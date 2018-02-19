package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// An example of using Hasura's auth UI kit & API gateway
// to avoid writing auth and session handling code.
//
// Hasura's auth microservice provides login/signup functionality
// and can redirect back to your own app.
//
// Hasura's API gateway will automatically resolve session tokens into
// HTTP headers containing user-id and roles.
func userInfo(c *gin.Context) {
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
		c.HTML(http.StatusOK, "auth_anonymous.html", gin.H{
			"baseDomain": c.GetHeader("X-Hasura-Base-Domain"),
		})
	} else {
		c.HTML(http.StatusOK, "auth_user.html", gin.H{
			"baseDomain": c.GetHeader("X-Hasura-Base-Domain"),
			"userId":     c.GetHeader("X-Hasura-User-Id"),
			"roles":      c.GetHeader("X-Hasura-Allowed-Roles"),
		})
	}
}
