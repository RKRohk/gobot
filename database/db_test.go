package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/rkrohk/gobot/database"
)

func TestDatabase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.Client.Ping(ctx, nil); err != nil {
		t.Error("Connection failed", err)
	} else {
		t.Log("Connection Successful")
	}
}
