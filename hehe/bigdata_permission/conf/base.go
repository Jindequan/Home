package conf

import (
	"bigdata_permission/pkg"
	"encoding/json"
	"io/ioutil"
)

type BaseConfig struct {
	App      string `json:"APP"`
	DBConfig struct {
		Dsn          string `json:"Dsn"`
		MaxIdleConns int    `json:"MaxIdleConns"`
		MaxOpenConns int    `json:"MaxOpenConns"`
	} `json:"DbConfig"`
	Domain struct {
		BigDataPermission string `json:"BigDataPermission"`
		Sso string `json:"Sso"`
	} `json:"Domain"`
	ExternalApiDomain struct {
		BigData string `json:"BigData"`
	} `json:"ExternalApiDomain"`
	HTTPClient struct {
		Dial      int `json:"Dial"`
		KeepAlive int `json:"KeepAlive"`
	} `json:"HttpClient"`
	HTTPServerConfig struct {
		Addr         string `json:"Addr"`
		ReadTimeout  int    `json:"ReadTimeout"`
		WriteTimeout int    `json:"WriteTimeout"`
	} `json:"HttpServerConfig"`
	LogPath      string `json:"LogPath"`
	TempFilePath string `json:"TempFilePath"`
	RedisConfig  struct {
		Addr     string `json:"Addr"`
		Password string `json:"Password"`
		Db       int    `json:"Db"`
	} `json:"RedisConfig"`
	SsoConfig struct {
		PlatformID string `json:"PlatformId"`
		SecretKey  string `json:"SecretKey"`
	} `json:"SsoConfig"`
	SwaagerToken string `json:"SwaagerToken"`
	Token        struct {
		ExpireTime int    `json:"ExpireTime"`
		Salt       string `json:"Salt"`
	} `json:"Token"`
	LogCenterConfig struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	} `json:"LogCenterConfig"`
	QueueConfig struct{
		Host string `json:"Host"`
		Ip string `json:"Ip"`
	}
	WeChat *WeChat
}

type WeChat struct {
	ApiHost string `json:"ApiHost"`
	ApprovalTemplateId struct {
		ModuleApprovalTemplateId  string `json:"ModuleApprovalTemplateId"`
	} `json:"ApprovalTemplateId"`
	ApprovalApp *WeChatAppWithToken `json:"ApprovalApp"`
}

type WeChatAppWithToken struct {
	AppID          string
	AgentID        int64
	CorpSecret     string
	Token          string `json:"Token"`
	EncodingAesKey string `json:"EncodingAesKey"`
}

var BaseConf *BaseConfig

func init() {
	path := "./settings/base_config_" + pkg.GetEnv() + ".json"
	pkg.DoNothing()
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bs, &BaseConf)
}
