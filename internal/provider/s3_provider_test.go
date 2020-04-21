package provider

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3iface/
// https://github.com/johannesboyne/gofakes3

func TestGetS3_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockS3API(ctrl)
	mock.EXPECT().GetObject(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

	prov := S3Provider{client: mock}
	defer prov.Close()
	_, err := prov.Get("blabla")
	assert.NotNil(t, err)
}

func TestGetS3_Nominal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockS3API(ctrl)
	ret := &s3.GetObjectOutput{
		Body: ioutil.NopCloser(strings.NewReader("fileContent")),
	}
	mock.EXPECT().GetObject(gomock.Any()).Return(ret, nil).AnyTimes()

	prov := S3Provider{client: mock}
	defer prov.Close()
	checkKeyValue(t, prov, "blabla", "fileContent")
}

func TestGetS3_UnknownKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockS3API(ctrl)
	err := awserr.New(s3.ErrCodeNoSuchKey, "mocked error", nil)
	mock.EXPECT().GetObject(gomock.Any()).Return(nil, err).AnyTimes()

	prov := S3Provider{client: mock}
	defer prov.Close()
	_, err2 := prov.Get("blabla")
	assert.NotNil(t, err2)
	assert.True(t, errors.Is(err2, ErrKeyNotFound))
}

type uploaderMock struct {
	err error
}

func (um uploaderMock) Upload(input *s3manager.UploadInput, options ...func(uploader *s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return nil, um.err
}

func TestSetS3_NewKey(t *testing.T) {
	mock := uploaderMock{}
	prov := S3Provider{}
	err := prov.set(mock, "k", strings.NewReader("content"))
	assert.Nil(t, err)
}

func TestSetS3_Error(t *testing.T) {
	mock := uploaderMock{fmt.Errorf("error")}
	prov := S3Provider{}
	err := prov.set(mock, "k", strings.NewReader("content"))
	assert.NotNil(t, err)
}