package util

import (
	"fmt"
	"github.com/dromara/carbon/v2"
	"time"
)

var Loc = time.Local

func Init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(fmt.Errorf("时区错误: %w", err))
	}
	Loc = loc
}

func GetTime() time.Time {
	return time.Now().In(Loc)
}
func GetTimeString() string {
	return carbon.Now().ToDateTimeString()
}
