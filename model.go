package main

type Change struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

type Label struct {
	Title string `json:"title,omitempty"`
}

type ObjectAttributes struct {
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
	ClosedAt    string  `json:"closed_at,omitempty"`
	ID          string  `json:"id,omitempty"`
	Labels      []Label `json:"labels,omitempty"`
	Description string  `json:"description,omitempty"`
}

type Event struct {
	ObjectKind       string           `json:"object_kind"`
	Before           string           `json:"before,omitempty"`
	CheckoutSHA      string           `json:"checkout_sha,omitempty"`
	Commits          []Change         `json:"commits,omitempty"`
	Status           string           `json:"status,omitempty"`
	StatusChangedAt  string           `json:"status_changed_at,omitempty"`
	DeploymentID     string           `json:"deployment_id,omitempty"`
	CommitURL        string           `json:"commit_url,omitempty"`
	ObjectAttributes ObjectAttributes `json:"object_attributes,omitempty"`
}
