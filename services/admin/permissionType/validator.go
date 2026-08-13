package permissionType

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	fastHttp "github.com/valyala/fasthttp"
)

type CreateInput struct {
	Name     map[string]string `validate:"required" json:"name"`
	Order    int               `validate:"min=1" json:"order"`
	UserId   string            `validate:"required" json:"user_id"`
	LangCode string            `validate:"" json:"lang_code"`
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
	Order    int    `validate:"min=1" json:"order"`
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
		typeDetail := permissionModel.GetTypeById(regRequest.Id, true)
		if typeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
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
		typeDetail := permissionModel.GetTypeById(regRequest.ID, true)
		if typeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchPermissionTypeInput struct {
	Key      string `validate:"" json:"key"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchPermissionTypeInput(ctx *fastHttp.RequestCtx) (regRequest SearchPermissionTypeInput, status bool) {
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
