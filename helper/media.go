package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func detectCategory(mime string) string {
	if strings.HasPrefix(mime, "image/") ||
		strings.HasPrefix(mime, "video/") ||
		strings.HasPrefix(mime, "audio/") {
		return "media"
	}

	if strings.HasPrefix(mime, "application/") ||
		strings.HasPrefix(mime, "text/") {
		return "document"
	}

	return "unknown"
}

func convertMedia(input, targetExt, mime string, logFile *os.File) error {
	base := strings.TrimSuffix(input, filepath.Ext(input))
	output := base + "." + targetExt

	// Detect if final name equals original name
	sameName := strings.EqualFold(input, output)

	var temp string
	if sameName {
		temp = base + "_tmp." + targetExt
	} else {
		temp = output
	}

	logLine(logFile, "Media temp: "+temp)

	// Determine source type
	srcIsVideo := strings.HasPrefix(mime, "video/")
	srcIsAudio := strings.HasPrefix(mime, "audio/")

	// Common extension maps
	audioExts := map[string]bool{"mp3": true, "m4a": true, "aac": true, "wav": true, "flac": true, "ogg": true, "opus": true}
	videoExts := map[string]bool{"mp4": true, "mkv": true, "webm": true, "avi": true}

	runFF := func(args []string) ([]byte, error) {
		full := append([]string{"-y"}, args...)
		cmd := exec.Command("ffmpeg", full...)
		return cmd.CombinedOutput()
	}

	// Helper to choose audio encoder when copy fails
	chooseAudioArgs := func(ext string, input string, output string) []string {
		switch ext {
		case "mp3":
			return []string{"-i", input, "-vn", "-c:a", "libmp3lame", "-b:a", "192k", output}
		case "wav":
			return []string{"-i", input, "-vn", "-c:a", "pcm_s16le", output}
		case "flac":
			return []string{"-i", input, "-vn", "-c:a", "flac", output}
		case "ogg":
			return []string{"-i", input, "-vn", "-c:a", "libvorbis", "-b:a", "192k", output}
		case "opus":
			return []string{"-i", input, "-vn", "-c:a", "libopus", "-b:a", "128k", output}
		case "m4a", "aac":
			return []string{"-i", input, "-vn", "-c:a", "aac", "-b:a", "192k", output}
		default:
			return []string{"-i", input, output}
		}
	}

	var out []byte
	var err error

	// CASE: video -> audio
	if srcIsVideo && audioExts[strings.ToLower(targetExt)] {
		logLine(logFile, "Detected video -> audio conversion; attempting audio extract")

		ext := strings.ToLower(targetExt)

		// WAV must always be PCM
		if ext == "wav" {
			logLine(logFile, "WAV target detected; forcing PCM encode")
			args := []string{"-i", input, "-vn", "-c:a", "pcm_s16le", temp}
			out, err = runFF(args)
		} else {
			// Try stream copy first
			args := []string{"-i", input, "-vn", "-c:a", "copy", temp}
			out, err = runFF(args)
			if err != nil {
				logLine(logFile, "Stream copy audio failed, trying re-encode:")
				logLine(logFile, string(out))
				args = chooseAudioArgs(ext, input, temp)
				out, err = runFF(args)
			}
		}

	} else if srcIsAudio && videoExts[strings.ToLower(targetExt)] {
		// CASE: audio -> video container
		logLine(logFile, "Detected audio -> video conversion; creating audio-only container")

		// Try copy audio into container
		args := []string{"-i", input, "-vn", "-c:a", "copy", temp}
		out, err = runFF(args)
		if err != nil {
			logLine(logFile, "Stream copy into container failed, re-encoding audio:")
			logLine(logFile, string(out))

			ext := strings.ToLower(targetExt)
			if ext == "webm" {
				args = []string{"-i", input, "-vn", "-c:a", "libvorbis", "-b:a", "192k", temp}
			} else {
				args = []string{"-i", input, "-vn", "-c:a", "aac", "-b:a", "192k", temp}
			}
			out, err = runFF(args)
		}

	} else {
		// Default behavior: stream copy for media types
		isStreamCopy := srcIsVideo || srcIsAudio
		if isStreamCopy {
			out, err = runFF([]string{"-i", input, "-c", "copy", temp})
		} else {
			out, err = runFF([]string{"-i", input, temp})
		}
	}

	if err != nil {
		logLine(logFile, "FFmpeg failed:")
		logLine(logFile, string(out))
		return errors.New("FFmpeg failed")
	}

	logLine(logFile, "FFmpeg conversion successful")

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
			return errors.New("Rename failed")
		}
	}

	return nil
}
