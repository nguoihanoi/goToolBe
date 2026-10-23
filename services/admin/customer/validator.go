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
	Key     string `validate:"" json:"key"`
	GroupID string `validate:"" json:"group_id"`
	Status  int64  `validate:"min=-1,max=1" json:"status"`
	Page    int64  `validate:"min=1" json:"page"`
	Limit   int64  `validate:"min=0" json:"limit"`
}

func ValidateSearchCustomerInput(ctx *fastHttp.RequestCtx) (regRequest SearchInput, status bool) {
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
			userDetail = userModel.GetUserById(userId, true)
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
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if hasGroupID && groupDetail.ID == "" {
			response.SendError(ctx, "This group's information does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type CreateInput struct {
	GroupID   string `validate:"" json:"group_id"`
	Email     string `validate:"required,email" json:"email"`
	FirstName string `validate:"required" json:"first_name"`
	LastName  string `validate:"required" json:"last_name"`
	Gender    int    `validate:"min=0" json:"gender"`
	Password  string `validate:"required" json:"password"`
	BirthDay  string `validate:"required,datetime=2006-01-02" json:"birth_day"`
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
			wg             sync.WaitGroup
			userDetail     userModel.User
			customerDetail customerModel.Customer
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
			customerDetail = customerModel.GetCustomerByEmail(regRequest.Email, "")
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if customerDetail.ID != "" {
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
	ID        string `validate:"required" json:"_id"`
	GroupID   string `validate:"" json:"group_id"`
	Email     string `validate:"required,email" json:"email"`
	FirstName string `validate:"required" json:"first_name"`
	LastName  string `validate:"required" json:"last_name"`
	Gender    int    `validate:"min=0" json:"gender"`
	BirthDay  string `validate:"required,datetime=2006-01-02" json:"birth_day"`
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
			wg                sync.WaitGroup
			oldCustomerDetail customerModel.Customer
			userDetail        userModel.User
			customerDetail    customerModel.Customer
		)
		wg.Add(3)
		go func() {
			defer wg.Done()
			oldCustomerDetail = customerModel.GetCustomerById(regRequest.ID, true)
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
			customerDetail = customerModel.GetCustomerByEmail(regRequest.Email, regRequest.ID)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			response.SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if oldCustomerDetail.ID == "" {
			response.SendError(ctx, "This customer does not exist in the system.", nil, 206)
			return
		}
		if customerDetail.ID != "" {
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
			wg                sync.WaitGroup
			oldCustomerDetail customerModel.Customer
			userDetail        userModel.User
		)
		wg.Add(2)
		go func() {
			defer wg.Done()
			oldCustomerDetail = customerModel.GetCustomerById(regRequest.ID, true)
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
		if oldCustomerDetail.ID == "" {
			response.SendError(ctx, "This customer does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type UpdatePasswordByIdInput struct {
	ID       string `bson:"_id" json:"_id"`
	Password string `validate:"required" json:"password"`
}

func ValidateUpdatePasswordByIdInput(ctx *fastHttp.RequestCtx) (oldCustomerDetail customerModel.Customer, regRequest UpdatePasswordByIdInput, status bool) {
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
			oldCustomerDetail = customerModel.GetCustomerById(regRequest.ID, true)
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
		if oldCustomerDetail.ID == "" {
			response.SendError(ctx, "This customer does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
	})
	return oldCustomerDetail, regRequest, status
}
