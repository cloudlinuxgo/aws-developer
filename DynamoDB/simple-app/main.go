package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	petname "github.com/dustinkirkland/golang-petname"
)

type Pet struct {
	UUID string
	Age  int
	Name string
}

const S3Endpoint = "http://localhost:4566"

var table string

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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				fmt.Println(err, table)
			}
		} else {
			log.Fatalf("Oh no! %s", err)
		}
	} else {
		fmt.Println("Created the table", table)
	}

	item := Pet{
		UUID: uuid.NewString(),
		Age:  rand.Intn(100),
		Name: randomPetName(),
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalln("Could not marshall the pet:", err)
	}

	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: &table,
	})

	if err != nil {
		log.Fatalln("Could not put the item:", err)
	}

	fmt.Printf("A new pet was added\n")
	readeable(item)
}

func randomPetName() string {
	return petname.Generate(2, "-")
}

func readeable(i interface{}) {
	out, _ := json.MarshalIndent(i, "", " ")
	fmt.Println(string(out))
}
