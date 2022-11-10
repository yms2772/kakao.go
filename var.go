package kakaogo

import (
	"log"
	"net/http"
	"os"
)

const (
	appVerSion = "3.4.3"
	osVersion  = "10.0"
	agent      = "win32"
	lang       = "ko"
	prtVersion = "1"
	deviceName = "kakao.go"
	deviceUUID = "a2FrYW9nb2RldjIwMjI="
	dtype      = 1
	ntype      = 0
	mccmnc     = "999"

	authHeader   = agent + "/" + appVerSion + "/" + lang
	uthUserAgent = "KT/" + appVerSion + " Wd/" + osVersion + " " + lang
)

const (
	requestPasscodeURL = "https://ac-sb-talk.kakao.com/win32/account/request_passcode.json"
	loginUrl           = "https://ac-sb-talk.kakao.com/win32/account/login.json"
	registerDeviceUrl  = "https://ac-sb-talk.kakao.com/win32/account/register_device.json"
	mediaURL           = "https://up-m.talk.kakao.com/upload"
)

var (
	httpClient = &http.Client{}
	logger     = log.New(os.Stdout, "INFO: ", log.LstdFlags)
)
