package main

import (
	"fmt"
	"log"
)

func main() {
	l := log.Logger{}
	err := UploadFileToDrive()
	if err != nil {
		l.Println("error: ", err)
		msg := fmt.Sprintf("fail: %s", err.Error())
		err = UploadDataToSplunk(msg)
		if err != nil {
			l.Println("error: ", err)
		}
	} else {
		err = UploadDataToSplunk("success")
		if err != nil {
			l.Println("error: ", err)
		}
	}
}
