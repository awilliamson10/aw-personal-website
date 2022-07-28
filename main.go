package main

import (
	modules "github.com/awilliamson10/blog/modules"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", modules.Index)
	router.GET("/:cid", modules.GetPost)
	router.Run()
}
