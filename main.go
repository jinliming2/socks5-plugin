package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	"golang.org/x/net/proxy"
)

// VERSION should be set when build
var VERSION = "unknown"

var (
	remoteAddress = flag.String("remoteAddr", "[::1]", "Remote address to forward")
	remotePort    = flag.Int("remotePort", 1080, "Remote port to forward")
	socks5Address = flag.String("socks5Addr", "[::1]", "Socks 5 address")
	socks5Port    = flag.Int("socks5Port", 1080, "Socks 5 port")
	localAddress  = flag.String("localAddr", "[::1]", "Local address to listen")
	localPort     = flag.Int("localPort", 1080, "Local port to listen")
	logFileName   = flag.String("log", "", "Log file name, default to stdout")
	version       = flag.Bool("version", false, "Show current version and exit")
	_             = flag.Bool("fast-open", false, "Enable TCP Fast Open (Not implement yet.)")
)

func main() {
	flag.Parse()

	if len(*logFileName) > 0 && *logFileName != "-" {
		logFile, err := os.OpenFile(*logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln("Cannot open log file!")
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	if *version {
		fmt.Println("Socks5-Plugin", VERSION)
		fmt.Println("Golang", runtime.Version())
		return
	}

	log.Println("Socks5 Plugin", VERSION)

	parseEnv()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(*localAddress),
		Port: *localPort,
	})
	if err != nil {
		log.Fatalf("Can't listen on local: %v\n", err)
	}
	defer listener.Close()
	log.Printf("Listen on %s:%d\n", *localAddress, *localPort)

	dailer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", *socks5Address, *socks5Port), nil, proxy.Direct)
	if err != nil {
		log.Fatalf("Can't connect to the proxy: %v\n", err)
	}

	remote := fmt.Sprintf("%s:%d", *remoteAddress, *remotePort)

	for {
		localConn, err := listener.Accept()
		log.Println("Accept", localConn.RemoteAddr().String())
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			remoteConn, err := dailer.Dial("tcp", remote)
			if err == nil {
				Pipe(localConn, remoteConn)
				remoteConn.Close()
			}
			localConn.Close()
		}()
	}
}
