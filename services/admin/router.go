package main

import (
	//commonMiddleware "auth-service/middlewares"
	fastHttpRouter "github.com/buaazp/fasthttprouter"
	"github.com/nguoihanoi/goToolBe/services/admin/accountType"
	"github.com/nguoihanoi/goToolBe/services/admin/customer"
	"github.com/nguoihanoi/goToolBe/services/admin/file"
	"github.com/nguoihanoi/goToolBe/services/admin/language"
	"github.com/nguoihanoi/goToolBe/services/admin/permission"
	"github.com/nguoihanoi/goToolBe/services/admin/permissionType"
	"github.com/nguoihanoi/goToolBe/services/admin/user"
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	fastHttp "github.com/valyala/fasthttp"
)

func Init(router *fastHttpRouter.Router, inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	user.InitController(inDb, inRedisClient, inJwtToken)
	language.InitController(inDb, inRedisClient, inJwtToken)
	file.InitController(inDb, inRedisClient, inJwtToken)
	customer.InitController(inDb, inRedisClient, inJwtToken)
	permission.InitController(inDb, inRedisClient, inJwtToken)
	permissionType.InitController(inDb, inRedisClient, inJwtToken)
	accountType.InitController(inDb, inRedisClient, inJwtToken)
	router.GET("/health", func(ctx *fastHttp.RequestCtx) {
		ctx.SetStatusCode(fastHttp.StatusOK)
		ctx.SetBodyString("Auth service is up and running")
	})
	router.GET("/api/v1/user", user.User)
	router.GET("/api/v1/file", file.File)
	router.POST("/api/v1/file/upload", file.UploadFile)
	router.GET("/api/v1/language", language.Language)
	router.GET("/api/v1/customer", customer.Customer)
	router.GET("/api/v1/customer/group", customer.CustomerGroup)
	router.GET("/api/v1/permission", permission.Permssion)
	router.GET("/api/v1/permission/type", permissionType.PermssionType)
	router.GET("/api/v1/accountType", accountType.AccountType)
}
