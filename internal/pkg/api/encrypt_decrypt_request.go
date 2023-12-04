package api

import (
	"net/http"
	"rip/internal/pkg/repo"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

		requests, err := r.GetEncryptDecryptRequests(status, startDate, endDate)

		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, requests)
	}
}

func getEncryptDecryptRequestsByID(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		req, dataServices, err := r.GetEncryptDecryptRequestWithDataByID(uint(id))
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
		draftID, err := r.CreateEncryptDecryptDraft(creatorID)

		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"draftID": draftID})
	}
}

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

func formEncryptDecryptRequest(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		err := r.FormEncryptDecryptRequestByID(uint(id))
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		req, dataServices, err := r.GetEncryptDecryptRequestWithDataByID(uint(id))
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
			if err := r.RejectEncryptDecryptRequestByID(uint(id), moderatorID); err != nil {
				respMessageAbort(c, http.StatusBadRequest, err.Error())
				return
			}
			respMessage(c, http.StatusOK, "rejected")
		} else if actionReq.Action == "finish" {
			if err := r.FinishEncryptDecryptRequestByID(uint(id), moderatorID); err != nil {
				respMessageAbort(c, http.StatusBadRequest, err.Error())
				return
			}
			respMessage(c, http.StatusOK, "finished")
		}
	}
}

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
