package main

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Downloader struct {
	l           sync.Mutex
	concurrency int
	bar         *progressbar.ProgressBar
}

func NewDownloader(concurrency int) *Downloader {
	return &Downloader{concurrency: concurrency}
}

func (d *Downloader) Download(url, filename string, numRoutines int) error {
	file, _ := os.Create(filename)
	length, err := getLength(url)
	if err != nil {
		return err
	}
	file.Truncate(length)
	d.setBar(length)
	//分割任务
	rangeSize := length / int64(numRoutines)

	var wg sync.WaitGroup
	log.Println("并发下载数: ", numRoutines)
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		startRange := int64(i) * rangeSize
		endRange := startRange + rangeSize

		if i == numRoutines-1 {
			endRange = length // 最后一片将结束字节设为文件大小
		}
		go func(start, end int64) {
			defer wg.Done()
			d.downloadRange(file, url, start, end) // 并发下载
		}(startRange, endRange)
	}
	wg.Wait()
	return nil
}

func (d *Downloader) downloadRange(w *os.File, url string, startRange, endRange int64) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	rangeHeader := "bytes=" + fmt.Sprintf("%d-%d", startRange, endRange-1)
	req.Header.Add("Range", rangeHeader)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	buffer := make([]byte, 10240)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read buffer: ", err)
			}
			w.WriteAt(buffer[:n], int64(startRange)) // 保存文件
			d.bar.Add(n)
			break
		}

		w.WriteAt(buffer[:n], int64(startRange)) // 保存文件
		startRange += int64(n)
		d.bar.Add(n)
	}
}
func (d *Downloader) setBar(length int64) {
	d.bar = progressbar.NewOptions64(
		length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("downloading..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func getLength(url string) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	length, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	if err != nil {
		return 0, err
	}
	return length, err
}
