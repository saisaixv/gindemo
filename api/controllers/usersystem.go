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
	"github.com/saisai/gindemo/common"
	"github.com/saisai/gindemo/models"
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

	_, errCode := models.Register(req)
	if errCode != 0 {
		rsp.Error_code = errCode
	}
}

func checkToken(head map[string]interface{}) int {

	token := head["x-us-token"]

	if token == nil {
		return msg.ErrUnauthorized
	}

	str := strings.Split(token.(string), common.SPLIT)
	if len(str) != 2 {
		return msg.ErrUnauthorized
	}

	key := str[0] + common.KEY_TOKEN

	has, value := cache.DoStrGet(key)
	if !has {
		return msg.ErrUnauthorized
	}

	if value != token {
		return msg.ErrUnauthorized
	}

	cache.DoExpire(key, common.ONE_MINUTE)

	return msg.OK
}

func Login(ctx *gin.Context) {

	req := new(msg.LoginReq)
	rsp := new(msg.LoginRsp)
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
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	err2 := bindBody(ctx, req)
	if err2 != nil {
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	models.Login(req, rsp)

	//	token := head["x-us-token"]
	//	str := strings.Split(token.(string), common.SPLIT)

	//	b, err2 := models.LoginCydexManager(str[0], "token", rsp.Token)
	//	if err2 != nil {
	//		rsp.Error_code = msg.ErrCydexManagerAuthError
	//		return
	//	}
	//	if !b {
	//		rsp.Error_code = msg.ErrCydexManagerAuthError
	//		return
	//	}

	fmt.Println(rsp)

}

func Logout(ctx *gin.Context) {
	rsp := new(msg.BaseRsp)
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
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	ret := checkToken(head)

	if ret != msg.OK {
		rsp.Error_code = ret
		return
	}
	//	//当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1 。 否则，以毫秒为单位，返回 key 的剩余生存时间。
	//	if ret == -2 {
	//		fmt.Println("ret == -2")
	//		rsp.Error_code = msg.ErrUnauthorized
	//		return
	//	}

	token := head["x-us-token"]
	str := strings.Split(token.(string), common.SPLIT)

	key := str[0]
	has := cache.DoDel(key)
	if !has {
		fmt.Println("clear cache error")
	}
}

func Info(ctx *gin.Context) {
	rsp := new(msg.InfoRsp)

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

	ret := checkToken(head)

	if ret != msg.OK {
		rsp.Error_code = ret
		return
	}

	token := head["x-us-token"]
	str := strings.Split(token.(string), common.SPLIT)

	err2 := models.UserInfo(str[0], rsp)
	if err2 != nil {
		fmt.Println(err2.Error())
		rsp.Error_code = msg.ErrServerInternalError
		return
	}

}

func AddIdentifyType(ctx *gin.Context) {

	req := new(msg.AddIdentifyTypeReq)
	rsp := new(msg.AddIdentifyTypeRsp)
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
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	ret := checkToken(head)

	if ret != msg.OK {
		rsp.Error_code = ret
		return
	}

	err2 := bindBody(ctx, req)
	if err2 != nil {
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	fmt.Println(req)

	rsp.Error_code = models.AddIdentifyType(req)
}

func Authentication(ctx *gin.Context) {

	req := new(msg.AuthenticationReq)
	rsp := new(msg.AuthenticationRsp)
	rsp.Error_code = msg.OK
	defer func() {
		if rsp.Error_code != msg.OK {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusOK, rsp)
	}()

	err := bindBody(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	if ret := models.Authentication(req); ret != msg.OK {
		fmt.Println("models.Authentication error")
		rsp.Error_code = ret
		return
	}

	code := models.GetUerInfo(req.Token, rsp)
	if code != msg.OK {
		fmt.Println("models.GetUerInfo error")
		rsp.Error_code = code
		return
	}

}
