package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/visualfc/goqt/ui"
)

const (
	seglen = 1452
	N      = 32
)

var chRow chan *seginfo
var chURL chan string
var chDir chan urldir
var chMsg chan string
var chTask chan *taskinfo
var chMrg chan string

var avidemux = [...]string{"avidemux_cli", "avidemux2.7_cli", "avidemux3_cli"}
var xdown bool = false
var threads int32 = 5
var automerge bool = true
var merger string
var autodel bool = true
var container = "Original"
var cindex int32 = 0
var dirs DirArray

func loadConfig(file string) {
	var config Config
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}
	xdown = config.Xdown
	threads = config.Threads
	automerge = config.Automerge
	autodel = config.Autodel
	cindex = config.CIndex
	dirs = config.Dirs
}

func dumpConfig(file string) {
	config := Config{xdown, threads, automerge, autodel, cindex, dirs}
	data, _ := json.MarshalIndent(config, "", "    ")
	ioutil.WriteFile(file, data, 0644)
}

func init() {
	chRow = make(chan *seginfo, N)
	chURL = make(chan string, N)
	chDir = make(chan urldir, N)
	chMsg = make(chan string, N)
	chTask = make(chan *taskinfo, N)
	chMrg = make(chan string, N)
}

func main() {
	var url string
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "bigrats://") {
		url = "http://www.flvcd.com/diy/" + os.Args[1][10:] + ".htm"
	}
	uid := os.Getuid()
	addr := "/tmp/bigrats:" + strconv.Itoa(uid)
listen:
	conn, err := net.ListenUnixgram("unixgram", &net.UnixAddr{addr, "unixgram"})
	if err != nil {
		conn, err = net.DialUnix("unixgram", nil, &net.UnixAddr{addr, "unixgram"})
		if err != nil {
			log.Println(err)
			os.Remove(addr)
			goto listen
		}
		_, err = conn.Write([]byte(url))
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer os.Remove(addr)

	user, err := user.Current()
	if err == nil {
		cfgfile := filepath.Join(user.HomeDir, ".bigrats")
		loadConfig(cfgfile)
		defer dumpConfig(cfgfile)
	}

	go func() {
		for {
			msg := make([]byte, 256)
			n, err := conn.Read(msg)
			if err != nil {
				log.Println(err)
				continue
			}
			runTask(string(msg[:n]))
		}
	}()

	if url != "" {
		go runTask(url)
	}

	go scheduler()

	ui.RunEx(os.Args, gui)
}

func runTask(url string) {
	chURL <- url
	tmp := <-chDir
	url = tmp.url
	dir := tmp.dir
	if url == "" || dir == "" {
		return
	}
	task, err := parseURL(url)
	if err != nil {
		chMsg <- err.Error()
		return
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	task.dir = dir
	chTask <- task
}

func scheduler() {
	var tasks []*taskinfo
	var active int32

	chseg := make(chan *seginfo, N)
	sched := func() {
		if active >= threads {
			return
		}
		for _, task := range tasks {
			for _, seg := range task.segs {
				if seg.status != READY {
					continue
				}
				go fetchSegment(seg, chseg)
				active++
				if active >= threads {
					return
				}
			}
		}
	}
	for {
	sel:
		select {
		case t := <-chTask:
			if t != nil {
				for _, task := range tasks {
					if t.tid == task.tid {
						chMsg <- "This video is already in task list."
						break sel
					}
				}
				for _, seg := range t.segs {
					chRow <- seg
				}
				tasks = append(tasks, t)
			}
			sched()
		case seg := <-chseg:
			chRow <- seg
			if seg.status == DONE || seg.status == ERROR {
				active--
				sched()
			}
			if seg.status == DONE && seg.task.Done() && automerge {
				chMrg <- "Merging " + seg.task.title
				go func() {
					err := seg.task.mergeSegs(container, autodel)
					if err != nil {
						chMrg <- err.Error()
					} else {
						chMrg <- ""
					}
				}()
			}
		}
	}
}
