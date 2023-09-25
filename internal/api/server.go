package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

type PlainData struct {
	ID   int
	Name string
	Blob string
}

func StartServer() {
	log.Println("Server start up")

	PlainDatas := []PlainData{
		PlainData{
			ID:   1,
			Name: "Encode your secrets",
			Blob: "secret",
		},
		PlainData{
			ID:   2,
			Name: "Decode your life",
			Blob: "01001001000100",
		},
	}

	r := gin.Default()

	/* r.SetFuncMap(template.FuncMap{
		"contains": strings.Contains,
	}) */

	r.LoadHTMLGlob("templates/*")

	r.GET("/", showAllPlainData(PlainDatas))

	r.GET("/filter", filterPlainDatas(PlainDatas))

	r.GET("/service/*id", showPlainData(PlainDatas))

	r.Static("/image", "./resources")

	r.Run()

	log.Println("Server down")
}
