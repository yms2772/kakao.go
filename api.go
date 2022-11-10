package kakaogo

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (k *Kakao) getXVC() string {
	sha := sha512.Sum512([]byte("HEATH|" + uthUserAgent + "|DEMIAN|" + k.email + "|" + deviceUUID))

	return fmt.Sprintf("%x", sha[:16])
}

func (k *Kakao) defaultData() *url.Values {
	httpData := &url.Values{}
	httpData.Add("email", k.email)
	httpData.Add("password", k.password)
	httpData.Add("device_name", deviceName)
	httpData.Add("device_uuid", deviceUUID)
	httpData.Add("os_version", osVersion)

	return httpData
}

func (k *Kakao) defaultHeader() http.Header {
	httpHeader := http.Request{
		Header: map[string][]string{},
	}
	httpHeader.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpHeader.Header.Add("A", authHeader)
	httpHeader.Header.Add("X-VC", k.getXVC())
	httpHeader.Header.Add("User-Agent", uthUserAgent)
	httpHeader.Header.Add("Accept", "*/*")
	httpHeader.Header.Add("Accept-Language", lang)

	return httpHeader.Header
}

func (k *Kakao) login() (data KakaoLogin, err error) {
	httpData := k.defaultData()
	httpData.Add("permanent", "true")
	httpData.Add("forced", "true")

	req, err := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(httpData.Encode()))
	if err != nil {
		return data, err
	}

	req.Header = k.defaultHeader()

	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return data, err
	}

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (k *Kakao) requestPasscode() (err error) {
	httpData := k.defaultData()
	httpData.Add("permanent", "true")
	httpData.Add("once", "false")

	req, err := http.NewRequest(http.MethodPost, requestPasscodeURL, strings.NewReader(httpData.Encode()))
	if err != nil {
		return err
	}

	req.Header = k.defaultHeader()

	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}

func (k *Kakao) registerDevice() (err error) {
	httpData := k.defaultData()
	httpData.Add("permanent", "true")
	httpData.Add("once", "false")
	httpData.Add("passcode", k.passcode)

	req, err := http.NewRequest(http.MethodPost, registerDeviceUrl, strings.NewReader(httpData.Encode()))
	if err != nil {
		return err
	}

	req.Header = k.defaultHeader()

	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}

func upload(file string, dataType string, userID int64) (path, key, urlStr string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return "", "", "", err
	}
	defer f.Close()

	httpData := &bytes.Buffer{}

	writer := multipart.NewWriter(httpData)
	_ = writer.WriteField("attachment_type", dataType)
	_ = writer.WriteField("user_id", strconv.FormatInt(userID, 10))
	part, _ := writer.CreateFormFile("attachment", filepath.Base(file))

	if _, err = io.Copy(part, f); err != nil {
		return "", "", "", err
	}

	if err = writer.Close(); err != nil {
		return "", "", "", err
	}

	req, err := http.NewRequest(http.MethodPost, mediaURL, httpData)
	if err != nil {
		return "", "", "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("A", authHeader)

	resp, err := httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", "", "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	path = string(body)
	key = strings.ReplaceAll(path, "/talkm", "")
	urlStr = "https://dn-m.talk.kakao.com" + path

	return path, key, urlStr, err
}
