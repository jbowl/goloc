package main

//https://github.com/golang-standards/project-layout

import (
	"context"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jbowl/goloc/internal/pkg/geoloc"
	"github.com/jbowl/goloc/internal/pkg/postgres"
	"github.com/pkg/errors" // TODO research this
)

// sudo su - postgres
// psql -U postgres -d loc_db

// https://blog.cloudflare.com/exposing-go-on-the-internet/
// fuser 8080/tcp
// fuser -k 8080/tcp
func main() {
	key := os.Getenv("MQ_CONSUMER_KEY")
	ls := &postgres.Locator{Db: nil, Mq: &geoloc.MqAPI{Consumerkey: key}}

	db, err := ls.OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	ls.Db = db

	err = ls.CreateUsersTable()
	if err != nil {
		log.Fatal(err)
	}

	err = ls.CreateLocationsTable()
	if err != nil {
		log.Fatal(err)
	}
	err = ls.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	SetAPI(ls)
	router := NewRouter()

	/*

		tlsConfig := &tls.Config{
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
	*/
	server := http.Server{Addr: "localhost:8080",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		//		TLSConfig:    tlsConfig,
		Handler: router}

	errors.New("test")

	go func() {
		//		log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
		log.Fatal(server.ListenAndServe())
	}()

	<-sigChannel

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	server.Shutdown(ctx)
}
