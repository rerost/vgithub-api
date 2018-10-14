package dgraph_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rerost/vgithub-api/infra/dgraph"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	_, err := dgraph.NewClient()
	if err != nil {
		t.Errorf("Can not create dgraph client: %v", err)
	}

	// TODO(@rerost) Implement other test.
}

func TestClose(t *testing.T) {
	t.Parallel()
	client, err := dgraph.NewClient()
	if err != nil {
		t.Errorf("Can not create dgraph client: %v", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Can not close dgraph client: %v", err)
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	inOutPairs := []struct {
		name string
		in   map[string]interface{}
		out  string
	}{
		{
			name: "test with empty",
			in: map[string]interface{}{
				"name": "test",
			},
		},
	}

	for _, inOutPair := range inOutPairs {
		t.Run(inOutPair.name, func(t *testing.T) {
			t.Parallel()

			client, err := dgraph.NewClient()
			if err != nil {
				t.Errorf("Can not create dgraph client: %v", err)
			}

			defer client.Close()

			ctx := context.Background()
			fmt.Println(inOutPair.in)
			body, err := json.Marshal(inOutPair.in)
			if err != nil {
				t.Error("Failed to marshal")
			}
			fmt.Println(string(body))
			_, err = client.Insert(ctx, body)
			if err != nil {
				t.Errorf("Failed to insert %v", err)
			}
		})
	}
}
