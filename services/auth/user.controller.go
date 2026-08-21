package main

import (
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

func login(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Login false!", 206)
	userDetail, regRequest, status := ValidateLoginInput(ctx)
	if status {
		_, resultStatus := userModel.Login(userDetail, regRequest.Password, false)
		if resultStatus {
			resp.Status = true
			resp.Message = "Login success!"
			resp.Data = userDetail
		}
		response.SendOutput(ctx, resp)
	}
}

var authCmdMap = map[string]CommandHandler{
	"login": login,
}

func Auth(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := authCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
