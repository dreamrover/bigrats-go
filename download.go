package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	chs "golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	READY = iota
	DOWN
	PAUSE
	ERROR
	DONE
)

func parseURL(url string) (*taskinfo, error) {
	var name, site, link string
	var task taskinfo

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reader := transform.NewReader(resp.Body, chs.GBK.NewDecoder())
	rd := bufio.NewReader(reader)
	htmldata, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	regR := regexp.MustCompile(`<R>.+`) //"?" represents lazy match
	task.title = string(bytes.TrimPrefix(regR.Find(htmldata), []byte("<R>")))
	if task.title == "" {
		return nil, errors.New("Script error: title not found")
	}

	regF := regexp.MustCompile(`<F>.+`)
	task.play = string(bytes.TrimPrefix(regF.Find(htmldata), []byte("<F>")))
	task.tid = md5.Sum([]byte(task.play))

	regQX := regexp.MustCompile(`<QX>.+`)
	task.quality = string(bytes.TrimPrefix(regQX.Find(htmldata), []byte("<QX>")))

	regS := regexp.MustCompile(`(?s)<\$>.+?<&>`) // (?s): let . include newline
	slices := regS.FindAll(htmldata, -1)
	if slices == nil {
		return nil, errors.New("Script error: video not found")
	}

	regN := regexp.MustCompile(`<N>.+`)
	regP := regexp.MustCompile(`<P>.+`)
	regU := regexp.MustCompile(`<U>.+`)
	for _, slice := range slices {
		name = string(bytes.TrimPrefix(regN.Find(slice), []byte("<N>")))
		site = string(bytes.TrimPrefix(regP.Find(slice), []byte("<P>")))
		link = string(bytes.TrimPrefix(regU.Find(slice), []byte("<U>")))
		name += filepath.Ext(link)
		sid := md5.Sum([]byte(link))
		task.segs = append(task.segs, &seginfo{name, site, link, sid, READY, -1, 0, -1, "-", &task})
	}
	task.suffix = filepath.Ext(task.segs[0].link)
	if task.suffix == "" {
		return nil, errors.New("Script error: unsupported video format")
	}
	return &task, nil
}

func fetchSegment(seg *seginfo, back chan *seginfo) (n int64, err error) {
	defer func() {
		back <- seg
	}()

	seg.status = DOWN
	seg.size = 0
	seg.down = 0
	seg.speed = 0
	back <- seg

	resp, err := http.Get(seg.link)
	if err != nil {
		log.Println(err)
		seg.status = ERROR
		return
	}
	/*fmt.Println("=============================================================")
	fmt.Println(seg.name, resp.Status, resp.Proto, resp.ContentLength)
	for k, v := range resp.Header {
		fmt.Println(k, v)
	}*/
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
		seg.status = ERROR
		return
	}
	length := resp.ContentLength
	seg.size = size(length)

	info, err := os.Stat(seg.task.dir + seg.name)
	if err == nil && info.Size() == length {
		n = length
		seg.down = size(n)
		seg.speed = -1
		seg.status = DONE
		seg.eta = "0s"
		return
	}
	file, err := os.Create(seg.task.dir + seg.name)
	if err != nil {
		log.Println(err)
		seg.status = ERROR
		return
	}
	defer file.Close()
	var down, n0, speed int64
	var eta string
	var dura time.Duration
	ticker := time.NewTicker(time.Second)
	for {
		down, err = io.CopyN(file, resp.Body, seglen)
		n += down
		select {
		case <-ticker.C:
			//fmt.Printf("\rDownloading %d / %d", n, length)
			//os.Stdout.Sync()
			speed = n - n0
			if n == n0 {
				eta = "∞"
			} else {
				t := strconv.FormatInt((length-n)/(n-n0), 10) + "s"
				dura, _ = time.ParseDuration(t)
				eta = dura.String()
			}
			seg.down = size(n)
			seg.speed = rate(speed)
			seg.eta = eta
			back <- seg
			n0 = n
		default:
		}
		if n == length || err == io.EOF {
			//fmt.Printf("\rFinished %d.             \n", length)
			ticker.Stop()
			seg.down = size(n)
			seg.speed = -1
			seg.eta = "0s"
			seg.status = DONE
			break
		}
	}
	return
}

func (task *taskinfo) Done() bool {
	for _, seg := range task.segs {
		if seg.status != DONE {
			return false
		}
	}
	return true
}

func (task *taskinfo) mergeSegs(format string, delsegs bool) error {
	var args []string

	if len(task.segs) == 1 && format == "Original" {
		if task.segs[0].name != task.title+task.suffix {
			return os.Rename(task.dir+task.segs[0].name, task.dir+task.title+task.suffix)
		}
		return nil
	}
	args = append(args, "--load")
	args = append(args, task.dir+task.segs[0].name)
	for _, seg := range task.segs[1:] {
		args = append(args, "--append")
		args = append(args, task.dir+seg.name)
	}
	args = append(args, "--output-format")
	if format == "Original" {
		format = task.suffix
	}
	suffix := format
	format = format[1:]
	if format == "mp4" {
		format = "mp4v2"
	} else if format == "flv" {
		format = "flv1"
	}
	args = append(args, format)
	args = append(args, "--save")
	args = append(args, task.dir+task.title+suffix)

	cmd := exec.Command(avidemux, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	if delsegs {
		for _, seg := range task.segs {
			e := os.Remove(task.dir + seg.name)
			if e != nil {
				err = e
			}
		}
	}
	return err
}
