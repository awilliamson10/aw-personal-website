package modules

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	CID      string        `json:"cid"`
	Title    string        `json:"title"`
	Date     string        `json:"date"`
	Content  template.HTML `json:"content"`
	Previous int           `json:"previous"`
}

var PostList = []Post{}

func CreatePost(title string, body []byte) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	date := time.Now().Format("2006-01-02")
	content := template.HTML(markdown.ToHTML(body, parser, nil))
	hash := sha256.Sum256([]byte(date + title + string(content)))
	cid := hex.EncodeToString(hash[:])
	PostList = append(PostList, Post{
		CID:      cid,
		Title:    title,
		Date:     date,
		Content:  content,
		Previous: len(PostList),
	})
}

func GetPosts(CID string) []Post {
	return getAllPosts(CID)
}

func getPostByCID(cid string) (Post, error) {
	for _, post := range PostList {
		if post.CID == cid {
			return post, nil
		}
	}
	return Post{}, errors.New("Post not found")
}

func getIPFSPost(CID string) Post {
	// HTTP GET https://ipfs.io/ipfs/<CID>
	resp, err := http.Get("https://ipfs.io/ipfs/" + CID)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var P Post
	json.Unmarshal(body, &P)
	P.CID = CID
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	P.Content = template.HTML(markdown.ToHTML([]byte(P.Content), parser, nil))
	return P
}

func getAllPosts(CID string) []Post {
	// write a for loop while getIPFSPost.Previous != 0
	for post := getIPFSPost(CID); ; {
		fmt.Println(post)
		PostList = append(PostList, post)
		if post.Previous == 0 {
			break
		}
	}
	return PostList
}
