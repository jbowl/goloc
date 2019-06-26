package main

import (
	"context"
	//	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	//	_ "github.com/lib/pq" //
	"net/http"

	"github.com/jbowl/gopostgresql1/mqapi"
)

// sudo su - postgres
// psql -U postgres -d loc_db

func dsn() *string {

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"loc_db")

	//	os.Getenv("DB_HOST"),
	//	os.Getenv("DB_PORT"),
	//	os.Getenv("DB_USERNAME"),
	//	os.Getenv("DB_PASSWORD"),
	//	os.Getenv("DB_NAME"))
	return &dsn
}

//Location define table for GORM
type Location struct {
	gorm.Model
	Date    time.Time
	Address string
	Lat     float32
	Lng     float32
}

func init() {
	db, err := gorm.Open("postgres", *dsn())

	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&Location{})
}

type geoDB struct {
	dbs *gorm.DB
	mq  *mqapi.MqAPI
}

// fuser 8080/tcp
// fuser -k 8080/tcp
func main() {

	db, err := gorm.Open("postgres", *dsn())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	key := os.Getenv("MQ_CONSUMER_KEY")

	router := newRouter(&geoDB{dbs: db, mq: &mqapi.MqAPI{Consumerkey: key}})

	server := http.Server{Addr: "localhost:8080", Handler: router}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-sigChannel

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	server.Shutdown(ctx)
}
