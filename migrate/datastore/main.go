package main

import (
	"context"
	"log"
	"reflect"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
)

var (
	clientSrc, clientDst *ds.Client
	ctx                  = context.Background()
)

// Migrate data betweet 2 datastore projects
func main() {
	var err error
	const projectSrc = "acourse-d9d0a"
	const projectDst = "acourse-156413"

	clientSrc, err = ds.NewClient(ctx, projectSrc)
	if err != nil {
		log.Fatal(err)
	}
	clientDst, err = ds.NewClient(ctx, projectDst)
	if err != nil {
		log.Fatal(err)
	}

	// kinds to migrate
	kinds := []struct {
		Type reflect.Type
		Name string
	}{
		{reflect.TypeOf(model.User{}), "User"},
		{reflect.TypeOf(model.Course{}), "Course"},
		{reflect.TypeOf(model.Role{}), "Role"},
		{reflect.TypeOf(model.Enroll{}), "Enroll"},
		{reflect.TypeOf(model.Payment{}), "Payment"},
		{reflect.TypeOf(model.Attend{}), "Attend"},
		{reflect.TypeOf(model.Assignment{}), "Assignment"},
		{reflect.TypeOf(model.UserAssignment{}), "UserAssignment"},
	}
	var xs interface{}
	for _, kind := range kinds {
		log.Println("Start migrate", kind.Name)
		xs = reflect.New(reflect.SliceOf(reflect.PtrTo(kind.Type))).Interface()
		err = clientSrc.Query(ctx, kind.Name, xs)
		err = ds.IgnoreFieldMismatch(err)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Number of %s: %d", kind.Name, reflect.ValueOf(xs).Elem().Len())

		migrate(reflect.ValueOf(xs).Elem().Interface())
		log.Println("Finish migrate", kind.Name)
	}
}

func migrate(xs interface{}) {
	xf := reflect.ValueOf(xs)
	if xf.Len() > 500 {
		migrate(xf.Slice(0, 500).Interface())
		migrate(xf.Slice(500, xf.Len()).Interface())
		return
	}
	log.Println("Migrate", xf.Len())
	err := clientDst.PutModels(ctx, xf.Interface())
	if err != nil {
		log.Fatal(err)
	}
}
