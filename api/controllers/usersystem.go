package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/saisai/utils/cache"

	"github.com/saisai/gindemo/api/msg"
	"github.com/saisai/gindemo/service"

	"github.com/gin-gonic/gin"
)

func getHeaders(ctx *gin.Context) map[string]interface{} {

	head := make(map[string]interface{})

	timeZone := ctx.Request.Header.Get("time-zone")
	if timeZone != "" {
		head["time-zone"] = timeZone
	}
	acceptLanguage := ctx.Request.Header.Get("accept-language")
	if timeZone != "" {
		head["accept-language"] = acceptLanguage
	}
	authtype := ctx.Request.Header.Get("x-us-authtype")
	if authtype != "" {
		head["x-us-authtype"] = authtype
	}
	token := ctx.Request.Header.Get("x-us-token")
	if token != "" {
		head["x-us-token"] = token
	}
	return head
}

func headCheck(head map[string]interface{}) *service.Error {

	fmt.Println(head)

	timeZone, err := strconv.Atoi(head["time-zone"].(string))
	if err != nil {
		fmt.Println(err.Error())
	}

	language := head["accept-language"].(string)
	//	authtype := head["x-us-authtype"].(int)

	if math.Abs(float64(timeZone)) > 12 {
		return service.ErrInvalidParam
	}

	if strings.ToUpper(language) != "EN" &&
		strings.ToUpper(language) != "ZH" {
		return service.ErrInvalidParam
	}
	return nil
}

func authCheck(ctx *gin.Context) *service.Error {
	head := getHeaders(ctx)
	err := headCheck(head)
	if err != nil {
		return err
	}
	token := head["x-us-token"]
	if token == "" {
		return service.ErrUnauthorized
	}

	b, _ := cache.DoStrGet(token.(string))
	if !b {
		return service.ErrUnauthorized
	}
	b = cache.DoExpire(token.(string), service.Redis_key_token_expire)
	if !b {
		return service.ErrUnauthorized
	}
	return nil
}

func bindBody(ctx *gin.Context, req interface{}) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, req)
}

func Register(ctx *gin.Context) {

	req := new(msg.RegisterReq)
	rsp := new(msg.RegisterRsp)

	rsp.Error_code = msg.OK

	defer func() {
		if rsp.Error_code != msg.OK {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusOK, rsp)
	}()

	head := getHeaders(ctx)
	err := headCheck(head)
	if err != nil {
		fmt.Println(err.Error())
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	err2 := bindBody(ctx, req)
	if err2 != nil {
		fmt.Println(err2.Error())
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	fmt.Println(req)

}
