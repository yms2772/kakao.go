package kakaogo

import (
	"crypto/tls"

	"go.mongodb.org/mongo-driver/bson"
	bson2 "gopkg.in/mgo.v2/bson"
)

func getBookingData() *Packet {
	hostname := "booking-loco.kakao.com"

	conn, err := tls.Dial("tcp", hostname+":443", &tls.Config{})
	if err != nil {
		logger.Printf("tls.Dial Err: %+v", err)
	}
	defer conn.Close()

	encoded, err := bson2.Marshal(bson.M{
		"os":     "win32",
		"model":  "",
		"MCCMNC": "",
	})

	b := Packet{
		packetID:   1000,
		statusCode: 0,
		packetName: "GETCONF",
		bodyType:   0,
		body:       encoded,
	}

	message := b.toLocoPacket()
	_, err = conn.Write(message)
	if err != nil {
		logger.Fatalf("client: write: %s", err)
	}

	reply := make([]byte, 4096)
	_, err = conn.Read(reply)
	if err != nil {
		logger.Fatalf("client: write: %s", err)
	}

	recvPacket := &Packet{}
	recvPacket.readLocoPacket(reply)

	return recvPacket
}
