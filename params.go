package main

import (
	"fmt"
	"path/filepath"
)

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	// control params
	Action      string `json:"action,omitempty"`
	Bucket      string `json:"bucket,omitempty"`
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	ACL         string `json:"acl,omitempty"`
	Recursive   *bool  `json:"recursive,omitempty"`
	Compress    *bool  `json:"compress,omitempty"`
	Parallel    *bool  `json:"parallel,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults() {
	if p.Recursive == nil {
		trueValue := true
		p.Recursive = &trueValue
	}
	if p.Compress == nil {
		trueValue := true
		p.Compress = &trueValue
	}
	if p.Parallel == nil {
		trueValue := true
		p.Parallel = &trueValue
	}
}

// Validate checks whether all parameters have valid values
func (p *Params) Validate(allowedBuckets []string) (bool, []error) {

	errors := []error{}

	sourcePath, err := filepath.Abs(p.Source)
	if err != nil {
		errors = append(errors, err)
	}

	if sourcePath == "/key-file.json" || sourcePath == "/" {
		errors = append(errors, fmt.Errorf("Source '%v' is not allowed; use a source inside the working directory", p.Source))
	}

	allowedBucket := false
	for _, b := range allowedBuckets {
		if p.Bucket == b {
			allowedBucket = true
		}
	}
	if !allowedBucket {
		errors = append(errors, fmt.Errorf("Bucket '%v' is not allowed for this credential", p.Bucket))
	}

	return len(errors) == 0, errors
}
