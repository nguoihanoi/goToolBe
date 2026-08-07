package main

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"

	//userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

type SearchCustomerInput struct {
	Key      string `validate:"" json:"key"`
	GroupID  string `validate:"" json:"group_id"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	LangCode string `validate:"" json:"lang_code"`
}

type SearchUserInput struct {
	Key           string `validate:"" json:"key"`
	AccountTypeId string `validate:"" json:"account_type"`
	Page          int64  `validate:"min=1" json:"page"`
	Limit         int64  `validate:"min=0" json:"limit"`
	LangCode      string `validate:"" json:"lang_code"`
}

type SearchInput struct {
	Key      string `validate:"" json:"key"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	LangCode string `validate:"" json:"lang_code"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=2"`
}

func ValidateSearchCustomerGroupInput(ctx *fastHttp.RequestCtx) (regRequest SearchInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if regRequest.LangCode == "" {
			regRequest.LangCode = "vi"
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

func ValidateSearchCustomerInput(ctx *fastHttp.RequestCtx) (regRequest SearchCustomerInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if regRequest.LangCode == "" {
			regRequest.LangCode = "vi"
		}
		if regRequest.GroupID != "" {
			groupDetail := customerModel.GetGroupById(regRequest.GroupID, true)
			if groupDetail.ID == "" {
				libUtilities.Response().SendError(ctx, "This group's information does not exist in the system.", nil, 206)
			}
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
func ValidateSearchUserInput(ctx *fastHttp.RequestCtx) (regRequest SearchUserInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if regRequest.LangCode == "" {
			regRequest.LangCode = "vi"
		}
		if regRequest.AccountTypeId != "" {
			groupDetail := customerModel.GetGroupById(regRequest.AccountTypeId, true)
			if groupDetail.ID == "" {
				libUtilities.Response().SendError(ctx, "This group's information does not exist in the system.", nil, 206)
			}
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
