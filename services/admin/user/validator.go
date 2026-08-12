package user

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	fastHttp "github.com/valyala/fasthttp"
)

type SearchUserInput struct {
	Key           string `validate:"" json:"key"`
	AccountTypeId string `validate:"" json:"account_type"`
	Page          int64  `validate:"min=1" json:"page"`
	Limit         int64  `validate:"min=0" json:"limit"`
	LangCode      string `validate:"" json:"lang_code"`
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
