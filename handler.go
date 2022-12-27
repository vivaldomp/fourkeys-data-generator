package main

import (
	log "github.com/sirupsen/logrus"
)

type generatorHandler struct {
	service EventsService
}

type GeneratorHandler interface {
	Run(int, int, int)
}

func NewGeneratorHandler(eventsService EventsService) GeneratorHandler {
	return &generatorHandler{
		eventsService,
	}
}

func (g *generatorHandler) Run(eventTimespan int, numEvents int, numIssues int) {
	changesSent, err := g.service.Generate(eventTimespan, numEvents, numIssues)
	if err != nil {
		log.Errorf("Error on generate events: %s", err)
		return
	}
	log.Infof("%d changes successfully sent to event-handler", changesSent)
}
