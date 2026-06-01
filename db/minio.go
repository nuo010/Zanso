package db

import (
	"context"
	"fmt"
	"io"
	"strings"
	"zanso/util"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

var (
	MinioClient *minio.Client
	MinioBucket string
)

func UseMinioStorage() bool {
	provider := strings.TrimSpace(strings.ToLower(viper.GetString("storage.provider")))
	return provider == "minio" || viper.GetBool("minio.enabled")
}

func MinioInit() {
	if !UseMinioStorage() {
		return
	}

	endpoint := strings.TrimSpace(viper.GetString("minio.endpoint"))
	accessKey := strings.TrimSpace(viper.GetString("minio.access_key"))
	secretKey := strings.TrimSpace(viper.GetString("minio.secret_key"))
	bucket := strings.TrimSpace(viper.GetString("minio.bucket"))
	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		util.Log().Panic("MinIO 配置不完整，请检查 minio.endpoint/access_key/secret_key/bucket")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: viper.GetBool("minio.secure"),
	})
	if err != nil {
		util.Log().Panic("初始化 MinIO 客户端失败: %v", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		util.Log().Panic("检查 MinIO bucket 失败: %v", err)
	}
	if !exists {
		if !viper.GetBool("minio.auto_create_bucket") {
			util.Log().Panic("MinIO bucket 不存在: %s", bucket)
		}
		if err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
			Region: strings.TrimSpace(viper.GetString("minio.region")),
		}); err != nil {
			util.Log().Panic("创建 MinIO bucket 失败: %v", err)
		}
	}

	MinioClient = client
	MinioBucket = bucket
}

func UploadMinioObject(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if MinioClient == nil || MinioBucket == "" {
		return "", fmt.Errorf("MinIO 未初始化")
	}
	objectName = strings.TrimLeft(strings.ReplaceAll(objectName, "\\", "/"), "/")
	if objectName == "" {
		return "", fmt.Errorf("MinIO objectName 不能为空")
	}
	_, err := MinioClient.PutObject(ctx, MinioBucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func DeleteMinioObject(ctx context.Context, objectName string) error {
	if MinioClient == nil || MinioBucket == "" {
		return fmt.Errorf("MinIO 未初始化")
	}
	objectName = strings.TrimLeft(strings.ReplaceAll(objectName, "\\", "/"), "/")
	if objectName == "" {
		return nil
	}
	return MinioClient.RemoveObject(ctx, MinioBucket, objectName, minio.RemoveObjectOptions{})
}

func BuildMinioObjectURL(objectName string) string {
	objectName = strings.TrimLeft(strings.ReplaceAll(objectName, "\\", "/"), "/")
	publicBaseURL := strings.TrimRight(strings.TrimSpace(viper.GetString("minio.public_base_url")), "/")
	if publicBaseURL != "" {
		return publicBaseURL + "/" + objectName
	}

	scheme := "http"
	if viper.GetBool("minio.secure") {
		scheme = "https"
	}
	endpoint := strings.TrimRight(strings.TrimSpace(viper.GetString("minio.endpoint")), "/")
	return fmt.Sprintf("%s://%s/%s/%s", scheme, endpoint, strings.Trim(MinioBucket, "/"), objectName)
}
