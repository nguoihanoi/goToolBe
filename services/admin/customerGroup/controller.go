package customerGroup

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var response *libUtilities.ResponseClass

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	languageModel.InitModel(inDb, inRedisClient)
	customerModel.InitModel(inDb, inRedisClient, inJwtToken)
	response = libUtilities.Response(inRedisClient, "codes")
}

func search(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchCustomerGroupInput(ctx)
	if status {
		inLangCode := libUtilities.GetLangCode(ctx)
		filter := bSon.M{"delete": 0}
		inSortOrder := bSon.D{{Key: "delete", Value: 1}, {Key: "order", Value: 1}}
		if regRequest.Key != "" {
			regexValue := bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
			filter["$or"] = bSon.A{
				bSon.D{{Key: "name." + inLangCode, Value: regexValue}},
			}
		}
		if regRequest.Status > -1 {
			filter["status"] = regRequest.Status
		}
		inSortOrder = append(inSortOrder, bSon.E{Key: "name." + inLangCode, Value: 1})
		results, total := customerModel.SearchGroups(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		response.SendOutput(ctx, resp)
	}
}

func create(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Create false!", 206)
	regRequest, userId, status := ValidateCreateInput(ctx)
	if status {
		result := customerModel.CreateGroup(customerModel.CustomerGroup{
			Name:     regRequest.Name,
			Order:    regRequest.Order,
			Status:   regRequest.Status,
			AuthorId: userId,
		})
		if result != "" {
			resp.Status = true
			resp.Message = "Create success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func update(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Update false!", 206)
	regRequest, status := ValidateUpdateInput(ctx)
	if status {
		updateOption := bSon.M{"name": regRequest.Name, "order": regRequest.Order, "status": regRequest.Status}
		result := customerModel.UpdateGroup(regRequest.Id, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Update success!"
		}
		response.SendOutput(ctx, resp)
	}
}
func delete(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Delete false!", 206)
	regRequest, status := ValidateDeleteInput(ctx)
	if status {
		result := customerModel.DeleteGroup(regRequest.Id)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		response.SendOutput(ctx, resp)
	}
}

func gets(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(true, "Get success!", 200)
	results := customerModel.GetGroups()
	resp.Data = results
	response.SendOutput(ctx, resp)
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

var customerGroupCmdMap = map[string]CommandHandler{
	"search": search,
	"create": create,
	"update": update,
	"delete": delete,
	"gets":   gets,
}

func CustomerGroup(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := customerGroupCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
