package osop

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"
	"github.com/yarenhere/tools/mfs"
	"io/fs"
	"net/http"
	"testing"
)

func TestUploadMultiPartFromReader(t *testing.T) {
	reader := mfs.NewMockRandomFile("/abc/d.txt", fs.ModePerm, 18*1024*1024)
	fileHash, err := mfs.GetFileMd5Hash(reader)
	assert.Equal(t, nil, err)
	reader.Close()

	s3Cfg := aws.NewConfig()
	s3Cfg.WithDisableSSL(true)
	s3Cfg.WithEndpoint("http://127.0.0.1:9000")
	s3Cfg.WithCredentials(credentials.NewStaticCredentials("32jDAV1Cy8vVhGRI", "WWtiiYyx2QLDYVByn7N9GXUhV3VMqAJT", ""))
	s3Cfg.WithRegion("us-east-1")
	s3Cfg.WithS3ForcePathStyle(true)

	session, err := session.NewSession(s3Cfg)
	assert.Equal(t, nil, err)
	svc := s3.New(session)

	bucket := "test-bucket"
	key := "abc/d.txt"
	logrus.SetLevel(logrus.TraceLevel)
	err = UploadMultiPartFromReader(svc, reader, "test-bucket", "abc/d.txt")
	assert.Equal(t, nil, err)
	presignUrl, err := GetDownloadUrl(svc, bucket, key)
	assert.Equal(t, nil, err)

	resp, err := http.Get(presignUrl)
	assert.Equal(t, nil, err)
	uploadFileHash, err := mfs.GetFileMd5Hash(resp.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, fileHash, uploadFileHash)
}
