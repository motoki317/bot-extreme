package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	bot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"

	"github.com/motoki317/bot-extreme/evaluate"
	"github.com/motoki317/bot-extreme/handler"
	"github.com/motoki317/bot-extreme/repository"
)

const (
	dbInitDirectory = "./mysql/init"
)

var (
	accessToken = os.Getenv("ACCESS_TOKEN")
)

func main() {
	log.SetFlags(log.LstdFlags)

	// connect to db
	db := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("MARIADB_USERNAME"),
		os.Getenv("MARIADB_PASSWORD"),
		os.Getenv("MARIADB_HOSTNAME"),
		os.Getenv("MARIADB_DATABASE"),
	))
	// db connection for batch executing, allowing multi statements
	dbForBatch := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?multiStatements=true&parseTime=true",
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
	b, err := bot.NewBot(&bot.Options{
		AccessToken: accessToken,
	})
	b.OnMessageCreated(handler.MessageReceived(repo))
	b.OnStampCreated(func(p *payload.StampCreated) {
		evaluate.AddStamp(p.Name, p.ID)
	})

	// Start (blocks on success)
	if err := b.Start(); err != nil {
		panic(err)
	}
}
