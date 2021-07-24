package network

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
	"banson.moe/t2f-bot/config"
)

const (
	S3Region = config.S3Region
	S3Acl    = "public-read"
)

func UploadFileToS3(bucket, filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	// Create a session which contains the default configurations for the SDK.
	// Use the session to create the service clients to make API calls to AWS.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3Region),
		Credentials: credentials.NewStaticCredentials(
			config.AwsID,     // id
			config.AwsSecret, // secret
			"") ,// token can be left blank for now
	})

	// Create S3 Uploader manager to concurrently upload the file
	svc := s3manager.NewUploader(sess)

	//fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(filename)),
		ACL:    aws.String(S3Acl),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	//fmt.Printf("Successfully uploaded %s to %s\n", filename, result.Location)
	return result.Location
}
