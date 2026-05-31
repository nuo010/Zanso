package db

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"log"
)

var MinioClient *minio.Client

func MinioInit() {
	// 初始化 Minio 客户端
	endpoint := viper.GetString("minio.endpoint") + ":" + viper.GetString("minio.port")
	accessKeyID := viper.GetString("minio.accessKey")
	secretAccessKey := viper.GetString("minio.secretKey")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	MinioClient = minioClient
}
