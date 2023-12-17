package s3

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type S3MinioAvatar struct {
	avatarBucket string
	cl           *minio.Client
}

func NewS3MinioAvatarSaver(avatarBucket string, client *minio.Client) *S3MinioAvatar {
	return &S3MinioAvatar{
		avatarBucket: avatarBucket,
		cl:           client,
	}
}

func (s *S3MinioAvatar) Put(ctx context.Context, avatar io.Reader, size int64) (uuid.UUID, error) {
	objUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	_, err = s.cl.PutObject(ctx, s.avatarBucket, objUUID.String(),
		avatar, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return uuid.Nil, err
	}

	return objUUID, nil
}

func (s *S3MinioAvatar) Delete(ctx context.Context, uuid uuid.UUID) error {
	return s.cl.RemoveObject(ctx, s.avatarBucket, uuid.String(), minio.RemoveObjectOptions{})
}
