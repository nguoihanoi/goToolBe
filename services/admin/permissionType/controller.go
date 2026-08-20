package permissionType

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var response *libUtilities.ResponseClass

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	languageModel.InitModel(inDb, inRedisClient)
	permissionModel.InitModel(inDb, inRedisClient)
	response = libUtilities.Response(inRedisClient, "codes")
}

func search(ctx *fastHttp.RequestCtx) {
	resp := response.GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchPermissionTypeInput(ctx)
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
		inSortOrder = append(inSortOrder, bSon.E{Key: "name." + inLangCode, Value: 1})
		results, total := permissionModel.SearchTypes(filter, inSortOrder, regRequest.Page, regRequest.Limit)
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
	regRequest, status := ValidateCreateInput(ctx)
	if status {
		result := permissionModel.CreateType(permissionModel.PermissionType{
			Name:     regRequest.Name,
			Order:    regRequest.Order,
			AuthorId: regRequest.UserId,
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
		updateOption := bSon.M{"name": regRequest.Name, "order": regRequest.Order}
		result := permissionModel.UpdateType(regRequest.Id, updateOption)
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
		result := permissionModel.DeleteType(regRequest.Id)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		response.SendOutput(ctx, resp)
	}
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

var permissionTypeCmdMap = map[string]CommandHandler{
	"search": search,
	"create": create,
	"update": update,
	"delete": delete,
}

func PermssionType(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := permissionTypeCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
