// http_sync_client project main.go
package main

import (
	_ "encoding/json"
	"flag"
	"http_sync_client/synclib"
	_ "httpsynccommon"
	_ "io/ioutil"
	_ "net/http"
	_ "path/filepath"
	_ "time"

	"log"
	_ "net/url"
	"os"
)

var (
	remote_server *string
	dir           *string
)

func main() {
	//set server url and local file path
	remote_server = flag.String("url", "http://localhost:8181/fs", "set server url")
	td := os.TempDir()
	//td := "/home/kovalev/12345678"
	dir = flag.String("dir", td, "Directory to file download")
	flag.Parse()
	log.Println("Sync content from ", *remote_server, " to : ", *dir)
	synclib.Sync(*remote_server, *dir)
}

//	//connect to server
//	r, err := http.Get(*remote_server + "/index.go")
//	if err != nil {
//		log.Println("could not get server directory : ", err)
//		return
//	}
//	defer r.Body.Close()
//	result, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Println(" response body could not read ", err)
//	} else {
//		//log.Println("server file content :", string(result))
//		if r.StatusCode == 200 {
//			log.Println("content ok, sync with it ", string(result))
//			content := map[string]httpsynccommon.FileInfo{}
//			sc := json.Unmarshal(result, &content)
//			if sc != nil {
//				log.Println("could not unmarshall server content, ", sc)
//			} else {
//				log.Println("server content  OK: ", content)
//				syncContent(content, *dir)
//			}

//		} else {
//			log.Println("bad status while read content ", r.StatusCode, string(result))
//		}
//	}

//}
//func syncContent(remote map[string]httpsynccommon.FileInfo, dir string) {

//	local := httpsynccommon.BuildFileStructure(dir)
//	log.Println("local content : ", local)
//	log.Println("get file not existing in local")
//	var files_to_download []string
//	var files_to_remove []string
//	for k, v := range remote {
//		lv, ok := local[k]
//		if ok {
//			//log.Println("file found :validate it by size", k, v)
//			if lv.Size != v.Size || lv.ModTime != v.ModTime {
//				log.Println("file differs by time or last access: Download", k)
//				files_to_download = append(files_to_download, k)

//			}
//		} else {

//			files_to_download = append(files_to_download, k)
//			log.Println("Append file to download list : ", k)
//		}
//	}
//	log.Println("NEED Download files  : ", len(files_to_download))
//	for _, f := range files_to_download {
//		downloadFile(f, remote[f].ModTime)
//	}

//	//remove files existing locally
//	for k, v := range local {
//		_, ok := remote[k]
//		if ok {

//		} else {

//			files_to_remove = append(files_to_remove, k)
//			log.Println("Append file to remove list : ", k, v)
//		}
//	}
//	log.Println("NEED Remove files  : ", len(files_to_remove))
//	for _, v := range files_to_remove {
//		fp := dir + v
//		e := os.Remove(fp)
//		if e != nil {
//			log.Println("Could not remove file ", fp)
//		} else {
//			log.Println("file removed : ", fp)
//		}

//	}

//}

//func downloadFile(srvpath string, mtime int64) bool {
//	u, _ := url.Parse(*remote_server + srvpath)
//	log.Println("download from network by url", u.String())
//	r, e := http.Get(u.String())
//	if e != nil {
//		log.Println("could not download file : ", e)
//		return false
//	}
//	defer r.Body.Close()
//	b, e := ioutil.ReadAll(r.Body)
//	if e != nil {
//		log.Println("could not read file content : ", e)
//		return false
//	} else {
//		//make dir for file
//		file_path := *dir + srvpath
//		file_dir := filepath.Dir(file_path)
//		if _, err := os.Stat(file_dir); os.IsNotExist(err) {
//			e := os.MkdirAll(file_dir, 0777)
//			if e != nil {
//				log.Println("could not create dir ", file_dir, "Error : ", e)
//				return false

//			} else {
//				log.Println("created directory  ", file_dir)
//			}
//		}
//		//write to file :
//		err := ioutil.WriteFile(file_path, b, 0666)
//		if err != nil {
//			log.Println("could not write file : ", file_path, " : ", err)
//			return false
//		} else {
//			log.Println("File ok downloaded : ", file_path)
//			e := os.Chtimes(file_path, time.Now(), time.Unix(mtime, 0))
//			if e != nil {
//				log.Println("could not chtime of file ", file_path, " error ", e)
//			} else {
//				log.Println("Chtime OK ", file_path)
//			}
//			return true
//		}
//	}
//	return true
//}
