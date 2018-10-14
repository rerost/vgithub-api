package dgraph

import (
	"context"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// Dgraph TODO
type Dgraph interface {
	Close() error
	Insert(ctx context.Context, json []byte) (*api.Assigned, error)
	QueryWithValues(ctx context.Context, query string, values map[string]string) (*api.Response, error)
}

type dgraphImp struct {
	client *dgo.Dgraph
	conn   *grpc.ClientConn
}

// TODO(@rerost) Move to env.
const dgraphHost = "127.0.0.1:9080"

// NewClient is create dgraph client. It use grpc
func NewClient() (Dgraph, error) {
	conn, err := grpc.Dial(dgraphHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	apiClient := api.NewDgraphClient(conn)
	dgraphClient := dgo.NewDgraphClient(apiClient)

	return &dgraphImp{
		client: dgraphClient,
		conn:   conn,
	}, nil
}

func (d *dgraphImp) Close() error {
	err := d.conn.Close()
	if err != nil {
		return err
	}

	d.client = nil

	return nil
}

func (d *dgraphImp) Insert(ctx context.Context, body []byte) (*api.Assigned, error) {
	mutation := &api.Mutation{
		CommitNow: true,
	}
	mutation.SetJson = body
	tx := d.client.NewTxn()
	assigned, err := tx.Mutate(ctx, mutation)
	return assigned, err
}

func (d *dgraphImp) QueryWithValues(ctx context.Context, query string, values map[string]string) (*api.Response, error) {
	tx := d.client.NewTxn()
	res, err := tx.QueryWithVars(ctx, query, values)
	return res, err
}
