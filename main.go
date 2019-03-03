package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/IceSandwich/Bilivedown/Logging"
	"github.com/IceSandwich/Bilivedown/M3U"
	"github.com/IceSandwich/Bilivedown/Net"
)

type BiliveInfo struct {
	Code    int
	Msg     string
	Message string
	Data    struct {
		CurrentQuality int
		AcceptQuality  []string
		Durl           []struct {
			Url        string
			Length     int
			Order      int
			StreamType int
		}
	}
}

const INF = int(^uint(0) >> 1)

var ( //Global settings
	PageId     = 0
	DataBase   = "Sequences"
	AutoRecord = false
	ATime      = 0
	ARunFail   = false
	ARunC      = ""
	ARunA      = ""
	ThreadNum  = 5
	SelDurl    = 0
	RetryTimes = 3
	Frequency  = 4000
	TooFreq    = 1000
	ReConn     = 1000
	ReTrd      = 1000
	ThreadPool []bool
)

func ReadSetting() {
	fh, err := os.Open("Setting.ini")
	if err != nil {
		return //Use default settings
	}
	defer fh.Close()
	sScan := bufio.NewScanner(fh)
	for sScan.Scan() {
		sstr := strings.TrimSpace(sScan.Text())
		if sstr != "" && sstr[0] != '#' {
			s := strings.IndexByte(sstr, '=')
			if s == -1 {
				continue
			}
			sk, sv := strings.TrimSpace(sstr[:s]), strings.TrimSpace(sstr[s+1:])
			switch sk {
			case "Roomid":
				PageId, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "DataBase":
				DataBase = sv
				Logging.ReadArg(sk, sv)
			case "AutoRecord":
				AutoRecord = (strings.ToLower(sv) == "true")
				Logging.ReadArg(sk, sv)
			case "AutoRecord_RunFail":
				ARunFail = (strings.ToLower(sv) == "true")
				Logging.ReadArg(sk, sv)
			case "AutoRecord_Time":
				ATime, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "AutoRecord_Run":
				extract2 := func(s string, t *string) int {
					var x int
					if s[0] == '"' {
						if x = strings.IndexByte(s[1:], '"'); x != -1 {
							*t = s[1 : x+1]
							return x + 1 + 1
						}
					} else if x = strings.IndexByte(s, ' '); x != -1 {
						*t = s[:x]
						return x + 1
					} else {
						*t = s
					}
					return x
				}
				if x := extract2(sv, &ARunA); x != -1 && x != len(sv) {
					extract2(strings.TrimSpace(sv[x:]), &ARunC)
					Logging.ReadArg(sk, sv)
				}
			case "MaxThread":
				ThreadNum, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "Retry":
				RetryTimes, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "Timeout":
				if t, err := strconv.Atoi(sv); err == nil {
					Net.SetTimeout(t)
					Logging.ReadArg(sk, sv)
				}
			case "UserAgent":
				if sv[0] == '"' && sv[len(sv)-1] == '"' {
					Net.UserAgent = sv
					Logging.ReadArg(sk, sv)
				}
			case "Durl":
				SelDurl, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "Frequency":
				Frequency, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "TooFreq":
				TooFreq, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "ReConn":
				ReConn, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			case "ReTrd":
				ReTrd, _ = strconv.Atoi(sv)
				Logging.ReadArg(sk, sv)
			}
		}
	}
}

func main() {
	Logging.PrintVersion()
	if err := Logging.Init(true); err != nil {
		log.Panic("Cannot initialize logging module.", err)
	}
	defer Logging.Close()
	if ReadSetting(); PageId == 0 {
		fmt.Println("Please set room id value in setting.ini")
		return
	}
	defer func() { /* cope with unexpected exceptions */
		if err := recover(); err != nil {
			Logging.Error("Unexpected signal.", err.(error))
			Logging.Dump("debug.Stack", string(debug.Stack()))
			if ARunFail == false {
				return
			}
		} else {
			Logging.Note("Exit normally.")
		}
		if AutoRecord && ARunA != "" {
			if err := exec.Command(ARunA, ARunC).Start(); err != nil {
				Logging.Error(fmt.Sprintf("Cannot execute command \" %s %s \".", ARunA, ARunC), err)
			}
		}
	}()

	/* create folders */
	if _, err := os.Stat(DataBase); err != nil {
		if err = os.Mkdir(DataBase, os.ModePerm); err != nil {
			Logging.Error("Cannot create database folder.", err)
			return
		}
	}
	DataBase = path.Join(DataBase, time.Now().Format("2006-01-02-15_04_05"))
	if err := os.Mkdir(DataBase, os.ModePerm); err != nil {
		Logging.Error("Cannot create current database folder.", err)
		return
	}

	/* fetch m3u url */
	infurl := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/playUrl?cid=%d&platform=h5&otype=json&quality=0", PageId)
	infpage, err := Net.FetchData(infurl)
	if err != nil {
		Logging.Error("Fetch video information failed.", err)
		return
	}
	var blinf BiliveInfo
	if err = json.Unmarshal(infpage, &blinf); err != nil {
		Logging.Dump("infpage", string(infpage))
		Logging.Error("Cannot parse video information. Maybe Bilibili changed the api. Please contact with the developers.", err)
		return
	}
	if len(blinf.Data.Durl) == 0 {
		Logging.ErrorS("No m3u found. Maybe Bilibili changed the api. Please contact with the developers.")
		return
	}
	if SelDurl >= len(blinf.Data.Durl) {
		Logging.Warning(fmt.Sprintf("SelDurl(%d) out of range(%d), use default value 0.", SelDurl, len(blinf.Data.Durl)))
		SelDurl = 0
	}
	durl := blinf.Data.Durl[SelDurl].Url
	domain := Net.GetDomain(durl)
	protocol := Net.GetProtocol(durl)

	/* download ts files */
	ThreadPool = make([]bool, ThreadNum)
	lastseq := INF
	rundownload := func(url string, seq int, last int) {
		for j := 0; j < 10; j++ { //this two loops just for seeking a free position in thread pool
			for i := 0; i < ThreadNum; i++ {
				if ThreadPool[i] == false {
					ThreadPool[i] = true
					go func(url string, seq int, tid int) {
						filename := path.Join(DataBase, fmt.Sprintf("%d.ts", seq))
						for retry := 0; retry < RetryTimes; retry++ {
							if err := Net.Download(url, filename); err != nil {
								switch err {
								case Net.Forbidden403:
									Logging.Error(fmt.Sprintf("Download seq %d got 403. Abort.", seq), err)
									log.Panic("Got 403. Please check the information you provided.")
								case Net.Notfound404:
									Logging.Error(fmt.Sprintf("Download seq %d got 404. Abort.", seq), err)
									log.Panic("Got 404. Maybe the live video has finished.")
								}
								Logging.Error(fmt.Sprintf("Download seq %d failed. Try again.", seq), err)
								time.Sleep(time.Millisecond * time.Duration(ReConn))
							} else {
								if last == -1 {
									Logging.RePrint(fmt.Sprintf("Got %d seq... ", seq))
								} else {
									Logging.RePrint(fmt.Sprintf("Got %d seq... (%ds) ", seq, last))
								}
								ThreadPool[tid] = false
								break
							}
						}
					}(Net.CombUrl(protocol, domain, url), seq, i)
					return
				}
			}
			Logging.Info("Wait for a free thread... Maybe you need to increase the thread num.")
			time.Sleep(time.Millisecond * time.Duration(ReTrd))
		}
		log.Panic("Too many times waiting. Maybe network traffic jam. Abort.")
	}
	fetchm3u := func() *M3U.EXTM3U {
		for retry := 0; retry < RetryTimes; retry++ {
			m3upage, err := Net.FetchContext(durl)
			if err == nil {
				if m3ulst, err := M3U.ParseM3u(m3upage); err == nil {
					return m3ulst
				} else {
					Logging.Dump("m3upage", m3upage)
					Logging.Dump("err", err.Error())
					Logging.ErrorS("Cannot parse m3u. Please contact with the developer.")
					return nil
				}
			}
			if err == Net.Notfound404 {
				Logging.Error(fmt.Sprintf("Fetch m3u got 404. Maybe the live has finished. Abort.\n - Url: %s", durl), err)
				return nil
			} else {
				Logging.Error(fmt.Sprintf("Cannot fetch m3u list. Retry again.\n - Url: %s", durl), err)
				time.Sleep(time.Millisecond * time.Duration(ReConn))
			}
		}
		Logging.ErrorS("Run out of retry times. Abort.")
		return nil
	}
	signal2continue := true
	go func() { //this func is to control when to stop.
		if AutoRecord && ATime == 0 {
			Logging.Warning("autorecord on but recording time is zero. now turn off autorecord.")
			AutoRecord = false
		}
		if AutoRecord {
			fmt.Printf("\n[AutoRecord] Record %d seconds.\n\n", ATime)
			time.Sleep(time.Second * time.Duration(ATime))
			fmt.Printf("\nTimes up. Wait for threads to exit.\n")
		} else {
			iReader := bufio.NewReader(os.Stdin)
			fmt.Printf("\nPress [Enter] to stop recording. Please don't press Ctrl-C to stop unless it's out-of-control.\n\n")
			iReader.ReadString('\n')
			fmt.Printf("Receive signal to quit. Wait for threads to exit.\n")
		}
		signal2continue = false
	}()
	for signal2continue {
		m3ulst := fetchm3u()
		if m3ulst == nil {
			break
		}
		if m3ulst.MediaSequence == lastseq {
			Logging.Info("Too frequent. Maybe you need to decrease the waiting time.")
			Logging.Dump("TargetDuration", strconv.Itoa(m3ulst.TargetDuration))
			if TooFreq == 0 {
				time.Sleep(time.Second * time.Duration(math.Ceil(float64(m3ulst.TargetDuration)/2.0)))
			} else {
				time.Sleep(time.Millisecond * time.Duration(TooFreq))
			}
			continue
		}
		if m3ulst.MediaSequence-1 > lastseq { /* m3ulst.MediaSequence > lastseq + 1 */
			/* it means that we miss some ts file, we try to rescue them  */
			miss := m3ulst.MediaSequence - lastseq - 1
			Logging.Info(fmt.Sprintf("Rescue %d pack(s).", miss))
			for ; miss > 0; miss-- {
				if m3ulst.LenInf-miss-1 < 0 {
					Logging.ErrorS(fmt.Sprintf("Cannot rescue %d pack(s) due to network problem.", miss))
					break
				}
				rundownload(m3ulst.Inf[m3ulst.LenInf-1-miss].File, m3ulst.MediaSequence-miss, -1)
			}
		}
		rundownload(m3ulst.Inf[m3ulst.LenInf-1].File, m3ulst.MediaSequence, m3ulst.TargetDuration)
		lastseq = m3ulst.MediaSequence
		if Frequency == 0 {
			time.Sleep(time.Second * time.Duration(m3ulst.TargetDuration-1))
		} else {
			time.Sleep(time.Millisecond * time.Duration(Frequency))
		}
	}
}
