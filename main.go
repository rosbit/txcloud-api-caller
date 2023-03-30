package main

import (
	sslcert "txcloud-api-caller/ssl-cert"
	"txcloud-api-caller/live"
	"txcloud-api-caller/conf"
	"fmt"
	"os"
)

func init() {
	if err := conf.CheckGlobalConf(); err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <command> <options>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Where <command> are:\n")
	fmt.Fprintf(os.Stderr, "  upload-ssl-certs\n")
	fmt.Fprintf(os.Stderr, "  bind-live-ssl-certs\n")
}

func main() {
	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "upload-ssl-certs":
		sslcert.UploadSslCerts()
	case "bind-live-ssl-certs":
		live.BindSSLCerts()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %s\n", os.Args[1])
		os.Exit(2)
	}
}

