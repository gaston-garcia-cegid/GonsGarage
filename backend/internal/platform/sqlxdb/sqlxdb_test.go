package sqlxdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen_rejectsEmptyDSN(t *testing.T) {
	_, err := Open("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty dsn")
}
