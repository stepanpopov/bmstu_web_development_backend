package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rip/internal/pkg/repo"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sethvargo/go-retry"
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
		draftID, err := r.CreateEncryptDecryptDraft(getUserUUIDFromCtx(c))

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

func formEncryptDecryptRequest(
	r repo.Repository,
	makeCalculationRequest func(uint, []repo.DataService) (int, error),
) func(c *gin.Context) {
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

		if err = retry.Do(context.Background(), retry.WithMaxRetries(3, retry.NewExponential(1*time.Second)), func(ctx context.Context) error {
			status, err := makeCalculationRequest(req.RequestID, dataServices)
			if err != nil {
				if status/100 == 5 {
					return retry.RetryableError(err)
				}

				return err
			}

			return nil
		}); err != nil {
			log.Printf("Error during calling async service: %v", err)
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

type CalculatedReq struct {
	CalculatedList []repo.Calculated `json:"calculated"`
	Token          string            `json:"token"`
	ReqID          int               `json:"req_id"`
}

func calculated(r repo.Repository, secret string) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("in calc")

		var req CalculatedReq
		if err := c.BindJSON(&req); err != nil {
			respMessage(c, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println(req.Token, secret)
		if req.Token != secret {
			respMessageAbort(c, http.StatusForbidden, "invalid token")
			return
		}

		if err := r.UpdateCalculated(uint(req.ReqID), req.CalculatedList); err != nil {
			respMessageAbort(c, http.StatusInternalServerError, "failed to update")
			return
		}

		respMessage(c, http.StatusOK, "updated")
	}
}
