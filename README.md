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
go install github.com/yourusername/fastfetch
```

## Usage
To use FastFetch, you can import it in your Go code:

```bash
fastfetch -o ds.zip -n=100 -u="http://example.com/file"
```
## License
FastFetch is released under the MIT License. See the LICENSE file for more details.