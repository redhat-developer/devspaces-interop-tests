package workspaces

// Workspace object obtained from che server
type Workspace struct {
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
	Status     string     `json:"status"`
}

// Workspace attributes
type Attributes struct {
	InfrastructureNamespace string `json:"infrastructureNamespace"`
}
