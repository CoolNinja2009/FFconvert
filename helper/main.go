package main

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
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
	f.WriteString(time.Now().Format("15:04:05") + " | " + text + "\n")
}

func mimeToExt(mime string) string {
	switch mime {
	// Images
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/webp":
		return "webp"
	case "image/avif":
		return "avif"
	case "image/gif":
		return "gif"
	case "image/bmp":
		return "bmp"
	case "image/tiff":
		return "tiff"

	// Video
	case "video/mp4":
		return "mp4"
	case "video/webm":
		return "webm"
	case "video/quicktime":
		return "mov"
	case "video/x-msvideo":
		return "avi"
	case "video/x-matroska":
		return "mkv"
	case "video/ogg":
		return "ogv"

	default:
		return ""
	}
}

func getType(ext string) string {
	ext = strings.ToLower(ext)

	image := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "webp": true,
		"avif": true, "gif": true, "bmp": true, "tiff": true,
	}

	video := map[string]bool{
		"mp4": true, "mkv": true, "webm": true,
		"mov": true, "avi": true, "ogv": true,
	}

	if image[ext] {
		return "image"
	}
	if video[ext] {
		return "video"
	}
	return "unknown"
}

func main() {
	logFile, _ := os.OpenFile("C:\\ffconvert\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	// Detect actual file type
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

	logLine(logFile, "Detected MIME: "+mtype.String())

	realExt := mimeToExt(mtype.String())
	if realExt == "" {
		logLine(logFile, "Unsupported file type")
		sendResponse(Response{"error", "Unsupported file type"})
		return
	}

	realType := getType(realExt)
	targetType := getType(targetExt)

	// Block cross-type conversions
	if realType != targetType && !req.Force {
		logLine(logFile, "Blocked cross-type conversion")
		sendResponse(Response{"error", "Illegal conversion"})
		return
	}

	// Skip if already correct
	if realExt == targetExt {
		logLine(logFile, "No conversion needed")
		sendResponse(Response{"ok", "No conversion needed"})
		return
	}

	base := strings.TrimSuffix(input, filepath.Ext(input))
	output := base + "." + targetExt

	var temp string
	if input == output {
		temp = base + "_tmp." + targetExt
	} else {
		temp = output
	}

	logLine(logFile, "Temp: "+temp)

	// Run FFmpeg
	cmd := exec.Command("ffmpeg", "-y", "-i", input, temp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logLine(logFile, "FFmpeg failed:")
		logLine(logFile, string(out))
		sendResponse(Response{"error", "FFmpeg failed"})
		return
	}

	logLine(logFile, "FFmpeg conversion successful")

	// Replace original if needed
	if input == output {
		logLine(logFile, "Starting replace sequence")

		success := false
		for i := 0; i < 20; i++ {
			os.Remove(input)
			err = os.Rename(temp, input)
			if err == nil {
				success = true
				logLine(logFile, "Rename succeeded")
				break
			}
			time.Sleep(200 * time.Millisecond)
		}

		if !success {
			logLine(logFile, "Final rename failed")
			sendResponse(Response{"error", "Rename failed"})
			return
		}
	}

	logLine(logFile, "Conversion complete")
	sendResponse(Response{"ok", "Conversion complete"})
}

func sendResponse(resp Response) {
	data, _ := json.Marshal(resp)
	length := uint32(len(data))
	binary.Write(os.Stdout, binary.LittleEndian, length)
	os.Stdout.Write(data)
}
