package kakaogo

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	bson2 "gopkg.in/mgo.v2/bson"
)

func (c *Channel) sendPacket(command string, data bson.M) {
	encoded, _ := bson2.Marshal(data)

	packet := &Packet{
		packetID:   0,
		statusCode: 0,
		packetName: command,
		bodyType:   0,
		body:       encoded,
	}

	c.writer.sendPacket(packet)
}

func (c *Channel) NotiRead(watermark int64) {
	c.sendPacket("NOTIREAD", bson.M{
		"chatId":    c.chatID,
		"watermark": watermark,
	})
}

func (c *Channel) sendChat(msg string, extra string, t int) {
	c.sendPacket("WRITE", bson.M{
		"chatId": c.chatID,
		"extra":  extra,
		"type":   t,
		"msgId":  float64(time.Now().Unix()),
		"msg":    msg,
		"noSeen": false,
	})
}

func (c *Channel) send(msg string) {
	c.sendChat(msg, "{}", 1)
}

func (c *Channel) forwardChat(msg, extra string, t int) {
	c.sendPacket("FORWARD", bson.M{
		"chatId": c.chatID,
		"extra":  extra,
		"type":   t,
		"msgId":  float64(time.Now().Unix()),
		"msg":    msg,
		"noSeen": false,
	})
}

func (c *Channel) hideMessage(logID int64, logType int) {
	c.sendPacket("REWRITE", bson.M{
		"c":     c.chatID,
		"li":    c.li,
		"logId": logID,
		"t":     logType,
	})
}

func (c *Channel) kickMember(author int64) {
	c.sendPacket("REWRITE", bson.M{
		"c":   c.chatID,
		"li":  c.li,
		"mid": author,
	})
}
