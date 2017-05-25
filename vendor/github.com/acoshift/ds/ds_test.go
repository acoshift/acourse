package ds

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"math/rand"
	"strconv"

	"cloud.google.com/go/datastore"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type ExampleModel struct {
	Model
	StampModel
	Name  string
	Value int
}

var kind string

func init() {
	rand.Seed(time.Now().Unix())
	kind = "Test_" + strconv.Itoa(rand.Int())
}

func (x *ExampleModel) NewKey() {
	x.NewIncomplateKey(kind, nil)
}

type ExampleNotModel struct {
	Name string
}

var ctx = context.Background()

func skipShort(t *testing.T, name string) {
	if testing.Short() {
		t.Skip("skipping", name)
	}
}

func initClient() (*Client, error) {
	// load service account from env
	serviceAccountStr := os.Getenv("service_account")
	opts := []option.ClientOption{}
	if serviceAccountStr != "" {
		serviceAccount, err := base64.StdEncoding.DecodeString(serviceAccountStr)
		if err != nil {
			return nil, err
		}
		cfg, err := google.JWTConfigFromJSON(serviceAccount, datastore.ScopeDatastore)
		if err != nil {
			return nil, err
		}
		opts = append(opts, option.WithTokenSource(cfg.TokenSource(ctx)))
	}
	projectID := os.Getenv("project_id")
	if projectID == "" {
		projectID = "acoshift-test"
	}
	return NewClient(ctx, projectID, opts...)
}

func TestInvalidNewClient(t *testing.T) {
	client, err := NewClient(ctx, "invalid-project-id", option.WithServiceAccountFile("invalid-file"))
	if err == nil {
		t.Fatalf("expected error not nil")
	}
	if client != nil {
		t.Fatalf("expected client to be nil")
	}
}

func prepareData(client *Client) []*datastore.Key {
	xs := []*ExampleModel{
		&ExampleModel{Name: "name1", Value: 1},
		&ExampleModel{Name: "name2", Value: 2},
		&ExampleModel{Name: "name3", Value: 3},
		&ExampleModel{Name: "name4", Value: 4},
		&ExampleModel{Name: "name5", Value: 5},
		&ExampleModel{Name: "name6", Value: 6},
		&ExampleModel{Name: "name7", Value: 7},
	}
	client.SaveModels(ctx, xs)
	return ExtractKeys(xs)
}

func removeData(client *Client) {
	keys, _ := client.QueryKeys(ctx, kind)
	client.DeleteMulti(ctx, keys)
}
