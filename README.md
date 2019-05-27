# Bigrats for Linux

* [**中文版**](https://github.com/dreamrover/bigrats-go/blob/master/README_CN.md)

The official site of Bigrats is [flvcd.com](http://www.flvcd.com), which only provides Windows version.

This project is the Linux version of Bigrats video downloader written in [**Golang**](https://golang.org/) with GUI by deploying [GoQt](https://github.com/visualfc/goqt), and uses [avidemux](http://fixounet.free.fr/avidemux) to merge video segments. It will support MacOS very soon. This project has nothing to do with flvcd.com.

## Installation

### 1.Install flvcd-helper extension for Firefox
Install the extension from [https://addons.mozilla.org/zh-CN/firefox/addon/flvcd-helper](https://addons.mozilla.org/zh-CN/firefox/addon/flvcd-helper)

### 2.Install [jq](https://stedolan.github.io/jq/) and Qt5 runtime
* sudo apt-get install jq libqt5printsupport5 libqt5widgets5

### 3.Install [avidemux](http://fixounet.free.fr/avidemux/) for merging video segments
This procedure is different for Debian and Ubuntu.

#### Debian:
Add following line into /etc/apt/sources.list<br>
`deb http://ftp.kaist.ac.kr/debian-multimedia/ buster main`<br>
(replace _buster_ with your Debian codename, you can get it by `lsb_release -cs`)
* sudo apt-get update -oAcquire::AllowInsecureRepositories=true
* sudo apt-get install deb-multimedia-keyring
* sudo apt-get update
* sudo apt-get install avidemux-cli

#### Ubuntu:
* sudo add-apt-repository ppa:ubuntuhandbook1/avidemux
* sudo apt-get update
* sudo apt-get install avidemux2.6-cli<br>
**or**
* sudo apt-get install avidemux2.7-plugins-cli avidemux2.7-cli<br>

(you can also install avidemux-qt to merge video segments manually)

### 4.Install bigrats-go-x.x.x-amd64.deb
Download bigrats-go installer from [releases](https://github.com/dreamrover/bigrats-go/releases).<br>
* dpkg -i bigrats-go-x.x.x-amd64.deb<br>
(replace "x.x.x" with release version)<br>

Now you can use it as if you were using the Windows version!

## Build from source
If you are interest in building bigrats-go from source code, please refer to: <br>
[How to build bigrats-go from source on Debian/Ubuntu](https://github.com/dreamrover/bigrats-go/wiki/Build-Bigrats-on-Debian-and-Ubuntu)

## Screenshots
### Ubuntu
![image](https://github.com/dreamrover/screenshots/blob/master/bigrats-ubuntu-19.04.png)
### Debian
![image](https://github.com/dreamrover/screenshots/blob/master/bigrats-debian-buster.png)
