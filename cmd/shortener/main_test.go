package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createLogger(t *testing.T) {
	logger, err := createLogger("xd")
	assert.NoError(t, err)
	require.NotNil(t, logger)
	err = logger.Sync()
	assert.Error(t, err)
}
