package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"zanso/db"
	"zanso/model"
	"zanso/result"
	"zanso/util"
)

func GetTblFileList(c *gin.Context) {
	var fileList []model.TblFile
	var count int64
	var pageQuery = model.PageQuery{}
	c.Bind(&pageQuery)
	db.DB.Model(&model.TblFile{}).Where("type = ?", pageQuery.Type).Count(&count).Order("creat_time DESC").Offset((pageQuery.PageNum - 1) * pageQuery.PageSize).Limit(pageQuery.PageSize).Find(&fileList)
	result.OkSetData(c, &model.PageResult{List: fileList, Total: count})
}

func UploadFile(c *gin.Context) {
	// 1. 打开文件流
	file, _ := c.FormFile("file")
	reader, err := file.Open()
	if err != nil {
		log.Println("Error opening file", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "file not exists",
		})
		return
	}
	defer reader.Close()

	//bucketName := "lgl"
	// 获取文件扩展名

	var filelj = "/root/data/i/"

	ext := filepath.Ext(file.Filename)

	//contentType := file.Header.Get("Content-Type")
	minioFileName := getFilePath("", ext)
	// 3. 通过 minio 客户端上传文件
	//_, err = db.MinioClient.PutObject(context.Background(), bucketName, minioFileName, reader, file.Size, minio.PutObjectOptions{
	//	ContentType: contentType,
	//})
	var fileName = filelj + "/" + minioFileName

	createFolderPath(filelj)

	c.SaveUploadedFile(file, fileName)

	//httpUrl := "https://ligl.top/file/" + bucketName + "/" + minioFileName
	httpUrl := "https://ligl.top/i/" + minioFileName

	db.DB.Save(&model.TblFile{ID: util.GetUuid(), Key: httpUrl, CreatTime: util.GetTime(), Fsize: file.Size, Type: "电脑保存图片"})
	result.OkSetData(c, httpUrl)
}
func createFolderPath(path string) error {
	// 使用 MkdirAll 递归创建多级文件夹
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("创建文件夹路径失败: %v", err)
	}
	return nil
}

func UploadFile2(c *gin.Context) {
	// 1. 打开文件流
	file, _ := c.FormFile("file")
	reader, err := file.Open()
	if err != nil {
		log.Println("Error opening file", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "file not exists",
		})
		return
	}
	defer reader.Close()

	bucketName := "camera"
	// 获取文件扩展名
	ext := filepath.Ext(file.Filename)

	contentType := file.Header.Get("Content-Type")
	minioFileName := getFilePath("", ext)
	// 3. 通过 minio 客户端上传文件
	_, err = db.MinioClient.PutObject(context.Background(), bucketName, minioFileName, reader, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	httpUrl := "https://ligl.top/file/" + bucketName + "/" + minioFileName
	db.DB.Save(&model.TblFile{ID: util.GetUuid(), Key: httpUrl, CreatTime: util.GetTime(), Fsize: file.Size, Type: "手动保存图片"})
	result.OkSetData(c, "https://ligl.top/file/"+bucketName+"/"+minioFileName)
}

// getFilePath 生成文件路径，包括年月日文件夹和文件名（如果文件名为空，则生成UUID作为文件名）
func getFilePath(fileName, fileType string) string {
	// 获取当前时间
	now := time.Now()
	year, month, _ := now.Date()

	// 格式化月份（Go中月份是从0开始的，所以需要+1）
	monthStr := fmt.Sprintf("%02d", month)

	// 构建路径
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d/%s/", year, monthStr))

	// 如果文件名为空，则生成UUID作为文件名
	if fileName == "" || strings.TrimSpace(fileName) == "" {
		fileName = util.GetUuid()[:5]
	}

	// 添加文件类型和文件名到路径
	sb.WriteString(fmt.Sprintf("%s%s", fileName, fileType))
	return sb.String()
}
