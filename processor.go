package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getURIDirectory(URI string) (error, string) {
	URL, err := url.Parse(URI)
	if err != nil {
		return err, ""
	}
	return nil, fmt.Sprintf("%s%s", URL.Hostname(), strings.ToLower(strings.Replace(URL.Path, `/`, "-", -1)))
}

func newProcess(URI string) (*exec.Cmd, string) {
	err, dirPath := getURIDirectory(URI)
	if err != nil {
		fmt.Println("Erorr happened while getting directory name", dirPath)
		return nil, ""
	}

	newPath := filepath.Join("./videos", dirPath)
	if err = os.MkdirAll(newPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-fflags",
		"nobuffer",
		"-rtsp_transport",
		"tcp",
		"-i",
		URI,
		"-vsync",
		"0",
		"-copyts",
		"-vcodec",
		"copy",
		"-movflags",
		"frag_keyframe+empty_moov",
		"-an",
		"-hls_flags",
		"delete_segments+append_list",
		"-f",
		"segment",
		"-segment_list_flags",
		"live",
		"-segment_time",
		"1",
		"-segment_list_size",
		"3",
		"-segment_format",
		"mpegts",
		"-segment_list",
		fmt.Sprintf("%s/index.m3u8", newPath),
		"-segment_list_type",
		"m3u8",
		"-segment_list_entry_prefix",
		fmt.Sprintf("/stream/%s/", dirPath),
		newPath+"/%d.ts",
	)
	return cmd, newPath
}
