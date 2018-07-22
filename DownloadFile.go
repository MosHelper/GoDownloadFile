package main

import (
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"os"
	"io"
)

func DownloadFile(filepath string, url string) error {
	log.Println("Crawler:",url)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Println("Error!:",err)
		return err
	}
	log.Println("Create file success. => ",filepath)
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error!:",err)
		return err
	}
	defer resp.Body.Close()
	log.Println("Init download ...")

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Error!:",err)
		return err
	}
	log.Println("Download success!")
	return nil
}

func GetDLLinkInUserscloud(url string) string {
	// Get key
	log.Println("Crawler:",url)
	key := strings.Split(url, "userscloud.com/")
	payload := strings.NewReader("op=download2&id=" + key[1])
	log.Println("key:",key[1])

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	link := ""
	doc.Find(".btn-icon-stacked").Each(func(i int, selection *goquery.Selection) {
		link, _ = selection.Attr("href")
		//log.Println("Found!:",link)
	})
	return link
}

func GetLinkUserscloudInFreefileload(url string) string {
	urlPost := "https://freefileload.tk/7x/"
	log.Println("Crawler:",url)

	key := strings.Split(url,"/lh/")
	payload := strings.NewReader("URL="+key[1])
	log.Println("key:",key[1])
	req, _ := http.NewRequest("POST", urlPost, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	link := ""
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() != "" {
			linkUserscloud := strings.Split(selection.Text(),"\"")

			//log.Println("Found!:",linkUserscloud[1])
			link = linkUserscloud[1]
		}
	})
	return link
}

