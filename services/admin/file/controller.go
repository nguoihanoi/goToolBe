package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	libUtilities "github.com/nguoihanoi/golang_shared/libs/utilities"
	fileModel "github.com/nguoihanoi/golang_shared/warehouses/files"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	userModel "github.com/nguoihanoi/golang_shared/warehouses/users"
	fastHttp "github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	bSon "go.mongodb.org/mongo-driver/v2/bson"
)

var storageDir string

func InitController(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string, inFolder string) {
	userModel.InitModel(inDb, inRedisClient, inJwtToken)
	languageModel.InitModel(inDb, inRedisClient)
	fileModel.InitModel(inDb, inRedisClient)
	storageDir = inFolder
}

func search(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Search false!", 206)
	regRequest, status := ValidateSearchFileInput(ctx)
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
		results, total := fileModel.Search(filter, inSortOrder, regRequest.Page, regRequest.Limit)
		if total > 0 {
			resp.Status = true
			resp.Message = "Search success!"
		}
		resp.Data = map[string]any{"list": results, "total": total}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

func UploadFile(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Upload false!", 206)
	// Process multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		libUtilities.Response().SendError(ctx, "Failed to parse form.", nil, 206)
		return
	}
	// Get uploaded file from form
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		libUtilities.Response().SendError(ctx, "No file uploaded.", nil, 206)
		return
	}
	fileHeader := fileHeaders[0]
	file, err := fileHeader.Open()
	if err != nil {
		libUtilities.Response().SendError(ctx, "Failed to open uploaded file.", nil, 206)
		return
	}
	defer file.Close()
	fileExt := filepath.Ext(fileHeader.Filename)

	// Create file metadata
	curTime := time.Now()
	fileMetadata := fileModel.File{
		ID:            primitive.NewObjectID().Hex(),
		Name:          fileHeader.Filename,
		Size:          fileHeader.Size,
		Type:          fileHeader.Header.Get("Content-Type"),
		StoredLocally: true,
	}
	// Store locally
	newFolderName := strings.Join(strings.Split(libUtilities.Time().FormatDate(curTime), "-"), "_")
	localDir := storageDir + "/" + newFolderName
	// Ensure uploads directory exists
	if err := os.MkdirAll(localDir, 0755); err != nil {
		libUtilities.Response().SendError(ctx, "Failed to create upload1s directory.", nil, 206)
		return
	}
	localFilePath := fmt.Sprintf("%s/%s%s", localDir, fileMetadata.ID, fileExt)
	fileMetadata.LocalPath = localFilePath
	// Save file to local path
	dst, err := os.Create(localFilePath)
	if err != nil {
		libUtilities.Response().SendError(ctx, "Failed to create local file.", nil, 206)
		return
	}
	defer dst.Close()
	// Reset file position to beginning
	if _, err = file.Seek(0, 0); err != nil {
		libUtilities.Response().SendError(ctx, "Failed to reset file position.", nil, 206)
		return
	}
	// Copy file content to destination
	if _, err = io.Copy(dst, file); err != nil {
		libUtilities.Response().SendError(ctx, "Failed to save file locally.", nil, 206)
		return
	}
	// Parse description and tags if provided
	if description := form.Value["description"]; len(description) > 0 {
		fileMetadata.Description = description[0]
	}
	if tags := form.Value["tags"]; len(tags) > 0 {
		var tagList []string
		if err := json.Unmarshal([]byte(tags[0]), &tagList); err == nil {
			fileMetadata.Tags = tagList
		}
	}
	//
	result := fileModel.Create(fileMetadata)
	if result != "" {
		resp.Status = true
		resp.Message = "Upload success!"
		resp.Data = fileMetadata.ID
	}
	libUtilities.Response().SendOutput(ctx, resp)
}
func DownloadFile(ctx *fastHttp.RequestCtx) {
	fileName := ctx.UserValue("id").(string)
	cleanFileName := filepath.Base(fileName)
	fileId := strings.Split(cleanFileName, ".")[0]
	fileDetail := fileModel.GetById(fileId, true)
	if fileDetail.ID != "" {
		filePath := fileDetail.LocalPath
		// 2. Kiểm tra xem File có tồn tại hay không
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			ctx.SetStatusCode(fastHttp.StatusNotFound)
			ctx.SetBodyString("File not found")
			return
		}
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileDetail.Name))
		ctx.Response.Header.Set("Content-Type", fileDetail.Type)
		ctx.Response.Header.Set("Content-Length", strconv.FormatInt(fileDetail.Size, 10))
		fastHttp.ServeFile(ctx, filePath)
	} else {
		ctx.SetStatusCode(fastHttp.StatusNotFound)
		ctx.SetBodyString("File not found")
		return
	}
}
func ViewImage(ctx *fastHttp.RequestCtx) {
	fileName := ctx.UserValue("id").(string)
	cleanFileName := filepath.Base(fileName)
	fileId := strings.Split(cleanFileName, ".")[0]
	fileDetail := fileModel.GetById(fileId, true)
	if fileDetail.ID != "" {
		filePath := fileDetail.LocalPath
		imageData, err := os.ReadFile(filePath)
		if err != nil {
			ctx.SetStatusCode(fastHttp.StatusNotFound)
			ctx.SetBodyString("File not found")
			return
		}
		// Set appropriate Content-Type header
		ctx.SetContentType(fileDetail.Type)

		// Set cache control for browser optimization (optional)
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")

		// Write raw image bytes
		ctx.SetBody(imageData)
	} else {
		ctx.SetStatusCode(fastHttp.StatusNotFound)
		ctx.SetBodyString("File not found")
		return
	}
}
func update(ctx *fastHttp.RequestCtx) {
	resp := libUtilities.Response().GetOutput(false, "Update false!", 206)
	regRequest, status := ValidateUpdateInput(ctx)
	if status {
		updateOption := bSon.M{"name": regRequest.Name, "order": regRequest.Order}
		result := fileModel.Update(regRequest.Id, updateOption)
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
		result := fileModel.Delete(regRequest.Id)
		if result == true {
			resp.Status = true
			resp.Message = "Delete success!"
		}
		libUtilities.Response().SendOutput(ctx, resp)
	}
}

type CommandHandler func(ctx *fastHttp.RequestCtx)

var fileCmdMap = map[string]CommandHandler{
	"search": search,
	"update": update,
	"delete": delete,
}

func File(ctx *fastHttp.RequestCtx) {
	cmd := ctx.Request.Header.Peek("X-API-Cmd")

	// Ép kiểu string ở đây nếu dùng Map, nhưng tra cứu O(1) rất nhanh
	if handler, exists := fileCmdMap[string(cmd)]; exists {
		handler(ctx)
		return
	}

	ctx.SetStatusCode(fastHttp.StatusBadRequest)
	ctx.SetBodyString("Invalid or missing X-API-Cmd header")
}
