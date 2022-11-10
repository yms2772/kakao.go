package kakaogo

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/lunixbochs/struc"
	"go.mongodb.org/mongo-driver/bson"
	bson2 "gopkg.in/mgo.v2/bson"
)

func pack(size string, data int) []byte {
	var buf bytes.Buffer
	var err error
	switch size {
	// big-endian
	case ">I":
		err = struc.Pack(&buf, &struc_big_I{Data: data})

	// little-endian
	case "<H":
		err = struc.Pack(&buf, &struc_little_H{Data: data})
	case "<i":
		err = struc.Pack(&buf, &struc_little_i{Data: data})
	case "<I":
		err = struc.Pack(&buf, &struc_little_I{Data: data})
	case "<b":
		err = struc.Pack(&buf, &struc_little_b{Data: data})
	case "<B":
		err = struc.Pack(&buf, &struc_little_B{Data: data})
	}

	if err != nil {
		logger.Fatalf("%s Pack Err: %+v", size, err)
		return nil
	}

	return buf.Bytes()
}

func unpack(size string, data []byte) int {
	buf := bytes.NewBuffer(data)

	var err error
	var result int
	switch size {
	// big-endian
	case ">I":
		item := &struc_big_I{}
		err = struc.Unpack(buf, item)
		result = item.Data

	// little-endian
	case "<H":
		item := &struc_little_H{}
		err = struc.Unpack(buf, item)
		result = item.Data
	case "<i":
		item := &struc_little_i{}
		err = struc.Unpack(buf, item)
		result = item.Data
	case "<I":
		item := &struc_little_I{}
		err = struc.Unpack(buf, item)
		result = item.Data
	case "<b":
		item := &struc_little_b{}
		err = struc.Unpack(buf, item)
		result = item.Data
	case "<B":
		item := &struc_little_B{}
		err = struc.Unpack(buf, item)
		result = item.Data
	}

	if err != nil {
		logger.Fatalf("%s Unpack Err: %+v", size, err)
		return 0
	}

	return result
}

func (k *Kakao) recvPacket() {
	defer func() {
		logger.Printf("Closing tcp connection")
		k.stream.Close()
	}()

	var encryptedBuffer []byte
	recv := make([]byte, 1<<8)
	currentPacketSize := 0

	console := bufio.NewReader(k.stream)

	for {
		n, err := console.Read(recv)
		if err != nil {
			if io.EOF == err {
				logger.Printf("connection is closed from client: %v | %v", k.stream.RemoteAddr().String(), err)
				return
			}

			logger.Printf("fail to receive data; err: %v", err)
			return
		}

		logger.Printf("packet received: %d bytes", n)
		message := recv[:n]
		encryptedBuffer = append(encryptedBuffer, message...)

		if currentPacketSize == 0 && len(encryptedBuffer) >= 4 {
			currentPacketSize = unpack("<I", encryptedBuffer[0:4])
		}

		if currentPacketSize != 0 {
			encryptedPacketSize := currentPacketSize + 4

			if len(encryptedBuffer) >= encryptedPacketSize {
				go k.processingPacket(encryptedBuffer[0:encryptedPacketSize])

				encryptedBuffer = encryptedBuffer[encryptedPacketSize:]
				currentPacketSize = 0
			}
		}
	}
}

func (k *Kakao) processingPacket(encryptedPacket []byte) {
	iv := encryptedPacket[4:20]
	body := encryptedPacket[20:]

	k.processingBuffer = append(k.processingBuffer, k.cryptoManager.aesDecrypt(body, iv)...)

	if len(k.processingHeader) == 0 && len(k.processingBuffer) >= 22 {
		k.processingHeader = k.processingBuffer[0:22]
		k.processingSize = unpack("<i", k.processingHeader[18:22]) + 22
	}

	if len(k.processingHeader) != 0 {
		if len(k.processingBuffer) >= k.processingSize {
			p := &Packet{}
			p.readLocoPacket(k.processingBuffer[:k.processingSize])

			go k.onPacket(p)

			k.processingBuffer = k.processingBuffer[k.processingSize:]
			k.processingHeader = []byte{}
		}
	}
}

func (k *Kakao) onPacket(packet *Packet) {
	body := packet.toJsonBody()

	switch strings.Trim(packet.packetName, "\x00") {
	case "MSG":
		chatLog := body["chatLog"].(map[string]interface{})
		chatID := chatLog["chatId"].(int64)

		var li int64
		if _, ok := body["li"]; ok {
			li = body["li"].(int64)
		} else {
			li = 0
		}

		channel := &Channel{
			chatID: chatID,
			li:     li,
			writer: k.writer,
		}

		message := &Message{
			channel: channel,
			rawBody: body,
			logID:   chatLog["logId"].(int64),
			logType: chatLog["type"].(int),
			Message: chatLog["message"].(string),
			id:      chatLog["msgId"].(int),
			Author:  chatLog["authorId"].(int64),
		}

		k.OnMessage(message)
	case "NEWMEM":
		chatLog := body["chatLog"].(map[string]interface{})
		chatID := chatLog["chatId"].(int64)

		var li int64
		if _, ok := body["li"]; ok {
			li = body["li"].(int64)
		} else {
			li = 0
		}

		channel := &Channel{
			chatID: chatID,
			li:     li,
			writer: k.writer,
		}

		k.OnJoin(packet, channel)
	case "DELMEM":
		chatLog := body["chatLog"].(map[string]interface{})
		chatID := chatLog["chatId"].(int64)

		var li int64
		if _, ok := body["li"]; ok {
			li = body["li"].(int64)
		} else {
			li = 0
		}

		channel := &Channel{
			chatID: chatID,
			li:     li,
			writer: k.writer,
		}

		k.OnQuit(packet, channel)
	case "DECUNREAD":
		chatID := body["chatId"].(int64)

		channel := &Channel{
			chatID: chatID,
			li:     0,
			writer: k.writer,
		}

		k.OnRead(channel, body)
	}
}

func (k *Kakao) __login(key ...string) (forever chan error) {
	forever = make(chan error)

	if len(key) == 0 {
		loginData, _ := k.login()

		var loginMsg string
		switch loginData.Status {
		case -101:
			loginMsg = "카카오톡에 로그인 되어있는 PC에서 로그아웃 해주세요"
		case -100:
			loginMsg = "인증이 되어있지 않습니다"
		case 12, 30:
			loginMsg = "이메일 또는 비밀번호를 확인해주세요"
		}

		if loginData.Status != 0 {
			forever <- errors.New(loginMsg)
			return
		}

		k.accessKey = loginData.AccessToken
	} else {
		k.accessKey = key[0]
	}

	logger.Printf("Access key: %s", k.accessKey)

	bookingData := getBookingData().toJsonBody()

	ticket := bookingData["ticket"].(map[string]interface{})
	wifi := bookingData["wifi"].(map[string]interface{})

	lsl := ticket["lsl"].([]interface{})
	ports := wifi["ports"].([]interface{})

	if len(lsl) == 0 || len(ports) == 0 {
		forever <- errors.New("error")
	}

	checkInData := getCheckInData(lsl[0].(string), strconv.Itoa(ports[0].(int))).toJsonBody()

	var err error
	k.stream, err = net.Dial("tcp", checkInData["host"].(string)+":"+strconv.Itoa(checkInData["port"].(int)))
	if err != nil {
		forever <- err
		return
	}

	k.cryptoManager = &CryptoManager{}
	k.cryptoManager.aesKey = make([]byte, 16)

	_, _ = rand.Read(k.cryptoManager.aesKey)

	k.writer = &Writer{
		cryptoManager: k.cryptoManager,
		stream:        k.stream,
		packetID:      0,
	}

	encoded, _ := bson2.Marshal(bson.M{
		"appVer":      appVerSion,
		"prtVer":      prtVersion,
		"os":          agent,
		"lang":        lang,
		"duuid":       deviceUUID,
		"oauthToken":  k.accessKey,
		"dtype":       dtype,
		"ntype":       ntype,
		"MCCMNC":      mccmnc,
		"revision":    0,
		"chatIds":     []string{},
		"maxIds":      []string{},
		"lastTokenId": 0,
		"lbk":         0,
		"bg":          false,
	})

	loginListPacket := &Packet{
		packetID:   0,
		statusCode: 0,
		packetName: "LOGINLIST",
		bodyType:   0,
		body:       encoded,
	}

	_, err = k.stream.Write(k.cryptoManager.getHandshakePacket())
	if err != nil {
		forever <- err
		return
	}

	go k.writer.sendPacket(loginListPacket)
	go k.recvPacket()
	go k.heartbeat()
	go k.OnReady()

	<-forever

	return
}

func (k *Kakao) heartbeat() {
	ticker := time.NewTicker(180 * time.Second)

	for range ticker.C {
		pingPacket := &Packet{
			packetID:   0,
			statusCode: 0,
			packetName: "PING",
			bodyType:   0,
			body:       nil,
		}

		k.writer.sendPacket(pingPacket)
	}
}

func (k *Kakao) Run() {
	fmt.Println(<-k.__login())
}

func (k *Kakao) RunWithKey(key string) {
	fmt.Println(<-k.__login(key))
}
