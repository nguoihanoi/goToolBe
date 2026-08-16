package permissionType

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
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
		var (
			wg             sync.WaitGroup
			userDetail     userModel.User
			resultValidate any
			statusValidate int
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		go func() {
			defer wg.Done()
			langCode := languageModel.GetCodes()
			resultValidate, statusValidate = libUtilities.ValidateLangValue(regRequest.Name, false, langCode)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		switch statusValidate {
		case 1, 2:
			libUtilities.Response().SendError(ctx, "Invalid input data!", resultValidate, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type UpdateInput struct {
	Id       string            `validate:"required" json:"_id"`
	Name     map[string]string `validate:"required" json:"name"`
	Order    int               `validate:"min=1" json:"order"`
	UserId   string            `validate:"required" json:"user_id"`
	LangCode string            `validate:"" json:"lang_code"`
}

func ValidateUpdateInput(ctx *fastHttp.RequestCtx) (regRequest UpdateInput, status bool) {
	status = false
	libProcess.Try(func() {
		// 1. Unmarshal & Validate Struct cơ bản
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if regRequest.LangCode == "" {
			regRequest.LangCode = "vi"
		}

		var (
			wg             sync.WaitGroup
			userDetail     userModel.User
			typeDetail     permissionModel.PermissionType
			resultValidate any
			statusValidate int
		)
		wg.Add(3)
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		go func() {
			defer wg.Done()
			typeDetail = permissionModel.GetTypeById(regRequest.Id, true)
		}()
		go func() {
			defer wg.Done()
			langCode := languageModel.GetCodes()
			resultValidate, statusValidate = libUtilities.ValidateLangValue(regRequest.Name, false, langCode)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if typeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
			return
		}
		switch statusValidate {
		case 1, 2:
			libUtilities.Response().SendError(ctx, "Invalid input data!", resultValidate, 206)
			return
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
		var (
			wg         sync.WaitGroup
			userDetail userModel.User
			typeDetail permissionModel.PermissionType
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		go func() {
			defer wg.Done()
			typeDetail = permissionModel.GetTypeById(regRequest.Id, true)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
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
	UserId   string `validate:"required" json:"user_id"`
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
		userDetail := userModel.GetUserById(regRequest.UserId, true)
		if userDetail.AccountType != "1" {
			userDetail.ID = ""
		}
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
