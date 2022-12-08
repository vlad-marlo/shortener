package filebased

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_Ping(t *testing.T) {
	filename := "file"
	defer assert.NoError(t, os.Remove(filename))
	store, err := New(filename)
	require.NoError(t, err, fmt.Sprintf("init storage: %v", err))
	err = store.Ping(context.Background())
	require.NoError(t, err, fmt.Sprintf("ping: %v", err))
}

func TestStore_Close(t *testing.T) {
	filename := "file"
	defer assert.NoError(t, os.Remove(filename))
	store, err := New(filename)
	require.NoError(t, err, fmt.Sprintf("init storage: %v", err))
	err = store.Close()
	require.NoError(t, err, fmt.Sprintf("ping: %v", err))
}
