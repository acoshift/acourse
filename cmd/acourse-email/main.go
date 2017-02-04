package main

import (
	"context"
	"log"
	"net"
	"os"

	"cloud.google.com/go/logging"
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/service/email"
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

	tokenSource, err := app.MakeTokenSource([]byte(cfg.ServiceAccount))
	if err != nil {
		log.Fatal(err)
	}

	loggerClient, err := logging.NewClient(context.Background(), cfg.ProjectID, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatal(err)
	}

	if !cfg.Debug {
		app.SetLogger(loggerClient.Logger("acourse-email"))
	}

	lis, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	creds, err := credentials.NewServerTLSFromFile(cfg.TLSCert, cfg.TLSKey)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(app.UnaryInterceptors), grpc.Creds(creds))
	acourse.RegisterEmailServiceServer(s, email.New(email.Config{
		From:     cfg.Email.From,
		Server:   cfg.Email.Server,
		Port:     cfg.Email.Port,
		User:     cfg.Email.User,
		Password: cfg.Email.Password,
	}))
	log.Fatal(s.Serve(lis))
}
