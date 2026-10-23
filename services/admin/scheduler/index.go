package scheduler

import (
	"fmt"
	"time"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	"github.com/robfig/cron/v3"
)

func InitPerSecond(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	languageModel.InitModel(inDb, inRedisClient)
	processUpdateCodes()
}
func processUpdateCodes() {
	c := cron.New() // Mặc định dùng chuẩn 5-field (phút, giờ, ngày, tháng, thứ)

	// Chạy mỗi phút
	c.AddFunc("* * * * *", func() {
		fmt.Println("Chạy mỗi phút:", time.Now())
	})

	c.Start()
	select {} // Giữ tiến trình không bị tắt
}
