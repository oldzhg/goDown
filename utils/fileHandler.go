package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"

func NewFile(downloadUrl string, md5 string, partNumber int) *fileInfo {
	return &fileInfo{
		URL:       downloadUrl,
		MD5:       md5,
		DoneParts: make([]filePart, partNumber),
	}
}

func (f *fileInfo) getHeader() error {
	resp, err := http.Head(f.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return fmt.Errorf("%s does not support range requests", f.URL)
	}
	if resp.Header.Get("Content-Disposition") != "" {
		contentDisposition := resp.Header.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			return err
		}
		f.Name = params["filename"]
	} else {
		f.Name = filepath.Base(f.URL)
	}
	f.Size, _ = strconv.Atoi(resp.Header.Get("Content-Length"))
	return nil
}

func (f fileInfo) downloadPart(c filePart) error {
	req, err := http.NewRequest("GET", f.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", c.Start, c.End))
	req.Header.Set("User-Agent", userAgent)
	//log.Printf("开始[%d]下载from:%d to:%d\n", c.Index, c.Start, c.End)
	log.Printf("开始下载分片%d\n", c.Index)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 206 {
		return errors.New(fmt.Sprintf("服务器状态码有误: %v", resp.StatusCode))
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(bs) != c.End-c.Start+1 {
		return errors.New("下载文件分片长度不正确")
	}
	c.Data = bs
	f.DoneParts[c.Index] = c
	return nil
}

func (f *fileInfo) mergeFileParts() error {
	path := filepath.Join(f.Path, f.Name)
	log.Println("所有分片下载完成, 开始合并文件")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileMd5 := sha256.New()
	totalSize := 0
	for _, part := range f.DoneParts {
		fileMd5.Write(part.Data)
		totalSize += len(part.Data)
		_, err := file.Write(part.Data)
		if err != nil {
			fmt.Printf("error when merge file: %v\n", err)
		}
	}
	if f.Size != totalSize {
		return errors.New("文件大小不正确")
	}
	if f.MD5 != "" {
		md5 := fmt.Sprintf("%x", fileMd5.Sum(nil))
		if f.MD5 != md5 {
			return errors.New("文件MD5不正确")
		}
		log.Println("文件SHA-256校验成功")
	}
	return nil
}
