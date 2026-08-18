package language

import (
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func searchGroup(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchGroupInput(ctx)
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
		results, total := languageModel.SearchGroupCode(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func createGroup(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Create false!", 206)
	regRequest, status := ValidateCreateGroupInput(ctx)
	if status {
		result := languageModel.CreateGroupCode(languageModel.GroupCode{
			Name:     regRequest.Name,
			Code:     regRequest.Code,
			AuthorId: regRequest.UserId,
		})
		if result != "" {
			resp.Status = true
			resp.Message = "Create success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func updateGroup(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Update false!", 206)
	regRequest, status := ValidateUpdateInput(ctx)
	if status {
		updateOption := bSon.M{"name": regRequest.Name, "code": regRequest.Code}
		result := languageModel.UpdateGroupCode(regRequest.Id, updateOption)
		if result == true {
			resp.Status = true
			resp.Message = "Update success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func deleteGroup(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Delete false!", 206)
	regRequest, status := ValidateDeleteGroupInput(ctx)
	if status {
		result := languageModel.DeleteGroupCode(regRequest.Id)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}
func getGroups(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(true, "Get success!", 200)
	results := languageModel.GetGroupCodes()
	resp.Data = results
	libUtilities.Response().SendOutput(ctx, resp)
}

var groupCmdMap = map[string]CommandHandler{
	"search": searchGroup,
	"create": createGroup,
	"update": updateGroup,
	"delete": deleteGroup,
	"gets":   getGroups,
}

func Group(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := groupCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
