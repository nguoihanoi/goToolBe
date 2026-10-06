package language

import (
	"time"

	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	fileModel "github.com/nguoihanoi/golang_shared/warehouses/files"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	fastHttp "github.com/valyala/fasthttp"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

type LanguageItem struct {
	ID        string    `bson:"_id" json:"_id"`
	Delete    int       `bson:"delete" json:"delete"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Name      string    `bson:"name" json:"name"`
	Code      string    `bson:"code" json:"code"`
	Image     string    `bson:"image" json:"image"`
	ImageLink string    `bson:"image_link" json:"image_link"`
	Order     int       `bson:"order" json:"order"`
	Status    int       `bson:"status" json:"status"`
	AuthorId  string    `bson:"author_id" json:"author_id"`
}

func formatLanguageItem(ctx *fastHttp.RequestCtx, item languageModel.Language) LanguageItem {
	temImageLink := ""
	if item.Image != "" {
		fileDetail := fileModel.GetById(item.Image, true)
		if fileDetail.ID != "" {
			temImageLink = libUtilities.Request().GetBaseURL(ctx, "image/"+fileDetail.S3Key)
		}
	}
	return LanguageItem{
		ID:        item.ID,
		Delete:    item.Delete,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		Name:      item.Name,
		Code:      item.Code,
		Image:     item.Image,
		ImageLink: temImageLink,
		Order:     item.Order,
		Status:    item.Status,
		AuthorId:  item.AuthorId,
	}
}

func doSearch(ctx *fastHttp.RequestCtx, regRequest SearchInput) ([]LanguageItem, int64) {
	filter := bSon.M{"delete": 0}
	inSortOrder := bSon.D{{Key: "delete", Value: 1}}
	if regRequest.Key != "" {
		regexValue := bSon.D{{Key: "$regex", Value: regRequest.Key}, {Key: "$options", Value: "i"}}
		filter["$or"] = bSon.A{
			bSon.D{{Key: "name", Value: regexValue}},
			bSon.D{{Key: "code", Value: regexValue}},
		}
	}
	if regRequest.Status > -1 {
		filter["status"] = regRequest.Status
	}
	inSortOrder = append(inSortOrder, bSon.E{Key: "name", Value: 1})
	output := []LanguageItem{}
	results, total := languageModel.Search(filter, inSortOrder, regRequest.Page, regRequest.Limit)
	for i := range results {
		temLanguageItem := formatLanguageItem(ctx, results[i])
		output = append(output, temLanguageItem)
	}
	return output, total
}
