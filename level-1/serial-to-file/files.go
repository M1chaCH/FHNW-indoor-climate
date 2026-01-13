package main

import (
	"bufio"
	"log"
	"os"
)

var writersCache = make(map[string]*bufio.Writer)
var openedFiles = make([]os.File, 0)

func AppendToCsvFile(path string, filename string, line string) error {
	filePath := path + "/" + filename + ".csv"

	writer, err := getCsvWriter(filePath)
	if err != nil {
		return err
	}

	_, err = writer.WriteString(line + "\n")
	if err != nil {
		return err
	}

	// writes to file immediately
	// doing this less often would improve performance but increase the risk of data loss
	err = writer.Flush()
	return err
}

func getCsvWriter(name string) (*bufio.Writer, error) {
	writer := writersCache[name]
	if writer == nil {
		log.Println("Opening writer for file: ", name)
		f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writersCache[name] = bufio.NewWriter(f)
		openedFiles = append(openedFiles, *f)
	}

	return writersCache[name], nil
}

func ForceCloseAllFiles() {
	for _, f := range openedFiles {
		f.Close()
	}
}
