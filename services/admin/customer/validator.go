package customer

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

type SearchInput struct {
	Key      string `validate:"" json:"key"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchCustomerGroupInput(ctx *fastHttp.RequestCtx) (regRequest SearchInput, status bool) {
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

type SearchCustomerInput struct {
	Key      string `validate:"" json:"key"`
	GroupID  string `validate:"" json:"group_id"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchCustomerInput(ctx *fastHttp.RequestCtx) (regRequest SearchCustomerInput, status bool) {
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
			wg          sync.WaitGroup
			userDetail  userModel.User
			groupDetail customerModel.CustomerGroup
		)

		hasGroupID := regRequest.GroupID != ""

		if hasGroupID {
			wg.Add(2)
		} else {
			wg.Add(1)
		}

		// Goroutine 1: Query User
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()

		// Goroutine 2: Query Group (Chỉ chạy khi có GroupID)
		if hasGroupID {
			go func() {
				defer wg.Done()
				groupDetail = customerModel.GetGroupById(regRequest.GroupID, true)
			}()
		}
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if hasGroupID && groupDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This group's information does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
