package main

import (
	"context"
	"log"
	"net"
	"os"

	"cloud.google.com/go/logging"
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/service/user"
	"github.com/acoshift/ds"
	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {
	grpclog.SetLogger(app.NewNoFatalLogger())

	configFile := os.Getenv("CONFIG")
	if configFile == "" {
		configFile = "config.yaml"
	}

	cfg, err := app.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	firApp, err := admin.InitializeApp(admin.AppOptions{
		ProjectID:      cfg.Firebase.ProjectID,
		ServiceAccount: []byte(cfg.Firebase.ServiceAccount),
	})
	if err != nil {
		log.Fatal(err)
	}
	firAuth, err := firApp.Auth()
	if err != nil {
		log.Fatal(err)
	}

	tokenSource, err := app.MakeTokenSource([]byte(cfg.ServiceAccount))
	if err != nil {
		log.Fatal(err)
	}

	client, err := ds.NewClient(context.Background(), cfg.ProjectID, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatal(err)
	}

	loggerClient, err := logging.NewClient(context.Background(), cfg.ProjectID, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatal(err)
	}

	if !cfg.Debug {
		app.SetLogger(loggerClient.Logger("acourse-user"))
	}

	app.InitService(firAuth)

	lis, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	creds, err := credentials.NewServerTLSFromFile(cfg.TLSCert, cfg.TLSKey)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(app.UnaryInterceptors), grpc.Creds(creds))
	acourse.RegisterUserServiceServer(s, user.New(client))
	log.Fatal(s.Serve(lis))
}
