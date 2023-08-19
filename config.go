package file_save

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type S3 struct {
	Credentials *Credentials
	Endpoint    string

	// BucketName and ObjectName define the object basic info
	BucketName string
	ObjectName string

	// S3 compatible object storage region
	// Default value is us-east-1.
	// Other valid values are listed below.
	// us-east-1 us-east-2 us-west-1 us-west-2 us-gov-west-1 us-gov-east-1
	// ca-central-1
	// eu-west-1 eu-west-2 eu-west-3 eu-central-1 eu-north-1
	// ap-east-1 ap-south-1 ap-southeast-1 ap-southeast-2 ap-northeast-1 ap-northeast-2 ap-northeast-3
	// me-south-1
	// sa-east-1
	// cn-north-1 cn-northwest-1
	Region string

	Secure bool
}

// Config define the config for middleware
type Config struct {
	s3 *S3

	// Key define the file's key
	Key string

	// Dst define the destination when save file locally
	Dst string

	// ErrorHandler defines a function which is executed when an error occurs.
	// Default: OutPut log and response 500.
	ErrorHandler app.HandlerFunc

	// SuccessHandler define a function which is executed when file save successfully
	// Default: Response 200
	SuccessHandler app.HandlerFunc
}

type Credentials struct {
	// AWS Access key ID
	AccessKeyID string

	// AWS Secret Access Key
	SecretAccessKey string

	// AWS Session Token
	SessionToken string
}

type Option func(cfg *Config)

var DefaultConfig = Config{
	ErrorHandler: func(ctx context.Context, c *app.RequestContext) {
		hlog.Error(c.Errors.Last())
		c.String(consts.StatusInternalServerError, "save file failed")
		c.Abort()
	},
	SuccessHandler: func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "success")
		c.Abort()
	},
	s3: &S3{},
}

func NewConfig(opts []Option) *Config {
	cfg := DefaultConfig
	cfg.Apply(opts)
	return &cfg
}

func (cfg *Config) Apply(options []Option) {
	for _, v := range options {
		v(cfg)
	}
}

func WithCredentials(id, secret, token string) Option {
	return func(cfg *Config) {
		cfg.s3.Credentials = &Credentials{
			AccessKeyID:     id,
			SecretAccessKey: secret,
			SessionToken:    token,
		}
	}
}

func WithEndpoint(endpoint string) Option {
	return func(cfg *Config) {
		cfg.s3.Endpoint = endpoint
	}
}

func WithRegion(region string) Option {
	return func(cfg *Config) {
		cfg.s3.Region = region
	}
}

func WithBucket(name string) Option {
	return func(cfg *Config) {
		cfg.s3.BucketName = name
	}
}

func WithObject(name string) Option {
	return func(cfg *Config) {
		cfg.s3.ObjectName = name
	}
}

func WithDestination(dst string) Option {
	return func(cfg *Config) {
		cfg.Dst = dst
	}
}

func WithErrorHandler(f app.HandlerFunc) Option {
	return func(cfg *Config) {
		cfg.ErrorHandler = f
	}
}

func WithKey(key string) Option {
	return func(cfg *Config) {
		cfg.Key = key
	}
}

func WithSuccessHandler(f app.HandlerFunc) Option {
	return func(cfg *Config) {
		cfg.SuccessHandler = f
	}
}
