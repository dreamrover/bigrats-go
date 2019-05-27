# 硕鼠网页视频下载器Linux版

* [**English version**](https://github.com/dreamrover/bigrats-go/edit/master/README.md)

硕鼠下载器官网[flvcd.com](http://www.flvcd.com)仅提供Windows版客户端，对于Linux用户，只好自己动手丰衣足食，于是有了本项目。

本项目是用[**Golang**](https://golang.org/)开发的，图形界面基于[GoQt](https://github.com/visualfc/goqt)，并调用[avidemux](http://fixounet.free.fr/avidemux)来合并所下载的视频片段。本项目很快将支持苹果Mac OSX系统。本项目的开发与硕鼠官网无关。

## 安装指南

### 1. 安装硕鼠下载器插件flvcd-helper
从[https://addons.mozilla.org/zh-CN/firefox/addon/flvcd-helper](https://addons.mozilla.org/zh-CN/firefox/addon/flvcd-helper)安装插件

### 2. 安装[jq](https://stedolan.github.io/jq/)和Qt5运行时环境
* sudo apt-get install jq libqt5printsupport5 libqt5widgets5

### 3. 安装用于合并视频片段的[avidemux](http://fixounet.free.fr/avidemux/)
这一步对于Debian和Ubuntu略有不同：
#### Debian
在文件/etc/apt/sources.list中添加下面这行：<br>
`deb http://ftp.kaist.ac.kr/debian-multimedia/ buster main`<br>
（注意将 _buster_ 替换成你所安装的Debian版本的开发代号，可通过命令`lsb_release -cs`获得）
* sudo apt-get update -oAcquire::AllowInsecureRepositories=true
* sudo apt-get install deb-multimedia-keyring
* sudo apt-get update
* sudo apt-get install avidemux-cli<br>
#### Ubuntu
* sudo add-apt-repository ppa:ubuntuhandbook1/avidemux
* sudo apt-get update
* sudo apt-get install avidemux2.7-plugins-common avidemux2.7-plugins-cli avidemux2.7-cli<br>

（可以同时安装带图形界面的avidemux2.7-qt5用于手动合并或编辑视频）

### 4. 安装bigrats-go-x.x.x-amd64.deb
从[releases](https://github.com/dreamrover/bigrats-go/releases)下载bigrats-go安装包，然后：
* dpkg -i bigrats-go-x.x.x-amd64.deb<br>
(将“x.x.x”替换为软件版本号)<br>

现在可以像使用Windows版的硕鼠下载器一样使用了，解析视频之后，点击“硕鼠专用链下载”，即可自动启动客户端，就像在Windows上一样。

## 从源码编译bigrats-go
如果你对编译bigrats-go感兴趣，请参考[在Debian和Ubuntu上从源码编译硕鼠下载器](https://github.com/dreamrover/bigrats-go/wiki/%E5%9C%A8Debian%E5%92%8CUbuntu%E4%B8%8A%E4%BB%8E%E6%BA%90%E7%A0%81%E7%BC%96%E8%AF%91%E7%A1%95%E9%BC%A0%E4%B8%8B%E8%BD%BD%E5%99%A8)

## 屏幕截图
### Ubuntu
![image](https://github.com/dreamrover/screenshots/blob/master/bigrats-ubuntu-19.04.png)
### Debian
![image](https://github.com/dreamrover/screenshots/blob/master/bigrats-debian-buster.png)
