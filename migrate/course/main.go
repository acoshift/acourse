package main

import (
	"log"

	"github.com/acoshift/acourse/pkg/store"
)

func main() {
	db := store.NewDB(store.ProjectID("acourse-d9d0a"))

	xs, err := db.CourseList()
	if err != nil {
		panic(err)
	}

	for _, x := range xs {
		err = db.CourseSave(x)
		if err != nil {
			panic(err)
		}
	}

	log.Println("Completed")
}
