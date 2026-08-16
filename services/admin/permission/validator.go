package permission

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
	Name             map[string]string `validate:"required" json:"name"`
	Code             string            `validate:"required" json:"code"`
	PermissionTypeID string            `validate:"required" json:"type_id"`
	Order            int               `validate:"min=1" json:"order"`
	UserId           string            `validate:"required" json:"user_id"`
	LangCode         string            `validate:"" json:"lang_code"`
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
			wg               sync.WaitGroup
			userDetail       userModel.User
			typeDetail       permissionModel.PermissionType
			permissionDetail permissionModel.Permission
			resultValidate   any
			statusValidate   int
		)
		wg.Add(4)
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		go func() {
			defer wg.Done()
			permissionDetail = permissionModel.CheckPermissionCode(regRequest.Code, "")
		}()
		go func() {
			defer wg.Done()
			typeDetail = permissionModel.GetTypeById(regRequest.PermissionTypeID, true)
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
		if permissionDetail.ID != "" {
			libUtilities.Response().SendError(ctx, "This permission code already exists in the system.", nil, 206)
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

type UpdateInput struct {
	ID               string            `validate:"required" json:"_id"`
	Name             map[string]string `validate:"required" json:"name"`
	Code             string            `validate:"required" json:"code"`
	PermissionTypeID string            `validate:"required" json:"type_id"`
	Order            int               `validate:"min=1" json:"order"`
	UserId           string            `validate:"required" json:"user_id"`
	LangCode         string            `validate:"" json:"lang_code"`
}

func ValidateUpdateInput(ctx *fastHttp.RequestCtx) (regRequest UpdateInput, status bool) {
	status = false

	libProcess.Try(func() {
		// 1. Validate struct input từ request
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		if regRequest.LangCode == "" {
			regRequest.LangCode = "vi"
		}
		// 2. Tối ưu xử lý song song 2 query độc lập
		var (
			wg                    sync.WaitGroup
			userDetail            userModel.User
			permissionDetail      permissionModel.Permission
			typeDetail            permissionModel.PermissionType
			permissionOtherDetail permissionModel.Permission
			resultValidate        any
			statusValidate        int
		)
		wg.Add(5)
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		go func() {
			defer wg.Done()
			permissionOtherDetail = permissionModel.CheckPermissionCode(regRequest.Code, regRequest.ID)
		}()
		go func() {
			defer wg.Done()
			permissionDetail = permissionModel.GetPermissionById(regRequest.ID, true)
		}()
		go func() {
			defer wg.Done()
			typeDetail = permissionModel.GetTypeById(regRequest.PermissionTypeID, true)
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
		if permissionOtherDetail.ID != "" {
			libUtilities.Response().SendError(ctx, "This permission code already exists in the system.", nil, 206)
			return
		}
		// 3. Kiểm tra kết quả thu được
		if permissionDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission information does not exist in the system.", nil, 206)
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
	ID       string `validate:"required" json:"_id"`
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
			wg               sync.WaitGroup
			userDetail       userModel.User
			permissionDetail permissionModel.Permission
		)
		wg.Add(5)
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		go func() {
			defer wg.Done()
			permissionDetail = permissionModel.GetPermissionById(regRequest.ID, true)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if permissionDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission information does not exist in the system.", nil, 206)
			return
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchPermissionInput struct {
	Key      string `validate:"" json:"key"`
	TypeId   string `validate:"" json:"type_id"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchPermissionInput(ctx *fastHttp.RequestCtx) (regRequest SearchPermissionInput, status bool) {
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
			wg         sync.WaitGroup
			userDetail userModel.User
			typeDetail permissionModel.PermissionType
		)

		hasTypeId := regRequest.TypeId != ""
		if hasTypeId {
			wg.Add(2)
		} else {
			wg.Add(1)
		}
		go func() {
			defer wg.Done()
			userDetail = userModel.GetUserById(regRequest.UserId, true)
			if userDetail.AccountType != "1" {
				userDetail.ID = ""
			}
		}()
		if hasTypeId {
			go func() {
				defer wg.Done()
				typeDetail = permissionModel.GetTypeById(regRequest.TypeId, true)
			}()
		}
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if hasTypeId && typeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
			return
		}

		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})

	return regRequest, status
}
