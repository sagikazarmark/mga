package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScaffold(t *testing.T) {
	expected := `package pkggen

// +kit:endpoint

// Service <insert your description>.
type Service interface {
	// Insert your operations here
}

// NewService returns a new Service.
func NewService() Service {
	return service{}
}

type service struct{}
`

	actual, err := Scaffold("pkggen")
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual), "the scaffolded code does not match the expected one")
}
