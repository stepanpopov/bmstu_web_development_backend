package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"rip/internal/pkg/api/consts"
	"rip/internal/pkg/repo"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

const s3Url = "http://localhost:9000/avatars/"

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

func toViewSlice(dd []repo.DataService) []DataServiceView {
	var view []DataServiceView
	for _, d := range dd {
		view = append(view, toView(d))
	}
	return view
}

func toViewWithOptResult(d repo.DataServiceWithOptResult) DataServiceView {
	return DataServiceView{
		DataID:   d.DataID,
		DataName: d.DataName,
		Encode:   d.Encode,
		Blob:     d.Blob,
		Active:   d.Active,
		ImageURL: s3Url + d.ImageUUID.String(),
		Result:   d.Result,
		Success:  d.Success,
	}
}

func toViewWithOptResultSlice(dd []repo.DataServiceWithOptResult) []DataServiceView {
	var view []DataServiceView
	for _, d := range dd {
		view = append(view, toViewWithOptResult(d))
	}
	return view
}

type DataServiceView struct {
	DataID   uint    `json:"data_id"`
	DataName string  `json:"data_name"`
	Encode   bool    `json:"encode"`
	Blob     string  `json:"blob"`
	Active   bool    `json:"active"`
	ImageURL string  `json:"image_url,omitempty"`
	Result   *string `json:"result,omitempty"`
	Success  *bool   `json:"success,omitempty"`
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getUserUUIDFromCtx(c *gin.Context) uuid.UUID {
	userID, _ := c.Get(consts.UserUUIDCtxParam)
	if userID == nil {
		return uuid.Nil
	}
	userIDCasted := userID.(uuid.UUID)
	return userIDCasted
}

func getJWTStr(gCtx *gin.Context) string {
	jwtStr := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, consts.JwtPrefix) {
		return ""
	}
	// отрезаем префикс
	return jwtStr[len(consts.JwtPrefix):]
}

type Calculate struct {
	ID   uint   `json:"id"`
	Data string `json:"data"`
}

type CalculateRequest struct {
	ReqID uint        `json:"req_id"`
	Calc  []Calculate `json:"calc"`
}

func (s *Server) makeCalculationRequest(reqID uint, dataServices []repo.DataServiceWithOptResult) (int, error) {
	// Define the data you want to send
	calc := make([]Calculate, 0, len(dataServices))
	for _, ds := range dataServices {
		calc = append(calc, Calculate{ID: ds.DataID, Data: ds.Blob})
	}

	calcReq := CalculateRequest{
		Calc:  calc,
		ReqID: reqID,
	}

	// Marshal the data into JSON
	json_data, err := json.Marshal(&calcReq)
	if err != nil {
		return 0, err
	}

	// Create a new HTTP request
	resp, err := http.Post(s.calculateCallback, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBuf := new(bytes.Buffer)
		respBuf.ReadFrom(resp.Body)
		return resp.StatusCode, errors.New(respBuf.String())
	}

	return 0, nil
}
