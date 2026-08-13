package language

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	fastHttp "github.com/valyala/fasthttp"
)

type CreateInput struct {
	Name     string `validate:"required" json:"name"`
	Code     string `validate:"required" json:"code"`
	Image    string `validate:"required" json:"image"`
	Order    int    `validate:"min=1" json:"order"`
	Status   int    `validate:"min=0" json:"status"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateCreateInput(ctx *fastHttp.RequestCtx) (regRequest CreateInput, status bool) {
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

type UpdateInput struct {
	Id       string `validate:"required" json:"_id"`
	Name     string `validate:"required" json:"name"`
	Code     string `validate:"required" json:"code"`
	Image    string `validate:"required" json:"image"`
	Order    int    `validate:"min=1" json:"order"`
	Status   int    `validate:"min=0" json:"status"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateUpdateInput(ctx *fastHttp.RequestCtx) (regRequest UpdateInput, status bool) {
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
		languageDetail := languageModel.GetById(regRequest.Id, true)
		if languageDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This language information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type DeleteInput struct {
	ID       string `validate:"required" json:"_id"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateDeleteInput(ctx *fastHttp.RequestCtx) (regRequest DeleteInput, status bool) {
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
		languageDetail := languageModel.GetById(regRequest.ID, true)
		if languageDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This language information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchInput struct {
	Key      string `validate:"" json:"key"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchInput(ctx *fastHttp.RequestCtx) (regRequest SearchInput, status bool) {
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
