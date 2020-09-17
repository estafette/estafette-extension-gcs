package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {

	t.Run("ReturnsFalseIfSourceIsRootKeyFileJson", func(t *testing.T) {

		allowedBuckets := []string{"bucket"}
		params := Params{
			Source: "/key-file.json",
		}

		// act
		valid, errors := params.Validate(allowedBuckets)

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfSourceIsRoot", func(t *testing.T) {

		allowedBuckets := []string{"bucket"}
		params := Params{
			Source: "/",
		}

		// act
		valid, errors := params.Validate(allowedBuckets)

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfSourceIsWorkDir", func(t *testing.T) {

		allowedBuckets := []string{"bucket"}
		params := Params{
			Bucket: "bucket",
			Source: "/estafette-work",
		}

		// act
		valid, errors := params.Validate(allowedBuckets)

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfBucketNotAllowed", func(t *testing.T) {

		allowedBuckets := []string{"bucket"}
		params := Params{
			Bucket: "notbucket",
			Source: "/estafette-work",
		}

		// act
		valid, errors := params.Validate(allowedBuckets)

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAllowedBucketsEmpty", func(t *testing.T) {

		allowedBuckets := []string{}
		params := Params{
			Bucket: "notbucket",
			Source: "/estafette-work",
		}

		// act
		valid, errors := params.Validate(allowedBuckets)

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})
}
