package util

import (
	"github.com/google/uuid"
	"strings"
)

func GetUuid() string {
	uuidStr := uuid.New().String() // 生成UUID作为文件名
	return strings.ReplaceAll(uuidStr, "-", "")
}
