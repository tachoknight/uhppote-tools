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

	conn.Close()

	data := hex.EncodeToString(message)
	return data
}

// Converts the scanned tag number into the format
// the board itself wants, which, shockingly, is
// not the same thing
func convertTagNum(tagSN int) (string, error) {
	bins := strconv.FormatInt(int64(tagSN), 2)
	bins = fmt.Sprintf("%024s", bins)

	frontb := bins[0:8]
	backb := bins[len(bins)-16:]

	f, err := strconv.ParseInt(frontb, 2, 32)
	if err != nil {
		return "", err
	}

	b, err := strconv.ParseInt(backb, 2, 32)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d%d", f, b), nil
}

// Need to twist the tag number into the format
// the board expects it
func formatTagNum(tagSN string) string {
	var fixedTagNum = ""
	tagNum, _ := strconv.Atoi(tagSN)
	base := strconv.FormatInt(int64(tagNum), 16)

	fixedTagNum = base
	c := len(fixedTagNum)

	for i := 8 - c; i > 0; i-- {
		fixedTagNum = "0" + fixedTagNum
	}

	fixedTagNum = flipBytes(fixedTagNum)

	return fixedTagNum
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

// AddUser adds a user's tag to the board
func AddUser(tagSN string) (bool, error) {
	success := false
	const auVerb = "50" // 0x50 is what the board uses as "Add User"

	payload := preamble
	payload += auVerb
	payload += pad("0", 4)
	payload += flipBytes(boardSerialNum)
	payload += formatTagNum(tagSN)

	// Now we're going to enable the user for today until a long
	// time from now
	currentTime := time.Now()
	payload += currentTime.Format("20060102")
	// Now lets add ten years
	futureTime := time.Now()
	futureTime = futureTime.AddDate(10, 0, 0)
	payload += futureTime.Format("20060102")

	// Turn the user on for all enabled systems (e.g. doors, machines...)
	payload += "01010101"

	// Now we have to pad the end so we get 64 bytes
	payload += pad("0", 128-len(payload))

	// Now we're gonna transmit and parse the result
	result := transmit(payload)

	// Parse the prelude in case we need it
	getPrelude(result)

	// And the code that determines whether we were successful
	// or not
	retSuccess := result[16:18]
	fmt.Println("Success: ", retSuccess)
	if retSuccess == "01" {
		success = true
	} else {
		success = false
	}

	return success, nil
}

// GetUser gets the information about a user's tag from the board
func GetUser(tagSN string) (bool, error) {
	success := false
	const guVerb = "5A" // 0x5A is what the board uses as "Get User"

	payload := preamble
	payload += guVerb
	payload += pad("0", 4)
	payload += flipBytes(boardSerialNum)
	payload += formatTagNum(tagSN)

	// Now we have to pad the end so we get 64 bytes
	payload += pad("0", 128-len(payload))

	result := transmit(payload)

	// Parse the prelude in case we need it
	getPrelude(result)

	// Now let's look at the info about this user

	// We need the hex value of the tag...
	retTagSN := result[16:24]
	// ... but we also want the number as we understand it
	retTagNum := hexToDec(retTagSN)
	fmt.Println("Num: ", retTagNum)
	// And we also get the start and end dates
	retStartDate := result[24:32]
	retEndDate := result[32:40]
	fmt.Println("Start: ", retStartDate, " End: ", retEndDate)
	// Now we find out what enabled systems (e.g. doors, machines...) the user
	// has access to
	retSystems := result[40:48]
	fmt.Println("Enabled systems: ", retSystems)

	// Success if the tag isn't 0x0 or 0xf
	if retTagSN != "00000000" && retTagSN != "ffffffff" {
		success = true
	}

	return success, nil
}

// DelUser deletes a user's tag from the board
func DelUser(tagSN string) (bool, error) {
	success := false
	const duVerb = "52" // 0x52 is what the board uses as "Delete User"

	payload := preamble
	payload += duVerb
	payload += pad("0", 4)
	payload += flipBytes(boardSerialNum)
	payload += formatTagNum(tagSN)

	// Now we have to pad the end so we get 64 bytes
	payload += pad("0", 128-len(payload))

	result := transmit(payload)

	// Parse the prelude in case we need it
	getPrelude(result)

	// And the code that determines whether we were successful
	// or not
	retSuccess := result[16:18]
	fmt.Println("Success: ", retSuccess)
	if retSuccess == "01" {
		success = true
	} else {
		success = false
	}

	return success, nil
}

func main() {

	// This is the tag we want to work with, as reported by
	// the RFID reader
	origTagNum := 10978235
	fixedTag, err := convertTagNum(origTagNum)
	if err != nil {
		fmt.Println("Couldn't convert tag: ", origTagNum, err)
	}

	fmt.Printf("Going to work with tag %s\n", fixedTag)

	// Now add the user's tag to the board
	AddUser(fixedTag)
	// Verify the user's tag is on the board
	GetUser(fixedTag)
	// And delete the user's tag from the board
	DelUser(fixedTag)
}
