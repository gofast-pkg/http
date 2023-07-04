package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		client := NewClient()
		assert.NotNil(t, client)
	})
}
