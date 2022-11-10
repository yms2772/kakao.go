package kakaogo

import (
	"fmt"
	"time"
)

func New(email, password string, skip ...bool) (kakao *Kakao, err error) {
	kakao = &Kakao{
		email:    email,
		password: password,

		OnReady:   func() {},
		OnJoin:    func(packet *Packet, channel *Channel) {},
		OnQuit:    func(packet *Packet, channel *Channel) {},
		OnMessage: func(chat *Message) {},
		OnRead:    func(channel *Channel, body map[string]interface{}) {},
	}

	if len(skip) == 0 {
		var checkReg string
		fmt.Print("If you have not authenticated, please press 'y' within 10 seconds: ")

		ch := make(chan int)

		go func() {
			_, _ = fmt.Scanln(&checkReg)
			ch <- 1
		}()

		select {
		case <-ch:
			if checkReg == "y" {
				if err := kakao.requestPasscode(); err != nil {
					return nil, err
				}

				fmt.Print("Check the KakaoTalk verification code with your mobile phone and enter it: ")
				_, _ = fmt.Scanln(&kakao.passcode)

				if err := kakao.registerDevice(); err != nil {
					return nil, err
				}

				fmt.Println("Verification completed!")
			} else {
				fmt.Println("Verification passed.")
			}
		case <-time.After(10 * time.Second):
			fmt.Println("Timed out, continue.")
		}
	} else {
		fmt.Println("Verification skipped.")
	}

	return kakao, nil
}
