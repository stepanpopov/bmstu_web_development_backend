package api

import (
	"net/http"
	"rip/internal/pkg/api/consts"
	"rip/internal/pkg/repo"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary      Filter EncryptDecryptRequests
// @Tags         EncryptDecryptRequest
// @Description  Get requests filtered by status, start and end date
// @Produce      json
// @Param		 status query string false 					"Status"
// @Param		 start_date query string false 				"Start Date"
// @Param		 end_date query string false 				"End Date"
// @Success      200    {object}  []repo.EncryptDecryptRequestView	   "Data Service Filtered"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/encryptDecryptRequest/filter [get]
func getEncryptDecryptRequests(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		const layoutDate = "2006-01-02"
		statusQuery, hasStatus := c.GetQuery("status")
		startDateQuery, hasStDate := c.GetQuery("start_date")
		endDateQuery, hasEnDate := c.GetQuery("end_date")

		var status repo.Status
		if !hasStatus {
			status = repo.UnknownStatus
		} else {
			var err error
			status, err = repo.FromString(statusQuery)
			if err != nil || status == repo.Deleted {
				respMessageAbort(c, http.StatusBadRequest, "невалидный статус")
				return
			}
		}

		var startDate, endDate time.Time
		if hasStDate {
			var err error
			startDate, err = time.Parse(layoutDate, startDateQuery)
			if err != nil {
				respMessageAbort(c, http.StatusBadRequest, "start_date invlaid format")
				return
			}
		}

		if hasEnDate {
			var err error
			endDate, err = time.Parse(layoutDate, endDateQuery)
			if err != nil {
				respMessageAbort(c, http.StatusBadRequest, "end_date invlaid format")
				return
			}

			if endDate.Before(startDate) {
				respMessageAbort(c, http.StatusBadRequest, "end_date должна быть позже start_date")
				return
			}
		}

		isModerator := c.GetBool(consts.ModeratorCtxParam)
		requests, err := r.GetEncryptDecryptRequests(status, startDate, endDate, getUserUUIDFromCtx(c), isModerator)

		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, requests)
	}
}

// @Summary      Get EncryptDecryptRequest by id
// @Tags         EncryptDecryptRequest
// @Description  Get EncryptDecryptRequest by id
// @Produce      json
// @Param		 req_id path int true 				"EncryptDecryptRequest ID"
// @Success      200    {object}  map[string]any	   "Draft EncryptDecryptRequest and Data Services"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/encryptDecryptRequest/{req_id} [get]
func getEncryptDecryptRequestsByID(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		req, dataServices, err := r.GetEncryptDecryptRequestWithDataByID(uint(id), getUserUUIDFromCtx(c), c.GetBool(consts.ModeratorCtxParam))
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"encDecReq":    req,
			"dataServices": toViewSlice(dataServices),
		})
	}
}

func createDraft(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		draftID, err := r.CreateEncryptDecryptDraft(getUserUUIDFromCtx(c))

		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"draftID": draftID})
	}
}

// @Summary      Delete EncryptDecryptRequest by id
// @Tags         EncryptDecryptRequest
// @Description  Delete EncryptDecryptRequest by id
// @Produce      json
// @Param		 id path int true 				"EncryptDecryptRequest ID"
// @Success      200    {object}  string	   "Delete"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/encryptDecryptRequest/{req_id} [delete]
func deleteEncryptDecryptRequest(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("req_id"), 10, 64)

		err := r.DeleteEncryptDecryptRequestByID(uint(id))
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
		}

		respMessage(c, http.StatusOK, "deleted")
	}
}

// @Summary      Form EncryptDecryptRequest by moderator
// @Tags         EncryptDecryptRequest
// @Description  Form EncryptDecryptRequest by moderator
// @Produce      json
// @Param		 id path int true 						"EncryptDecryptRequest ID"
// @Success      200    {object}  map[string]any	   	"EncryptDecryptRequest and DataServices"
// @Failure      400    {object}  error  			   "Bad Request"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/encryptDecryptRequest/form/{id} [put]
func formEncryptDecryptRequest(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		err := r.FormEncryptDecryptRequestByID(uint(id))
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		req, dataServices, err := r.GetEncryptDecryptRequestWithDataByID(uint(id), getUserUUIDFromCtx(c), false)
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"encDecReq":    req,
			"dataServices": toViewSlice(dataServices),
		})
	}
}

// @Summary      Update EncryptDecryptRequest by moderator
// @Tags         EncryptDecryptRequest
// @Description  Update EncryptDecryptRequest by moderator
// @Produce      json
// @Accept		 json
// @Param		 id path int true 				"EncryptDecryptRequest ID"
// @Param		 action body int true 				"Action: finish or reject"
// @Success      200    {object}  string	   "Finished/Rejected"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/encryptDecryptRequest/update_moderator/{id} [put]
func updateModeratorEncryptDecryptRequest(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		type Action struct {
			Action string
		}
		actionReq := Action{}

		if err := c.BindJSON(&actionReq); err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		if actionReq.Action == "reject" {
			if err := r.RejectEncryptDecryptRequestByID(uint(id), getUserUUIDFromCtx(c)); err != nil {
				respMessageAbort(c, http.StatusBadRequest, err.Error())
				return
			}
			respMessage(c, http.StatusOK, "rejected")
		} else if actionReq.Action == "finish" {
			if err := r.FinishEncryptDecryptRequestByID(uint(id), getUserUUIDFromCtx(c)); err != nil {
				respMessageAbort(c, http.StatusBadRequest, err.Error())
				return
			}
			respMessage(c, http.StatusOK, "finished")
		}
	}
}

// @Summary      Delete DataService from EncryptDecryptRequest
// @Tags         EncryptDecryptRequest
// @Description  Delete DataService from EncryptDecryptRequest
// @Produce      json
// @Param		 req_id path int true 				"EncryptDecryptRequest ID"
// @Param		 data_id path int true 				"DataService ID"
// @Success      200    {object}  string	   "Delete"
// @Failure      404    {object}  error  			   "Not found"
// @Failure      500    {object}  error  			   "Server error"
// @Router       /api/encryptDecryptRequest/{req_id}/delete/{data_id} [delete]
func deleteDataFromEncryptDecryptRequest(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		reqID, _ := strconv.ParseUint(c.Param("req_id"), 10, 64)
		dataID, _ := strconv.ParseUint(c.Param("data_id"), 10, 64)

		if err := r.DeleteDataServiceFromEncryptDecryptRequest(uint(dataID), uint(reqID)); err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		respMessage(c, http.StatusOK, "deleted")
	}
}
