package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
)

// This route renders a file-upload page for users.
//
// Not logged-in users (anonymous):
//   > request that they login
//
// logged-in users (anonymous):
//   > list files they own
//   > show file-upload box
//
// The file-upload and download API uses Hasura's filestore APIs
// so you will notice that this code has no file-handling code!
func userFiles(c *gin.Context) {
	// Important for local development, ignore otherwise.
	// Local development will need you to add custom headers to simulate the API gateway.
	if len(c.GetHeader("X-Hasura-Allowed-Roles")) == 0 {
		c.Data(http.StatusOK, "text/html", []byte(`This route can only be accessed via the Hasura API gateway.
Add headers via your browser if you're testing locally.<br/>
<a href="https://docs.hasura.io/0.15/manual/gateway/session-middleware.html"
target="_blank">Read the docs.</a>`))
		return
	}
	// End local development directives

	// If user is not logged in, render HTML page with login link
	if strings.Contains(c.GetHeader("X-Hasura-Allowed-Roles"), "anonymous") {
		c.HTML(http.StatusOK, "filestore_anonymous.html", gin.H{
			"baseDomain": c.GetHeader("X-Hasura-Base-Domain"),
		})
	} else {
		resp, err := grequests.Post(dataUrl,
			&grequests.RequestOptions{
				JSON: map[string]interface{}{
					"type": "select",
					"args": map[string]interface{}{
						"table": map[string]string{
							"name":   "hf_file",
							"schema": "hf_catalog",
						},
						"columns": []string{"file_id", "content_type", "file_size"},
						"where": map[string]string{
							"user_id": c.GetHeader("X-Hasura-User-Id"),
						},
					},
				},
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if !resp.Ok {
			c.JSON(resp.StatusCode, gin.H{
				"error": fmt.Sprintf("code: %d, data: %s", resp.StatusCode, string(resp.Bytes())),
			})
			return
		}

		var files []fileResponse
		err = resp.JSON(&files)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "filestore_user.html", gin.H{
			"baseDomain": c.GetHeader("X-Hasura-Base-Domain"),
			"files":      files,
		})
	}
}

type fileResponse struct {
	FileId      string `json:"file_id"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
}
