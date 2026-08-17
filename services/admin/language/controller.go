package language

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	languageModel.InitModel(inDb, inRedisClient)
}

func search(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchInput(ctx)
	if status {
		filter := bSon.M{"delete": 0}
		inSortOrder := bSon.D{{Key: "delete", Value: 1}}
		if regRequest.Key != "" {
			regexValue := bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
			filter["$or"] = bSon.A{
				bSon.D{{Key: "name", Value: regexValue}},
				bSon.D{{Key: "code", Value: regexValue}},
			}
		}
		inSortOrder = append(inSortOrder, bSon.E{Key: "name", Value: 1})
		results, total := languageModel.Search(filter, inSortOrder, regRequest.Page, regRequest.Limit)
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
		result := languageModel.Create(languageModel.Language{
			Name:     regRequest.Name,
			Code:     regRequest.Code,
			Image:    regRequest.Image,
			Order:    regRequest.Order,
			Status:   regRequest.Status,
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
		updateOption := bSon.M{"name": regRequest.Name, "code": regRequest.Code, "image": regRequest.Image, "order": regRequest.Order, "status": regRequest.Status}
		result := languageModel.Update(regRequest.Id, updateOption)
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
		result := languageModel.Delete(regRequest.Id)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

func gets(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(true, "Get success!", 200)
	results := languageModel.Gets()
	resp.Data = results
	libUtilities.Response().SendOutput(ctx, resp)
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

var languageCmdMap = map[string]CommandHandler{
	"search": search,
	"create": create,
	"update": update,
	"delete": delete,
	"gets":   gets,
}

func Language(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := languageCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
