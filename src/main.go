package main

// #cgo LDFLAGS: -lpam
// #include <security/pam_modules.h>
// #include <security/pam_appl.h>
import "C"

import (
	"fmt"
	"github.com/google/uuid"
	"rsc.io/qr"
	"unsafe"
)

type eezeMessage struct {
	Typ   string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
	Token *uuid.UUID  `json:"token,omitempty"`
	Url   string      `json:"url,omitempty"`
}

//export goAuthenticate
func goAuthenticate(handle *C.pam_handle_t, flags C.int, argv []string) C.int {
	hdl := Handle{unsafe.Pointer(handle)}
	fmt.Println("authenticate:", argv)

	usr, err := hdl.GetUser()
	if err != nil {
		return C.PAM_AUTH_ERR
	}

	fmt.Println("User:", usr)
	if usr != "tester" {
		return C.PAM_USER_UNKNOWN
	}

	_, err = hdl.Conversation(
		Message{
			Style: MessageTextInfo,
			Msg:   "Made it this far",
		})

	//c, _, err := websocket.DefaultDialer.Dial("wss://eeze.io/api/v1/did-auth/ws?clientId=ffaa8b2d-1f7a-4297-8fc0-89ac4743639b", nil)
	//if err != nil {
	//	log.Fatal("dial:", err)
	//}
	//defer c.Close()
	//
	//for {
	//	_, message, err := c.ReadMessage()
	//	if err != nil {
	//		log.Println("read:", err)
	//		return C.PAM_AUTH_ERR
	//	}
	//	msg := eezeMessage{}
	//	json.Unmarshal(message, &msg)
	//
	//	switch msg.Typ {
	//	case "did-auth":
	//		qrcode, err := encodeQrCode(string(message))
	//		if err != nil {
	//			return C.PAM_AUTH_ERR
	//		}
	//
	//		_, err = hdl.Conversation(
	//			Message{
	//				Style: MessageTextInfo,
	//				Msg:   qrcode,
	//			})
	//
	//		break
	//	}
	//}
	//
	//resps, err := hdl.Conversation(
	//	Message{
	//		Style: MessageEchoOff,
	//		Msg:   "Password: ",
	//	},
	//)
	//
	//if err != nil {
	//	fmt.Println("Error: ", err)
	//	return C.PAM_CONV_ERR
	//}
	//
	//if resps[0] != "cake" {
	//	return C.PAM_AUTH_ERR
	//}
	//
	//resps, err = hdl.Conversation(
	//	Message{
	//		Style: MessageEchoOn,
	//		Msg:   "Favourite colour: ",
	//	})
	//
	//if err != nil {
	//	fmt.Println("Error: ", err)
	//	return C.PAM_CONV_ERR
	//}
	//
	//fmt.Println("I can't believe you like the colour", resps[0])
	//hdl.SetModuleData("fav-colour", resps[0])

	return C.PAM_SUCCESS
}

const ANSI_RESET = "\x1B[0m"
const ANSI_BLACKONGREY = "\x1B[30;47;27m"
const ANSI_WHITE = "\x1B[27m"
const ANSI_BLACK = "\x1B[7m"
const UTF8_BOTH = "\xE2\x96\x88"
const UTF8_TOPHALF = "\xE2\x96\x80"
const UTF8_BOTTOMHALF = "\xE2\x96\x84"

// Display QR code visually. If not possible, return 0.
func encodeQrCode(val string) (string, error) {

	code, err := qr.Encode(val, qr.M)
	if err != nil {
		return "", err
	}

	// Output QRCode using ANSI colors. Instead of black on white, we
	// output black on grey, as that works independently of whether the
	// user runs their terminal in a black on white or white on black color
	// scheme.
	// But this requires that we print a border around the entire QR Code.
	// Otherwise readers won't be able to recognize it.

	qrcode := ""
	counter := 0

	for i := 0; i < 2; i++ {
		qrcode += ANSI_BLACKONGREY
		for x := 0; x < code.Size+4; x++ {
			qrcode += "  "
		}
		qrcode += ANSI_RESET
		qrcode += "\n"
	}
	for y := 0; y < code.Size; y++ {
		qrcode += ANSI_BLACKONGREY + "    "
		isBlack := 0
		for x := 0; x < code.Size; x++ {
			if code.Bitmap[counter] == 0 {
				if isBlack == 0 {
					qrcode += ANSI_BLACK
				}
				isBlack = 1
			} else {
				if isBlack > 0 {
					qrcode += ANSI_WHITE
				}
				isBlack = 0
			}
			qrcode += "  "
			counter++
		}
		if isBlack > 0 {
			qrcode += ANSI_WHITE
		}
		qrcode += "    "
		qrcode += ANSI_RESET
		qrcode += "\n"
	}
	for i := 0; i < 2; i++ {
		qrcode += ANSI_BLACKONGREY
		for x := 0; x < code.Size+4; x++ {
			qrcode += "  "
		}
		qrcode += ANSI_RESET
		qrcode += "\n"
	}

	return qrcode, nil
}

//export setCred
func setCred(handle *C.pam_handle_t, flags C.int, argv []string) C.int {
	fmt.Println("setcred: ", argv)
	return C.PAM_SUCCESS
}

// Does Nothing?
func main() {}
