package main

import (
	"log"
	"os"
	"strconv"
	"time"

	fastHttpRouter "github.com/buaazp/fasthttprouter"
	"github.com/joho/godotenv"
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libCrypto "github.com/nguoihanoi/golang_shared/libs/crypto"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	"github.com/redis/go-redis/v9"
	fastHttp "github.com/valyala/fasthttp"
)

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

func corsMiddleware(next fastHttp.RequestHandler) fastHttp.RequestHandler {
	output := func(ctx *fastHttp.RequestCtx) {
		// Set CORS headers
		start := time.Now()
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Expose-Headers", "Authorization")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, Cache-Control")
		// Handle preflight (OPTIONS) requests
		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(fastHttp.StatusOK)
			return
		}

		// Do middleware things
		defer func() {
			log.Println(string(ctx.Path()), time.Since(start))
		}()
		// Call the next handler
		next(ctx)
	}
	return fastHttp.CompressHandler(output)
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
	PREFIX_TOKEN := os.Getenv("PREFIX_TOKEN")

	//Todo: init db
	mgoDB1 := initDb(MONGO_COMMON_URI, MONGO_COMMON_DB)
	redisDb1 := initRedisDb(REDIS_HOST, REDIS_DB, REDIS_PASSWORD, REDIS_PREFIX)

	//Todo: init  crypto
	libCrypto.JWT().SetToken(JWT_SECRET, PREFIX_TOKEN)

	//Todo: init router
	mainRouter := fastHttpRouter.New()

	//Todo: init db
	Init(mainRouter, mgoDB1, redisDb1)

	// Create a custom logger
	myLogger := log.New(log.Writer(), "FasthttpServer: ", log.LstdFlags)
	fastHttpServer := &fastHttp.Server{
		Handler:            corsMiddleware(mainRouter.Handler),
		Logger:             myLogger,
		Name:               SERVER_NAME,
		MaxRequestBodySize: 512 * 1024 * 1024,
	}
	log.Fatal(fastHttpServer.ListenAndServe(PORT))
}
