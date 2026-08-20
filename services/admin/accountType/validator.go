package accountType

import (
	"sync"
	"sync/atomic"

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
}

func ValidateCreateInput(ctx *fastHttp.RequestCtx) (regRequest CreateInput, status bool) {
	status = false

	libProcess.Try(func() {
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if len(regRequest.Permissions) > 0 {
			regRequest.Permissions = libUtilities.Array().RemoveDuplicates(regRequest.Permissions, false)
		}
		var (
			wg               sync.WaitGroup
			userDetail       userModel.User
			invalidPermFound atomic.Bool
			sem              = make(chan struct{}, 4)
		)
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		for _, permID := range regRequest.Permissions {
			if invalidPermFound.Load() {
				break
			}
			wg.Add(1)
			sem <- struct{}{}
			go func(id string) {
				defer wg.Done()
				defer func() { <-sem }()
				if invalidPermFound.Load() {
					return
				}
				permissionDetail := permissionModel.GetPermissionById(id, true)
				if permissionDetail.ID == "" {
					invalidPermFound.Store(true)
				}
			}(permID)
		}
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if invalidPermFound.Load() {
			response.SendError(ctx, "This permission information does not exist in the system.", nil, 206)
			return
		}

		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
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
}

func ValidateUpdateInput(ctx *fastHttp.RequestCtx) (regRequest UpdateInput, status bool) {
	status = false
	libProcess.Try(func() {
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if len(regRequest.Permissions) > 0 {
			regRequest.Permissions = libUtilities.Array().RemoveDuplicates(regRequest.Permissions, false)
		}
		var (
			wg                sync.WaitGroup
			userDetail        userModel.User
			accountTypeDetail permissionModel.AccountType
			invalidPermFound  atomic.Bool
			sem               = make(chan struct{}, 4)
		)
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
		}()
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			accountTypeDetail = permissionModel.GetAccountTypeById(regRequest.Id, true)
		}()
		for _, permID := range regRequest.Permissions {
			if invalidPermFound.Load() {
				break
			}
			wg.Add(1)
			sem <- struct{}{}
			go func(id string) {
				defer wg.Done()
				defer func() { <-sem }()
				if invalidPermFound.Load() {
					return
				}
				permissionDetail := permissionModel.GetPermissionById(id, true)
				if permissionDetail.ID == "" {
					invalidPermFound.Store(true)
				}
			}(permID)
		}
		wg.Wait()
		if userDetail.AccountType != "1" {
			userDetail.ID = ""
		}
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if invalidPermFound.Load() {
			response.SendError(ctx, "This permission information does not exist in the system.", nil, 206)
			return
		}
		if accountTypeDetail.ID == "" {
			response.SendError(ctx, "This account type information does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type DeleteInput struct {
	Id     string `validate:"required" json:"_id"`
	UserId string `validate:"required" json:"user_id"`
}

func ValidateDeleteInput(ctx *fastHttp.RequestCtx) (regRequest DeleteInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg                sync.WaitGroup
			userDetail        userModel.User
			accountTypeDetail permissionModel.AccountType
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
			accountTypeDetail = permissionModel.GetAccountTypeById(regRequest.Id, true)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if accountTypeDetail.ID == "" {
			response.SendError(ctx, "This account type information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchAccountTypeInput struct {
	Key    string `validate:"" json:"key"`
	Status int    `validate:"min=1" json:"status"`
	Page   int64  `validate:"min=1" json:"page"`
	Limit  int64  `validate:"min=0" json:"limit"`
	UserId string `validate:"required" json:"user_id"`
}

func ValidateSearchAccountTypeInput(ctx *fastHttp.RequestCtx) (regRequest SearchAccountTypeInput, status bool) {
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
