package object_storage

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"com.jadud.search.six/pkg/vcap"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Bucket struct {
	Name            string
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	Endpoint        string
	client          *s3.S3
	//session           *session.Session
}

type Buckets struct {
	Ephemeral Bucket
	Databases Bucket
}

func get_creds(vc *vcap.VcapServices, bucket string, field string) string {
	query := fmt.Sprintf("s3.#(instance_name==\"%s\").credentials.%s", bucket, field)
	s := vc.VCAP.Get(query).String()
	log.Println(bucket, field, s)
	return s
}

func InitBuckets(vc *vcap.VcapServices) *Buckets {
	Buckets := Buckets{}
	Buckets.Ephemeral = Bucket{}
	Buckets.Databases = Bucket{}
	Buckets.Ephemeral.Name = "ephemeral-storage"
	Buckets.Ephemeral.Region = get_creds(vc, Buckets.Ephemeral.Name, "region")
	Buckets.Ephemeral.SecretAccessKey = get_creds(vc, Buckets.Ephemeral.Name, "secret_access_key")
	Buckets.Ephemeral.AccessKeyId = get_creds(vc, Buckets.Ephemeral.Name, "access_key_id")
	Buckets.Ephemeral.Endpoint = get_creds(vc, Buckets.Ephemeral.Name, "uri")
	Buckets.Ephemeral.initClient()
	Buckets.Databases.Name = "database-storage"
	Buckets.Databases.Region = get_creds(vc, Buckets.Databases.Name, "region")
	Buckets.Databases.SecretAccessKey = get_creds(vc, Buckets.Databases.Name, "secret_access_key")
	Buckets.Databases.AccessKeyId = get_creds(vc, Buckets.Databases.Name, "access_key_id")
	Buckets.Databases.Endpoint = get_creds(vc, Buckets.Databases.Name, "uri")
	Buckets.Databases.initClient()
	// Create the buckets if they don't exist
	Buckets.Ephemeral.CreateBucket()
	Buckets.Databases.CreateBucket()

	return &Buckets
}

func s3_client(b *Bucket) *s3.S3 {
	// https://stackoverflow.com/questions/41544554/how-to-run-aws-sdk-with-credentials-from-variables
	creds := credentials.NewStaticCredentials(b.AccessKeyId, b.SecretAccessKey, "")

	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(b.Endpoint),
		Region:      aws.String(b.Region),
		Credentials: creds,
	})
	if err != nil {
		log.Fatal("CANNOT INIT AWS SESSION")
	}
	svc := s3.New(sess)
	return svc
}

func minio_client(b *Bucket) *s3.S3 {
	// Configure to use MinIO Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(b.AccessKeyId, b.SecretAccessKey, ""),
		Endpoint:         aws.String(b.Endpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal("CANNOT INIT AWS SESSION")
	}
	s3Client := s3.New(sess)
	return s3Client
}

// https://github.com/nitisht/cookbook/blob/master/docs/aws-sdk-for-go-with-minio.md

func (b *Bucket) initClient() {
	switch os.Getenv("ENV") {
	case "LOCAL":
		fallthrough
	case "DOCKER":
		b.client = minio_client(b)
	default:
		b.client = s3_client(b)
	}
}

// This will attempt to create the bucket.
// If it already exists, or is already owned by us, then
// nothing will happen. It silently moves on.
// We want to try and create buckets every time we start up.
func (b *Bucket) CreateBucket() {
	cparams := &s3.CreateBucketInput{
		Bucket: &b.Name,
	}
	_, err := b.client.CreateBucket(cparams)
	if err != nil {
		// Casting to the awserr.Error type will allow you to inspect the error
		// code returned by the service in code. The error code can be used
		// to switch on context specific functionality. In this case a context
		// specific error message is printed to the user based on the bucket
		// and key existing.
		//
		// For information on other S3 API error codes see:
		// http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				// pass
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				// pass
			default:
				log.Fatal(aerr)
			}
		}
	}
}

func (b *Bucket) ListObjects() {
	resp, err := b.client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(b.Name)})
	if err != nil {
		log.Println(err)
		log.Fatal("COULD NOT LIST OBJECTS IN BUCKET ", b.Endpoint, b.Name)
	}
	for _, item := range resp.Contents {
		log.Println("Name:         ", *item.Key)
		log.Println("Last modified:", *item.LastModified)
		log.Println("Size:         ", *item.Size)
		log.Println("Storage class:", *item.StorageClass)
		log.Println("")
	}
}

func (b *Bucket) PutObject(path []string, object []byte) {
	key := strings.Join(path, "/")
	log.Printf("storing object at %s", key)
	_, err := b.client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(object),
		Bucket: &b.Name,
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal(err)
	}
}
