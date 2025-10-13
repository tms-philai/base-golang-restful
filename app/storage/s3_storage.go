package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucket     string
	region     string
	baseURL    string
}

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Endpoint        string
	UsePathStyle    bool
}

func NewS3Storage(config S3Config) (*S3Storage, error) {
	awsConfig := &aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	}

	if config.Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.Endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(config.UsePathStyle)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)

	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", config.Bucket, config.Region)
	if config.Endpoint != "" {
		baseURL = config.Endpoint
	}

	return &S3Storage{
		client:     client,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		bucket:     config.Bucket,
		region:     config.Region,
		baseURL:    baseURL,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, input UploadInput) (*UploadResult, error) {
	uploadInput := &s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(input.Key),
		Body:        input.Content,
		ContentType: aws.String(input.ContentType),
	}

	if input.ACL != "" {
		uploadInput.ACL = aws.String(input.ACL)
	}

	if len(input.Metadata) > 0 {
		metadata := make(map[string]*string)
		for k, v := range input.Metadata {
			metadata[k] = aws.String(v)
		}
		uploadInput.Metadata = metadata
	}

	result, err := s.uploader.UploadWithContext(ctx, uploadInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	return &UploadResult{
		Key:      input.Key,
		URL:      result.Location,
		ETag:     aws.StringValue(result.ETag),
		Size:     input.Size,
		Metadata: input.Metadata,
	}, nil
}

func (s *S3Storage) Download(ctx context.Context, key string) ([]byte, error) {
	buffer := aws.NewWriteAtBuffer([]byte{})

	_, err := s.downloader.DownloadWithContext(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	return buffer.Bytes(), nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

func (s *S3Storage) GetURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}

func (s *S3Storage) GetSignedURL(ctx context.Context, key string, expiration int64) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	duration := time.Duration(expiration) * time.Second
	urlStr, err := req.Presign(duration)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return urlStr, nil
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if isNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}

func (s *S3Storage) List(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	var keys []string

	err := s.client.ListObjectsV2PagesWithContext(ctx, input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			keys = append(keys, aws.StringValue(obj.Key))
		}
		return !lastPage
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	return keys, nil
}

func (s *S3Storage) CopyObject(ctx context.Context, sourceKey, destKey string) error {
	copySource := fmt.Sprintf("%s/%s", s.bucket, sourceKey)

	_, err := s.client.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destKey),
	})

	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}

func (s *S3Storage) GetObjectMetadata(ctx context.Context, key string) (map[string]string, error) {
	result, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	metadata := make(map[string]string)
	for k, v := range result.Metadata {
		metadata[k] = aws.StringValue(v)
	}

	return metadata, nil
}

func isNotFoundError(err error) bool {
	return err != nil && (err.Error() == "NotFound" || err.Error() == "NoSuchKey")
}
