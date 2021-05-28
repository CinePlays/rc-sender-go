package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var protocols []Protocol;
var pin uint8;
var code int64;
var length int64;
var currProtocol Protocol;
var pulseLength int64;
var nRepeatTransmit int64;
var gpio rpio.Pin;

type Protocol struct {
	pulseLength int;

	syncFactor HighLow;
	zero HighLow;
	one HighLow;

	invertedSignal bool;
}

type HighLow struct {
	high int;
	low int;
}

func init() {
	protocols = []Protocol{
		{ pulseLength: 350, syncFactor: HighLow{ high: 1, low: 31 }, zero: HighLow{ high: 1, low: 3 }, one: HighLow{ high: 3, low: 1 }, invertedSignal: false },     // protocol 1
		{ pulseLength: 650, syncFactor: HighLow{ high: 1, low: 10 }, zero: HighLow{ high: 1, low: 2 }, one: HighLow{ high: 2, low: 1 }, invertedSignal: false },     // protocol 2
		{ pulseLength: 100, syncFactor: HighLow{ high: 30, low: 71 }, zero: HighLow{ high: 4, low: 11 }, one: HighLow{ high: 9, low: 6 }, invertedSignal: false },   // protocol 3
		{ pulseLength: 380, syncFactor: HighLow{ high: 1, low: 6 }, zero: HighLow{ high: 1, low: 3 }, one: HighLow{ high: 3, low: 1 }, invertedSignal: false },      // protocol 4
		{ pulseLength: 500, syncFactor: HighLow{ high: 6, low: 14 }, zero: HighLow{ high: 1, low: 2 }, one: HighLow{ high: 2, low: 1 }, invertedSignal: false },     // protocol 5
		{ pulseLength: 450, syncFactor: HighLow{ high: 23, low: 1 }, zero: HighLow{ high: 1, low: 2 }, one: HighLow{ high: 2, low: 1 }, invertedSignal: true },      // protocol 6 HT6P20B
		{ pulseLength: 150, syncFactor: HighLow{ high: 2, low: 62 }, zero: HighLow{ high: 1, low: 6 }, one: HighLow{ high: 6, low: 1 }, invertedSignal: false },     // protocol 7 HS2303-PT - AUKEY REMOTE
		{ pulseLength: 200, syncFactor: HighLow{ high: 3, low: 130 }, zero: HighLow{ high: 7, low: 16 }, one: HighLow{ high: 3, low: 16 }, invertedSignal: false },  // protocol 8 Conrad RS-200 RX
		{ pulseLength: 200, syncFactor: HighLow{ high: 130, low: 7 }, zero: HighLow{ high: 16, low: 7 }, one: HighLow{ high: 16, low: 3 }, invertedSignal: true },   // protocol 9 Conrad RS-200 TX
		{ pulseLength: 365, syncFactor: HighLow{ high: 18, low:1 }, zero: HighLow{ high: 3, low: 1 }, one: HighLow{ high: 1, low: 3 }, invertedSignal: true },      // protocol 10 1ByOne DOORBELL
		{ pulseLength: 270, syncFactor: HighLow{ high: 36, low: 1 }, zero: HighLow{ high: 1, low: 2 }, one: HighLow{ high: 2, low: 1 }, invertedSignal: true },      // protocol 11 HT12E
		{ pulseLength: 320, syncFactor: HighLow{ high: 36, low: 1 }, zero: HighLow{ high: 1, low: 2 }, one: HighLow{ high: 2, low: 1 }, invertedSignal: true },      // protocol 12 SM5212
	};
}

func main() {
	if len(os.Args) != 7 {
		fmt.Println("Syntax: <pin> <code> <length> <protocol> <pulse_length> <repeat_transmit>");
		os.Exit(0);
	}

	rawPin, err := strconv.ParseInt(os.Args[1], 10, 8);
	code, err = strconv.ParseInt(os.Args[2], 10, 64);
	length, err = strconv.ParseInt(os.Args[3], 10, 64);
	protocol, err := strconv.ParseInt(os.Args[4], 10, 64);
	pulseLength, err = strconv.ParseInt(os.Args[5], 10, 64);
	nRepeatTransmit, err = strconv.ParseInt(os.Args[6], 10, 64);

	if err != nil {
		fmt.Println("Arguments have to be numbers");
		os.Exit(0);
	}
	if rawPin < 0 || code < 0 || length < 0 || protocol < 0 || pulseLength < 0 || nRepeatTransmit < 0 {
		fmt.Println("Arguments have to be positive");
		os.Exit(0);
	}

	if protocol > int64(len(protocols)) {
		fmt.Println("Invalid Protocol");
		os.Exit(0);
	}

	pin = uint8(rawPin);
	currProtocol = protocols[protocol-1];

	if pulseLength > 1 {
		currProtocol.pulseLength = int(pulseLength);
	}

	if nRepeatTransmit < 1 {
		nRepeatTransmit = 10;
	}

	err = rpio.Open();
	if err != nil {
		fmt.Printf("Error accessing gpio mem: %v\n", err);
		os.Exit(0);
	}

	defer rpio.Close();

	gpio = rpio.Pin(pin);
	gpio.Output();

	fmt.Printf("Pin: %v\n", pin);
	fmt.Printf("Code: %v\n", code);
	fmt.Printf("Length: %v\n", length);
	fmt.Printf("Protocol: %v\n", protocol);
	fmt.Printf("Pulselength: %v\n", currProtocol.pulseLength);
	fmt.Printf("InvertedSignal: %v\n", currProtocol.invertedSignal);
	fmt.Printf("RepeatTransmit: %v\n", nRepeatTransmit);

	send();
}

func send() {
	for repeat := 0; repeat < int(nRepeatTransmit); repeat++ {
		for i := length-1; i >= 0; i-- {
			if (code & (1 << i)) > 0 {
				transmit(currProtocol.one);
			} else {
				transmit(currProtocol.zero);
			}
		}
		transmit(currProtocol.syncFactor);
	}
	gpio.Low()
}

func transmit(pulse HighLow) {
	firstSignal := rpio.High;
	secondSignal := rpio.Low;

	if currProtocol.invertedSignal {
		firstSignal = rpio.Low;
		secondSignal = rpio.High;
	}
	gpio.Write(firstSignal);
	time.Sleep(time.Duration(currProtocol.pulseLength * pulse.high) * time.Microsecond);
	gpio.Write(secondSignal);
	time.Sleep(time.Duration(currProtocol.pulseLength * pulse.low) * time.Microsecond);
}