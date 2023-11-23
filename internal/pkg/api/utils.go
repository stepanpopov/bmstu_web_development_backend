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

func toView(d repo.DataService) repo.DataServiceView {
	return repo.DataServiceView{
		DataID:   d.DataID,
		DataName: d.DataName,
		Encode:   d.Encode,
		Blob:     d.Blob,
		Active:   d.Active,
		ImageURL: s3Url + d.ImageUUID.String(),
	}
}

func toViewSlice(dd []repo.DataService) []repo.DataServiceView {
	var view []repo.DataServiceView
	for _, d := range dd {
		view = append(view, toView(d))
	}
	return view
}
