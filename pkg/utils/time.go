package utils

import (
	"time"
)

func Today() string {
	// Load múi giờ UTC+7 (có thể dùng location cố định hoặc IANA)
	location := time.FixedZone("UTC+7", 7*60*60)

	// Lấy thời gian hiện tại theo UTC+7
	now := time.Now().In(location)

	// Định dạng chính xác đến phút: YYYYMMDDHHMM
	return now.Format("20060102")
}
