package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ERROR   = -1
	SUCCESS = 1
)

type Result struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Page   *Page       `json:"page"`
}

type Page struct {
	Count int `json:"count"`
	Index int `json:"index"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

func (result *Result) Response(c *gin.Context) {
	c.JSON(http.StatusOK, result)
}

func (result *Result) IsOK() bool {
	return result.Status == 1
}

func Success() Result {
	result := new(Result)
	result.Status = SUCCESS
	result.Msg = "ok"
	return *result
}

func SuccessWithMsg(msg string) Result {
	result := new(Result)
	result.Status = SUCCESS
	result.Msg = msg
	return *result
}

func SuccessWithData(data interface{}) Result {
	result := new(Result)
	result.Status = SUCCESS
	result.Msg = "ok"
	result.Data = data
	return *result
}

func SuccessWithDetail(data interface{}, page *Page) Result {
	result := new(Result)
	result.Status = SUCCESS
	result.Msg = "ok"
	result.Data = data
	result.Page = page
	return *result
}

func Fail() Result {
	result := new(Result)
	result.Status = ERROR
	result.Msg = "fail"
	return *result
}

func FailWithMsg(msg string) Result {
	result := new(Result)
	result.Status = ERROR
	result.Msg = msg
	return *result
}

func FailWithErr(msg string, err error) Result {
	result := new(Result)
	result.Status = ERROR
	result.Msg = msg + ":" + err.Error()
	return *result
}
