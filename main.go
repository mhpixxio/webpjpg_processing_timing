package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/h2non/bimg"
)

func main() {

	//set true, if files are only present in the "originals" folder. All files will get created before the comparison starts.
	create_files := false

	//define file types
	var args []string
	original_file_endnung := ".jpg"
	number_of_files := 42
	m := make(map[int]string)
	m[0] = ".jpg"
	m[1] = ".png"
	m[2] = ".webp"
	m[3] = ".tiff"
	m[4] = ".pdf"

	//create all files necessary
	if create_files == true {
		for _, file_endung := range m {
			for i := 0; i < number_of_files; i++ {
				if file_endung == ".jpg" {
					args = []string{"./originals/" + strconv.Itoa(i) + original_file_endnung + "[0]", "-background", "white", "-alpha", "remove", "./files_for_comparison/" + strconv.Itoa(i) + file_endung}
				} else {
					args = []string{"./originals/" + strconv.Itoa(i) + original_file_endnung + "[0]", "./files_for_comparison/" + strconv.Itoa(i) + file_endung}
				}
				cmd := exec.Command("magick", args...)
				_, err := cmd.Output()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	//start comparison
	for _, file_endung := range m {
		for i := 0; i < number_of_files; i++ {
			//resize
			buffer, err := bimg.Read("./files_for_comparison/" + strconv.Itoa(i) + file_endung)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			newImage, err := bimg.NewImage(buffer).Resize(800, 600)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			size, err := bimg.NewImage(newImage).Size()
			if size.Width == 800 && size.Height == 600 {
				fmt.Println("The image size is valid")
			}

			bimg.Write("./Zwischenspeicher/"+strconv.Itoa(i)+file_endung, newImage)

			//rotate

			//force resize

			//watermark

			//convert
		}
	}

}
