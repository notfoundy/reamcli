package player

import (
	"encoding/hex"
	"net/url"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/notfoundy/reamcli/internal/ani"
)

func Launch(ep *ani.Episode) {
	urls := generateUrls(ep)
	ok := false
	for _, u := range urls {
		if !ok {
			play(u)
			ok = true
		}
	}
}

func play(url string) error {
	binary := "mpv"
	if runtime.GOOS == "windows" {
		binary = "mpv.exe"
	}

	cmd := exec.Command(binary, url)
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func generateUrls(ep *ani.Episode) []string {
	var urls []string

	sort.Slice(ep.Sources, func(i, j int) bool {
		return ep.Sources[i].Priority < ep.Sources[j].Priority
	})

	for _, src := range ep.Sources {
		if src == nil || src.SourceUrl == "" {
			continue
		}

		var urlStr string
		if strings.HasPrefix(src.SourceUrl, "--") {
			decoded, err := decryptSourceURL(strings.TrimPrefix(src.SourceUrl, "--"))
			if err != nil {
				continue
			}
			urlStr = decoded
		} else {
			decoded, err := url.QueryUnescape(src.SourceUrl)
			if err != nil {
				urlStr = src.SourceUrl
			} else {
				urlStr = decoded
			}
		}

		urls = append(urls, urlStr)
	}

	return urls
}

func decryptSourceURL(encoded string) (string, error) {
	const key = "1234567890abcdef"
	keyLen := len(key)

	data, err := hex.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	decoded := make([]byte, len(data))
	for i := range data {
		decoded[i] = data[i] ^ key[i%keyLen]
	}

	return string(decoded), nil
}
