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
		var request struct {
			Status    string `json:"status"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
		}
		const layoutDate = "2006-01-02"

		if err := c.BindJSON(&request); err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		status, err := repo.FromString(request.Status)
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, "невалидный статус")
			return
		}

		startDate, err := time.Parse(layoutDate, request.StartDate)
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, "start_date invlaid format")
			return
		}

		endDate, err := time.Parse(layoutDate, request.EndDate)
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, "end_date invlaid format")
			return
		}

		if endDate.Before(startDate) {
			respMessageAbort(c, http.StatusBadRequest, "end_date должна быть позже start_date")
			return
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
			"dataServices": dataServices,
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
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

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
			"dataServices": dataServices,
		})
	}
}

func rejectEncryptDecryptRequest(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		if err := r.RejectEncryptDecryptRequestByID(uint(id), moderatorID); err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		respMessage(c, http.StatusOK, "rejected")
	}
}

func finishEncryptDecryptRequest(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		if err := r.FinishEncryptDecryptRequestByID(uint(id), moderatorID); err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		respMessage(c, http.StatusOK, "finished")
	}
}
