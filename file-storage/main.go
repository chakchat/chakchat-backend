package main

import (
	"github.com/chakchat/chakchat/backend/file-storage/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.POST("/upload", handlers.Upload(&handlers.UploadConfig{
		FileSizeLimit: 20 << 10,
	}, nil))
	r.Run(":5005")
}
