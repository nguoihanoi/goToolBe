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
	Name   map[string]string `validate:"required" json:"name"`
	Order  int               `validate:"min=1" json:"order"`
	Status int               `validate:"min=0" json:"status"`
}

func ValidateCreateInput(ctx *fastHttp.RequestCtx) (regRequest CreateInput, userId string, status bool) {
	status = false
	userId = libUtilities.GetUserId(ctx)
	if userId == "" {
		response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
		return
	}
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
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
			userDetail = userModel.GetUserById(userId, true)
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
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		switch statusValidate {
		case 1, 2:
			response.SendError(ctx, "Invalid input data!", resultValidate, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, userId, status
}

type UpdateInput struct {
	Id     string            `validate:"required" json:"_id"`
	Name   map[string]string `validate:"required" json:"name"`
	Order  int               `validate:"min=1" json:"order"`
	Status int               `validate:"min=0" json:"status"`
}

func ValidateUpdateInput(ctx *fastHttp.RequestCtx) (regRequest UpdateInput, status bool) {
	status = false
	userId := libUtilities.GetUserId(ctx)
	if userId == "" {
		response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
		return
	}
	libProcess.Try(func() {
		// 1. Unmarshal & Validate Struct cơ bản
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
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
			userDetail = userModel.GetUserById(userId, true)
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
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if typeDetail.ID == "" {
			response.SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
			return
		}
		switch statusValidate {
		case 1, 2:
			response.SendError(ctx, "Invalid input data!", resultValidate, 206)
			return
		}

		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})

	return regRequest, status
}

type DeleteInput struct {
	Id string `validate:"required" json:"_id"`
}

func ValidateDeleteInput(ctx *fastHttp.RequestCtx) (regRequest DeleteInput, status bool) {
	status = false
	userId := libUtilities.GetUserId(ctx)
	if userId == "" {
		response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
		return
	}
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg         sync.WaitGroup
			userDetail userModel.User
			typeDetail permissionModel.PermissionType
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(userId, true)
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
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if typeDetail.ID == "" {
			response.SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchPermissionTypeInput struct {
	Key    string `validate:"" json:"key"`
	Status int64  `validate:"min=-1" json:"status"`
	Page   int64  `validate:"min=1" json:"page"`
	Limit  int64  `validate:"min=0" json:"limit"`
}

func ValidateSearchPermissionTypeInput(ctx *fastHttp.RequestCtx) (regRequest SearchPermissionTypeInput, status bool) {
	status = false
	userId := libUtilities.GetUserId(ctx)
	if userId == "" {
		response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
		return
	}
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		userDetail := userModel.GetUserById(userId, true)
		if userDetail.AccountType != "1" {
			userDetail.ID = ""
		}
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
