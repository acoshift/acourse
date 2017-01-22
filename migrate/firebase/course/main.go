package main

import (
	"encoding/json"
	"log"
	"time"

	"context"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
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

	// db.CoursePurge()
	// log.Println("Purged Courses")
	// db.EnrollPurge()
	// log.Println("Purged Enrolls")

	client, _ := google.DefaultClient(context.Background())
	resp, _ := client.Get("https://acourse-d9d0a.firebaseio.com/course.json")
	courses := map[string]firCourse{}
	json.NewDecoder(resp.Body).Decode(&courses)
	resp.Body.Close()

	resp, _ = client.Get("https://acourse-d9d0a.firebaseio.com/content.json")
	contents := map[string][]firContent{}
	json.NewDecoder(resp.Body).Decode(&contents)
	resp.Body.Close()

	es := []*model.Enroll{}
	for k, v := range courses {
		start, _ := time.Parse("2006-01-02", v.Start)

		x := model.Course{
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

		// contents
		cs := contents[k]
		if len(cs) > 0 {
			x.Contents = make([]model.CourseContent, len(cs))
			for i, c := range cs {
				x.Contents[i] = model.CourseContent{
					Title:       c.Title,
					Description: c.Content,
				}
			}
		}

		// err := db.CourseSave(&x)
		// log.Println("Migrated Course:", x.Title)
		// if err != nil {
		// 	log.Println(err)
		// }

		for uid, s := range v.Student {
			if s {
				es = append(es, &model.Enroll{UserID: uid, CourseID: mapCourseID(k)})
			}
		}
	}
	log.Println("Migrated Courses")

	// db.EnrollCreateAll(es[:400])
	// db.EnrollCreateAll(es[400:])
	log.Println(len(es))
	for _, e := range es {
		x, _ := db.EnrollFind(e.UserID, e.CourseID)
		if x == nil {
			log.Println("Not found, " + e.UserID + " " + e.CourseID)
		}
	}
	log.Println("Migrated Enrolls")

	log.Println("Completed")
}

func mapCourseID(id string) string {
	switch id {
	case "-KS2ivTwZpPIQiV3Sqm-":
		return "5751646390845440"
	case "-KS6g5qpLq8I0MPdF7oc":
		return "5701751084679168"
	case "-KUNcswtqfZHIAzWHccQ":
		return "5644101080842240"
	case "-KYnk-y_6KoFDSsOVJex":
		return "5671617594130432"
	case "-KZf-Y5eG97WlFmv87g1":
		return "5766596232478720"
	default:
		return ""
	}
}
