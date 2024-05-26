package config

import (
	"encoding/base64"
	"os"

	"github.com/spf13/cast"

	"nh-downloader/consts"
)

type ConfigModel struct {
	port         string
	cachePath    string
	downloadPath string
	logsPath     string
	proxy        string
	retried      int
	lanraragi    LanraragiConf
}

type LanraragiConf struct {
	Url       string
	KeyBase64 string
}

var config *ConfigModel

// Init initialize config with env.
func Init() {
	config = &ConfigModel{}
	config.port = getEnv(consts.NHDL_PORT_KEY, consts.NHDL_PORT_DEFAULT)
	config.cachePath = getEnv(consts.NHDL_CACHE_PATH_KEY, consts.NHDL_CACHE_PATH_DEFAULT)
	config.downloadPath = getEnv(consts.NHDL_DOWNLOAD_PATH_KEY, consts.NHDL_DOWNLOAD_PATH_DEFAULT)
	config.logsPath = getEnv(consts.NHDL_LOGS_PATH_KEY, consts.NHDL_LOGS_PATH_DEFAULT)
	config.proxy = getEnv(consts.NHDL_PROXY_KEY, consts.NHDL_PROXY_DEFAULT)
	config.retried = cast.ToInt(getEnv(consts.NHDL_RETRIED_KEY, consts.NHDL_RETRIED_DEFAULT))

	lanraragiUrl := getEnv(consts.NHDL_LANRARAGI_URL_KEY, consts.NHDL_LANRARAGI_URL_DEFAULT)
	if lanraragiUrl[len(lanraragiUrl)-1] != '/' {
		lanraragiUrl += "/"
	}
	config.lanraragi.Url = lanraragiUrl
	lanraragiKey := getEnv(consts.NHDL_LANRARAGI_KEY_KEY, "")
	if lanraragiKey != "" {
		config.lanraragi.KeyBase64 = base64.StdEncoding.EncodeToString([]byte(lanraragiKey))
	}
}

// getEnv get value from env, or default.
func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}

func Port() string {
	return config.port
}

func Lanraragi() LanraragiConf {
	return config.lanraragi
}

func CachePath() string {
	return config.cachePath
}

func DownloadPath() string {
	return config.downloadPath
}

func LogsPath() string {
	return config.logsPath
}

func Proxy() string {
	return config.proxy
}

func Retried() int {
	return config.retried
}
