package main

// GCSCredentials represents the credentials of type cloud-storage as defined in the server config and passed to this trusted image
type GCSCredentials struct {
	Name                 string                            `json:"name,omitempty"`
	Type                 string                            `json:"type,omitempty"`
	AdditionalProperties GCSCredentialAdditionalProperties `json:"additionalProperties,omitempty"`
}

// GCSCredentialAdditionalProperties contains the non standard fields for this type of credentials
type GCSCredentialAdditionalProperties struct {
	AllowedBuckets        []string `json:"allowedBuckets,omitempty"`
	Project               string   `json:"project,omitempty"`
	ServiceAccountKeyfile string   `json:"serviceAccountKeyfile,omitempty"`
}

// GetCredentialsByName returns a credential if the name exists
func GetCredentialsByName(c []GCSCredentials, credentialName string) *GCSCredentials {

	for _, cred := range c {
		if cred.Name == credentialName {
			return &cred
		}
	}

	return nil
}
