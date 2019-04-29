package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	path := flag.String("path", ".", "-path=<PATH>")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(`Going to reduce image sizes in the folder : "`, *path, `" .Do you want to continue (y/n)?`, "\n")
	scanner.Scan()

	if scanner.Text() != "y" {
		return
	} else {
		err := start(*path)

		if err != nil {
			fmt.Println(err)
		}
	}

}

func start(p string) error {

	err := filepath.Walk(p,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			} else {
				return compressMe(path, info)
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

func compressMe(path string, info os.FileInfo) error {

	if info.Mode().IsRegular() {

		f, err := os.Open(path)

		if err != nil {
			return err
		}

		c, err := getFileContentType(f)

		if err != nil {
			return err
		}

		err = f.Close()

		if err != nil {
			return err
		}

		if c == "image/jpeg" {
			f, err := os.Open(path)

			if err != nil {
				return err
			}

			i, err := jpeg.Decode(f)

			if err != nil {
				return err
			}

			err = f.Close()

			if err != nil {
				return err
			}

			o, err := os.Create(fmt.Sprintf("%v", path))

			if err != nil {
				return err
			}

			err = jpeg.Encode(o, i, &jpeg.Options{Quality: 30})

			if err != nil {
				return err
			}

			s, _ := o.Stat()
			ns := s.Size()
			err = o.Close()

			if err != nil {
				return err
			}

			fmt.Printf("%v is compressed.\t Before: %v KB \t After: %v KB\t Comprasion Level: %%%v \n", path, info.Size()/1024, ns/1024, (ns/1024)*100/(info.Size()/1024))
		}
	}

	return nil
}
