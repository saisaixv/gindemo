package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/saisai/gindemo/api/msg"
	"github.com/saisai/gindemo/common"

	"github.com/saisai/gindemo/utils"
	"github.com/saisai/gindemo/utils/cache"
	"github.com/saisai/gindemo/utils/captcha"
)

func Register(req *msg.RegisterReq) (string, int) {
	if req.Nickname == "" ||
		req.Credential == "" {
		return "", msg.ErrInvalidParam
	}

	if req.Email == "" && req.Phone == "" {
		return "", msg.ErrInvalidParam
	}

	userId := utils.GetMongoObjectId()

	user := User{Id: userId, Nickname: req.Nickname, Avatar: req.Avatar, Sex: req.Sex}
	newUser := new(User)
	has, err := DB().Where("nickname=?", req.Nickname).Get(newUser)
	if err != nil {
		fmt.Println(err.Error())
		return "", msg.ErrInvalidParam
	}
	if has {
		return "", msg.ErrNicknameIsExist
	}

	affected, err := DB().Insert(user)
	if err != nil {
		fmt.Println(err.Error())
		return "", msg.ErrInvalidParam
	}

	fmt.Println(affected)

	if req.Email != "" {
		auths := UserAuths{UserId: userId, IdentifyType: "email",
			Identifier: req.Email, Credential: req.Credential, Latestlogintime: "1970-1-1 0:0:0"}
		_, err := DB().Insert(auths)
		if err != nil {
			fmt.Println(err.Error())
			return "", msg.ErrInvalidParam
		}
	}

	if req.Phone != "" {
		auths := UserAuths{UserId: userId, IdentifyType: "phone",
			Identifier: req.Phone, Credential: req.Credential, Latestlogintime: "1970-1-1 0:0:0"}
		_, err := DB().Insert(auths)
		if err != nil {
			fmt.Println(err.Error())
			return "", msg.ErrInvalidParam
		}
	}

	return userId, msg.OK
}

func Login(req *msg.LoginReq, rsp *msg.LoginRsp) {
	if req.Identify_type == "" || req.Identifier == "" || req.Credential == "" {
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	auth := new(UserAuths)

	has, err := DB().Where("identify_type = ? and identifier = ?", req.Identify_type, req.Identifier).Get(auth)
	if err != nil {
		fmt.Println(err.Error())
		rsp.Error_code = msg.ErrInvalidParam
		return
	}

	if !has {
		rsp.Error_code = msg.ErrAccountNotExist
		return
	}

	key_login_err := auth.UserId + common.KEY_LOGIN_ERROR_COUNT

	var errCount int
	has, strCount := cache.DoStrGet(key_login_err)

	if has {
		errCount, err = strconv.Atoi(strCount)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if errCount > 10 {
		rsp.Error_code = msg.ErrTooManyLoginError
		return
	}
	//	else if errCount >= 3 {

	//		if req.CaptchaId == "" || req.Value == "" {
	//			rsp.Error_code = msg.ErrInvalidParam
	//			return
	//		}

	//		has = captcha.VerifyCaptcha(req.CaptchaId, req.Value)
	//		if !has {
	//			captchaId := captcha.GenerateCaptcha()
	//			rsp.CaptchaId = captchaId
	//			rsp.Error_code = msg.ErrCaptchaError
	//			rsp.CaptchaUrl = "http://192.168.150.130:9007/usersystem/api/v1/captcha/" + captchaId + ".png"
	//			return
	//		}
	//	}

	if auth.Credential != req.Credential {
		errCount = errCount + 1
		cache.DoStrSet(key_login_err, strconv.Itoa(errCount), common.FIVE_MINUTE)
		rsp.Error_code = msg.ErrPasswordError
		rsp.ErrCount = errCount
		captchaId := captcha.NewLen(4)
		rsp.CaptchaId = captchaId
		rsp.CaptchaUrl = "http://192.168.150.130:9007/usersystem/api/v1/captcha/" + captchaId + ".png"
		return
	}

	token := auth.UserId + common.SPLIT + utils.GetToken()
	has = cache.DoStrSet(auth.UserId+common.KEY_TOKEN, token, common.FIVE_MINUTE)
	//	has = cache.DoExpire(token, common.ONE_MINUTE)
	if !has {
		rsp.Error_code = msg.ErrServerInternalError
		captchaId := captcha.NewLen(4)
		rsp.CaptchaId = captchaId
		rsp.CaptchaUrl = "http://192.168.150.130:9007/usersystem/api/v1/captcha/" + captchaId + ".png"

		return
	}
	rsp.Token = token

}

func UserInfo(userId string, rsp *msg.InfoRsp) error {
	user := new(User)
	has, err := DB().Where("id = ?", userId).Get(user)
	if err != nil {
		return err
	}
	if !has {
		err2 := fmt.Errorf("user not exist!")
		return err2
	}

	fmt.Println(user)

	auths := make([]UserAuths, 0)
	err = DB().Where("user_id=?", userId).Find(&auths)
	if err != nil {
		return err
	}

	fmt.Println(auths)

	rsp.Id = userId
	rsp.Nickname = user.Nickname
	rsp.Avatar = user.Avatar
	rsp.Sex = user.Sex
	rsp.CreateTime = user.CreateTime

	for _, auth := range auths {
		if auth.IdentifyType == "email" {
			rsp.Email = auth.Identifier
		} else if auth.IdentifyType == "phone" {
			rsp.Phone = auth.Identifier
		}
	}

	return nil
}

func GetUerInfo(token string, rsp *msg.AuthenticationRsp) (error_code int) {

	str := strings.Split(token, common.SPLIT)
	//	fmt.Println(str)
	if len(str) != 2 {
		return msg.ErrUnauthorized
	}

	userId := str[0]

	user := new(User)
	has, err := DB().Where("id = ?", userId).Get(user)
	if err != nil {
		return msg.ErrServerInternalError
	}
	if !has {
		//		err2 := fmt.Errorf("user not exist!")
		return msg.ErrServerInternalError
	}

	fmt.Println(user)

	auths := make([]UserAuths, 0)
	err = DB().Where("user_id=?", userId).Find(&auths)
	if err != nil {
		return msg.ErrServerInternalError
	}

	fmt.Println(auths)

	rsp.Id = userId
	rsp.Nickname = user.Nickname
	rsp.Avatar = user.Avatar
	rsp.Sex = user.Sex
	rsp.CreateTime = user.CreateTime

	for _, auth := range auths {
		if auth.IdentifyType == "email" {
			rsp.Email = auth.Identifier
		} else if auth.IdentifyType == "phone" {
			rsp.Phone = auth.Identifier
		}
	}

	return msg.OK
}

func AddIdentifyType(req *msg.AddIdentifyTypeReq) int {
	auth := new(UserAuths)
	has, err := DB().Where("user_id=? and identify_type=?", req.User_id, req.Identify_type).Get(auth)
	if err != nil {
		fmt.Println(err.Error())
		return msg.ErrServerInternalError
	}

	if has {
		return msg.ErrIdentifyTypeExist
	}

	auth.IdentifyType = req.Identify_type
	auth.UserId = req.User_id
	auth.Identifier = req.Identifier
	auth.Credential = req.Credential
	auth.Latestlogintime = "1970-1-1 0:0:0"

	_, err = DB().Insert(auth)
	if err != nil {
		fmt.Println(err.Error())
		return msg.ErrServerInternalError
	}
	return msg.OK
}
