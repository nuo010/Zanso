package db

import (
	"context"
	"fmt"
	"io"
	"net/url"
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

	endpoint := normalizeMinioEndpoint(strings.TrimSpace(viper.GetString("minio.endpoint")))
	accessKey := strings.TrimSpace(viper.GetString("minio.accessKey"))
	if accessKey == "" {
		accessKey = strings.TrimSpace(viper.GetString("minio.access_key"))
	}
	secretKey := strings.TrimSpace(viper.GetString("minio.secretKey"))
	if secretKey == "" {
		secretKey = strings.TrimSpace(viper.GetString("minio.secret_key"))
	}
	bucket := strings.TrimSpace(viper.GetString("minio.bucket"))
	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		util.Log().Panic("MinIO 配置不完整，请检查 minio.endpoint/accessKey/secretKey/bucket")
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
	objectName = TrimMinioBucketPrefix(objectName)
	if objectName == "" {
		return nil
	}
	return MinioClient.RemoveObject(ctx, MinioBucket, objectName, minio.RemoveObjectOptions{})
}

func BuildMinioStoragePath(objectName string) string {
	objectName = strings.TrimLeft(strings.ReplaceAll(objectName, "\\", "/"), "/")
	bucket := strings.Trim(MinioBucket, "/")
	if bucket == "" || strings.HasPrefix(objectName, bucket+"/") {
		return objectName
	}
	return bucket + "/" + objectName
}

func TrimMinioBucketPrefix(objectName string) string {
	objectName = strings.TrimLeft(strings.ReplaceAll(objectName, "\\", "/"), "/")
	bucket := strings.Trim(MinioBucket, "/")
	return strings.TrimPrefix(objectName, bucket+"/")
}

func BuildMinioObjectURL(objectName string) string {
	objectName = BuildMinioStoragePath(objectName)
	downloadBase := strings.TrimRight(strings.TrimSpace(viper.GetString("minio.downloadBase")), "/")
	if downloadBase == "" {
		downloadBase = strings.TrimRight(strings.TrimSpace(viper.GetString("minio.public_base_url")), "/")
	}
	if downloadBase != "" {
		return downloadBase + "/" + objectName
	}

	scheme, endpoint := minioSchemeAndEndpoint()
	return fmt.Sprintf("%s://%s/%s/%s", scheme, endpoint, strings.Trim(MinioBucket, "/"), objectName)
}

func normalizeMinioEndpoint(endpoint string) string {
	if endpoint == "" {
		return ""
	}
	if parsed, err := url.Parse(endpoint); err == nil && parsed.Host != "" {
		return parsed.Host
	}
	return strings.TrimRight(endpoint, "/")
}

func minioSchemeAndEndpoint() (string, string) {
	rawEndpoint := strings.TrimSpace(viper.GetString("minio.endpoint"))
	if parsed, err := url.Parse(rawEndpoint); err == nil && parsed.Host != "" {
		return parsed.Scheme, parsed.Host
	}
	scheme := "http"
	if viper.GetBool("minio.secure") {
		scheme = "https"
	}
	return scheme, strings.TrimRight(rawEndpoint, "/")
}
