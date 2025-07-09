package main

import (
	"os"

	"faissal.com/blogSpace/internal/db"
	"faissal.com/blogSpace/internal/repository"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file", "msg", err)
	}

	dbConfig := DBConf{
		Addr:        os.Getenv("DB_ADDR"),
		MaxOpenConn: 30,
		MaxIdleConn: 30,
		MaxIdleTime: "15m",
	}

	db, err := db.New(dbConfig.Addr, dbConfig.MaxOpenConn, dbConfig.MaxIdleConn, dbConfig.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Info("database connection pool has established")

	// TODO: Inject to service layer
	repository := repository.NewRepostory(db)

	application := Application{
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
		Env:  os.Getenv("ENV"),

		// what if this is not used in app dependency ?
		DbConfig: dbConfig,
	}

	mux := application.Mux()

	err = application.Run(mux)
	if err != nil {
		log.Fatal(err)
	}

}
