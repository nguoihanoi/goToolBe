package permission

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
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
		typeDetail := permissionModel.GetTypeById(regRequest.PermissionTypeID, true)
		if typeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
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
			wg               sync.WaitGroup
			permissionDetail permissionModel.Permission     // Thay bằng struct tương ứng của bạn
			typeDetail       permissionModel.PermissionType // Thay bằng struct tương ứng
		)
		wg.Add(2)
		// Query 1: Lấy chi tiết Permission
		go func() {
			defer wg.Done()
			permissionDetail = permissionModel.GetPermissionById(regRequest.ID, true)
		}()
		// Query 2: Lấy chi tiết Permission Type
		go func() {
			defer wg.Done()
			typeDetail = permissionModel.GetTypeById(regRequest.PermissionTypeID, true)
		}()
		// Chờ cả 2 goroutine hoàn thành
		wg.Wait()
		// 3. Kiểm tra kết quả thu được
		if permissionDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission information does not exist in the system.", nil, 206)
			return
		}
		if typeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
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
		permissionDetail := permissionModel.GetPermissionById(regRequest.ID, true)
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
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchPermissionInput(ctx *fastHttp.RequestCtx) (regRequest SearchPermissionInput, status bool) {
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
		if regRequest.TypeId != "" {
			typeDetail := permissionModel.GetTypeById(regRequest.TypeId, true)
			if typeDetail.ID == "" {
				libUtilities.Response().SendError(ctx, "This permission type information does not exist in the system.", nil, 206)
			}
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
