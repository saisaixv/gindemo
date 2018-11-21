package models

import (
	"encoding/json"
	//	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/saisai/gindemo/api/msg"
	"github.com/saisai/gindemo/common"
	"github.com/saisai/gindemo/utils/cache"
	ss_http "github.com/saisai/gindemo/utils/http"
)

func Authentication(req *msg.AuthenticationReq) int {
	str := strings.Split(req.Token, common.SPLIT)
	//	fmt.Println(str)
	if len(str) != 2 {
		return msg.ErrUnauthorized
	}

	ret, token := cache.DoStrGet(str[0] + common.KEY_TOKEN)
	if !ret {

		return msg.ErrUnauthorized
	}

	if token != req.Token {
		return msg.ErrUnauthorized
	}
	return msg.OK
}

func LoginCydexManager(user_id, authtype, auth string) (bool, error) {
	req := new(msg.LoginCydexManagerReq)
	rsp := new(msg.LoginCydexManagerRsp)

	header := http.Header{
		"x-us-authtype":   {"1"},
		"accept-language": {"zh"},
		"time-zone":       {"8"},
	}

	req.AuthType = authtype
	req.Auth = auth
	content, err := json.Marshal(req)
	if err != nil {
		return false, err
	}
	ret, err := ss_http.CallAPI("POST", "http://127.0.0.1:9005/cydex/api/v1/thirdparty_auth", content, header, 5*time.Second)
	if err != nil {
		return false, err
	}

	body, err := ss_http.ResponseBody(ret)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(body, rsp)
	if err != nil {
		return false, err
	}

	cache.DoStrSet(user_id+common.KEY_CYDEX_AUTH, rsp.Token, common.FIVE_MINUTE)

	return true, nil
}
