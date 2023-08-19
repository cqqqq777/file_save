package file_save

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func New(opts ...Option) app.HandlerFunc {
	config := NewConfig(opts)
	return func(ctx context.Context, c *app.RequestContext) {
		extractor, err := NewExtractor(config.Key)
		if err != nil {
			panic(err)
		}
		if config.s3.Credentials == nil && config.Dst == "" {
			panic("s3 save func and local save func can not be empty in the same time")
		}

		saver := make([]SaveFunc, 0, 2)

		if config.Dst != "" {
			saver = append(saver, NewLocalSaveFunc(config.Dst))
		}

		if config.s3.Credentials != nil {
			if config.s3.Region == "" {
				config.s3.Region = "us-east-1"
			}
			cli, err := NewS3Client(config.s3)
			if err != nil {
				hlog.Error("new s3 save func failed err:", err)
			}
			saver = append(saver, NewS3SaveFunc(ctx, cli, config.s3.BucketName, config.s3.ObjectName))
		}

		if len(saver) == 0 {
			panic("s3 save func and local save func can not be empty in the same time")
		}

		file, err := extractor(c)
		if err != nil {
			_ = c.Error(err)
			config.ErrorHandler(ctx, c)
			return
		}
		for _, v := range saver {
			if err = v(file); err != nil {
				_ = c.Error(err)
				config.ErrorHandler(ctx, c)
				return
			}
		}
		config.SuccessHandler(ctx, c)
	}
}
