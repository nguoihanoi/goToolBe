package language

import (
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func searchCode(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchCodeInput(ctx)
	if status {
		filter := bSon.M{"delete": 0}
		inSortOrder := bSon.D{{Key: "delete", Value: 1}, {Key: "type", Value: 1}}
		if regRequest.Type > 0 {
			filter["type"] = regRequest.Type
		}
		if regRequest.GroupId != "" {
			filter["group_id"] = regRequest.GroupId
		}
		if regRequest.Key != "" {
			regexValue := bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
			filter["$or"] = bSon.A{
				bSon.D{{Key: "name", Value: regexValue}},
			}
		}
		inSortOrder = append(inSortOrder, bSon.E{Key: "name", Value: 1})
		results, total := languageModel.SearchCode(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func createCode(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Create false!", 206)
	regRequest, groupId, status := ValidateCreateCodeInput(ctx)
	if status {
		result := languageModel.CreateCode(languageModel.LanguageCode{
			Name:     regRequest.Name,
			Value:    regRequest.Value,
			Type:     regRequest.Type,
			GroupId:  groupId,
			AuthorId: regRequest.UserId,
		})
		if result != "" {
			resp.Status = true
			resp.Message = "Create success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func updateCode(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Update false!", 206)
	regRequest, groupId, status := ValidateUpdateCodeInput(ctx)
	if status {
		updateOption := bSon.M{"name": regRequest.Name, "value": regRequest.Value, "type": regRequest.Type, "group_id": groupId}
		result := languageModel.UpdateCode(regRequest.Id, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Update success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func deleteCode(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Delete false!", 206)
	regRequest, status := ValidateDeleteCodeInput(ctx)
	if status {
		result := languageModel.DeleteCode(regRequest.Id)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

var codeCmdMap = map[string]CommandHandler{
	"search": searchCode,
	"create": createCode,
	"update": updateCode,
	"delete": deleteCode,
}

func Code(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := codeCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
