package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	concurrency int
	bar         *progressbar.ProgressBar
}
type Data struct {
	Bytes []byte
	Start int
}

func NewDownloader(concurrency int) *Downloader {
	return &Downloader{concurrency: concurrency}
}

func (d *Downloader) Download(strURL, filename string) error {
	// if you use http.Head to access oss links with expiration time in alibaba cloud,
	// you will report 403, so use http.Get
	resp, err := http.Get(strURL)
	if err != nil {
		return err
	}

	length, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// preallocate file
	out.Truncate(int64(length))
	// create progress bar
	d.setBar(length)

	partSize := length / d.concurrency

	dataCh := make(chan Data, d.concurrency)

	// start a goroutine to write data to file
	var wgp sync.WaitGroup
	wgp.Add(1)
	go func() {
		defer wgp.Done()
		for data := range dataCh {
			out.Seek(int64(data.Start), 0)
			out.Write(data.Bytes)
			d.bar.Add(len(data.Bytes))
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < d.concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			start := partSize * i
			end := start + partSize
			if i == d.concurrency-1 {
				end = length
			}

			req, err := http.NewRequest("GET", strURL, nil)
			if err != nil {
				log.Fatal(err)
			}
			// end needs to be reduced by 1 because the data is downloaded starting from 0 bytes
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			// read the data and send it to the channel
			bytes := make([]byte, end-start)
			_, err = io.ReadFull(resp.Body, bytes)
			if err != nil {
				panic(err)
			}
			dataCh <- Data{Bytes: bytes, Start: start}
		}(i)
	}
	wg.Wait()
	close(dataCh)
	wgp.Wait()
	log.Println("下载完成")
	return nil
}

func (d *Downloader) setBar(length int) {
	d.bar = progressbar.NewOptions(
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
