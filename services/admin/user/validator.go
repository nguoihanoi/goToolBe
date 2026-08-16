package user

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

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
