package polyreduce_test

import (
	"context"
	"testing"

	"changkun.de/x/infloop/tools/polyreduce-sdk-go"
)

func TestPolyreduce_Ping(t *testing.T) {
	c := polyreduce.NewClient()
	r, err := c.Ping(context.Background())
	if err != nil {
		t.Fatalf("failed to ping polyreduce.com: %v", err)
	}

	t.Log(r.Version)
}
