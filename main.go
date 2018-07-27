package main

import (
	"os"

	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"./uploadtodrive"
	"github.com/MosHelper/Userscloud-Download"
	"github.com/PuerkitoBio/goquery"
	"github.com/mholt/archiver"
)

var basePath = "/home/" + os.Getenv("USER") + "/_karanCrawler"

func main() {
	os.Mkdir(basePath, os.ModePerm)
	saveFile("https://freefileload.tk/lh/zjNU1z0ElQuy0a64J7O080x0Q5or_x97khONxS4bM9t12gi6xGPGDWNNTpo0_8inI5ziCxdSS49YrD371QD5_g")
	uploadtodrive.Upload(basePath+"/upload.zip", "sssssss.zip")
}

func saveFile(url string) {
	s := getLinkUserscloudInFreefileload(url)
	userscloud.DownloadFile(basePath+"/temp.rar", s)
	archiver.Rar.Open(basePath+"/temp.rar", basePath+"/temp")
	os.Remove(basePath + "/temp.rar")

	readme := []byte("Download software free => https://loadhit.com\nGoogle Drive Fast download Unlimited.")
	ioutil.WriteFile(basePath+"/temp/Readme.txt", readme, 0644)
	zipscript := basePath + "/zip.sh"
	ioutil.WriteFile(zipscript, []byte("#!/bin/bash\ncd "+basePath+"/temp\nzip -P loadhit.com -r ../upload.zip ./"), 0644)

	exec.Command("/bin/sh", zipscript).Run()
	os.RemoveAll(basePath + "/temp")
	os.RemoveAll(basePath + "/zip.sh")
}

func getLinkUserscloudInFreefileload(url string) string {
	urlPost := "https://freefileload.tk/7x/"
	log.Println("Crawler:", url)

	key := strings.Split(url, "/lh/")
	payload := strings.NewReader("URL=" + key[1])
	log.Println("key:", key[1])
	req, _ := http.NewRequest("POST", urlPost, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	link := ""
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() != "" {
			linkUserscloud := strings.Split(selection.Text(), "\"")
			log.Println("Found!:", linkUserscloud[1])
			link = linkUserscloud[1]
		}
	})
	return link
}
