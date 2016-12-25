package main

import (
	"acourse/store"
	"log"
)

func main() {
	db := store.NewDB(store.ProjectID("acourse-d9d0a"))

	courses, err := db.CourseList()
	if err != nil {
		panic(err)
	}

	for _, course := range courses {
		err = db.CourseSave(course)
		if err != nil {
			panic(err)
		}
	}

	log.Println("Completed")
}
