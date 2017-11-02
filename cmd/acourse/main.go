package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acoshift/configfile"
	_ "github.com/lib/pq"

	"github.com/acoshift/acourse/pkg/app"
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

	go func() {
		log.Printf("acourse: start server at %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("acourse: server shutdown error: %v\n", err)
		return
	}
	log.Println("acourse: server shutdown")
}
