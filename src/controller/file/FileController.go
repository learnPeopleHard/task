package file

import (
	"encoding/json"
	"fmt"
	"gotest1/src/common"
	"gotest1/src/logic/check"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const maxUploadSize = 2 * 1024 * 2014 // 2 MB
const uploadPath = "/Users/weihua.tian/GolandProjects/gotest1/tmp"

func UploadFileHandler(w http.ResponseWriter, r *http.Request)  {
	resp := UploadFile(w,r)
	//向客户端发送json数据
	bytes, _ := json.Marshal(resp)
	fmt.Fprint(w, string(bytes))
}

func UploadFile(w http.ResponseWriter, r *http.Request) common.BaseResponseDTO{
		tokenObj,err := check.CheckLoginAuthPost(w,r)
		if err!=nil{
			return common.BaseResponseDTO{Code: 5,Message: err.Error()}
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			return common.BaseResponseDTO{101,"FILE_TOO_BIG"}
		}
		fileType := r.PostFormValue("type")
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			return common.BaseResponseDTO{101,"INVALID_FILE"}
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return common.BaseResponseDTO{101,"INVALID_FILE"}
		}
		filetype := http.DetectContentType(fileBytes)
		if filetype != "image/jpeg" && filetype != "image/jpg" &&
			filetype != "image/gif" && filetype != "image/png" &&
			filetype != "application/pdf" {
			return common.BaseResponseDTO{101,"INVALID_FILE_TYPE"}
		}
		fileName := time.Now().UnixNano()
		fileEndings, err := mime.ExtensionsByType(fileType)
		if err != nil {
			return common.BaseResponseDTO{101,"CANT_READ_FILE_TYPE"}
		}

		newPath := filepath.Join(uploadPath,  strconv.FormatInt(fileName,10)+fileEndings[0])
		fmt.Printf("userName: %s FileType: %s, File: %s\n", tokenObj.UserName,fileType, newPath)
		newFile, err := os.Create(newPath)
		if err != nil {
			fmt.Println("CANT_WRITE_FILE error:", err.Error())
			return common.BaseResponseDTO{101,"CANT_WRITE_FILE"}
		}
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil {
			fmt.Println("CANT_WRITE_FILE error:", err.Error())
			return common.BaseResponseDTO{101,"CANT_WRITE_FILE"}
		}
		return common.BaseResponseDTO{100,newPath}
	}
