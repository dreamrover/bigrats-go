BIN=./deb/usr/bin/bigrats-go
DEB=./bigrats-go-0.2-amd64.deb

all: deb

deb: $(DEB)

$(DEB): $(BIN)
	dpkg -b ./deb $(DEB)

$(BIN): *.go
	go build -o $(BIN)
	strip $(BIN)

.PHONY:clean
clean:
	rm $(BIN) $(DEB)
