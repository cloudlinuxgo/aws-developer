package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	petname "github.com/dustinkirkland/golang-petname"
)

const S3Endpoint = "http://localhost:4566"

var table string

func randomPetName() string {
	return petname.Generate(2, "-")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	flag.StringVar(&table, "t", "pets", "Bucket name.")
	flag.Parse()

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("user", "secret", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(S3Endpoint),
	}))

	dynamo := dynamodb.New(sess)

	tableCreate := &dynamodb.CreateTableInput{
		TableName: &table,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("UUID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Age"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("UUID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Age"),
				KeyType:       aws.String("RANGE"),
			},
			// {
			// 	AttributeName: aws.String("Name"),
			// 	KeyType:       aws.String("RANGE"),
			// },
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(2),
			WriteCapacityUnits: aws.Int64(2),
		},
	}

	_, err := dynamo.CreateTable(tableCreate)
	if err != nil {
		log.Fatalf("Oh no! %s", err)
	}

	fmt.Println("Created the table", table)

}
