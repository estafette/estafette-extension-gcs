package main

import (
	"fmt"
	"path/filepath"
	"strings"
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
func (p *Params) Validate() (bool, []error) {

	errors := []error{}

	sourcePath, err := validPath(p.Source)
	if err != nil {
		errors = append(errors, err)
	}

	destPath, err := validPath(p.Destination)
	if err != nil {
		errors = append(errors, err)
	}

	if isGcsPath(sourcePath) && isGcsPath(destPath) {
		errors = append(errors, fmt.Errorf("Either Source or Destination might be a GCS path, not both"))
	}

	if sourcePath == "/key-file.json" || sourcePath == "/" {
		errors = append(errors, fmt.Errorf("Source '%v' is not allowed; use a source inside the working directory", p.Source))
	}

	return len(errors) == 0, errors
}

func isGcsPath(path string) bool {
	return strings.HasPrefix(path, "gs://")
}

func validPath(path string) (string, error) {
	if isGcsPath(path) {
		return path, nil
	}
	return filepath.Abs(path)
}
