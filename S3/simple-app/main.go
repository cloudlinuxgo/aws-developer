package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const S3Endpoint = "http://localhost:4566"

var bucket, fileName string

func main() {

	flag.StringVar(&bucket, "b", "mybucket", "Bucket name.")
	flag.StringVar(&fileName, "f", "myfile", "File name")
	flag.Parse()

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("user", "secret", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(S3Endpoint),
	}))

	svc := s3.New(sess)
	ctx := context.Background()

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: &bucket,
	})

	if err != nil {
		log.Fatalln("Could not create bucket:", err)
	}
	fmt.Printf("Bucket created %s\n", bucket)

	file, _ := os.Open("file-sample.json")

	_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		log.Fatalln("Could not upload file:", err)
	}

	fmt.Printf("File uploaded %s in the bucket %s\n", fileName, bucket)

}
