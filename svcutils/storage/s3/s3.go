package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/viper"
)

const (
	urlFormat = "https://%s.s3.%s.amazonaws.com"

	envNameAccessKey = "AWS_ACCESS_KEY_ID"
	envNameSecretKey = "AWS_SECRET_ACCESS_KEY"
	envNameRegion    = "AWS_REGION"
)

type BucketClient struct {
	region string
	bucket string
	client *s3.Client
}

type ClientInterface interface {
	Upload(context.Context, *File) (string, error)
	// Get fetches a file from s3 based on the file name and returns a reader for that file
	Get(context.Context, string) (io.ReadCloser, error)
	GetURLFromFileName(name string) string
}

type File struct {
	Name          string
	FullPath      string
	Content       io.ReadCloser
	IsPublic      bool
	ContentLength int64
}

func NewS3BucketClient(cfg *viper.Viper) (*BucketClient, error) {
	var (
		accessKey  = cfg.GetString("aws.accessKey")
		secretKey  = cfg.GetString("aws.secretKey")
		bucketName = cfg.GetString("aws.bucketName")
		region     = cfg.GetString("aws.region")
	)
	switch "" {
	case accessKey:
		return nil, errors.New("aws access key not set")
	case secretKey:
		return nil, errors.New("aws secret key not set")
	case bucketName:
		return nil, errors.New("s3 bucket name not set")
	case region:
		return nil, errors.New("aws region not set")
	}

	os.Setenv(envNameAccessKey, accessKey)
	os.Setenv(envNameSecretKey, secretKey)
	os.Setenv(envNameRegion, region)
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %v", err)
	}
	client := s3.NewFromConfig(awsConfig)

	return &BucketClient{
		region: region,
		bucket: bucketName,
		client: client,
	}, nil
}

func (c *BucketClient) Upload(ctx context.Context, file *File) (string, error) {
	// generate the new name from the current timestamp
	t := strconv.FormatInt(time.Now().Unix(), 10)
	fName := fmt.Sprintf("%s_%s", t, strings.ReplaceAll(filepath.Base(file.Name), " ", "+"))
	if file.FullPath != "" {
		fName = file.FullPath
	}

	defer file.Content.Close()

	input := s3.PutObjectInput{
		Bucket: &c.bucket,
		Body:   file.Content,
		Key:    &fName,
	}

	if file.ContentLength > 0 {
		input.ContentLength = &file.ContentLength
	}

	if file.IsPublic {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	_, err := c.client.PutObject(ctx, &input)
	if err != nil {
		return "", err
	}

	urlString := fmt.Sprintf(urlFormat, c.bucket, c.region) + "/" + fName
	return urlString, nil
}

func (c *BucketClient) GetFileNameFromURL(url string) string {
	urlStr := fmt.Sprintf(urlFormat, c.bucket, c.region)

	name := strings.TrimLeft(strings.ReplaceAll(url, urlStr, ""), "/")
	return name
}

func (c *BucketClient) GetURLFromFileName(name string) string {
	urlStr := fmt.Sprintf(urlFormat, c.bucket, c.region)

	return strings.Join([]string{urlStr, name}, "/")
}

// Get fetches a file from s3 based on the file name and returns a reader for that file
func (c *BucketClient) Get(ctx context.Context, fileName string) (io.ReadCloser, error) {
	input := s3.GetObjectInput{Bucket: &c.bucket, Key: &fileName}
	out, err := c.client.GetObject(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("error fetching object: %v", err)
	}
	return out.Body, nil
}
