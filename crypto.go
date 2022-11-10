package kakaogo

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"io"
	"math/big"
)

func (c *CryptoManager) getHandshakePacket() []byte {
	buf := &bytes.Buffer{}

	enced := c.rsaEncrypt()
	buf.Write(pack("<I", len(enced)))
	buf.Write(pack("<I", 12))
	buf.Write(pack("<I", 2))
	buf.Write(enced)

	return buf.Bytes()
}

func (c *CryptoManager) getRsaPublicKey() *rsa.PublicKey {
	n := "A44960441C7E83BB27898156ECB13C8AFAF05D284A4D1155F255CD22D3176CDE50482F2F27F71348E4D2EB5F57BF9671EF15C9224E042B1B567AC1066E06691143F6C50F88787F68CF42716B210CBEF0F59D53405A0A56138A6872212802BB0AEEA6376305DBD428831E8F61A232EFEDD8DBA377305EF972321E1352B5F64630993E5549C64FCB563CDC97DA2124B925DDEA12ADFD00138910F66937FAB68486AE43BFE203C4A617F9F232B5458A9AB409BAC8EDADEF685545F9B013986747737B3FD76A9BAC121516226981EA67225577D15D0F082B8207EAF7CDCB13123937CB12145837648C2F3A65018162315E77EAD2D2DD5986E46251764A43B9BA8F79"
	e := 3

	bigN := new(big.Int)
	bigN.SetString(n, 16)

	return &rsa.PublicKey{
		N: bigN,
		E: e,
	}
}

func (c *CryptoManager) rsaEncrypt() []byte {
	encrypted, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, c.getRsaPublicKey(), c.aesKey, nil)
	if err != nil {
		fmt.Println(err)
	}

	return encrypted
}

func (c *CryptoManager) aesEncrypt(data, iv []byte) []byte {
	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		fmt.Println(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext[aes.BlockSize:]
}

func (c *CryptoManager) aesDecrypt(data, iv []byte) []byte {
	decrypted := make([]byte, len(data))

	aesBlockDecrypt, err := aes.NewCipher(c.aesKey)
	if err != nil {
		fmt.Println(err)
	}

	aesDecrypt := cipher.NewCFBDecrypter(aesBlockDecrypt, iv)
	aesDecrypt.XORKeyStream(decrypted, data)

	return decrypted
}
