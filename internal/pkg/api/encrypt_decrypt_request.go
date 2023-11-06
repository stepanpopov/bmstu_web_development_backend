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
			Status    string    `json:"status"`
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}
		if err := c.BindJSON(&request); err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		status, err := repo.FromString(request.Status)
		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, "невалидный статус")
			return
		}

		if request.EndDate.Before(request.StartDate) {
			respMessageAbort(c, http.StatusBadRequest, "end_date должна быть позже start_date")
			return
		}

		requests, err := r.GetEncryptDecryptRequests(status, request.StartDate, request.EndDate)

		if err != nil {
			respMessageAbort(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, requests)
	}
}

func getEncryptDecryptRequestsByID(r repo.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id")[1:], 10, 64)

		// TODO: creator
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
