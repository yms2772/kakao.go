package kakaogo

import (
	"net"
)

type Kakao struct {
	email    string
	password string
	passcode string

	stream    net.Conn
	accessKey string

	cryptoManager *CryptoManager

	writer *Writer

	processingBuffer []byte
	processingHeader []byte
	processingSize   int

	OnReady   func()
	OnMessage func(chat *Message)
	OnRead    func(channel *Channel, body map[string]interface{})
	OnJoin    func(packet *Packet, channel *Channel)
	OnQuit    func(packet *Packet, channel *Channel)
}

type Packet struct {
	packetID   int
	statusCode int
	packetName string
	bodyType   int
	bodySize   int
	body       []byte
}

type CryptoManager struct {
	aesKey []byte
}

type Writer struct {
	cryptoManager *CryptoManager
	stream        net.Conn
	packetID      int
}

type Channel struct {
	chatID int64
	li     int64
	writer *Writer
}

type Message struct {
	channel    *Channel
	rawBody    map[string]interface{}
	logID      int64
	logType    int
	Message    string
	id         int
	Author     int64
	attachment map[string]interface{}
	nickname   int64
}

type KakaoLogin struct {
	UserID               int    `json:"userId"`
	CountryIso           string `json:"countryIso"`
	CountryCode          string `json:"countryCode"`
	AccountID            int    `json:"accountId"`
	ServerTime           int    `json:"server_time"`
	AccessToken          string `json:"access_token"`
	RefreshToken         string `json:"refresh_token"`
	TokenType            string `json:"token_type"`
	AutoLoginAccountID   string `json:"autoLoginAccountId"`
	DisplayAccountID     string `json:"displayAccountId"`
	MainDeviceAgentName  string `json:"mainDeviceAgentName"`
	MainDeviceAppVersion string `json:"mainDeviceAppVersion"`
	Status               int    `json:"status"`
}

// big-endian
type struc_big_I struct {
	Data int `struc:"big,uint32"`
}

// little-endian
type struc_little_i struct {
	Data int `struc:"little,uint32"`
}

type struc_little_I struct {
	Data int `struc:"little,uint32"`
}

type struc_little_H struct {
	Data int `struc:"little,uint16"`
}

type struc_little_b struct {
	Data int `struc:"little,int8"`
}

type struc_little_B struct {
	Data int `struc:"little,uint8"`
}
