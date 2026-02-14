package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func convertDocument(input, targetExt, mime string, logFile *os.File) error {
	base := strings.TrimSuffix(input, filepath.Ext(input))
	output := base + "." + targetExt

	sameName := strings.EqualFold(input, output)

	var temp string
	if sameName {
		temp = base + "_tmp." + targetExt
	} else {
		temp = output
	}

	logLine(logFile, "Document temp: "+temp)

	ext := strings.ToLower(targetExt)

	// Build pandoc command
	var cmd *exec.Cmd

	// Determine real input format from MIME
	inputFormat := ""
	if mime == "application/pdf" {
		inputFormat = "pdf"
	} else if mime == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		inputFormat = "docx"
	} else if strings.HasPrefix(mime, "text/") {
		inputFormat = "plain"
	}

	switch ext {
	case "txt":
		if inputFormat != "" {
			cmd = exec.Command("pandoc", "-f", inputFormat, input, "-t", "plain", "-o", temp)
		} else {
			cmd = exec.Command("pandoc", input, "-t", "plain", "-o", temp)
		}

	case "pdf":
		if inputFormat != "" {
			cmd = exec.Command("pandoc", "-f", inputFormat, input, "-o", temp)
		} else {
			cmd = exec.Command("pandoc", input, "-o", temp)
		}

	case "html":
		if inputFormat != "" {
			cmd = exec.Command("pandoc", "-f", inputFormat, input, "-o", temp)
		} else {
			cmd = exec.Command("pandoc", input, "-o", temp)
		}

	default:
		logLine(logFile, "Unsupported document target: "+ext)
		return errors.New("unsupported document target")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		logLine(logFile, "Pandoc failed:"+err.Error())
		logLine(logFile, string(out))
		return errors.New("pandoc failed")
	}

	logLine(logFile, "Pandoc conversion successful")

	// Replace original if needed
	if sameName {
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
			return errors.New("rename failed")
		}
	}

	return nil
}
