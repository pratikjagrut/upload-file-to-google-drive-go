package main

import "log"

func main() {
	err := UploadFileToDrive()
	if err != nil {
		log.Fatal(err)
	}
}
