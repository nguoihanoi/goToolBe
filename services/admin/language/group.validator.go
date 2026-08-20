package language

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

type CreateGroupInput struct {
	Name   string `validate:"required" json:"name"`
	Code   string `validate:"required" json:"code"`
	UserId string `validate:"required" json:"user_id"`
}

func ValidateCreateGroupInput(ctx *fastHttp.RequestCtx) (regRequest CreateGroupInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg               sync.WaitGroup
			userDetail       userModel.User
			otherGroupDetail languageModel.GroupCode
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
			otherGroupDetail = languageModel.CheckGroupCode(regRequest.Code, "")
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if otherGroupDetail.ID != "" {
			response.SendError(ctx, "The group code has already been used.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type UpdateGroupInput struct {
	Id     string `validate:"required" json:"_id"`
	Name   string `validate:"required" json:"name"`
	Code   string `validate:"required" json:"code"`
	UserId string `validate:"required" json:"user_id"`
}

func ValidateUpdateGroupInput(ctx *fastHttp.RequestCtx) (regRequest UpdateGroupInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg               sync.WaitGroup
			userDetail       userModel.User
			otherGroupDetail languageModel.GroupCode
			groupDetail      languageModel.GroupCode
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
			otherGroupDetail = languageModel.CheckGroupCode(regRequest.Code, regRequest.Id)
		}()
		go func() {
			defer wg.Done()
			groupDetail = languageModel.GetGroupCodeById(regRequest.Id, true)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if otherGroupDetail.ID != "" {
			response.SendError(ctx, "The group code has already been used.", nil, 206)
		}
		if groupDetail.ID == "" {
			response.SendError(ctx, "This group code information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type DeleteGroupInput struct {
	Id     string `validate:"required" json:"_id"`
	UserId string `validate:"required" json:"user_id"`
}

func ValidateDeleteGroupInput(ctx *fastHttp.RequestCtx) (regRequest DeleteGroupInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg          sync.WaitGroup
			userDetail  userModel.User
			groupDetail languageModel.GroupCode
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
			groupDetail = languageModel.GetGroupCodeById(regRequest.Id, true)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if groupDetail.ID == "" {
			response.SendError(ctx, "This group information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchGroupInput struct {
	Key    string `validate:"" json:"key"`
	Page   int64  `validate:"min=1" json:"page"`
	Limit  int64  `validate:"min=0" json:"limit"`
	UserId string `validate:"required" json:"user_id"`
}

func ValidateSearchGroupInput(ctx *fastHttp.RequestCtx) (regRequest SearchGroupInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		userDetail := userModel.GetUserById(regRequest.UserId, true)
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
