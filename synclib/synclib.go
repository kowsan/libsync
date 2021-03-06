package synclib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"log"
	"net/url"
	"os"
)

var (
	remote_server string
	dir           string
)

type FileInfo struct {
	//Name    string    `json:"name"`
	ModTime int64  `json:"modtime"`
	Size    int64  `json:"size"`
	Md5     string `json:"md5"`
}

func hash_file_md5(filePath string) string {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String

}

func BuildFileStructure(dir string, useCsumm bool) map[string]FileInfo {

	log.Println("Begin ", dir)
	fileList := map[string]FileInfo{}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		e := os.MkdirAll(dir, 0777)
		if e != nil {
			log.Println("could not create dir ", dir, "Error : ", e)
			return fileList
		}
	}

	err := filepath.Walk(dir, func(f_path string, f os.FileInfo, err error) error {
		if err == nil {
			if !f.IsDir() {
				var fi FileInfo

				//fi.Name = path
				fi.ModTime = f.ModTime().UTC().Unix()
				fi.Size = f.Size()
				if useCsumm {
					fi.Md5 = hash_file_md5(f_path)
				} else {
					fi.Md5 = "0"
				}

				//fileList = append(fileList, fi)

				p := strings.TrimPrefix(f_path, dir)
				p = strings.Replace(p, "\\", "/", -1)
				fileList[p] = fi
			}
		}
		return nil
	})
	if err != nil {
		log.Println("error file list", err)
	}

	return fileList
}
func Sync(server_url, directory string, useCsumm bool) {
	//	remote_server = flag.String("url", "http://localhost:8181/fs", "set server url")
	//	td := os.TempDir()
	//	//td := "/home/kovalev/12345678"
	//	dir = flag.String("dir", td, "Directory to file download")

	//flag.Parse()
	remote_server = server_url
	dir = directory
	log.Println("Sync content from ", remote_server, " to : ", dir)

	//connect to server
	r, err := http.Get(remote_server + "/index.go")
	if err != nil {
		log.Println("could not get server directory : ", err)
		return
	}
	defer r.Body.Close()
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(" response body could not read ", err)
	} else {
		//log.Println("server file content :", string(result))
		if r.StatusCode == 200 {
			log.Println("content ok, sync  it ")
			content := map[string]FileInfo{}
			sc := json.Unmarshal(result, &content)
			if sc != nil {
				log.Println("could not unmarshall server content, ", sc)
			} else {
				//log.Println("server content  OK: ", content)
				syncContent(content, dir, useCsumm)
			}

		} else {
			log.Println("bad status while read content ", r.StatusCode, string(result))
		}
	}
}

func syncContent(remote map[string]FileInfo, dir string, useCsumm bool) {

	local := BuildFileStructure(dir, useCsumm)
	//log.Println("local content : ", local)
	log.Println("get file not existing in local")
	var files_to_download []string
	var files_to_remove []string
	for k, v := range remote {
		lv, ok := local[k]
		if ok {
			//log.Println("file found :validate it by size", k, v)

			if lv.Size != v.Size || (useCsumm && lv.Md5 != v.Md5) {
				log.Println("file differs by size or md5: Download", k)
				files_to_download = append(files_to_download, k)

			}
		} else {

			files_to_download = append(files_to_download, k)
			//log.Println("Append file to download list : ", k)
		}
	}
	log.Println("NEED Download files  : ", len(files_to_download))
	for _, f := range files_to_download {
		downloadFile(f, remote[f].ModTime)
	}

	//remove files existing locally
	for k, _ := range local {
		_, ok := remote[k]
		if ok {

		} else {

			files_to_remove = append(files_to_remove, k)
			//log.Println("Append file to remove list : ", k, v)
		}
	}
	log.Println("NEED Remove files  : ", len(files_to_remove))
	for _, v := range files_to_remove {
		fp := dir + v
		e := os.Remove(fp)
		if e != nil {
			log.Println("Could not remove file ", fp, e)
		} else {
			log.Println("file removed : ", fp)
		}

	}

}

func downloadFile(srvpath string, mtime int64) bool {
	u, e := url.Parse(remote_server)
	u.Path = srvpath
	if e != nil {
		log.Fatalln("Could not download  ", srvpath, e)
	}
	log.Println("download from network by path", srvpath)
	log.Println("download from network by url", u.String())
	r, e := http.Get(u.String())
	if e != nil {
		log.Fatalln("could not download file : ", e)

	}

	//	b, e := ioutil.ReadAll(r.Body)
	//	if e != nil {
	//		log.Println("could not read file content : ", e)
	//		return false
	//	} else {
	//make dir for file
	file_path := dir + srvpath
	file_dir := filepath.Dir(file_path)

	if _, err := os.Stat(file_dir); os.IsNotExist(err) {
		e := os.MkdirAll(file_dir, 0777)
		if e != nil {
			log.Fatalln("could not create dir ", file_dir, "Error : ", e)

		} else {
			log.Println("created directory  ", file_dir)
		}
	}

	//write to file :
	fp, ex := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if ex != nil {
		log.Fatalln("Could not open file for write ", ex)
	}
	bytes_c, err := io.Copy(fp, r.Body)
	defer fp.Close()
	//err := ioutil.WriteFile(file_path, b, 0666)
	if err != nil {
		log.Fatalln("could not write file : ", file_path, " : ", err, bytes_c)

	} else {
		log.Println("File ok downloaded : ", file_path, bytes_c)

		return true
	}

	return true
}
