package dgraph_test

import (
	"context"
	"encoding/json"
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
			name: "test with simple",
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
			body, err := json.Marshal(inOutPair.in)
			if err != nil {
				t.Error("Failed to marshal")
			}
			_, err = client.Insert(ctx, body)
			if err != nil {
				t.Errorf("Failed to insert %v", err)
			}
		})
	}
}

func TestQueryWithValues(t *testing.T) {
	t.Parallel()

	var uid string
	// Before test
	{
		ctx := context.Background()
		client, _ := dgraph.NewClient()
		body, _ := json.Marshal(map[string]interface{}{"name": "test"})
		res, _ := client.Insert(ctx, body)
		uid = res.Uids["blank-0"]
	}
	inOutPairs := []struct {
		name string
		in   map[string]interface{}
		out  string
	}{
		{
			name: "test",
			in: map[string]interface{}{
				"query": `query Me($id: string){
					me(func: uid($id)) {
						name
					}
				}`,
				"values": map[string]string{"$id": uid},
			},
			out: `{"me":[{"name":"test"}]}`,
		},
	}

	for _, inOutPair := range inOutPairs {
		t.Run(inOutPair.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			client, err := dgraph.NewClient()
			if err != nil {
				t.Error(err)
			}

			res, err := client.QueryWithValues(ctx, inOutPair.in["query"].(string), inOutPair.in["values"].(map[string]string))
			if err != nil {
				t.Error(err)
			}

			resBody := res.GetJson()
			if string(resBody) != inOutPair.out {
				t.Errorf(`
					want: %v
					have: %v
				`, inOutPair.out, string(resBody))
			}
		})
	}
}
