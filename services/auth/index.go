package main

import (
	"log"
	"os"
	"strconv"

	fastHttpRouter "github.com/buaazp/fasthttprouter"
	"github.com/joho/godotenv"
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libCrypto "github.com/nguoihanoi/golang_shared/libs/crypto"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	mdWare "github.com/nguoihanoi/golang_shared/libs/middleware"
	"github.com/redis/go-redis/v9"
	fastHttp "github.com/valyala/fasthttp"
)

var libJwt *libCrypto.JwtClass

func initDb(urlConnect string, dbName string) *libDb.DatabaseClass {
	mongoClient := libDb.MongoConnect(urlConnect)
	newDb := libDb.NewDatabase(mongoClient, dbName)
	return newDb
}

func initRedisDb(inAddr string, inDB int, inPassword string, inPrefix string) *libCache.Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     inAddr,     // Redis server address
		Password: inPassword, // No password by default
		DB:       inDB,       // Default DB
	})
	return libCache.NewCache(rdb, inPrefix)
}

func main() {
	//
	log.Println("init started")
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	//
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi load file .env:", err)
	}

	//Todo: Đọc các giá trị bằng os.Getenv
	// redis
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	REDIS_DB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	REDIS_PREFIX := os.Getenv("REDIS_PREFIX")
	// mongodb
	MONGO_COMMON_URI := os.Getenv("MONGO_COMMON_URI")
	MONGO_COMMON_DB := os.Getenv("MONGO_COMMON_DB")
	// server
	SERVER_NAME := os.Getenv("SERVER_NAME")
	PORT := os.Getenv("PORT")
	// crypto
	JWT_SECRET := os.Getenv("JWT_SECRET")
	//PREFIX_TOKEN := os.Getenv("PREFIX_TOKEN")

	//Todo: init db
	mgoDB1 := initDb(MONGO_COMMON_URI, MONGO_COMMON_DB)
	redisDb1 := initRedisDb(REDIS_HOST, REDIS_DB, REDIS_PASSWORD, REDIS_PREFIX)

	//Todo: init router
	mainRouter := fastHttpRouter.New()

	//Todo: init db
	libJwt = libCrypto.JWT(JWT_SECRET)
	Init(mainRouter, mgoDB1, redisDb1, JWT_SECRET)
	corsMid := mdWare.Init("*", "POST", JWT_SECRET)

	newToken, nextTime, err := libJwt.CreateToken(`{"email":"playhard24h@gmail.com","password":"abc123!@#"}`)
	log.Println(newToken, nextTime, err)

	//Todo: Create a custom logger
	myLogger := log.New(log.Writer(), "FasthttpServer: ", log.LstdFlags)
	fastHttpServer := &fastHttp.Server{
		Handler:            corsMid.CorsMiddleware(mainRouter.Handler),
		Logger:             myLogger,
		Name:               SERVER_NAME,
		MaxRequestBodySize: 512 * 1024 * 1024,
	}
	log.Fatal(fastHttpServer.ListenAndServe(PORT))
}
