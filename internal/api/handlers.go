package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func showAllServices(services []Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "Services",
			"services": services,
			"contains": strings.Contains,
			"filtered": "",
		})
	}
}

func showService(services []Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": c.Param("id"),
		})
	}
}

type text struct {
	Text string
}

func filterServices(services []Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req text
		c.BindJSON(req)

		c.HTML(http.StatusOK, "index.tmpl",
			gin.H{
				"title":    "Services",
				"services": services,
				"contains": strings.Contains,
				"filtered": req.Text,
			})
	}
}
