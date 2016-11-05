package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		//fmt.Println(name, site, link)
		name += filepath.Ext(link)
		sid := md5.Sum([]byte(link))
		task.segs = append(task.segs, &seginfo{name, site, link, sid, READY, &task})
	}
	task.suffix = filepath.Ext(task.segs[0].link)
	if task.suffix == "" {
		return nil, errors.New("Script error: unsupported video format")
	}
	//fmt.Println(task.suffix)
	return &task, nil
}

func fetchSegment(seg *seginfo, back chan backinfo) (n int64, err error) {
	rinfo := rowinfo{seg, 0, 0, 0, "-", "Connecting"}
	chRow <- rinfo

	resp, err := http.Get(seg.link)
	if err != nil {
		log.Println(err)
		rinfo.status = "Error"
		chRow <- rinfo
		back <- backinfo{seg, ERROR}
		return
	}
	fmt.Println("=============================================================")
	fmt.Println(seg.name, resp.Status, resp.Proto, resp.ContentLength)
	for k, v := range resp.Header {
		fmt.Println(k, v)
	}
	defer resp.Body.Close()
	length := resp.ContentLength
	rinfo.size = size(length)
	rinfo.status = "Downloading"
	chRow <- rinfo

	info, err := os.Stat(seg.task.dir + seg.name)
	if err == nil && info.Size() == length {
		n = length
		rinfo.down = size(n)
		rinfo.speed = -1
		rinfo.eta = "0s"
		rinfo.status = "Finished"
		chRow <- rinfo
		back <- backinfo{seg, DONE}
		return
	}
	file, err := os.Create(seg.task.dir + seg.name)
	if err != nil {
		log.Println(err)
		rinfo.status = "Error"
		chRow <- rinfo
		back <- backinfo{seg, ERROR}
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
			rinfo.down = size(n)
			rinfo.speed = rate(speed)
			rinfo.eta = eta
			chRow <- rinfo
			n0 = n
		default:
		}
		if n == length || err == io.EOF {
			//fmt.Printf("\rFinished %d.             \n", length)
			ticker.Stop()
			rinfo.down = size(n)
			rinfo.speed = -1
			rinfo.eta = "0s"
			rinfo.status = "Finished"
			chRow <- rinfo
			back <- backinfo{seg, DONE}
			break
		}
	}
	return
}

func mergeSegs(task *taskinfo, format string, joiner *mrgtool, del bool) (err error) {
	var args []string

	if len(task.segs) == 1 && format == "Original" {
		if len(task.segs) == 1 && task.title+task.suffix != task.segs[0].name {
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
	format = format[1:]
	if format == "mp4" {
		format = "mp4v2"
	}
	args = append(args, format)
	args = append(args, "--save")
	args = append(args, task.dir+task.title)
	args = append(args, fmt.Sprintf("%v", del))
	for _, arg := range args {
		_, err = fmt.Fprintln(joiner.wr, arg)
		if err != nil {
			return
		}
	}
	_, err = fmt.Fprint(joiner.wr, "\n")
	return
}

/*
func delSegs(task *taskinfo, format string) (err error) {
	var size, sum int64
	var file string
	var info os.FileInfo

	if format == "Original" {
		file = task.dir + task.title + task.suffix
	} else {
		file = task.dir + task.title + format
	}
	for _, seg := range task.segs[:len(task.segs)-1] {
		info, err = os.Stat(task.dir + seg.name)
		if err != nil {
			return
		}
		sum += info.Size()
	}
	for i := 0; i < len(task.segs); i++ {
		time.Sleep(time.Second)
		info, err = os.Stat(file)
		if err != nil {
			continue
		}
		if size != 0 && size == info.Size() && size > sum {
			for _, seg := range task.segs {
				e := os.Remove(task.dir + seg.name)
				if e != nil {
					err = e
				}
			}
			return
		}
		size = info.Size()
	}
	return errors.New("Error deleting segments.")
}
*/
