package main

import (
	//commonMiddleware "auth-service/middlewares"
	fastHttpRouter "github.com/buaazp/fasthttprouter"
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	fastHttp "github.com/valyala/fasthttp"
)

func Init(router *fastHttpRouter.Router, inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	InitController(inDb, inRedisClient, inJwtToken)
	router.GET("/health", func(ctx *fastHttp.RequestCtx) {
		ctx.SetStatusCode(fastHttp.StatusOK)
		ctx.SetBodyString("Auth service is up and running")
	})
	router.GET("/api/v1/auth", Auth)
	router.POST("/api/v1/customer", Customer)
}
