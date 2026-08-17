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
	AccountType string `bson:"account_type" json:"account_type"`
	Email       string `bson:"email" json:"email"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Password    string `bson:"password" json:"password"`
	UserId      string `validate:"required" json:"user_id"`
	LangCode    string `validate:"" json:"lang_code"`
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
			wg              sync.WaitGroup
			userDetail      userModel.User
			userOtherDetail userModel.User
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
			userOtherDetail = userModel.GetUserByEmail(regRequest.Email, "")
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if userOtherDetail.ID != "" {
			libUtilities.Response().SendError(ctx, "This email address has already been registered on the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type UpdateInput struct {
	ID          string `bson:"_id" json:"_id"`
	AccountType string `bson:"account_type" json:"account_type"`
	Email       string `bson:"email" json:"email"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	UserId      string `validate:"required" json:"user_id"`
	LangCode    string `validate:"" json:"lang_code"`
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
			userDetail = userModel.GetUserById(regRequest.UserId, true)
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
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		if userOtherDetail.ID != "" {
			libUtilities.Response().SendError(ctx, "This email address has already been registered on the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type DeleteInput struct {
	ID       string `bson:"_id" json:"_id"`
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
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

func ValidateActiveInput(ctx *fastHttp.RequestCtx) (oldUserDetail userModel.User, status bool) {
	status = false
	regRequest := DeleteInput{}
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
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			oldUserDetail = userModel.GetUserById(regRequest.ID, true)
		}()
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return oldUserDetail, status
}

type UpdatePasswordInput struct {
	UserId   string `validate:"required" json:"user_id"`
	Password string `validate:"required" json:"password"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateUpdatePasswordInput(ctx *fastHttp.RequestCtx) (oldUserDetail userModel.User, regRequest UpdatePasswordInput, status bool) {
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
		oldUserDetail = userModel.GetUserById(regRequest.UserId, true)
		if oldUserDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return oldUserDetail, regRequest, status
}

type UpdatePasswordByIdInput struct {
	ID       string `bson:"_id" json:"_id"`
	UserId   string `validate:"required" json:"user_id"`
	Password string `validate:"required" json:"password"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateUpdatePasswordByIdInput(ctx *fastHttp.RequestCtx) (oldUserDetail userModel.User, regRequest UpdatePasswordByIdInput, status bool) {
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
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			oldUserDetail = userModel.GetUserById(regRequest.ID, true)
		}()
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldUserDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return oldUserDetail, regRequest, status
}

type SearchUserInput struct {
	Key           string `validate:"" json:"key"`
	AccountTypeId string `validate:"" json:"account_type"`
	Page          int64  `validate:"min=1" json:"page"`
	Limit         int64  `validate:"min=0" json:"limit"`
	UserId        string `validate:"required" json:"user_id"`
	LangCode      string `validate:"" json:"lang_code"`
}

func ValidateSearchUserInput(ctx *fastHttp.RequestCtx) (regRequest SearchUserInput, status bool) {
	status = false
	libProcess.Try(func() {
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if regRequest.LangCode == "" {
			regRequest.LangCode = "vi"
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
			userDetail = userModel.GetUserById(regRequest.UserId, true)
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
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if hasAccountTypeId && accountTypeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This account type information does not exist in the system.", nil, 206)
			return
		}

		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})

	return regRequest, status
}
