package main

//https://github.com/golang-standards/project-layout

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"sync"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jbowl/goloc/internal/pkg/geoloc"
	"github.com/jbowl/goloc/internal/pkg/goloc"
	"github.com/jbowl/goloc/internal/pkg/postgres"

	"github.com/jbowl/goloc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/jbowl/goloc/internal/pkg/mongoloc"

	"github.com/go-redis/redis"

	// mongo
	//	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors" // TODO research this
)

// sudo su - postgres
// psql -U postgres -d loc_db

// fuser 8080/tcp
// fuser -k 8080/tcp

var client *redis.Client

// NewTLSConfig - // https://blog.cloudflare.com/exposing-go-on-the-internet/
func NewTLSConfig() *tls.Config {
	return &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			//tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			//tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

}

// NewServer -
func NewServer(Addr string,
	TLSConfig *tls.Config,
	Handler http.Handler) *http.Server {
	return &http.Server{Addr: Addr,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    TLSConfig,
		Handler:      Handler}
}

// InterfaceStartup -
func InterfaceStartup(ls goloc.Locator) error {

	err := ls.CreateUsersTable()
	if err != nil {
		return err
	}

	err = ls.CreateLocationsTable()
	if err != nil {
		return err
	}
	err = ls.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// move some env strings to redis
	//client = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	db := flag.String("database", "postgresql2", "[postgresql | mongo]")
	flag.Parse()

	key := os.Getenv("MQ_KEY") // mapquest api key

	var ls goloc.Locator

	if *db == "postgresql" {
		ls = &postgres.Locator{Db: nil, Mq: &geoloc.MqAPI{Consumerkey: key}}
	} else {
		ls = &mongoloc.Locator{Client: nil, Mq: &geoloc.MqAPI{Consumerkey: key}}
	}

	err := ls.OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer ls.Close()

	err = InterfaceStartup(ls)
	if err != nil {
		log.Fatal(err)
	}

	SetAPI(ls)
	server := NewServer("localhost:8080", NewTLSConfig(), NewRouter())

	errors.New("test") // study this

	////////////////////// GRPC

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	s := api.Server{Locapi: ls}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile("./cert/server.crt", "./cert/server.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)
	api.RegisterGolocServer(grpcServer, &s)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	var waitgroup sync.WaitGroup

	go func() {
		waitgroup.Add(1)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	////////////////////// GRPC
	go func() {
		waitgroup.Add(1)
		//	log.Fatal(server.ListenAndServeTLS("./cert/server.crt", "./cert/server.key"))

		log.Fatal(server.ListenAndServe())
	}()

	<-sigChannel // kill signal

	waitgroup.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
