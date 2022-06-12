package utils

import (
	post "piepay/controllers/POST"
	"time"
)

var timer *time.Ticker

func Uploader() {
	if timer == nil {
		timer = time.NewTicker(10 * time.Second)
		go func() {
			if timer.C != nil {
				post.UploadVideoMetaData()
			}
		}()
	}
}
