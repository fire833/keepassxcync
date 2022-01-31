
all: 
	go build cmd/keepassxcync/main.go

install:
	install cmd/keepassxcync/keepassxcync /usr/bin/keepassxcync

clean:
	rm -rf /usr/bin/keepassxcync