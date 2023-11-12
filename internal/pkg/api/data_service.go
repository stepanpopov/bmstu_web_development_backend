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
			respMessageAbort(c, 500, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data_service": filt,
			"draft_id":     draftID,
		})
	}
}

func getDataServiceByID(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		print(id)
		d, err := r.GetDataServiceById(uint(id))
		if err != nil {
			respMessage(c, 404, "not found")
			return
		}

		c.JSON(http.StatusOK, d)
	}
}

func deleteDataService(r repo.Repository, a repo.Avatar) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id")[:], 10, 64)

		avatarUUID, err := r.DeleteDataService(uint(id))
		if err != nil {
			notFound(c)
			return
		}

		if err := a.Delete(c, avatarUUID); err != nil {
			respMessageAbort(c, http.StatusInternalServerError, "не получается удалить картинку")
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
		data.DataID = 0
		data.Active = true

		if strings.ReplaceAll(data.DataName, " ", "") == "" ||
			strings.ReplaceAll(data.Blob, " ", "") == "" ||
			len(data.DataName) > 30 {

			respMessage(c, http.StatusBadRequest, "invalid input")
			return
		}

		dataID, err := r.CreateDataService(data)
		if err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data_id": dataID,
		})
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

func putImage(r repo.Repository, a repo.Avatar) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		form, err := c.MultipartForm()
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, "не получается достать изображение")
			return
		}
		fileHeader := form.File["avatar"][0]
		// TODO: check fileHeader.Header
		avatar, err := form.File["avatar"][0].Open()
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, "не получается достать изображение:")
			return
		}

		avatarUUID, err := a.Put(c, avatar, fileHeader.Size)
		if err != nil {
			respMessageAbort(c, http.StatusInternalServerError, "не получается сохранить изображение: "+err.Error())
			return
		}

		if err := r.UpdateImageUUID(avatarUUID, uint(id)); err != nil {
			respMessageAbort(c, http.StatusInternalServerError, "не получается сохранить изображение")
			return
		}

		respMessage(c, http.StatusOK, "uploaded")
	}
}

func addToDraft(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		draftID, err := r.AddDataServiceToDraft(uint(id), creatorID)

		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"draftID": draftID,
		})
	}
}

func deleteFromDraft(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		if err := r.DeleteDataServiceFromDraft(uint(id), creatorID); err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		respMessage(c, http.StatusOK, "deleted")
	}
}
