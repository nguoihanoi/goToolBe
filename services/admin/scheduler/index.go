package scheduler

import (
	"fmt"
	"time"

	libCache "github.com/nguoihanoi/golang_shared/libs/cache"
	libDb "github.com/nguoihanoi/golang_shared/libs/database"
	languageModel "github.com/nguoihanoi/golang_shared/warehouses/languages"
	"github.com/robfig/cron/v3"
)

var newCron *cron.Cron

func InitPerMinute(inDb *libDb.DatabaseClass, inRedisClient *libCache.Cache, inJwtToken string) {
	newCron = cron.New()
	languageModel.InitModel(inDb, inRedisClient)

	// Chạy mỗi phút
	newCron.AddFunc("* * * * *", func() {
		processUpdateCodes(inRedisClient)
	})

	newCron.Start()
	select {} // Giữ tiến trình không bị tắt
}
func processUpdateCodes(inRedisClient *libCache.Cache) {
	fmt.Println("Chạy mỗi phút:", time.Now())
	queueNewCodes := inRedisClient.NewQueue("queue_new_codes")
	total := queueNewCodes.Count()
	if total > 0 {
		return
	}
	languages := languageModel.Gets()
	for i := 0; i < int(total); i += 1 {
		item := queueNewCodes.Pop()
		if item != nil {
			code := item.(string)
			processUpdateCode(code, languages)
		}
	}
}
func processUpdateCode(code string, languages []languageModel.Language) {
	codeDetail := languageModel.GetCodeByName(code)
	if codeDetail.ID == "" {
		temValue := make(map[string]string)
		for i := range languages {
			temValue[languages[i].Code] = code
		}
		languageModel.CreateCode(languageModel.LanguageCode{
			Name:  code,
			Value: temValue,
		})
	}
}
