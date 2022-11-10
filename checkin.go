package kakaogo

import (
	"crypto/rand"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	bson2 "gopkg.in/mgo.v2/bson"
)

func getCheckInData(host, port string) *Packet {
	cryptoManager := &CryptoManager{}
	cryptoManager.aesKey = make([]byte, 16)

	_, _ = rand.Read(cryptoManager.aesKey)

	sock, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("client: connect: %s", err)
	}
	defer sock.Close()

	_, err = sock.Write(cryptoManager.getHandshakePacket())
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

	encoded, err := bson2.Marshal(bson.M{
		"userId": 24951845,
		"os":     agent,
		"ntype":  ntype,
		"appVer": appVerSion,
		"MCCMNC": mccmnc,
		"lang":   "ko",
	})

	p := &Packet{
		packetID:   1,
		statusCode: 0,
		packetName: "CHECKIN",
		bodyType:   0,
		body:       encoded,
	}

	_, err = sock.Write(p.toEncryptedLocoPacket(cryptoManager))
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

	reply := make([]byte, 2048)
	_, err = sock.Read(reply)
	if err != nil {
		log.Fatalf("client: read: %s", err)
	}

	recvPacket := &Packet{}
	recvPacket.readEncryptedLocoPacket(reply, cryptoManager)

	return recvPacket
}
