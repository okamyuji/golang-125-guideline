package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestDeadlineExceeded(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(100*time.Millisecond))
		defer cancel()

		err := doWithDeadline(ctx, 200*time.Millisecond)
		synctest.Wait() // バブル内の待ちを収束させる

		if err != context.DeadlineExceeded {
			t.Fatalf("want DeadlineExceeded, got %v", err)
		}
	})
}
