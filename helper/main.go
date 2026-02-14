package main

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

type Request struct {
	Input     string `json:"input"`
	TargetExt string `json:"target_ext"`
	Force     bool   `json:"force"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func logLine(f *os.File, text string) {
	filename := f.Name()
	entry := time.Now().Format("15:04:05") + " | " + text + "\n"

	old, err := os.ReadFile(filename)
	if err != nil {
		// If we can't read (file may not exist yet), just write the entry
		_ = os.WriteFile(filename, []byte(entry), 0644)
		return
	}

	// Prepend new entry
	newContent := append([]byte(entry), old...)
	_ = os.WriteFile(filename, newContent, 0644)
}

func sendResponse(resp Response) {
	data, _ := json.Marshal(resp)
	length := uint32(len(data))
	binary.Write(os.Stdout, binary.LittleEndian, length)
	os.Stdout.Write(data)
}

func main() {
	logFile, _ := os.OpenFile("C:\\ffconvert\\debug.log", os.O_CREATE|os.O_RDWR, 0644)
	defer logFile.Close()

	logLine(logFile, "Helper started")

	// Read message length
	var length uint32
	err := binary.Read(os.Stdin, binary.LittleEndian, &length)
	if err != nil {
		logLine(logFile, "Failed to read message length")
		return
	}

	// Read JSON message
	msg := make([]byte, length)
	_, err = os.Stdin.Read(msg)
	if err != nil {
		logLine(logFile, "Failed to read message body")
		return
	}

	var req Request
	err = json.Unmarshal(msg, &req)
	if err != nil {
		logLine(logFile, "Failed to parse JSON")
		return
	}

	input := req.Input
	targetExt := strings.ToLower(req.TargetExt)

	logLine(logFile, "Input: "+input)
	logLine(logFile, "Target: "+targetExt)

	// Detect MIME
	file, err := os.Open(input)
	if err != nil {
		logLine(logFile, "Failed to open file")
		sendResponse(Response{"error", "Failed to open file"})
		return
	}

	mtype, err := mimetype.DetectReader(file)
	file.Close()
	if err != nil {
		logLine(logFile, "MIME detection failed")
		sendResponse(Response{"error", "MIME detection failed"})
		return
	}

	mime := mtype.String()
	logLine(logFile, "Detected MIME: "+mime)

	category := detectCategory(mime)

	if category == "media" {
		err = convertMedia(input, targetExt, mime, logFile)
	} else if category == "document" {
		err = convertDocument(input, targetExt, mime, logFile)
	} else {
		logLine(logFile, "Unsupported file type")
		sendResponse(Response{"error", "Unsupported file type"})
		return
	}

	if err != nil {
		sendResponse(Response{"error", err.Error()})
		return
	}

	logLine(logFile, "Conversion complete")
	sendResponse(Response{"ok", "Conversion complete"})
}
