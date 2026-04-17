package handler_test

import (
	"testing"

	apihandler "github.com/gaston-garcia-cegid/gonsgarage/internal/handler"
	"github.com/stretchr/testify/require"
)

func TestNewAuthHandler_nonNil(t *testing.T) {
	h := apihandler.NewAuthHandler(nil)
	require.NotNil(t, h)
}

func TestNewUserHandler_nonNil(t *testing.T) {
	h := apihandler.NewUserHandler()
	require.NotNil(t, h)
}
