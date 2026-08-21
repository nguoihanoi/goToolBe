package main

import (
	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

func ValidateLoginInput(ctx *fastHttp.RequestCtx) (userDetail userModel.User, regRequest LoginInput, status bool) {
	status = false
	libProcess.Try(func() {
		//Todo: get struct input
		err := libUtilities.Validate(ctx, &regRequest)
		if err != nil {
			libProcess.Throw(err)
		}
		userDetail = userModel.GetUserByEmail(regRequest.Email, "")
		if userDetail.ID != "" {
			status = true
		} else {
			response.SendError(ctx, "This email does not exist in the system.", nil, 206)
		}
	}).Catch(func(e libProcess.E) {
		response.SendError(ctx, "Invalid input data!", e, 206)
		status = false
	})
	return userDetail, regRequest, status
}
