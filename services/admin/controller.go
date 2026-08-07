package main

import (
	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	customerModel "github.com/nguoihanoi/golang_shared/warehouses/customers"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	customerModel.InitModel(inDb, inRedisClient, inJwtToken)
}

func SearchCustomerGroup(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchCustomerGroupInput(ctx)
	if status {
		filter := bSon.M{"delete": 0}
		if regRequest.Key != "" {
			filter["name."+regRequest.LangCode] = bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
		}
		inSortOrder := bSon.D{{Key: "delete", Value: 1}}
		inSortOrder = append(inSortOrder, bSon.E{Key: "first_name", Value: 1})
		inSortOrder = append(inSortOrder, bSon.E{Key: "last_name", Value: 1})
		results, total := customerModel.SearchGroups(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

func SearchCustomer(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
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
		inSortOrder = append(inSortOrder, bSon.E{Key: "first_name", Value: 1})
		inSortOrder = append(inSortOrder, bSon.E{Key: "last_name", Value: 1})
		results, total := customerModel.Search(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

func SearchUser(ctx *fastHttp.RequestCtx) {
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
