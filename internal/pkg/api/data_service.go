package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Filter data service
// @Tags         DataService
// @Description  Get data services filtered by name
// @Accept       json
// @Produce      json
// @Param		 dataname query string true 				"Name"
// @Success      200    {object}  map[string]any 		   "Data Service Filtered"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/ [get]
func filterDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		queryText, _ := c.GetQuery("dataname")

		filt, err := r.GetActiveDataServiceFilteredByName(queryText)
		if err != nil {
			notFound(c)
			return
		}

		var draftID *uint
		if userUUID := getUserUUIDFromCtx(c); userUUID != uuid.Nil {
			draftID, err = r.GetEncryptDecryptDraftID(getUserUUIDFromCtx(c))
			if err != nil {
				respMessageAbort(c, 500, err.Error())
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"data_service": toViewSlice(filt),
			"draft_id":     draftID,
		})
	}
}

// @Summary      Get data service by id
// @Tags         DataService
// @Description  Get data service by id
// @Accept       json
// @Produce      json
// @Param		 dataServiceID path int true 				"Data Service ID"
// @Success      200    {object}  DataServiceView 		   "Got Data Service"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/{dataServiceID} [get]
func getDataServiceByID(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		print(id)
		d, err := r.GetDataServiceById(uint(id))
		if err != nil {
			respMessage(c, 404, "not found")
			return
		}

		c.JSON(http.StatusOK, toView(*d))
	}
}

// @Summary      Delete data service
// @Tags         DataService
// @Description  Delete data service
// @Accept       json
// @Produce      json
// @Param		 dataServiceID path int true 				"Data Service ID"
// @Success      200    {object}  string 		   "Delete Data Service"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/{dataServiceID} [delete]
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

// @Summary      Create data service
// @Tags         DataService
// @Description  Create data services
// @Accept       json
// @Produce      json
// @Param		 data	body		repo.DataService    true "Data Service"
// @Success      200    {object}  map[string]int 		   "Create Data Service"
// @Failure      400    {object}  error  			   "Bad request"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/ [post]
func createDataService(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		data := repo.DataService{}
		if err := c.BindJSON(&data); err != nil {
			print("after bind")
			respMessageAbort(c, http.StatusBadRequest, err.Error())
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
			print("after orm")
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data_id": dataID,
		})
	}
}

// @Summary      Update data service
// @Tags         DataService
// @Description  Update data services
// @Accept       json
// @Produce      json
// @Param		 data	body		repo.DataService    true "Data Service"
// @Success      200    {object}  map[string]int 		   "Create Data Service"
// @Failure      400    {object}  error  			   "Bad request"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/ [put]
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

		c.JSON(http.StatusOK, toView(data))
	}
}

// @Summary      Put image for data service
// @Tags         DataService
// @Description  Put image for data service
// @Accept       multipart/form-data
// @Produce      json
// @Param		 avatar formData file true 				   "DataService avatar file"
// @Success      200    {object}  string 		   "Create Data Service"
// @Failure      400    {object}  error  			   "Bad request"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/{dataServiceID}/image [post]
func putImage(r repo.Repository, a repo.Avatar) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("In put image")
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		form, err := c.MultipartForm()
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, "не получается достать изображение: "+err.Error())
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

// @Summary      Add data service to draft
// @Tags         DataService
// @Description  Add data service to draft
// @Produce      json
// @Param		 id path int true 				"Data Service ID"
// @Success      200    {object}  string 		   "Data Service added to draft"
// @Failure      400    {object}  error  			   "Bad request"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/draft/{id} [post]
func addToDraft(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		draftID, err := r.AddDataServiceToDraft(uint(id), getUserUUIDFromCtx(c))

		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"draftID": draftID,
		})
	}
}

// @Summary      Delete data service from draft
// @Tags         DataService
// @Description  Delete data service from draft
// @Produce      json
// @Param		 id path int true 				"Data Service ID"
// @Success      200    {object}  string 		   "Data Service added to draft"
// @Failure      400    {object}  error  			   "Bad request"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/dataService/draft/{id} [delete]
func deleteFromDraft(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		if err := r.DeleteDataServiceFromDraft(uint(id), getUserUUIDFromCtx(c)); err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		respMessage(c, http.StatusOK, "deleted")
	}
}
