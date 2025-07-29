package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

type Stream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type FFProbeOutput struct {
	Streams []Stream `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %v", err)
	}

	var probeOutput FFProbeOutput
	if err := json.Unmarshal(stdout.Bytes(), &probeOutput); err != nil {
		return "", fmt.Errorf("could not parse ffprobe output: %v", err)
	}

	if len(probeOutput.Streams) == 0 {
		return "", errors.New("no video streams found")
	}

	width := probeOutput.Streams[0].Width
	height := probeOutput.Streams[0].Height

	if width == 16*height/9 {
		return "16:9", nil
	} else if height == 16*width/9 {
		return "9:16", nil
	}
	return "other", nil

}
