package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func showAllPlainData(p []PlainData) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "PlainData",
			"services": p,
		})
	}
}

func showPlainData(p []PlainData) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, _ := strconv.Atoi(c.Param("id")[1:])

		for _, v := range p {
			if v.ID == id {
				c.HTML(http.StatusOK, "service.tmpl", gin.H{
					"name": v.Name,
					"blob": v.Blob,
				})
			}
		}
	}
}

func filterPlainDatas(p []PlainData) func(c *gin.Context) {
	return func(c *gin.Context) {
		queryText, _ := c.GetQuery("text")

		var filtered []PlainData
		for _, val := range p {
			if strings.Contains(val.Name, queryText) {
				filtered = append(filtered, val)
			}
		}

		c.HTML(http.StatusOK, "index.tmpl",
			gin.H{
				"title":    "PlainData",
				"services": filtered,
				"filtered": queryText,
			})
	}
}
