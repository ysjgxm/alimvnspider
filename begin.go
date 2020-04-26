package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type FolderStruct struct {
	Object []Node
	Successful bool
}
type Node struct {
	NodeKey string
	NodeName string
	NodePath string
	NodeType string
	RepoId string
}
type VOStruct struct {
	ArtifactId   string
	Classifier   string
	GroupId      string
	Packaging    string
	RepositoryId string
	Version				string
}

type FileDetail struct {
	ArtifactDetailVO	VOStruct
	ContentLength		string
	DownloadUrl			string
	Exist				string
	LastModified		string
	Path				string
}

type FileStruct struct {
	Object FileDetail
	Successful bool
}

var httpcli = http.Client{}

func main() {
	fmt.Println("根URL是：   " + os.Args[1])
	fmt.Println("本地目录是： " + os.Args[2])
	c := colly.NewCollector(
		colly.Async(true))
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
	})

	c.OnRequest(func(r *colly.Request){
		fmt.Println("start visiting ", r.URL)
	})
	c.OnError(func(_ *colly.Response, e error) {
		fmt.Println("something error ", e)
	})
	c.OnResponse(func(r *colly.Response) {
		uriSuffix := strings.Split(r.Request.URL.RequestURI(), "?_input_charset=")[0]
		if uriSuffix == "/browse/tree" {
			o := FolderStruct{}
			json.Unmarshal(r.Body, &o)
			rand.Seed(time.Now().UnixNano())
			sleepTime := rand.Intn(100)
			for _, v := range o.Object {
				myv := v
				if myv.NodeType == "FOLDER" {
					os.MkdirAll(os.Args[2] + myv.NodePath, os.ModePerm)
					time.Sleep(time.Duration(200 * sleepTime) * time.Millisecond)
					c.Visit("https://maven.aliyun.com/browse/tree?_input_charset=utf-8&repoId=central&path="+myv.NodePath)
				} else if myv.NodeType == "FILE" {
					time.Sleep(time.Duration(80 * sleepTime) * time.Millisecond)
					c.Visit("https://maven.aliyun.com/browse/fileInfo?_input_charset=utf-8&repoId=central&path="+myv.NodePath)
				}

			}
		} else if uriSuffix == "/browse/fileInfo" {
			o := FileStruct{}
			json.Unmarshal(r.Body, &o)
			tmpPath := o.Object.Path
			downUrl := o.Object.DownloadUrl
			//下载文件
			downFile(downUrl, os.Args[2] + tmpPath)
		}

	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnScraped(func(r *colly.Response) {
		//fmt.Println("finished ", r.Request.URL)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(os.Args[1])
	c.Wait()
}

func downFile(link string, fullFileName string) error{
	resp,err := httpcli.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	allBody,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fullFileName,allBody,0666)
	if err != nil {
		return err
	}
	return nil
}




