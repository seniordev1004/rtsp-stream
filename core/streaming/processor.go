package streaming

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Roverr/rtsp-stream/core/config"
	"github.com/sirupsen/logrus"
)

// GetURIDirectory is a function to create a directory string from an URI
func GetURIDirectory(URI string) (string, error) {
	URL, err := url.Parse(URI)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", URL.Hostname(), strings.ToLower(strings.Replace(URL.Path, `/`, "-", -1))), nil
}

// NewProcess creates a new transcoding process for ffmpeg
func NewProcess(URI string, spec *config.Specification) (*exec.Cmd, string, string) {
	dirPath, err := GetURIDirectory(URI)
	if err != nil {
		logrus.Error("Error happened while getting directory name", dirPath)
		return nil, "", ""
	}

	newPath := filepath.Join(spec.StoreDir, dirPath)
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
	return cmd, filepath.Join("stream", dirPath), fmt.Sprintf("%s/index.m3u8", newPath)
}
