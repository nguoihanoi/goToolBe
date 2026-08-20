package main

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

type CustomerRegistrationInput struct {
	Email     string `json:"email" validate:"email"`
	Password  string `json:"password" validate:"min=2"`
	FirstName string `json:"first_name" validate:"min=2"`
	LastName  string `json:"last_name" validate:"min=2"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=2"`
}

func ValidateLoginInput(ctx *fastHttp.RequestCtx) (userDetail userModel.User, regRequest LoginInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		userDetail = userModel.GetUserByEmail(regRequest.Email, "")
		if userDetail.ID != "" {
			status = true
		} else {
			response.SendError(ctx, "This email does not exist in the system.", nil, 206)
		}
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
		status = false
	})
	return userDetail, regRequest, status
}

func ValidateCustomerLoginInput(ctx *fastHttp.RequestCtx) (customerDetail customerModel.Customer, regRequest LoginInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		customerDetail = customerModel.GetCustomerByEmail(regRequest.Email, "")
		if customerDetail.ID != "" {
			status = true
		} else {
			response.SendError(ctx, "This email does not exist in the system.", nil, 206)
		}
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
		status = false
	})
	return customerDetail, regRequest, status
}

func ValidateRegisterInput(ctx *fastHttp.RequestCtx) (regRequest CustomerRegistrationInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		//
		status = false
		customerDetail := customerModel.GetCustomerByEmail(regRequest.Email, "")
		if customerDetail.ID != "" {
			response.SendError(ctx, "This email is registered in the system.", nil, 206)
		} else {
			status = true
		}
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
