package accountType

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

type CreateInput struct {
	Name        map[string]string `validate:"required" json:"name"`
	Content     map[string]string `validate:"required" json:"content"`
	Permissions []string          `validate:"required" json:"permissions"`
	Status      int               `validate:"min=0" json:"status"`
	Order       int               `validate:"min=1" json:"order"`
	UserId      string            `validate:"required" json:"user_id"`
	LangCode    string            `validate:"" json:"lang_code"`
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
		userDetail := userModel.GetUserById(regRequest.UserId, true)
		if userDetail.AccountType != "1" {
			userDetail.ID = ""
		}
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if len(regRequest.Permissions) > 0 {
			regRequest.Permissions = libUtilities.Array().RemoveDuplicates(regRequest.Permissions, false)
			for i := 0; i < len(regRequest.Permissions); i += 1 {
				permissionDetail := permissionModel.GetPermissionById(regRequest.Permissions[i], true)
				if permissionDetail.ID == "" {
					libUtilities.Response().SendError(ctx, "This permission information does not exist in the system.", nil, 206)
					return
				}
			}
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type UpdateInput struct {
	Id          string            `validate:"required" json:"_id"`
	Name        map[string]string `validate:"required" json:"name"`
	Content     map[string]string `validate:"required" json:"content"`
	Permissions []string          `validate:"required" json:"permissions"`
	Status      int               `validate:"min=0" json:"status"`
	Order       int               `validate:"min=1" json:"order"`
	UserId      string            `validate:"required" json:"user_id"`
	LangCode    string            `validate:"" json:"lang_code"`
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
		if len(regRequest.Permissions) > 0 {
			regRequest.Permissions = libUtilities.Array().RemoveDuplicates(regRequest.Permissions, false)
			for i := 0; i < len(regRequest.Permissions); i += 1 {
				permissionDetail := permissionModel.GetPermissionById(regRequest.Permissions[i], true)
				if permissionDetail.ID == "" {
					libUtilities.Response().SendError(ctx, "This permission information does not exist in the system.", nil, 206)
					return
				}
			}
		}
		accountTypeDetail := permissionModel.GetAccountTypeById(regRequest.Id, true)
		if accountTypeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account type information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type DeleteInput struct {
	Id       string `validate:"required" json:"_id"`
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
		accountTypeDetail := permissionModel.GetAccountTypeById(regRequest.Id, true)
		if accountTypeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account type information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchAccountTypeInput struct {
	Key      string `validate:"" json:"key"`
	Status   int    `validate:"min=1" json:"status"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchAccountTypeInput(ctx *fastHttp.RequestCtx) (regRequest SearchAccountTypeInput, status bool) {
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
