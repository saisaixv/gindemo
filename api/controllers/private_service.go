package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/saisai/gindemo/api/msg"
	ss_http "github.com/saisai/utils/http"

	"github.com/gin-gonic/gin"
)

func Private_Register(ctx *gin.Context) {

	errNum := 0

	message := "success"

	defer func() {
		if errNum != 0 {
			message = fmt.Sprintf("{\"error\":%d}", errNum)
		}
		ctx.JSON(http.StatusOK, message)
	}()

	header := map[string][]string{
		"time-zone":       {"-8"},
		"x-us-authtype":   {"1"},
		"accept_language": {"zh"},
	}

	content := "{\"nickname\":\"saisai\",\"avatar\":\"www.baidu.com\",\"phone\":\"15033156272\",\"email\":\"478306328@qq.com\",\"sex\":1}"

	rsp, err := ss_http.CallAPI("POST",
		"http://192.168.150.130:9007/usersystem/api/v1/register",
		header,
		[]byte(content),
		5*time.Second)
	if err != nil {
		errNum = -1
		return
	}

	body, err := ss_http.ResponseBody(rsp)
	if err != nil {
		fmt.Println(err.Error())
		errNum = -2
		return
	}

	fmt.Println(body)

	bean := new(msg.RegisterRsp)

	err = json.Unmarshal(body, bean)
	if err != nil {
		fmt.Println(err.Error())
		errNum = -3
		return
	}

	fmt.Println(bean)

}
