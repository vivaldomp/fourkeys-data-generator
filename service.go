package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type eventsService struct {
	driver EventsDriver
}

type EventsService interface {
	Generate(eventTimespan int, numEvents int, numIssues int) (int, error)
}

func NewEventsService(driver EventsDriver) EventsService {
	return &eventsService{
		driver,
	}
}

func (e *eventsService) Generate(eventTimespan int, numEvents int, numIssues int) (int, error) {
	gitDeploymentIDMaxSize := math.Max(1000, math.Pow(10, math.Ceil(math.Log10(float64(numEvents)))))
	indexes := make(map[int]bool)

	for i := 0; i < numEvents; i++ {
		var r int
		for {
			r = int(rand.Int63n(int64(gitDeploymentIDMaxSize)))
			if indexes[r] {
				continue
			}
			break
		}

		indexes[r] = true
	}

	gitlabDeployIDs := make([]int, 0, numEvents)
	for k := range indexes {
		gitlabDeployIDs = append(gitlabDeployIDs, k)
	}

	changesSent := 0
	allChangeSets := makeAllChangeSets(numEvents, eventTimespan)
	ctx := context.Background()

	for i := 0; i < numEvents; i++ {
		changeset := allChangeSets[i]
		deployID := gitlabDeployIDs[i]
		indChanges := makeIndChangesFromChangeset(changeset)
		for _, currChange := range indChanges {
			changesSent++
			e.driver.MakeEventRequest(ctx, &currChange, "push")
		}

		e.driver.MakeEventRequest(ctx, &changeset, "push")

		if RandBool() {
			deploy := createDeployEvent(changeset, deployID)
			e.driver.MakeEventRequest(ctx, &deploy, "deployment")
		}
	}
	changesetsWithIssues := append([]Event(nil), allChangeSets...)
	rand.Shuffle(
		len(changesetsWithIssues),
		func(i, j int) {
			changesetsWithIssues[i], changesetsWithIssues[j] = changesetsWithIssues[j], changesetsWithIssues[i]
		},
	)
	changesetsWithIssues = changesetsWithIssues[:numIssues]
	for _, changeset := range changesetsWithIssues {
		issue := makeIssue(changeset)
		e.driver.MakeEventRequest(ctx, &issue, "issues")
	}
	return changesSent, nil
}

func makeAllChangeSets(numEvents int, eventTimespan int) []Event {
	allChangeSets := make([]Event, 0, numEvents)
	prevChangeSha := TokenHex(20)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numEvents; i++ {
		numChanges := rand.Intn(5) + 1
		changeSet := makeChanges(numChanges, eventTimespan, prevChangeSha)
		prevChangeSha = changeSet.CheckoutSHA
		allChangeSets = append(allChangeSets, changeSet)
	}

	return allChangeSets
}

func makeChanges(numChanges int, eventTimespan int, before string) Event {
	changes := make([]Change, 0, numChanges)
	maxTime := time.Now().Add(time.Duration(-eventTimespan) * time.Second)
	var headCommit Change
	if len(before) == 0 {
		before = TokenHex(20)
	}
	for i := 0; i < numChanges; i++ {
		changeId := TokenHex(20)
		unixTimeStamp := time.Now().Add(time.Duration(-rand.Intn(eventTimespan)) * time.Second)
		change := Change{
			ID:        changeId,
			Timestamp: unixTimeStamp.Format(time.RFC3339),
		}
		if unixTimeStamp.After(maxTime) {
			maxTime = unixTimeStamp
			headCommit = change
		}
		changes = append(changes, change)
	}

	return Event{
		ObjectKind:  "push",
		Before:      before,
		CheckoutSHA: headCommit.ID,
		Commits:     changes,
	}
}

func makeIndChangesFromChangeset(changeset Event) []Event {
	indChanges := make([]Event, 0, len(changeset.Commits))
	changesetSHA := changeset.CheckoutSHA
	prevChangeSHA := "0000000000000000000000000000000000000000"
	for _, c := range changeset.Commits {
		if c.ID != changesetSHA {
			currChange := Event{
				ObjectKind:  "push",
				Before:      prevChangeSHA,
				CheckoutSHA: c.ID,
				Commits:     []Change{c},
			}
			prevChangeSHA = c.ID
			indChanges = append(indChanges, currChange)
		}
	}
	return indChanges
}

func createDeployEvent(changes Event, deployID int) Event {
	var deployment Event
	checkoutSHA := changes.CheckoutSHA
	for _, c := range changes.Commits {
		if c.ID == checkoutSHA {
			deployment = Event{
				ObjectKind:      "deployment",
				Status:          "success",
				StatusChangedAt: c.Timestamp,
				DeploymentID:    fmt.Sprint(deployID),
				CommitURL:       fmt.Sprintf("http://example.com/root/test/commit/%s", checkoutSHA),
			}
			break
		}
	}
	return deployment
}

func makeIssue(changes Event) Event {
	var issue Event
	checkoutSHA := changes.CheckoutSHA
	for _, c := range changes.Commits {
		if c.ID == checkoutSHA {
			issue = Event{
				ObjectKind: "issue",
				ObjectAttributes: ObjectAttributes{
					CreatedAt: c.Timestamp,
					UpdatedAt: time.Now().Format(time.RFC3339),
					ClosedAt:  time.Now().Format(time.RFC3339),
					ID:        fmt.Sprint(rand.Intn(1000) + 1),
					Labels: []Label{
						{
							Title: "Incident",
						},
					},
					Description: fmt.Sprintf("root cause: %s", c.ID),
				},
			}
			break
		}
	}
	return issue
}
