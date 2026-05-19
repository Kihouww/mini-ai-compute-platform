package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Kihouww/mini-ai-compute-platform/internal/api"
)

func main() {
	r := gin.Default()

	api.RegisterRoutes(r)

	fmt.Println("server started at :8080")

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
