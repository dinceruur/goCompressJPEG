package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/jpeg"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

var totalSave float64
var totalJPEG float64

func main() {

	path := flag.String("path", ".", "-path=<PATH>")
	quality := flag.Int("quality", 5, "-quality=<0-100>")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(`Going to reduce image quality 100% to `, *quality, `% in the folder : "`, *path, `" .Do you want to continue (y/n)?`, "\n")
	scanner.Scan()

	if scanner.Text() != "y" {
		return
	} else {
		err := start(*path, *quality)
		fmt.Printf("Total JPEG Count is %v, saved %v MB \n", totalJPEG, totalSave)

		if err != nil {
			fmt.Println(err)
		}
	}

}

func start(p string, q int) error {

	err := filepath.Walk(p,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			} else {
				compressMe(path, info, q)
				return nil
			}

		})

	return err
}

func getFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func compressMe(path string, info os.FileInfo, quality int) {

	if info.Mode().IsRegular() {

		f, err := os.Open(path)

		if err != nil {
			fmt.Println(err)
		}

		c, err := getFileContentType(f)

		if err != nil {
			fmt.Println(err)
		}

		err = f.Close()

		if err != nil {
			fmt.Println(err)
		}

		if c == "image/jpeg" {
			f, err := os.Open(path)

			if err != nil {
				fmt.Println(err)
			}

			i, err := jpeg.Decode(f)

			if err != nil {
				fmt.Println(err)
			}

			err = f.Close()

			if err != nil {
				fmt.Println(err)
			}

			o, err := os.Create(fmt.Sprintf("%v", path))

			if err != nil {
				fmt.Println(err)
			}

			err = jpeg.Encode(o, i, &jpeg.Options{Quality: quality})

			if err != nil {
				fmt.Println(err)
			}

			s, _ := o.Stat()
			ns := s.Size()
			err = o.Close()

			if err != nil {
				fmt.Println(err)
			}

			sizeBefore := float64(info.Size())
			sizeAfter := float64(ns)
			totalSave += (sizeBefore - sizeAfter) / 1000000.0
			totalJPEG++
			fmt.Printf("%v is compressed.\t Before: %v KB \t After: %v KB\t Comprasion Level: %%%v \n", path, info.Size()/1000, ns/1000, math.Round(((sizeBefore-sizeAfter)/sizeBefore)*100))
		}
	}

}
