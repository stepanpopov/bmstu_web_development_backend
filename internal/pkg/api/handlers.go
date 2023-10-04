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

func servicesView(d []repo.DataService, opts ...func(d []repo.DataService)) []repo.DataService {
	copyD := make([]repo.DataService, len(d))
	copy(d, copyD)

	for _, o := range opts {
		o(copyD)
	}
	return copyD
}

func viewWithBlobLenCtrl(n int) func(d []repo.DataService) {
	return func(d []repo.DataService) {
		n := n
		for _, v := range d {
			if len(v.Blob) > n {
				v.Blob = v.Blob[:n]
			}
		}
	}
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
			"services": servicesView(d, viewWithBlobLenCtrl(10)),
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
		queryText, _ := c.GetQuery("dataname")

		filt, err := r.GetDataServiceFilteredByName(queryText)
		if err != nil {
			notFound(c)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl",
			gin.H{
				"title":    "Шифрование кодом для коррекции ошибок",
				"services": servicesView(filt, viewWithBlobLenCtrl(10)),
				"filtered": queryText,
			})
	}
}
