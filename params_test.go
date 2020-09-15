package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {

	t.Run("ReturnsFalseIfSourceIsRootKeyFileJson", func(t *testing.T) {

		params := Params{
			Source: "/key-file.json",
		}

		// act
		valid, errors := params.Validate()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsFalseIfSourceIsRoot", func(t *testing.T) {

		params := Params{
			Source: "/",
		}

		// act
		valid, errors := params.Validate()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfSourceIsWorkDir", func(t *testing.T) {

		params := Params{
			Source: "/estafette-work",
		}

		// act
		valid, errors := params.Validate()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfSourceDestinationAreGcs", func(t *testing.T) {

		params := Params{
			Source:      "gs://estafette-work",
			Destination: "gs://estafette-work2",
		}

		// act
		valid, errors := params.Validate()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})
}
