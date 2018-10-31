package msg

type BaseRsp struct {
	Error_code int `json:"error_code"`
}

type User struct {
	Id         string `json:"id" xorm:"varchar(24) pk"`
	Nickname   string `json:"nickname" xorm:"varchar(100)"`
	Avatar     string `json:"avatar" xorm:"varchar(100)"`
	Sex        int    `json:"sex" xorm:"int"`
	CreateTime string `json:"createtime" xorm:"DateTime created"`
}

type UserInfo struct {
	User
	UpdateTime string `json:"updatetime" xorm:"DateTime updated"`
	Phone      string `json:"phone" xorm:"varchar(100)"`
	Email      string `json:"email" xorm:"varchar(100)"`
}

type RegisterReq struct {
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Credential string `json:"credential"`
	Sex        int    `json:"sex"`
}

type RegisterRsp struct {
	BaseRsp
}

type LoginReq struct {
	Identify_type string `json:"identify_type"`
	Identifier    string `json:"identifier"`
	Credential    string `json:"credential"`
	CaptchaId     string `json:"captcha_id"`
	Value         string `json:"value"`
}

type LoginRsp struct {
	BaseRsp
	Token      string `json:"token"`
	ErrCount   int    `json:"err_count"`
	CaptchaId  string `json:"captcha_id"`
	CaptchaUrl string `json:"captcha_url"`
}

type InfoRsp struct {
	BaseRsp
	UserInfo
}

type AddIdentifyTypeReq struct {
	User_id       string `json:"user_id"`
	Identify_type string `json:"identify_type"`
	Identifier    string `json:"identifier"`
	Credential    string `json:"credential"`
}

type AddIdentifyTypeRsp struct {
	BaseRsp
}

type AuthenticationReq struct {
	Token string `json:"token"`
}

type AuthenticationRsp struct {
	BaseRsp
	UserInfo
}

type LoginCydexManagerReq struct {
	AuthType string `json:"authtype"`
	Auth     string `json:"auth"`
}

type LoginCydexManagerRsp struct {
	BaseRsp
	Token string `json:"token"`
}
