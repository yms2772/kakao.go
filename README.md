# kakao.go
A Simple kakaotalk chatbot using LOCO protocal for Go

# Introduction
A package that makes the LOCO protocol compatible with Go.

*When using this package, it must not be abused or used illegally, and the developer is not responsible for any disadvantages caused by using it.*

# Bot Start
```Go
package main

import (
	"fmt"
	"log"

	"github.com/yms2772/kakaogo"
)

func main() {
	kakao, err := kakaogo.New("example@email.com", "pAsSwORd")
	// kakao, err := kakaogo.New("example@email.com", "pAsSwORd", true) // If you are already authenticated, put 'true' in the parameter.
	if err != nil {
		log.Fatalln(err)
	}

	kakao.OnReady = func() {
		log.Printf("Logged On")
	}

	kakao.OnMessage = func(chat *kakaogo.Message) {
		log.Printf("Message: %s", chat.Message)
		chat.Send(chat.Message)
		chat.SendPhoto("test.jpg")
	}

	kakao.Run()
	//kakao.RunWithKey("AccessKey~~~")  // If you have an access key, use this function.
}


```

# Thanks to
[jhleekr/kakao.py](https://github.com/jhleekr/kakao.py)

# License
MIT Licence
