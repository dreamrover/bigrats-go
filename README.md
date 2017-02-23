# bigrats-go
硕鼠Linux版 (Bigrats for Linux)

The official site of Bigrats(硕鼠下载器) is http://www.flvcd.com/, which only provides Windows and Mac OSX version.

This project is the Linux version of Bigrats video downloader written in [**Golang**](https://golang.org/) with GUI by deploying [GoQt](https://github.com/visualfc/goqt). It will support Mac OSX in the future. This project has nothing to do with flvcd.com.

## usage
The usage of bigrats-go is quite similar with the Windows version except that it uses [**avidemux**](http://fixounet.free.fr/avidemux/) to merge video segments.

1. Download bigrats-go binaries(x86_64 only) from https://github.com/dreamrover/bigrats-go/releases and extract, e.g. to /opt. Change to the directory and:

    >chmod +x bigrats

    >Debian: sudo cp -d lib/* /usr/lib
    
    >CentOS/RHEL: sudo cp -d lib/* /usr/lib64

    If you build from source, building GoQt with **Qt5** is highly recommended to prevent crash on receiving signals.
2. Install Qt runtime libraries.
3. Install Firefox browser and flvcd-helper extension from https://addons.mozilla.org/zh-CN/firefox/addon/flvcd-helper/.
4. Edit ~/.mozilla/firefox/xxxxxxxx.default/mimeTypes.rdf, ('xxxxxxxx' may be deferent in deferent systems)

    add a new line     
    >`<RDF:li RDF:resource="urn:scheme:bigrats"/>` 
    
    after  
    >`<RDF:Seq RDF:about="urn:schemes:root">`
    
    
    And then add following lines (replace '/opt/bigrats-go' with the directory you extracted to):
    
     >`<RDF:Description RDF:about="urn:scheme:bigrats" NC:value="bigrats">`<br>
     >`<NC:handlerProp RDF:resource="urn:scheme:handler:bigrats"/>`<br>
     >`</RDF:Description>`<br>
     >`<RDF:Description RDF:about="urn:handler:local:/opt/bigrats-go/bigrats" NC:prettyName="bigrats" NC:path="/opt/bigrats-go/bigrats" />`    
    
    Open Firefox--Preferences--Applications, choose bigrats binary file in the 'Action' column:
    ![image](https://github.com/dreamrover/screenshots/blob/master/settings.png)
    
5. Install [**avidemux**](http://fixounet.free.fr/avidemux/) for merging video segments. Avidemux is not included in the default software repositories of Debian or CentOS/RHEL.

    >Debian users should add deb-multimedia repository by the instructions in https://deb-multimedia.org/ first, 
    and run: sudo apt-get install avidemux avidemux-cli
    
    >CentOS/RHEL users should add Nux repository by the instructions in http://li.nux.ro/repos.html first, 
    and run: sudo yum install avidemux avidemux-cli
    
6. Now you can use it like the official Windows version.
    
Screenshots:

![image](https://github.com/dreamrover/screenshots/blob/master/bigrats-go.png)
