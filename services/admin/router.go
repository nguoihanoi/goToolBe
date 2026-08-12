package main

import (
	//commonMiddleware "auth-service/middlewares"
	fastHttpRouter "github.com/buaazp/fasthttprouter"
	"github.com/nguoihanoi/goToolBe/services/admin/customer"
	"github.com/nguoihanoi/goToolBe/services/admin/user"
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	fastHttp "github.com/valyala/fasthttp"
)

func Init(router *fastHttpRouter.Router, inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	user.InitController(inDb, inRedisClient, inJwtToken)
	customer.InitController(inDb, inRedisClient, inJwtToken)
	router.GET("/health", func(ctx *fastHttp.RequestCtx) {
		ctx.SetStatusCode(fastHttp.StatusOK)
		ctx.SetBodyString("Auth service is up and running")
	})
	router.POST("/api/v1/user/search", user.User)
	router.POST("/api/v1/customer", customer.Customer)
	router.GET("/api/v1/customer/group", customer.CustomerGroup)
	//router.POST("/login", commonMiddleware.Post(Login))
	//router.POST("/api/v1/user/register", commonMiddleware.Post(Register))
}
