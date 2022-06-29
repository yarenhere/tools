package osop

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

func GetDownloadUrl(svc *s3.S3, bucket, key string) (presignUrl string, err error) {
	input := new(s3.GetObjectInput)
	input.SetBucket(bucket)
	input.SetKey(key)

	req, _ := svc.GetObjectRequest(input)
	return req.Presign(time.Hour)
}
