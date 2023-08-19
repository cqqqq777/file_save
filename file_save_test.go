package file_save

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"testing"
)

func TestNew(t *testing.T) {
	h := server.New(server.WithHostPorts(":8080"))
	handler := New(WithCredentials("minio", "123456", ""),
		WithEndpoint("localhost:9000"),
		WithDestination("./hello.jpg"),
		WithBucket("mytxt"),
		WithObject("hello.jpg"),
		WithKey("a"))
	h.POST("/file", handler)
	h.Spin()
}
