package object_storage

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"com.jadud.search.six/pkg/vcap"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Bucket struct {
	Name            string
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	Endpoint        string
	client          *s3.S3
	session         *session.Session
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

func s3_client(b *Bucket) (*s3.S3, *session.Session) {
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
	return svc, sess
}

func minio_client(b *Bucket) (*s3.S3, *session.Session) {
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

	return s3Client, sess
}

// https://github.com/nitisht/cookbook/blob/master/docs/aws-sdk-for-go-with-minio.md

func (b *Bucket) initClient() {
	switch os.Getenv("ENV") {
	case "LOCAL":
		fallthrough
	case "DOCKER":
		b.client, b.session = minio_client(b)
	default:
		b.client, b.session = s3_client(b)
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

func (b *Bucket) ListObjects(filter string) []*s3.Object {
	resp, err := b.client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(b.Name)})
	if err != nil {
		log.Println(err)
		log.Fatal("COULD NOT LIST OBJECTS IN BUCKET ", b.Endpoint, b.Name)
	}

	keys := make([]*s3.Object, 0)

	for _, item := range resp.Contents {
		log.Println("Name:         ", *item.Key)
		log.Println("Last modified:", *item.LastModified)
		log.Println("Size:         ", *item.Size)
		log.Println("Storage class:", *item.StorageClass)
		log.Println("")

		log.Printf("CHECKING OBJECT %s against filter %s\n", *item.Key, filter)
		if found, _ := regexp.MatchString(filter, *item.Key); found {
			keys = append(keys, item)
		}
	}
	return keys
}

func get_mime_type(path string) string {
	m := map[string]string{
		"json":    "application/json",
		"txt":     "text/plain",
		"md":      "text/plain",
		"pdf":     "application/pdf",
		"sqlite":  "application/x-sqlite3",
		"sqlite3": "application/x-sqlite3",
		// https://www.iana.org/assignments/media-types/application/zstd
		"zstd": "application/zstd",
	}
	for k, v := range m {
		if bytes.HasSuffix([]byte(path), []byte(k)) {
			return v
		}
	}
	return m["json"]
}

// https://github.com/nitisht/cookbook/blob/master/docs/aws-sdk-for-go-with-minio.md
func (b *Bucket) PutObject(path []string, object []byte) {
	key := strings.Join(path, "/")

	log.Printf("storing object at %s", key)
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#PutObjectInput
	_, err := b.client.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(object),
		Bucket:      &b.Name,
		Key:         aws.String(key),
		ContentType: aws.String(get_mime_type(path[len(path)-1])),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (b *Bucket) GetObject(key string) []byte {

	log.Printf("getting object at %s", key)
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#GetObjectInput
	goo, err := b.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal(err)
	}
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#GetObjectOutput
	buf := new(bytes.Buffer)
	buf.ReadFrom(goo.Body)
	defer goo.Body.Close()

	return buf.Bytes()
}

func (b *Bucket) DownloadFile(path []string, filename string) {
	key := strings.Join(path, "/")

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(b.session)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(b.Name),
			Key:    aws.String(key),
		})
	if err != nil {
		fmt.Println(err)
	}
	if numBytes == 0 {
		log.Printf("DownloadFile: %s -> %s was 0 bytes\n", key, filename)
	}
}

func (b *Bucket) UploadFile(path []string, filename string) {
	key := strings.Join(path, "/")

	log.Printf("uploading file to %s", key)

	// Create an uploader with the session and custom options
	uploader := s3manager.NewUploader(b.session, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
		u.Concurrency = 5            // default
	})

	//open the file
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to open file %q, %v", filename, err)
		log.Fatal(err)
	}
	defer f.Close()

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(b.Name),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(get_mime_type(path[len(path)-1])),
	})

	//in case it fails to upload
	if err != nil {
		log.Printf("failed to upload file, %v", err)
		log.Fatal(err)
	}
}
