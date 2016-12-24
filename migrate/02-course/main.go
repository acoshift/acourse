package main

import (
	"acourse/store"
	"encoding/json"
	"log"

	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

// migrate course from firebase to datastore

type firCourse struct {
	Title             string          `json:"title"`
	Description       string          `json:"description"`
	ShortDescription  string          `json:"shortDescription"`
	Start             string          `json:"start"`
	Video             string          `json:"video"`
	Timestamp         int64           `json:"timestamp"`
	Favorite          map[string]bool `json:"favorite"`
	HasAssignment     bool            `json:"hasAssignment"`
	Open              bool            `json:"open"`
	CanAttend         bool            `json:"canAttend"`
	CanQueueEnroll    bool            `json:"canQueueEnroll"`
	Owner             string          `json:"owner"`
	Photo             string          `json:"photo"`
	QueueEnrollDetail string          `json:"queueEnrollDetail"`
	Public            bool            `json:"public"`
	Student           map[string]bool `json:"student"`
}

type firContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	db := store.NewDB(store.ProjectID("acourse-d9d0a"))

	db.CoursePurge()
	log.Println("Purged Courses")
	db.EnrollPurge()
	log.Println("Purged Enrolls")

	client, _ := google.DefaultClient(context.Background())
	resp, _ := client.Get("https://acourse-d9d0a.firebaseio.com/course.json")
	courses := map[string]firCourse{}
	json.NewDecoder(resp.Body).Decode(&courses)
	resp.Body.Close()

	resp, _ = client.Get("https://acourse-d9d0a.firebaseio.com/content.json")
	contents := map[string][]firContent{}
	json.NewDecoder(resp.Body).Decode(&contents)
	resp.Body.Close()

	es := []*store.Enroll{}
	for k, v := range courses {
		start, _ := time.Parse("2006-01-02", v.Start)

		x := store.Course{
			Title:            v.Title,
			ShortDescription: v.ShortDescription,
			Description:      v.Description,
			Photo:            v.Photo,
			Owner:            v.Owner,
			Start:            start,
		}
		x.Options.Public = v.Public
		x.Options.Enroll = v.Open
		x.Options.Attend = v.CanAttend
		x.Options.Assignment = v.HasAssignment
		x.Options.Purchase = v.CanQueueEnroll

		// contents
		cs := contents[k]
		if len(cs) > 0 {
			x.Contents = make([]store.CourseContent, len(cs))
			for i, c := range cs {
				x.Contents[i] = store.CourseContent{
					Title:       c.Title,
					Description: c.Content,
				}
			}
		}

		err := db.CourseSave(&x)
		log.Println("Migrated Course:", x.Title)
		if err != nil {
			log.Println(err)
		}

		for uid, s := range v.Student {
			if s {
				es = append(es, &store.Enroll{UserID: uid, CourseID: x.ID})
			}
		}
	}
	log.Println("Migrated Courses")

	db.EnrollCreateAll(es)
	log.Println("Migrated Enrolls")

	log.Println("Completed")
}
