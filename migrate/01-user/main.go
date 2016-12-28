package main

import (
	"encoding/json"
	"log"

	"acourse/store"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

// migrate user from firebase to datastore

type firUser struct {
	Name    string `json:"name"`
	Photo   string `json:"photo"`
	AboutMe string `json:"aboutMe"`
}

func main() {
	db := store.NewDB(store.ProjectID("acourse-d9d0a"))

	// clean data
	db.UserPurge()
	log.Println("Purged Users")
	db.RolePurge()
	log.Println("Purged Roles")

	client, _ := google.DefaultClient(context.Background())
	resp, _ := client.Get("https://acourse-d9d0a.firebaseio.com/user.json")
	defer resp.Body.Close()
	users := map[string]firUser{}
	json.NewDecoder(resp.Body).Decode(&users)
	xs := make([]*model.User, len(users))
	uids := make([]string, len(users))
	i := 0
	for uid, u := range users {
		x := model.User{
			Name:     u.Name,
			Username: uid,
			Photo:    u.Photo,
			AboutMe:  u.AboutMe,
		}
		xs[i] = &x
		uids[i] = uid
		i++
	}
	err := db.UserCreateAll(uids[:400], xs[:400])
	if err != nil {
		log.Println(err)
	}
	err = db.UserCreateAll(uids[401:], xs[401:])
	if err != nil {
		log.Println(err)
	}
	log.Println("Migrated Users")

	// Instructor
	resp, _ = client.Get("https://acourse-d9d0a.firebaseio.com/instructor.json")
	roles := map[string]bool{}
	json.NewDecoder(resp.Body).Decode(&roles)
	resp.Body.Close()
	for uid, r := range roles {
		u, _ := db.UserGet(uid)
		x, _ := db.RoleFindByUserID(u.ID)
		x.Instructor = r
		err := db.RoleSave(x)
		log.Println(x)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Migrated Role Instructor")

	// Admin
	resp, _ = client.Get("https://acourse-d9d0a.firebaseio.com/admin.json")
	roles = map[string]bool{}
	json.NewDecoder(resp.Body).Decode(&roles)
	resp.Body.Close()
	for uid, r := range roles {
		u, _ := db.UserGet(uid)
		x, _ := db.RoleFindByUserID(u.ID)
		x.Admin = r
		err := db.RoleSave(x)
		log.Println(x)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Migrated Role Admin")

	log.Println("Completed")
}
