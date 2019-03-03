package M3U

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type M3UInf struct {
	Last float64
	File string
}
type EXTM3U struct {
	Version        int
	AllowCache     bool
	MediaSequence  int
	TargetDuration int
	Inf            []M3UInf
	LenInf         int
}

func ParseM3u(data string) (*EXTM3U, error) {
	var m3u EXTM3U
	sReader := strings.NewReader(data)
	bReader := bufio.NewScanner(sReader)
	bReader.Split(bufio.ScanLines)
	if bReader.Scan() == false {
		return nil, errors.New("Empty m3u file.")
	}
	flag := bReader.Text()
	if flag != "#EXTM3U" {
		return nil, errors.New(fmt.Sprintf("Tag %s is not a regular m3u file.", flag))
	}
	var err error
	for bReader.Scan() {
		line := bReader.Text()
		s := strings.Index(line, ":")
		switch line[:s] {
		case "#EXT-X-VERSION":
			m3u.Version, err = strconv.Atoi(line[s+1:])
		case "#EXT-X-ALLOW-CACHE":
			m3u.AllowCache = (line[s+1:] != "NO")
		case "#EXT-X-MEDIA-SEQUENCE":
			m3u.MediaSequence, err = strconv.Atoi(line[s+1:])
		case "#EXT-X-TARGETDURATION":
			m3u.TargetDuration, err = strconv.Atoi(line[s+1:])
		case "#EXTINF":
			var lst float64
			lst, err = strconv.ParseFloat(line[s+1:len(line)-1], 32)
			if bReader.Scan() == false {
				return nil, errors.New("Inf tag is not complete. This is not a regular m3u file.")
			}
			m3u.Inf = append(m3u.Inf, M3UInf{Last: lst, File: bReader.Text()})
			m3u.LenInf++
		}
		if err != nil {
			return nil, err
		}
	}
	return &m3u, nil
}
