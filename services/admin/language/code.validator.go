package language

import (
	"sync"

	libProcess "github.com/nguoihanoi/golang_shared/libs/process"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
)

type CreateCodeInput struct {
	Name     string            `validate:"required" json:"name"`
	Value    map[string]string `validate:"required" json:"value"`
	Code     string            `validate:"required" json:"code"`
	Type     int               `validate:"required" json:"type"`
	UserId   string            `validate:"required" json:"user_id"`
	LangCode string            `validate:"" json:"lang_code"`
}

func ValidateCreateCodeInput(ctx *fastHttp.RequestCtx) (regRequest CreateCodeInput, groupId string, status bool) {
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
			otherCodeDetail  languageModel.LanguageCode
			otherGroupDetail languageModel.GroupCode
			resultValidate   any
			statusValidate   int
		)
		groupId = ""
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
			langCode := languageModel.GetCodes()
			resultValidate, statusValidate = libUtilities.ValidateLangValue(regRequest.Value, false, langCode)
		}()
		go func() {
			defer wg.Done()
			otherCodeDetail = languageModel.GetCodeByName(regRequest.Name)
		}()
		go func() {
			defer wg.Done()
			otherGroupDetail = languageModel.CheckGroupCode(regRequest.Code, "")
			if otherGroupDetail.ID != "" {
				groupId = otherGroupDetail.ID
			}
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if otherCodeDetail.ID != "" {
			libUtilities.Response().SendError(ctx, "The lang code has already been used.", nil, 206)
		}
		switch statusValidate {
		case 1, 2:
			libUtilities.Response().SendError(ctx, "Invalid input data!", resultValidate, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, groupId, status
}

type UpdateCodeInput struct {
	Id       string            `validate:"required" json:"_id"`
	Name     string            `validate:"required" json:"name"`
	Value    map[string]string `validate:"required" json:"value"`
	Code     string            `validate:"required" json:"code"`
	Type     int               `validate:"required" json:"type"`
	UserId   string            `validate:"required" json:"user_id"`
	LangCode string            `validate:"" json:"lang_code"`
}

func ValidateUpdateCodeInput(ctx *fastHttp.RequestCtx) (regRequest UpdateCodeInput, groupId string, status bool) {
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
			otherCodeDetail  languageModel.LanguageCode
			codeDetail       languageModel.LanguageCode
			otherGroupDetail languageModel.GroupCode
			resultValidate   any
			statusValidate   int
		)
		groupId = ""
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
			langCode := languageModel.GetCodes()
			resultValidate, statusValidate = libUtilities.ValidateLangValue(regRequest.Value, false, langCode)
		}()
		go func() {
			defer wg.Done()
			otherCodeDetail = languageModel.CheckCodeByName(regRequest.Name, regRequest.Id)
		}()
		go func() {
			defer wg.Done()
			otherGroupDetail = languageModel.CheckGroupCode(regRequest.Code, "")
			if otherGroupDetail.ID != "" {
				groupId = otherGroupDetail.ID
			}
		}()
		go func() {
			defer wg.Done()
			codeDetail = languageModel.GetCodeById(regRequest.Id, true)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if otherCodeDetail.ID != "" {
			libUtilities.Response().SendError(ctx, "The code has already been used.", nil, 206)
		}
		if codeDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This code information does not exist in the system.", nil, 206)
		}
		switch statusValidate {
		case 1, 2:
			libUtilities.Response().SendError(ctx, "Invalid input data!", resultValidate, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, groupId, status
}

type DeleteCodeInput struct {
	Id       string `validate:"required" json:"_id"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateDeleteCodeInput(ctx *fastHttp.RequestCtx) (regRequest DeleteCodeInput, status bool) {
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
			wg          sync.WaitGroup
			userDetail  userModel.User
			groupDetail languageModel.LanguageCode
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
			groupDetail = languageModel.GetCodeById(regRequest.Id, true)
		}()
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if groupDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This group information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}

type SearchCodeInput struct {
	Key      string `validate:"" json:"key"`
	Type     int    `validate:"min=0" json:"type"`
	GroupId  string `validate:"" json:"group_id"`
	Page     int64  `validate:"min=1" json:"page"`
	Limit    int64  `validate:"min=0" json:"limit"`
	UserId   string `validate:"required" json:"user_id"`
	LangCode string `validate:"" json:"lang_code"`
}

func ValidateSearchCodeInput(ctx *fastHttp.RequestCtx) (regRequest SearchCodeInput, status bool) {
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
			wg          sync.WaitGroup
			userDetail  userModel.User
			groupDetail languageModel.GroupCode
		)
		hasGroupID := regRequest.GroupId != ""
		if hasGroupID {
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
		if hasGroupID {
			go func() {
				defer wg.Done()
				groupDetail = languageModel.GetGroupCodeById(regRequest.GroupId, true)
			}()
		}
		wg.Wait()
		if userDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "You do not have permission to perform this function.", nil, 206)
			return
		}
		if hasGroupID && groupDetail.ID == "" {
			libUtilities.Response().SendError(ctx, "This group information does not exist in the system.", nil, 206)
		}
		status = true
	}).Catch(func(e libProcess.E) {
		libUtilities.Response().SendError(ctx, "Invalid input data!", e, 206)
	})
	return regRequest, status
}
