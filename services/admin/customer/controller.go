package customer

import (
	"time"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var response *libUtilities.ResponseClass

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	customerModel.InitModel(inDb, inRedisClient, inJwtToken)
	response = libUtilities.Response(inRedisClient, "codes")
}
func createCustomer(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Create false!", 206)
	regRequest, status := ValidateCreateInput(ctx)
	if status {
		temBirthDay, _ := time.Parse("2006-01-02", regRequest.BirthDay)
		Password, PasswordHash := libUtilities.String().GetHashPassWord(regRequest.Password, "", false)
		result := customerModel.CreateCustomer(customerModel.Customer{
			CustomerGroupID: regRequest.GroupID,
			Email:           regRequest.Email,
			FirstName:       regRequest.FirstName,
			LastName:        regRequest.LastName,
			Gender:          regRequest.Gender,
			BirthDay:        temBirthDay,
			Password:        Password,
			PasswordHash:    PasswordHash,
			IsActive:        false,
		})
		if result != "" {
			resp.Status = true
			resp.Message = "Create success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func updateCustomer(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Update false!", 206)
	regRequest, status := ValidateUpdateInput(ctx)
	if status {
		temBirthDay, _ := time.Parse("2006-01-02", regRequest.BirthDay)
		updateOption := bSon.M{
			"customer_group_id": regRequest.GroupID,
			"email":             regRequest.Email,
			"first_name":        regRequest.FirstName,
			"last_name":         regRequest.LastName,
			"gender":            regRequest.Gender,
			"birth_day":         temBirthDay,
		}
		result := customerModel.UpdateCustomer(regRequest.ID, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Update success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func deleteCustomer(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Delete false!", 206)
	regRequest, status := ValidateDeleteInput(ctx)
	if status {
		result := customerModel.DeleteCustomer(regRequest.ID)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func updatePasswordById(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Update false!", 206)
	userDetail, regRequest, status := ValidateUpdatePasswordByIdInput(ctx)
	if status {
		Password, PasswordHash := libUtilities.String().GetHashPassWord(regRequest.Password, userDetail.PasswordHash, false)
		updateOption := bSon.M{"password_hash": PasswordHash, "password": Password}
		result := customerModel.UpdateCustomer(regRequest.ID, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Update success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func searchCustomer(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchCustomerInput(ctx)
	if status {
		filter := bSon.M{"delete": 0}
		inSortOrder := bSon.D{{Key: "delete", Value: 1}}
		if regRequest.GroupID != "" {
			filter["customer_group_id"] = regRequest.GroupID
			inSortOrder = append(inSortOrder, bSon.E{Key: "customer_group_id", Value: 1})
		}
		if regRequest.Key != "" {
			regexValue := bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
			filter["$or"] = bSon.A{
				bSon.D{{Key: "first_name", Value: regexValue}},
				bSon.D{{Key: "last_name", Value: regexValue}},
				bSon.D{{Key: "email", Value: regexValue}},
			}
		}
		if regRequest.Status > -1 {
			filter["status"] = regRequest.Status
		}
		inSortOrder = append(inSortOrder, bSon.E{Key: "first_name", Value: 1})
		inSortOrder = append(inSortOrder, bSon.E{Key: "last_name", Value: 1})
		results, total := customerModel.Search(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		response.SendOutput(ctx, resp)
	}
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

var customerCmdMap = map[string]CommandHandler{
	"search":             searchCustomer,
	"create":             createCustomer,
	"update":             updateCustomer,
	"delete":             deleteCustomer,
	"updatePasswordById": updatePasswordById,
}

func Customer(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := customerCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
