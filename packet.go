package kakaogo

import (
	"bytes"
	"crypto/rand"
	"log"
	"strings"

	bson2 "gopkg.in/mgo.v2/bson"
)

func (p *Packet) toJsonBody() (data map[string]interface{}) {
	data = make(map[string]interface{})
	if err := bson2.Unmarshal(p.body, &data); err != nil {
		log.Fatalf("toJsonBody: %s", err)
	}

	return data
}

func (p *Packet) readLocoPacket(packet []byte) {
	p.packetID = unpack("<I", packet[:4])
	p.statusCode = unpack("<H", packet[4:6])
	p.packetName = strings.ReplaceAll(string(packet[6:17]), `\0`, "")
	p.bodyType = unpack("<b", packet[17:18])
	p.bodySize = unpack("<i", packet[18:22])
	p.body = packet[22:]
}

func (p *Packet) toLocoPacket() []byte {
	buf := &bytes.Buffer{}
	buf.Write(pack("<I", p.packetID))
	buf.Write(pack("<H", p.statusCode))
	buf.Write([]byte(p.packetName))

	var null []byte
	for i := 0; i < 11-len(p.packetName); i++ {
		null = append(null, 0)
	}

	buf.Write(null)
	buf.Write(pack("<b", p.bodyType))
	buf.Write(pack("<i", len(p.body)))

	buf.Write(p.body)

	return buf.Bytes()
}

func (p *Packet) toEncryptedLocoPacket(cryptoManager *CryptoManager) []byte {
	iv := make([]byte, 16)
	_, _ = rand.Read(iv)

	encryptedPacket := cryptoManager.aesEncrypt(p.toLocoPacket(), iv)

	buf := &bytes.Buffer{}
	buf.Write(pack("<I", len(encryptedPacket)+len(iv)))
	buf.Write(iv)
	buf.Write(encryptedPacket)

	return buf.Bytes()
}

func (p *Packet) readEncryptedLocoPacket(packet []byte, cryptoManager *CryptoManager) {
	iv := packet[4:20]
	data := packet[20:]

	dec := cryptoManager.aesDecrypt(data, iv)

	p.readLocoPacket(dec)
}
