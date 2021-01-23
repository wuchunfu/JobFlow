package main

import (
	"flag"
	"gin-vue/middleware/rpc/auth"
	"gin-vue/middleware/rpc/server"
	"gin-vue/utils/fileUtils"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"runtime"
	"strings"
)

func main() {
	var serverAddr string
	var allowRoot bool
	var CAFile string
	var certFile string
	var keyFile string
	var enableTLS bool
	var logLevel string
	flag.BoolVar(&allowRoot, "allow-root", false, "./gocron-node -allow-root")
	flag.StringVar(&serverAddr, "s", "0.0.0.0:5922", "./gocron-node -s ip:port")
	flag.BoolVar(&enableTLS, "enable-tls", false, "./gocron-node -enable-tls")
	flag.StringVar(&CAFile, "ca-file", "", "./gocron-node -ca-file path")
	flag.StringVar(&certFile, "cert-file", "", "./gocron-node -cert-file path")
	flag.StringVar(&keyFile, "key-file", "", "./gocron-node -key-file path")
	flag.StringVar(&logLevel, "log-level", "info", "-log-level error")
	flag.Parse()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(level)

	if enableTLS {
		if !fileUtils.FileExist(CAFile) {
			log.Fatalf("failed to read ca cert file: %s", CAFile)
		}
		if !fileUtils.FileExist(certFile) {
			log.Fatalf("failed to read server cert file: %s", certFile)
			return
		}
		if !fileUtils.FileExist(keyFile) {
			log.Fatalf("failed to read server key file: %s", keyFile)
			return
		}
	}

	certificate := auth.Certificate{
		CAFile:   strings.TrimSpace(CAFile),
		CertFile: strings.TrimSpace(certFile),
		KeyFile:  strings.TrimSpace(keyFile),
	}
	if runtime.GOOS != "windows" && os.Getuid() == 0 && !allowRoot {
		log.Fatal("Do not run gocron-node as root user")
		return
	}

	server.Start(serverAddr, enableTLS, certificate)
}
