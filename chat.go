package kakaogo

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
)

func (m *Message) init() {
	//chatLog := m.rawBody["chatLog"].(map[string]interface{})
	//fmt.Println(chatLog)
	//
	//if _, ok := chatLog["attachment"]; ok {
	//	m.attachment = chatLog["attachment"].(map[string]interface{})
	//	if err := json.Unmarshal(chatLog["attachment"].([]byte), &m.attachment); err != nil {
	//		fmt.Println(err)
	//	}
	//} else {
	//	m.attachment = map[string]interface{}{}
	//}

	m.nickname = m.Author
}

func (m *Message) Read() {
	m.init()
	m.channel.NotiRead(m.logID)
}

func (m *Message) Send(msg string) {
	m.init()
	m.channel.send(msg)
}

func (m *Message) Reply(msg string) {
	m.init()

	type data struct {
		AttachOnly bool     `json:"attach_only"`
		AttachType int      `json:"attach_type"`
		Mentions   []string `json:"mentions"`
		SrcLinkId  int64    `json:"src_linkId"`
		SrcMessage string   `json:"src_message"`
		SrcType    int      `json:"src_type"`
		SrcUserId  int64    `json:"src_userId"`
	}

	marshaled, _ := json.Marshal(data{
		AttachOnly: false,
		AttachType: 1,
		Mentions:   []string{},
		SrcLinkId:  m.channel.li,
		SrcMessage: m.Message,
		SrcType:    m.logType,
		SrcUserId:  m.Author,
	})

	m.channel.sendChat(msg, string(marshaled), 26)
}

func (m *Message) SendPhoto(file string) (err error) {
	m.init()

	type data struct {
		ThumbnailURL    string `json:"thumbnailUrl"`
		ThumbnailHeight int    `json:"thumbnailHeight"`
		ThumbnailWidth  int    `json:"thumbnailWidth"`
		URL             string `json:"url"`
		K               string `json:"k"`
		Cs              string `json:"cs"`
		S               int64  `json:"s"`
		W               int    `json:"w"`
		H               int    `json:"h"`
		Mt              string `json:"mt"`
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fstat, err := f.Stat()
	if err != nil {
		return err
	}

	_, key, urlStr, err := upload(file, "image/jpeg", m.Author)
	if err != nil {
		return err
	}

	sha1Data := sha1.New()
	if _, err := io.Copy(sha1Data, f); err != nil {
		return err
	}

	marshaled, _ := json.Marshal(data{
		ThumbnailURL:    urlStr,
		ThumbnailHeight: 500,
		ThumbnailWidth:  500,
		URL:             urlStr,
		K:               key,
		Cs:              hex.EncodeToString(sha1Data.Sum(nil)),
		S:               fstat.Size(),
		W:               500,
		H:               500,
		Mt:              "image/jpeg",
	})

	m.channel.forwardChat("", string(marshaled), 2)

	return nil
}

func (m *Message) Hide() {
	m.channel.hideMessage(m.logID, m.logType)
}

func (m *Message) Kick() {
	m.channel.kickMember(m.Author)
}
