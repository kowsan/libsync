// http_sync_client project main.go
package main

import (
	"flag"
	"http_sync_client/synclib"
	"log"
	"os"
)

var (
	remote_server string
	dir           string
	use_csumm     bool
)

func main() {
	//set server url and local file path
	flag.StringVar(&remote_server, "url", "http://localhost:8181/fs", "set server url")
	td := os.TempDir()
	flag.StringVar(&dir, "dir", td, "Directory to file download")
	flag.BoolVar(&use_csumm, "csumm", false, "Use file control summ")
	flag.Parse()
	log.Println("Sync content from ", remote_server, " to : ", dir, "Csumm: ", use_csumm)
	synclib.Sync(remote_server, dir, use_csumm)
}
