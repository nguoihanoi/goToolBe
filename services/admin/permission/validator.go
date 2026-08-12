package permission

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	fastHttp "github.com/valyala/fasthttp"
)

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
