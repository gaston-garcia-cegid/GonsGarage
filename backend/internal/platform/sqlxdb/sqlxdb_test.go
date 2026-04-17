package sqlxdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapPostgres_nilDB(t *testing.T) {
	require.Nil(t, WrapPostgres(nil))
}
