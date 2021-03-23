package helpers_test

import (
	"context"
	"testing"
	"time"

	"github.com/rkrohk/gobot/helpers"
)

func TestDatabase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := helpers.Client.Ping(ctx, nil); err != nil {
		t.Error("Connection failed", err)
	} else {
		t.Log("Connection Successful")
	}
}
