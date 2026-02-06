package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var BucketName string

// InitStorage connects to Object Storage using env variables
func InitStorage() {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	BucketName = os.Getenv("STORAGE_BUCKET")
	region := os.Getenv("STORAGE_REGION")

	if endpoint == "" || accessKey == "" || secretKey == "" || BucketName == "" {
		log.Fatal("Error: Storage configuration missing in .env")
	}

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("Error: Failed to load S3 config: %v", err)
	}

	S3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(endpoint)
	})

	fmt.Println("Connected to Object Storage")
}

func UploadFile(file multipart.File, filename string, size int64) (string, error) {
	_, err := S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(BucketName),
		Key:           aws.String(filename),
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("application/pdf"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Construct the URL using the endpoint from env
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	return fmt.Sprintf("%s/%s/%s", endpoint, BucketName, filename), nil
}
