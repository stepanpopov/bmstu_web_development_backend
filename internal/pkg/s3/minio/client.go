package s3

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3AvatarSaver struct {
	avatarBucket string
	avatarFolder string
	cl           *minio.Client
}

func NewS3AvatarSaver(avatarBucket, avatarFolder string, client *minio.Client) *S3AvatarSaver {
	return &S3AvatarSaver{
		avatarBucket: avatarBucket,
		avatarFolder: avatarFolder,
		cl:           client,
	}
}

func MakeS3MinioClient(endpoint, accessKey, secret string) (*minio.Client, error) {
	if endpoint == "" || accessKey == "" || secret == "" {
		return nil, fmt.Errorf("invalid config")
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secret, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func (s *S3AvatarSaver) Save(ctx context.Context, avatar io.Reader, fileName string, size int64) error {

	objectPath := filepath.Join(s.avatarFolder, fileName)
	_, err := s.cl.PutObject(context.Background(), s.avatarBucket, objectPath,
		avatar, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})

	if err != nil {
		return err
	}

	return nil
}

/* s3Client, err := s3.MakeS3MinioClient(os.Getenv(config.S3HostParam), os.Getenv(config.S3AccessKeyParam), os.Getenv(config.S3SecretKeyParam))
if err != nil {
	logger.Errorf("Error while connecting to S3: %v", err)
	return
}
userS3 := userS3.NewS3AvatarSaver(os.Getenv(config.S3BucketParam), os.Getenv(config.S3AvatarFolderParam), s3Client)
*/
