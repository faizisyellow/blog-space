package main

import (
	"fmt"
	"net"
	"os"
	"time"

	_ "faissal.com/blogSpace/docs"
	"faissal.com/blogSpace/internal/auth"
	"faissal.com/blogSpace/internal/db"
	"faissal.com/blogSpace/internal/repository"
	"faissal.com/blogSpace/internal/services"
	"faissal.com/blogSpace/internal/uploader"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

//	@title			Blog Space Rest API
//	@version		1.0
//	@description	Rest API Documentation for Blog Space Services.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						Authorization

// @schemes	http https
// @host		localhost:8080
// @BasePath	/v1
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file", "msg", err)
	}

	dbConfig := DBConf{
		Addr: os.Getenv("DB_ADDR"),

		MaxOpenConn: 30,

		MaxIdleConn: 30,

		MaxIdleTime: "15m",
	}

	dbs, err := db.New(dbConfig.Addr, dbConfig.MaxOpenConn, dbConfig.MaxIdleConn, dbConfig.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer dbs.Close()

	log.Info("database connection pool has established")

	repository := repository.NewRepostory(dbs)

	services := services.NewServices(*repository, db.WithTx, dbs)

	jwtTokenConfig := JwtConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
		Iss:       "authentication",
		Sub:       "user",
		Exp:       time.Now().Add(time.Hour * 24 * 3).Unix(),
	}

	jwtAuthentication := auth.New(jwtTokenConfig.SecretKey, jwtTokenConfig.Iss, jwtTokenConfig.Sub)

	r2conf := R2Conf{
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
		AccountId:       os.Getenv("R2_ACCOUNT_ID"),
		AccessKeyId:     os.Getenv("R2_ACCESS_KEY"),
		AccessKeySecret: os.Getenv("R2_ACCESS_KEY_SECRET"),
	}

	r2Store := uploader.NewR2Client(r2conf.BucketName, r2conf.AccountId, r2conf.AccessKeyId, r2conf.AccessKeySecret)

	application := Application{
		Port: os.Getenv("PORT"),

		Host: os.Getenv("HOST"),

		Env: os.Getenv("ENV"),

		DbConfig: dbConfig,

		Services: *services,

		JwtAuth: jwtTokenConfig,

		Authentication: jwtAuthentication,

		//http:domain:port/version/swagger/*
		SwaggerUrl: fmt.Sprintf("http://%v/v%v/swagger/doc.json", net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT")), 1),

		R2Config: r2conf,

		Uploading: r2Store,
	}

	mux := application.Mux()

	err = application.Run(mux)
	if err != nil {
		log.Fatal(err)
	}

}
