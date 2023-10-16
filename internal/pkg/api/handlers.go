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
			"title":    "Шифрование кодом для коррекции ошибок",
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
			"name":   d.DataName,
			"blob":   d.Blob,
			"id":     d.DataID,
			"encode": d.Encode,
			"active": d.Active,
		})
	}
}

func filterDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		queryText, _ := c.GetQuery("dataname")

		filt, err := r.GetActiveDataServiceFilteredByName(queryText)
		if err != nil {
			notFound(c)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl",
			gin.H{
				"title":    "Шифрование кодом для коррекции ошибок",
				"services": filt,
				"filtered": queryText,
			})
	}
}

func deleteDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id")[1:], 10, 64)
		if r.DeleteDataService(uint(id)) != nil {
			notFound(c)
			return
		}

		show := showDataService(r)
		show(c)
	}
}
