package main

import (
	"os"

	"github.com/mholt/archiver"
	"io/ioutil"
	"os/exec"
)
var basePath = "/home/"+os.Getenv("USER")+"/_karanCrawler"
func main() {
	os.Mkdir(basePath,os.ModePerm)
	SaveFile("https://freefileload.tk/lh/zjNU1z0ElQuy0a64J7O080x0Q5or_x97khONxS4bM9t12gi6xGPGDWNNTpo0_8inI5ziCxdSS49YrD371QD5_g")
}

func SaveFile(url string)  {

	s := GetLinkUserscloudInFreefileload(url)
	DownloadFile(basePath+"/temp.rar", 	GetDLLinkInUserscloud(s))
	archiver.Rar.Open(basePath+"/temp.rar", basePath+"/temp")
	os.Remove(basePath+"/temp.rar")

	readme := []byte("Download software free => https://loadhit.com\nGoogle Drive Fast download Unlimited.")
	ioutil.WriteFile(basePath+"/temp/Readme.txt",readme, 0644)
	zipscript := basePath+"/zip.sh"
	ioutil.WriteFile(zipscript,[]byte("#!/bin/bash\ncd "+basePath+"/temp\nzip -P passw0rd -r ../secure.zip ./"), 0644)

	exec.Command("/bin/sh",zipscript).Run()
	os.RemoveAll(basePath+"/temp")
	os.RemoveAll(basePath+"/zip.sh")

}