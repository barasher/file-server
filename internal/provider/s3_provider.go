package provider

// https://github.com/minio/cookbook/blob/master/docs/aws-sdk-for-go-with-minio.md
// https://medium.com/@alexsante/serving-up-videos-from-s3-to-the-browser-using-go-974dfc11b738

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const defaultRegion = "us-east-1"

type S3Provider struct {
	client *s3.S3
	session *session.Session
	bucket *string
}

func (p S3Provider) Provide(key string) (io.ReadCloser, error) {
	obj, err := p.client.GetObject(&s3.GetObjectInput{
		Bucket: p.bucket,
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok &&  aerr.Code() == s3.ErrCodeNoSuchKey {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return obj.Body, nil
}

func (p S3Provider) Close() {
	// nothing
}

func NewS3Provider(conf S3Conf) (S3Provider, error) {
	prov := S3Provider{}
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.AccessId, conf.AccessSecret, ""),
		Endpoint:         aws.String(conf.URL),
		Region:           aws.String(defaultRegion),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return prov, err
	}
	prov.session = newSession
	prov.client = s3.New(newSession)
	prov.bucket = aws.String(conf.Bucket)
	return prov, nil
}

type S3Conf struct {
	AccessId     string
	AccessSecret string
	Bucket       string
	URL string
}