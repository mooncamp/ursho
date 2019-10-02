package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/douglasmakey/ursho/encoding"
	"github.com/douglasmakey/ursho/storage"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
	"mooncamp.com/dgx/gql"
	"mooncamp.com/dgx/render"
)

func New(host, port string, coder encoding.Coder) (storage.Service, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	schema := `
type ursho {
  ursho.createdAt:dateTime
  ursho.modifiedAt:dateTime
  ursho.url:string
  ursho.visited:bool
  ursho.count:int
}
`

	if err := dg.Alter(context.Background(), &api.Operation{Schema: schema}); err != nil {
		return nil, err
	}

	return &dgraph{
		ClientConn: conn,
		Dgraph:     dg,
		Coder:      coder,
	}, nil
}

type dgraph struct {
	ClientConn *grpc.ClientConn
	Dgraph     *dgo.Dgraph
	Coder      encoding.Coder
}

type dgraphItem struct {
	ID         string     `json:"uid"`
	CreatedAt  *time.Time `json:"ursho.createdAt,omitempty"`
	ModifiedAt *time.Time `json:"ursho.modifiedAt,omitempty"`

	URL     string `json:"ursho.url,omitempty"`
	Visited bool   `json:"ursho.visited,omitempty"`
	Count   int    `json:"ursho.count,omitempty"`

	Type string `json:"dgraph.type"`
}

func (d *dgraph) Save(url string) (string, error) {
	now := time.Now()
	item := dgraphItem{
		ID: "_:blank-0",

		CreatedAt:  &now,
		ModifiedAt: &now,

		URL:     url,
		Visited: false,
		Count:   0,

		Type: "ursho",
	}

	js, err := json.Marshal(item)
	if err != nil {
		return "", nil
	}

	assigned, err := d.Dgraph.NewTxn().Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: js})
	if err != nil {
		return "", err
	}

	id, err := strconv.ParseInt(assigned.Uids["blank-0"], 0, 64)
	if err != nil {
		return "", nil
	}

	return d.Coder.Encode(id)
}

func (d *dgraph) Load(code string) (string, error) {
	txn := d.Dgraph.NewTxn()
	defer txn.Discard(context.Background())

	item, err := d.loadInfo(txn, code)
	if err != nil {
		return "", err
	}

	now := time.Now()
	js, err := json.Marshal(dgraphItem{ID: item.ID, Visited: true, Count: item.Count + 1, ModifiedAt: &now, Type: "ursho"})
	if err != nil {
		return "", err
	}

	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: js})
	if err != nil {
		return "", err
	}

	return item.URL, txn.Commit(context.Background())
}

func (d *dgraph) loadInfo(txn *dgo.Txn, code string) (dgraphItem, error) {
	id, err := d.Coder.Decode(code)
	if err != nil {
		return dgraphItem{}, err
	}

	q := gql.GraphQuery{
		Alias: "items",
		UID:   []uint64{uint64(id)},
		Func:  &gql.Function{Name: "uid"},
		Children: []gql.GraphQuery{
			{Attr: "uid"},
			{Attr: "ursho.url"},
			{Attr: "ursho.count"},
			{Attr: "ursho.visited"},
		},
	}

	query, err := render.Render(render.Query{Queries: []gql.GraphQuery{q}})
	if err != nil {
		return dgraphItem{}, err
	}

	resp, err := txn.Query(context.Background(), query)
	if err != nil {
		return dgraphItem{}, err
	}

	var items map[string][]dgraphItem
	if err := json.Unmarshal(resp.Json, &items); err != nil {
		return dgraphItem{}, err
	}

	if len(items["items"]) == 0 {
		return dgraphItem{}, fmt.Errorf("no such url: %s", code)
	}

	return items["items"][0], nil
}

func (d *dgraph) LoadInfo(code string) (*storage.Item, error) {
	txn := d.Dgraph.NewTxn()
	defer txn.Discard(context.Background())

	item, err := d.loadInfo(txn, code)
	if err != nil {
		return nil, err
	}

	return &storage.Item{URL: item.URL, Visited: item.Visited, Count: item.Count}, txn.Commit(context.Background())
}

func (d *dgraph) Close() error {
	return d.ClientConn.Close()
}
