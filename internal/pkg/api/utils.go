package api

import (
	"rip/internal/pkg/repo"

	"github.com/gin-gonic/gin"
)

func notFound(c *gin.Context) {
	respMessage(c, 404, "not found")
}

func respMessage(c *gin.Context, code uint, message string) {
	c.JSON(int(code), gin.H{"message": message})
}

func respMessageAbort(c *gin.Context, code uint, message string) {
	c.AbortWithStatusJSON(int(code), gin.H{"message": message})
}

const s3Url = "http://localhost:9001/images/"

func toView(d repo.DataService) DataServiceView {
	return DataServiceView{
		DataID:   d.DataID,
		DataName: d.DataName,
		Encode:   d.Encode,
		Blob:     d.Blob,
		Active:   d.Active,
		ImageURL: s3Url + d.ImageUUID.String(),
	}
}

type DataServiceView struct {
	DataID   uint   `json:"data_id"`
	DataName string `json:"data_name"`
	Encode   bool   `json:"encode"`
	Blob     string `json:"blob"`
	Active   bool   `json:"active"`
	ImageURL string `json:"image_url,omitempty"`
}

func toViewSlice(dd []repo.DataService) []DataServiceView {
	var view []DataServiceView
	for _, d := range dd {
		view = append(view, toView(d))
	}
	return view
}
