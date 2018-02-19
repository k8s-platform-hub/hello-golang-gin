package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
)

// For local development,
// First: connect to Hasura Data APIs directly on port 9000
// $ hasura ms port-forward data -n hasura --local-port=9000
// Second: Uncomment the line below
// const dataUrl = "http://localhost:9000/v1/query"

// When deployed to your cluster, use this:
const dataUrl = "http://data.hasura/v1/query"

func getArticles(c *gin.Context) {
	if !(strings.Contains(c.Request.Host, "hasura-app.io") || !strings.Contains(dataUrl, "data.hasura")) {
		c.Data(http.StatusOK, "text/html", []byte(`Edit the <code>dataUrl</code> variable in <code>microservices/app/src/data.go</code> to test locally.`))
		return
	}

	resp, err := grequests.Post(dataUrl,
		&grequests.RequestOptions{
			JSON: map[string]interface{}{
				"type": "select",
				"args": map[string]interface{}{
					"table":   "article",
					"columns": []string{"title", "id", "author_id", "rating"},
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
			"error": fmt.Errorf("code: %d, data: %s", resp.StatusCode, string(resp.Bytes())),
		})
		return
	}

	var data interface{}
	err = resp.JSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	jsonString, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "data.html", gin.H{
		"data": string(jsonString),
	})
}
