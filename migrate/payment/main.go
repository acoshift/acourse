package main

import (
	"log"

	"github.com/acoshift/acourse/store"
)

func main() {
	db := store.NewDB(store.ProjectID("acourse-d9d0a"))

	xs, err := db.PaymentList()
	if err != nil {
		panic(err)
	}

	for _, x := range xs {
		log.Println(x.ID)
		err = db.PaymentSave(x)
		if err != nil {
			panic(err)
		}
	}

	log.Println("Completed")
}
