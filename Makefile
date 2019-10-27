DATE=`date +%Y%m%d%H%M%S`
LDFLAGS="-X main.VERSION=${DATE} -s -w"

.PHONY: all build clean

all: clean build

build:
	go build -v -ldflags ${LDFLAGS} -o socks5-plugin

clean:
	if [ -a socks5-plugin ] ; then rm socks5-plugin ; fi;
