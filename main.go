package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/motoki317/bot-extreme/handler"
	"github.com/motoki317/bot-extreme/repository"
	bot "github.com/traPtitech/traq-bot"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	dbInitDirectory = "./mysql/init"
)

var (
	port = os.Getenv("PORT")
)

func main() {
	log.SetFlags(log.LstdFlags)
	if port == "" {
		log.Println("Setting default port to 80")
		port = "80"
	}

	// connect to db
	db := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s",
		os.Getenv("MARIADB_USERNAME"),
		os.Getenv("MARIADB_PASSWORD"),
		os.Getenv("MARIADB_HOSTNAME"),
		os.Getenv("MARIADB_DATABASE"),
	))
	// db connection for batch executing, allowing multi statements
	dbForBatch := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?multiStatements=true",
		os.Getenv("MARIADB_USERNAME"),
		os.Getenv("MARIADB_PASSWORD"),
		os.Getenv("MARIADB_HOSTNAME"),
		os.Getenv("MARIADB_DATABASE"),
	))

	// create schema
	var paths []string
	err := filepath.Walk(dbInitDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		log.Printf("Executing file %s...", path)
		dbForBatch.MustExec(string(data))
	}
	log.Println("Successfully initialized DB schema!")

	// repository impl
	repo := repository.NewRepositoryImpl(db)

	// traq bot handlers
	handlers := bot.EventHandlers{}
	handlers.SetMessageCreatedHandler(handler.MessageReceived(repo))

	// traq bot server
	vt := os.Getenv("VERIFICATION_TOKEN")
	server := bot.NewBotServer(vt, handlers)
	log.Println("Listening on port " + port + "...")
	log.Fatal(server.ListenAndServe(":" + port))
}
