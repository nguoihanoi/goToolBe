package main

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	fileModel "github.com/nguoihanoi/golang_shared/warehouses/files"
	fastHttp "github.com/valyala/fasthttp"
)

type LoginInput struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=2"`
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

type CustomerRegistrationInput struct {
	Email     string `json:"email" validate:"email"`
	Password  string `json:"password" validate:"min=2"`
	FirstName string `json:"first_name" validate:"min=2"`
	LastName  string `json:"last_name" validate:"min=2"`
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

type UpdateAvatarInput struct {
	Image      string `validate:"" json:"image"`
	CustomerId string `validate:"required" json:"customer_id"`
}

func ValidateUpdateAvatarInput(ctx *fastHttp.RequestCtx) (regRequest UpdateAvatarInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg              sync.WaitGroup
			custommerDetail customerModel.Customer
			fileDetail      fileModel.File
		)
		hasImageID := regRequest.Image != ""
		if hasImageID {
			wg.Add(2)
		} else {
			wg.Add(1)
		}
		wg.Add(2)
		go func() {
			defer wg.Done()
			custommerDetail = customerModel.GetCustomerById(regRequest.CustomerId, true)
			if custommerDetail.IsActive == false {
				custommerDetail.ID = ""
			}
		}()
		if hasImageID {
			go func() {
				defer wg.Done()
				fileDetail = fileModel.GetById(regRequest.Image, true)
			}()
		}
		wg.Wait()
		if custommerDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if hasImageID && fileDetail.ID == "" {
			response.SendError(ctx, "This image does not exist in the system.", nil, 206)
		}
		status = false
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type UpdateProfileInput struct {
	FirstName  string `bson:"required" json:"first_name"`
	LastName   string `bson:"required" json:"last_name"`
	CustomerId string `validate:"required" json:"customer_id"`
}

func ValidateUpdateProfileInput(ctx *fastHttp.RequestCtx) (regRequest UpdateProfileInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		custommerDetail := customerModel.GetCustomerById(regRequest.CustomerId, true)
		if custommerDetail.IsActive == false {
			custommerDetail.ID = ""
		}
		if custommerDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		status = false
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
