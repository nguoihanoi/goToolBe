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

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	permissionModel.InitModel(inDb, inRedisClient)
}

func searchUser(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
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
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

var userCmdMap = map[string]CommandHandler{
	"search": searchUser,
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
