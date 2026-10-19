package user

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

type CreateInput struct {
	AccountType string `validate:"" json:"account_type"`
	Email       string `validate:"required,email" json:"email"`
	FirstName   string `validate:"required" json:"first_name"`
	LastName    string `validate:"required" json:"last_name"`
	Gender      int    `validate:"min=0" json:"gender"`
	Password    string `validate:"required" json:"password"`
	BirthDay    string `validate:"required,datetime=2006-01-02" json:"birth_day"`
}

func ValidateCreateInput(ctx *fastHttp.RequestCtx) (regRequest CreateInput, status bool) {
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
			wg              sync.WaitGroup
			userDetail      userModel.User
			userOtherDetail userModel.User
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
			userOtherDetail = userModel.GetUserByEmail(regRequest.Email, "")
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if userOtherDetail.ID != "" {
			response.SendError(ctx, "This email address has already been registered on the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type UpdateInput struct {
	ID          string `validate:"required" json:"_id"`
	AccountType string `validate:"" json:"account_type"`
	Email       string `validate:"required,email" json:"email"`
	FirstName   string `validate:"required" json:"first_name"`
	LastName    string `validate:"required" json:"last_name"`
	Gender      int    `validate:"min=0" json:"gender"`
	BirthDay    string `validate:"required,datetime=2006-01-02" json:"birth_day"`
}

func ValidateUpdateInput(ctx *fastHttp.RequestCtx) (regRequest UpdateInput, status bool) {
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
			wg              sync.WaitGroup
			oldUserDetail   userModel.User
			userDetail      userModel.User
			userOtherDetail userModel.User
		)
		wg.Add(3)
		go func() {
			defer wg.Done()
			oldUserDetail = userModel.GetUserById(regRequest.ID, true)
		}()
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(userId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		go func() {
			defer wg.Done()
			userOtherDetail = userModel.GetUserByEmail(regRequest.Email, regRequest.ID)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			response.SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		if userOtherDetail.ID != "" {
			response.SendError(ctx, "This email address has already been registered on the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type DeleteInput struct {
	ID string `validate:"required" json:"_id"`
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
			wg            sync.WaitGroup
			oldUserDetail userModel.User
			userDetail    userModel.User
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			oldUserDetail = userModel.GetUserById(regRequest.ID, true)
		}()
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(userId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			response.SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

func ValidateActiveInput(ctx *fastHttp.RequestCtx) (oldUserDetail userModel.User, status bool) {
	status = false
	userId := libUtilities.GetUserId(ctx)
	if userId == "" {
		response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
		return
	}
	regRequest := DeleteInput{}
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg         sync.WaitGroup
			userDetail userModel.User
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			oldUserDetail = userModel.GetUserById(regRequest.ID, true)
		}()
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(userId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			response.SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return oldUserDetail, status
}

type UpdatePasswordInput struct {
	Password string `validate:"required" json:"password"`
}

func ValidateUpdatePasswordInput(ctx *fastHttp.RequestCtx) (oldUserDetail userModel.User, regRequest UpdatePasswordInput, status bool) {
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
		oldUserDetail = userModel.GetUserById(userId, true)
		if oldUserDetail.ID == "" {
			response.SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return oldUserDetail, regRequest, status
}

type UpdatePasswordByIdInput struct {
	ID       string `bson:"_id" json:"_id"`
	Password string `validate:"required" json:"password"`
}

func ValidateUpdatePasswordByIdInput(ctx *fastHttp.RequestCtx) (oldUserDetail userModel.User, regRequest UpdatePasswordByIdInput, status bool) {
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
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			oldUserDetail = userModel.GetUserById(regRequest.ID, true)
		}()
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(userId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			response.SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return oldUserDetail, regRequest, status
}

type SearchUserInput struct {
	Key           string `validate:"" json:"key"`
	AccountTypeId string `validate:"" json:"account_type"`
	Page          int64  `validate:"min=1" json:"page"`
	Limit         int64  `validate:"min=0" json:"limit"`
}

func ValidateSearchUserInput(ctx *fastHttp.RequestCtx) (regRequest SearchUserInput, status bool) {
	status = false
	userId := libUtilities.GetUserId(ctx)
	if userId == "" {
		response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
		return
	}
	libProcess.Try(func() {
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		var (
			wg                sync.WaitGroup
			userDetail        userModel.User
			accountTypeDetail permissionModel.AccountType
		)
		hasAccountTypeId := regRequest.AccountTypeId != ""
		if hasAccountTypeId {
			wg.Add(2)
		} else {
			wg.Add(1)
		}
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(userId, true)
		}()
		if hasAccountTypeId {
			go func() {
				defer wg.Done()
				accountTypeDetail = permissionModel.GetAccountTypeById(regRequest.AccountTypeId, true)
			}()
		}
		wg.Wait()
		if userDetail.AccountType != "1" {
			userDetail.ID = ""
		}
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if hasAccountTypeId && accountTypeDetail.ID == "" {
			response.SendError(ctx, "This account type information does not exist in the system.", nil, 206)
			return
		}

		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})

	return regRequest, status
}
