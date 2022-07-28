package modules

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var firstLoad = false
var posts []Post

var Index = func(c *gin.Context) {

	if !firstLoad {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error loading .env file")
		}
		var LATESTCID = os.Getenv("LATEST_CID")
		fmt.Println("LATEST CID: " + LATESTCID)
		posts = GetPosts(LATESTCID)
		firstLoad = true
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Andrew Williamson's website",
		"payload": posts,
	})
}

var GetPost = func(c *gin.Context) {
	cid := c.Param("cid")
	post, err := getPostByCID(cid)
	if err != nil {
		fmt.Println("ERROR", err)
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.HTML(http.StatusOK, "post.html", gin.H{
		"title":   post.Title,
		"payload": post,
	})
}
