package api

import (
	"net/http"
	"strconv"
	"strings"

	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
)

func getDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		queryText, _ := c.GetQuery("dataname")

		filt, err := r.GetActiveDataServiceFilteredByName(queryText)
		if err != nil {
			notFound(c)
			return
		}

		c.JSON(http.StatusOK, filt)
	}
}

func getDataServiceByID(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id")[1:], 10, 64)

		d, err := r.GetDataServiceById(uint(id))
		if err != nil {
			respMessage(c, 404, "not found")
			return
		}

		c.JSON(http.StatusOK, d)
	}
}

func deleteDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request struct {
			ID uint `json:"id"`
		}
		if err := c.BindJSON(&request); err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		if r.DeleteDataService(request.ID) != nil {
			notFound(c)
			return
		}

		respMessage(c, http.StatusOK, "deleted")
	}
}

func createDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		data := repo.DataService{}
		if err := c.BindJSON(&data); err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		if strings.ReplaceAll(data.DataName, " ", "") == "" ||
			strings.ReplaceAll(data.Blob, " ", "") == "" ||
			len(data.DataName) > 30 ||
			data.DataID == 0 {

			respMessage(c, http.StatusBadRequest, "invalid input")
			return
		}

		if err := r.CreateDataService(data); err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		respMessage(c, http.StatusCreated, "created")
	}
}

func updateDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		data := repo.DataService{}
		if err := c.BindJSON(&data); err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		if strings.ReplaceAll(data.DataName, " ", "") == "" ||
			strings.ReplaceAll(data.Blob, " ", "") == "" ||
			len(data.DataName) > 30 ||
			data.DataID == 0 {

			respMessage(c, http.StatusBadRequest, "invalid input")
			return
		}

		if err := r.UpdateDataService(&data); err != nil {
			respMessage(c, http.StatusOK, "deleted")
		}

		c.JSON(http.StatusOK, data)
	}
}
