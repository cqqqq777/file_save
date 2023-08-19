package file_save

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type SaveFunc func(file *multipart.FileHeader) error

func NewLocalSaveFunc(dst string) SaveFunc {
	return func(file *multipart.FileHeader) error {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, src)
		return err
	}
}

func NewS3Client(s3 *S3) (*minio.Client, error) {
	client, err := minio.New(s3.Endpoint, &minio.Options{
		Region: s3.Region,
		Creds:  credentials.NewStaticV4(s3.Credentials.AccessKeyID, s3.Credentials.SecretAccessKey, s3.Credentials.SessionToken),
		Secure: s3.Secure,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewS3SaveFunc(ctx context.Context, cli *minio.Client, bucketName, objectName string) SaveFunc {
	return func(file *multipart.FileHeader) error {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = cli.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{})
		return err
	}
}
