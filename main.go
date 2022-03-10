package main

import (
	"flag"
	"fmt"
	"goDown/utils"
	"log"
	"os"
)

var (
	h         bool
	v         bool
	url       string
	md5       string
	threadNum int
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.StringVar(&url, "url", "", "下载链接")
	flag.StringVar(&md5, "sha256", "", "sha256校验码")
	flag.IntVar(&threadNum, "thread", 4, "下载线程数量")
	flag.Usage = usage
}

func usage() {
	_, err := fmt.Fprintf(os.Stderr, `goDown version: goDown/1.0
Usage: goDown [-hv] [-url URL] [-t thred] [-sha256]
Options:
`)
	if err != nil {
		return
	}
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
	}
	if v {
		fmt.Println("goDown/1.0")
		os.Exit(0)
	}
	url := url
	md5 := md5
	if url == "" {
		log.Println("请输入合法的下载链接！")
		return
	}
	downloader := utils.NewFile(url, md5, threadNum)
	if err := downloader.Run(threadNum); err != nil {
		log.Fatal(err)
	}
	log.Println("文件下载完成!")
}
