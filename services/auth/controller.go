package main

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	customerModel.InitModel(inDb, inRedisClient, inJwtToken)
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

func CustomerLogin(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Login false!", 206)
	customerDetail, regRequest, status := ValidateCustomerLoginInput(ctx)
	if status {
		_, resultStatus := customerModel.Login(customerDetail, regRequest.Password, false)
		if resultStatus {
			resp.Status = true
			resp.Message = "Login success!"
			resp.Data = customerDetail
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

func Login(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Login false!", 206)
	userDetail, regRequest, status := ValidateLoginInput(ctx)
	if status {
		_, resultStatus := userModel.Login(userDetail, regRequest.Password, false)
		if resultStatus {
			resp.Status = true
			resp.Message = "Login success!"
			resp.Data = userDetail
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

func Register(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Register false!", 206)
	regRequest, status := ValidateRegisterInput(ctx)
	if status {
		Password, PasswordHash := libUtilities.String().GetHashPassWord("", "", false)
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
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
