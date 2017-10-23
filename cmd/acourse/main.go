package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/configfile"
)

func main() {
	time.Local = time.UTC

	config := configfile.NewReader("config")

	err := app.Init(app.Config{
		ProjectID:      config.String("project_id"),
		ServiceAccount: config.Bytes("service_account"),
		BucketName:     config.String("bucket"),
		EmailServer:    config.String("email_server"),
		EmailPort:      config.Int("email_port"),
		EmailUser:      config.String("email_user"),
		EmailPassword:  config.String("email_password"),
		EmailFrom:      config.String("email_from"),
		BaseURL:        config.String("base_url"),
		XSRFSecret:     config.String("xsrf_key"),
		SQLURL:         config.String("sql_url"),
		RedisAddr:      config.String("redis_addr"),
		RedisPass:      config.String("redis_pass"),
		RedisPrefix:    config.String("redis_prefix"),
		SessionSecret:  config.Bytes("session_secret"),
		SlackURL:       config.String("slack_url"),
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	mux.Handle("/", app.Handler())

	// lets reverse proxy handle other settings
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Start server at :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
