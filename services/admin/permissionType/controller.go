package permissionType

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	permissionModel "github.com/nguoihanoi/golang_shared/warehouses/permissions"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	permissionModel.InitModel(inDb, inRedisClient)
}

func search(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchPermissionTypeInput(ctx)
	if status {
		filter := bSon.M{"delete": 0}
		inSortOrder := bSon.D{{Key: "delete", Value: 1}, {Key: "order", Value: 1}}
		if regRequest.Key != "" {
			regexValue := bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
			filter["$or"] = bSon.A{
				bSon.D{{Key: "name." + regRequest.LangCode, Value: regexValue}},
			}
		}
		inSortOrder = append(inSortOrder, bSon.E{Key: "name." + regRequest.LangCode, Value: 1})
		results, total := permissionModel.SearchTypes(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

func create(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Create false!", 206)
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
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func update(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Update false!", 206)
	regRequest, status := ValidateUpdateInput(ctx)
	if status {
		updateOption := bSon.M{"name": regRequest.Name, "order": regRequest.Order}
		result := permissionModel.UpdateType(regRequest.Id, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Update success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func delete(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Delete false!", 206)
	regRequest, status := ValidateDeleteInput(ctx)
	if status {
		result := permissionModel.DeleteType(regRequest.Id)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
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
