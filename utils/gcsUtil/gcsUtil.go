package gcsUtil

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

// GCSにファイルをアップロードする
func UploadGCS(file *multipart.FileHeader) (string, error) {
	//アップロード先バケット名取得
	bucketName := os.Getenv("GCS_BUCKET")

	//ファイル名取得
	fileName := file.Filename

	//ファイル名から拡張子を取得
	extension := filepath.Ext(fileName)

	//uuidでファイル名を生成
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	//保存ファイル名
	uploadFileName := uid.String() + extension

	//GCSクライアントの初期化
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	defer client.Close()
	if err != nil {
		return "", err
	}

	//GCSバケットオブジェクトの取得
	bucket := client.Bucket(bucketName)

	//ファイルを開く
	tempFile, err := file.Open()
	defer tempFile.Close()
	if err != nil {
		return "", err
	}

	//GCSにアップロード
	object := bucket.Object(uploadFileName)
	writer := object.NewWriter(ctx)
	if _, err := io.Copy(writer, tempFile); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	return uploadFileName, err
}

// GCSのファイルを削除する
func DeleteGCS(path string) error {
	//アップロード先バケット名取得
	bucketName := os.Getenv("GCS_BUCKET")

	//GCSクライアントの初期化
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	defer client.Close()
	if err != nil {
		return err
	}

	//GCSバケットオブジェクトの取得
	bucket := client.Bucket(bucketName)

	//削除
	err = bucket.Object(path).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetUrlGCS(path string, minute time.Duration) (string, error) {
	//アップロード先バケット名取得
	bucketName := os.Getenv("GCS_BUCKET")

	//GCSクライアントの初期化
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	defer client.Close()
	if err != nil {
		return "", err
	}

	//GCSバケットオブジェクトの取得
	bucket := client.Bucket(bucketName)

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(minute * time.Minute),
	}

	url, err := bucket.SignedURL(path, opts)
	if err != nil {
		return "", err
	}

	return url, nil
}
