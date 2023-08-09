# FastFetch

FastFetch is a concurrent file downloader written in Go. It splits a file into multiple parts and downloads them concurrently, which can significantly improve the download speed.

## Features

- Concurrent download: FastFetch splits a file into multiple parts and downloads them concurrently.
- Progress bar: FastFetch provides a progress bar to show the download progress.
- Separation of downloading and writing files: After concurrent downloading, data is written to the channel and written to the file according to the position of each block.
- Truncate: Using the Truncate method does not generate temporary files

## Installation

To install FastFetch, you need to have Go installed on your machine, then run the following command:

```bash
go install github.com/zhqqqy/fastfetch@latest
```

## Usage
To use FastFetch, you can import it in your Go code:

```bash
fastfetch -h
NAME:
   fastfetch - File concurrency download

USAGE:
   fastfetch [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url URL, -u URL               URL to download
   --output filename, -o filename  Output filename
   --max-connect value, -n value   Specify maximum number of connections (default: 10)
   --help, -h                      show help

fastfetch -o ds.zip -n=100 -u="http://example.com/file"
```
## License
FastFetch is released under the MIT License. See the LICENSE file for more details.