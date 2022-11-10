package kakaogo

import "fmt"

func (w *Writer) getPacketID() int {
	w.packetID++
	return w.packetID
}

func (w *Writer) sendPacket(packet *Packet) {
	pid := w.getPacketID()

	packet.packetID = pid
	if _, err := w.stream.Write(packet.toEncryptedLocoPacket(w.cryptoManager)); err != nil {
		fmt.Println(err)
	}
}
