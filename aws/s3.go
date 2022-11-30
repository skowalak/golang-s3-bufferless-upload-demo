package aws

import (
	"compress/gzip"
	"context"
	"io"
	"os"

	"file-upload-demo/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type S3 struct {
	config *config.AwsConfig
	logger *zap.SugaredLogger
	sess   *session.Session
}

func NewS3(cfg *config.Config, logger *zap.SugaredLogger) *S3 {
	// Set S3 envvars from config
	os.Setenv("AWS_ACCESS_KEY_ID", cfg.KeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.KeySecret)

	ret := &S3{
		config: &cfg.AwsConfig,
		logger: logger,
		sess: session.Must(session.NewSession(&aws.Config{
			Region: aws.String(cfg.Region),
		})),
	}
	return ret
}

func (s *S3) UploadStream(ctx context.Context, file io.Reader, fileName string, contentType string) (*s3manager.UploadOutput, error) {

	reader, writer := io.Pipe()

	go func() {
		gzipWriter := gzip.NewWriter(writer)
		_, err := io.Copy(gzipWriter, file)
		if err != nil {
			return
		}
		gzipWriter.Close()
		writer.Close()
	}()

	uploader := s3manager.NewUploader(s.sess)

	uploadInput := s3manager.UploadInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(fileName),
		Body:        reader,
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	}

	result, err := uploader.Upload(&uploadInput)
	if err != nil {
		s.logger.Errorw("error while uploading to s3", "error", err)
		return nil, err
	}
	s.logger.Infof("succesful upload: %s\n", aws.StringValue(&result.Location))
	return result, nil
}
