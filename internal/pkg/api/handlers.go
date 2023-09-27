package api

import (
	"net/http"
	"strconv"

	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
)

func notFound(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}

func showAllDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		d, err := r.GetDataServiceAll()
		if err != nil {
			notFound(c)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "DataService",
			"services": d,
		})
	}
}

func showDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, _ := strconv.ParseUint(c.Param("id")[1:], 10, 64)

		d, err := r.GetDataServiceById(uint(id))
		if err != nil {
			notFound(c)
			return
		}

		c.HTML(http.StatusOK, "service.tmpl", gin.H{
			"name": d.Name,
			"blob": d.Blob,
		})
	}
}

func filterDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		queryText, _ := c.GetQuery("text")

		filt, err := r.GetDataServiceFilteredByName(queryText)
		if err != nil {
			notFound(c)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl",
			gin.H{
				"title":    "DataService",
				"services": filt,
				"filtered": queryText,
			})
	}
}
