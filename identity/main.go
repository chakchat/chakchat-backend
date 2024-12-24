package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()

	r.Use(CheckIdempotencyKey(nil))

	r.Run()
}
