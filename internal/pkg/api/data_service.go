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

		draftID, err := r.GetEncryptDecryptDraftID(creatorID)
		if err != nil {
			respMessage(c, 500, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"data_service": filt,
			"draft_id":     draftID,
		})
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
		id, _ := strconv.ParseUint(c.Param("id")[1:], 10, 64)

		if r.DeleteDataService(uint(id)) != nil {
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

/*
TODO
func updateImage(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
*/

func addToDraft(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id")[1:], 10, 64)

		draftID, err := r.AddDataServiceToDraft(uint(id), creatorID)

		if err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"draftID": draftID,
		})
	}
}
