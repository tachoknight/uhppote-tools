package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// All commands *must* begin with this
const preamble = "17"

// This is the unique serial number of the board,
// typically the last four bytes of the MAC
// address
const boardSerialNum = ""

// The IP of the board
const boardIP = ""

// The port the board's listening on
const boardPort = ""

// Buffer size for return payload
const payloadBuffer = 2045

// The stuff at the beginning of the message that is common
// to all everything we send and receive
type prelude struct {
	preamble  string
	command   string
	buffer    string
	serialnum string
}

func getPrelude(payload string) prelude {
	var p prelude

	// The preamble of the command
	p.preamble = payload[0:2]
	// Now the command we sent
	p.command = payload[2:4]
	// And the buffer
	p.buffer = payload[4:8]
	// And the serial number
	p.serialnum = flipBytes(payload[8:16])

	return p
}

// Function for building the first part of every command we
// send to the board
func buildPrelude(commandVerb string) string {
	payload := preamble
	payload += commandVerb
	payload += pad("0", 4)
	payload += flipBytes(boardSerialNum)

	return payload
}

// Actually performs the transmission to the board, gets
// the reply, and returns the result to the caller
func transmit(payload string) string {
	decoded, _ := hex.DecodeString(payload)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, decoded)
	if err != nil {
		panic(err)
	}

	// Hello Central!
	conn, _ := net.Dial("udp", boardIP+":"+boardPort)
	conn.Write(buf.Bytes())

	message := make([]byte, payloadBuffer)
	conn.Read(message)

	data := hex.EncodeToString(message)
	return data
}

// Helper function for splitting a string into
// an array every n characters. Used in
// flipBytes()
func splitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

// We have to reverse the bytes because that's
// what the board wants. This function gets used
// a lot.
func flipBytes(byteString string) string {
	ba := splitSubN(byteString, 2)
	baSize := len(ba)

	outBytes := ""
	for i := (baSize - 1); i >= 0; i-- {
		outBytes += ba[i]
	}

	return outBytes
}

// We pad the payloads with zeros in several spots, so
// this is a simple helper function
func pad(padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen)
}

// Helper function as we convert from hex to decimal a lot
func hexToDec(hexVal string) int64 {
	decVal, err := strconv.ParseInt(flipBytes(hexVal), 16, 0)
	if err != nil {
		panic(err)
	}

	return decVal
}

// GetBoardTime queries the board for the current date and time
func GetBoardTime() string {
	const grcVerb = "32" // 0x32 is Get Time
	payload := buildPrelude(grcVerb)
	// Now we have to pad the end so we get 64 bytes
	payload += pad("0", 128-len(payload))

	result := transmit(payload)

	tsLayout := "20060102150405"
	ts, err := time.Parse(tsLayout, result[16:30])
	if err != nil {
		fmt.Println(err)
	}

	return ts.String()
}

// SetBoardTime sets the board's clock to the current
// date and time
func SetBoardTime() {
	const grcVerb = "30" // 0x30 is Set Time
	payload := buildPrelude(grcVerb)
	currentTime := time.Now()
	payload += currentTime.Format("20060102150405")
	// Now we have to pad the end so we get 64 bytes
	payload += pad("0", 128-len(payload))

	transmit(payload)
}

func main() {
	SetBoardTime()
	fmt.Println(GetBoardTime())
}
