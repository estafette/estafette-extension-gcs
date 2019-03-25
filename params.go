package main

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	// control params
	Action      string `json:"action,omitempty"`
	Bucket      string `json:"bucket,omitempty"`
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	ACL         string `json:"acl,omitempty"`
	Compress    *bool  `json:"compress,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults() {
	if p.Compress == nil {
		trueValue := true
		p.Compress = &trueValue
	}
}
