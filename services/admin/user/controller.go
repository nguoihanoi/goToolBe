package user

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var response *libUtilities.ResponseClass

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	permissionModel.InitModel(inDb, inRedisClient)
	response = libUtilities.Response(inRedisClient, "codes")
}

func createUser(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Create false!", 206)
	regRequest, status := ValidateCreateInput(ctx)
	if status {
		Password, PasswordHash := libUtilities.String().GetHashPassWord(regRequest.Password, "", false)
		result := userModel.CreateUser(userModel.User{
			AccountType:  regRequest.AccountType,
			Email:        regRequest.Email,
			FirstName:    regRequest.FirstName,
			LastName:     regRequest.LastName,
			Password:     Password,
			PasswordHash: PasswordHash,
			IsActive:     false,
		})
		if result != "" {
			resp.Status = true
			resp.Message = "Create success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func updateUser(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Update false!", 206)
	regRequest, status := ValidateUpdateInput(ctx)
	if status {
		updateOption := bSon.M{
			"account_type": regRequest.AccountType,
			"email":        regRequest.Email,
			"first_name":   regRequest.FirstName,
			"last_name":    regRequest.LastName,
		}
		result := userModel.UpdateUser(regRequest.ID, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Update success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func deleteUser(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Delete false!", 206)
	regRequest, status := ValidateDeleteInput(ctx)
	if status {
		result := userModel.DeleteUser(regRequest.ID)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func activeUser(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Active false!", 206)
	userDetail, status := ValidateActiveInput(ctx)
	if status {
		updateOption := bSon.M{"is_active": true}
		if userDetail.IsActive == true {
			updateOption["is_active"] = false
		}
		result := userModel.UpdateUser(userDetail.ID, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Active success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func updatePassword(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Active false!", 206)
	userDetail, regRequest, status := ValidateUpdatePasswordInput(ctx)
	if status {
		Password, PasswordHash := libUtilities.String().GetHashPassWord(regRequest.Password, userDetail.PasswordHash, false)
		updateOption := bSon.M{"password_hash": PasswordHash, "password": Password}
		result := userModel.UpdateUser(regRequest.UserId, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Active success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func updatePasswordById(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Active false!", 206)
	userDetail, regRequest, status := ValidateUpdatePasswordByIdInput(ctx)
	if status {
		Password, PasswordHash := libUtilities.String().GetHashPassWord(regRequest.Password, userDetail.PasswordHash, false)
		updateOption := bSon.M{"password_hash": PasswordHash, "password": Password}
		result := userModel.UpdateUser(regRequest.UserId, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Active success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func searchUser(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchUserInput(ctx)
	if status {
		filter := bSon.M{"delete": 0}
		inSortOrder := bSon.D{{Key: "delete", Value: 1}}
		if regRequest.AccountTypeId != "" {
			filter["account_type"] = regRequest.AccountTypeId
			inSortOrder = append(inSortOrder, bSon.E{Key: "account_type", Value: 1})
		}
		if regRequest.Key != "" {
			regexValue := bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
			filter["$or"] = bSon.A{
				bSon.D{{Key: "first_name", Value: regexValue}},
				bSon.D{{Key: "last_name", Value: regexValue}},
				bSon.D{{Key: "email", Value: regexValue}},
			}
		}
		inSortOrder = append(inSortOrder, bSon.E{Key: "first_name", Value: 1})
		inSortOrder = append(inSortOrder, bSon.E{Key: "last_name", Value: 1})
		results, total := userModel.Search(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		response.SendOutput(ctx, resp)
	}
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

var userCmdMap = map[string]CommandHandler{
	"search":             searchUser,
	"create":             createUser,
	"update":             updateUser,
	"delete":             deleteUser,
	"active":             activeUser,
	"updatePassword":     updatePassword,
	"updatePasswordById": updatePasswordById,
}

func User(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := userCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
