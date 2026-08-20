package main

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

var response *libUtilities.ResponseClass

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	customerModel.InitModel(inDb, inRedisClient, inJwtToken)
	response = libUtilities.Response(inRedisClient, "codes")
	listEmail := []string{"playhard24h@gmail.com"}
	listName := []string{"Thành Nguyễn Viết"}
	passwordDefault := "abc123!@#"
	for i := range listEmail {
		userDetail := userModel.GetUserByEmail(listEmail[i], "")
		if userDetail.ID == "" {
			Password, PasswordHash := libUtilities.String().GetHashPassWord(passwordDefault, "", false)
			inFirstName, inLastName := libUtilities.String().GetFromFullName(listName[i])
			insertData := userModel.User{
				AccountType:  "1",
				Email:        listEmail[i],
				FirstName:    inFirstName,
				LastName:     inLastName,
				PasswordHash: PasswordHash,
				Password:     Password,
				IsActive:     true,
			}
			userModel.CreateUser(insertData)
		}
	}
}

func customerLogin(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Login false!", 206)
	customerDetail, regRequest, status := ValidateCustomerLoginInput(ctx)
	if status {
		_, resultStatus := customerModel.Login(customerDetail, regRequest.Password, false)
		if resultStatus {
			resp.Status = true
			resp.Message = "Login success!"
			resp.Data = customerDetail
		}
		response.SendOutput(ctx, resp)
	}
}

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

func register(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Register false!", 206)
	regRequest, status := ValidateRegisterInput(ctx)
	if status {
		Password, PasswordHash := libUtilities.String().GetHashPassWord(regRequest.Password, "", false)
		customerId := customerModel.CreateCustomer(customerModel.Customer{
			Email:        regRequest.Email,
			FirstName:    regRequest.FirstName,
			LastName:     regRequest.LastName,
			PasswordHash: PasswordHash,
			Password:     Password,
			IsActive:     false,
		})
		if customerId != "" {
			resp.Status = true
			resp.Message = "Register success!"
		}
		response.SendOutput(ctx, resp)
	}
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

// Khai báo map handler 1 lần duy nhất lúc khởi chạy app
var customerCmdMap = map[string]CommandHandler{
	"login":    customerLogin,
	"register": register,
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
func Customer(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := customerCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
