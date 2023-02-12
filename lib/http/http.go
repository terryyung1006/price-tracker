package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpRequest(req *http.Request, result interface{}) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("[HttpRequest] request failed with error: %s", err.Error())
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[HttpPost] ioutil ReadAll resp body failed with error: %s", err.Error())
	}
	err = json.Unmarshal([]byte(respBody), result)
	if err != nil {
		return fmt.Errorf("[HttpPost] response unmarshal failed with error: %s", err.Error())
	}
	errStruct := ErrorStruct{}
	_ = json.Unmarshal([]byte(respBody), &errStruct)
	if errStruct.Error != "" {
		return fmt.Errorf("[HttpPost] request failed with error: %s", errStruct.Error)
	}
	return nil
}

type ErrorStruct struct {
	Error string
}

func ResponseJson(ctx *gin.Context, data interface{}, err error, httpErrorCode int) {
	if err != nil {
		ctx.JSON(httpErrorCode, gin.H{
			"message": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
