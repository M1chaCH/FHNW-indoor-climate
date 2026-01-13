package main

import (
	"bufio"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
)

const (
	PortName = "/dev/tty.usbmodem2102"
	BaudRate = 115200
	FilePath = "."
)

var invalidDataLine = errors.New("received line has invalid format")

func main() {
	log.Printf("Attempting connection to serial port: '%v' with baud rate: %d \n", PortName, BaudRate)

	config := &serial.Config{Name: PortName, Baud: BaudRate}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to serial port, reading data...")

	defer port.Close()
	defer ForceCloseAllFiles()

	for {
		key, value, err := readValue(port)
		if errors.Is(err, invalidDataLine) {
			log.Println("WARN: got invalid line", key)
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		now := time.Now().UTC().Format(time.RFC3339)

		log.Printf("Read key value pair: key='%v',value='%v'\n", key, value)
		err = AppendToCsvFile(FilePath, key, value+","+now)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readValue(port *serial.Port) (string, string, error) {
	scanner := bufio.NewScanner(port)
	scanner.Split(bufio.ScanLines) // this is the default, but keep it here as a reminder that this could be changed.
	scanner.Scan()

	if scanner.Err() != nil {
		return "", "", scanner.Err()
	}

	text := strings.TrimSpace(scanner.Text())
	split := strings.Split(text, ":")

	if len(split) != 2 {
		return text, "", invalidDataLine
	}

	return split[0], split[1], nil
}
