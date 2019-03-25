package main

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	// control params
	Action      string `json:"action,omitempty"`
	Bucket      string `json:"bucket,omitempty"`
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
}
