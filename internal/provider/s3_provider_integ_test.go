// +build integration_tests

package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewS3Provider_LifeCycle(t *testing.T) {
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	defer ts.Close()

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Endpoint:         aws.String(ts.URL),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	s3Session, err := session.NewSession(s3Config)
	assert.Nil(t, err)
	s3Client := s3.New(s3Session)

	cparams := &s3.CreateBucketInput{
		Bucket: aws.String("my-bucket"),
	}
	_, err = s3Client.CreateBucket(cparams)
	t.Logf("%v", err)
	assert.Nil(t, err)

	c := S3Conf{
		AccessId:     "bla",
		AccessSecret: "bla",
		Bucket:       "my-bucket",
		URL:          ts.URL,
	}
	prov, err := NewS3Provider(c)
	assert.Nil(t, err)
	defer prov.Close()

	err = prov.Set("k", strings.NewReader("content"))
	assert.Nil(t, err)
	checkKeyValue(t, prov, "k", "content")
}