package main

import (
	//commonMiddleware "auth-service/middlewares"
	fastHttpRouter "github.com/buaazp/fasthttprouter"
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	mdWare "github.com/nguoihanoi/golang_shared/libs/middleware"
	fastHttp "github.com/valyala/fasthttp"
)

func Init(router *fastHttpRouter.Router, inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	InitController(inDb, inRedisClient, inJwtToken)
	router.GET("/health", func(ctx *fastHttp.RequestCtx) {
		ctx.SetStatusCode(fastHttp.StatusOK)
		ctx.SetBodyString("Auth service is up and running")
	})
	router.POST("/api/v1/auth/login", mdWare.Post(Login))
	router.POST("/api/v1/customer/login", mdWare.Post(CustomerLogin))
	router.POST("/api/v1/customer/register", mdWare.Post(Register))
	//router.POST("/login", commonMiddleware.Post(Login))
	//router.POST("/api/v1/user/register", commonMiddleware.Post(Register))
}
