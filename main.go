package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/h2non/bimg"
)

func main() {

	//set true, if files are only present in the "originals" folder. All files will get created before the comparison starts.
	create_files := false

	//settings for the original files
	original_file_endnung := ".jpg"
	number_of_files := 3

	//define file types
	m := make(map[int]string) //file types for convertion tests
	m[0] = ".jpg"
	m[1] = ".png"
	m[2] = ".webp"
	m[3] = ".tiff"
	m_2 := make(map[int]string) //file types for manipulation tests
	m_2[0] = ".jpg"
	m_2[1] = ".webp"

	//definition of other necessary variables
	var args []string
	benchmark_time_convertion_jpg := make([][4]int, 0) //4 entries: file_endung, file_number, quali, time
	benchmark_time_convertion_webp := make([][4]int, 0)
	benchmark_time_manipulation_jpg := make([][4]int, 0) //4 entries: file_number, quali, test, time
	benchmark_time_manipulation_webp := make([][4]int, 0)

	//create all files necessary from the originals
	if create_files == true {
		for _, file_endung := range m {
			for i := 0; i < number_of_files; i++ {
				if file_endung == ".jpg" {
					args = []string{"./originals/" + strconv.Itoa(i) + original_file_endnung + "[0]", "-background", "white", "-alpha", "remove", "./files_for_comparison/" + strconv.Itoa(i) + file_endung}
				} else {
					args = []string{"./originals/" + strconv.Itoa(i) + original_file_endnung + "[0]", "./files_for_comparison/" + strconv.Itoa(i) + file_endung}
				}
				cmd := exec.Command("convert", args...)
				_, err := cmd.Output()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	//convert all to jpg in different qualities
	for key, file_endung := range m {
		for i := 0; i < number_of_files; i++ {
			for j := 0; j < 20; j++ {
				quali := (j + 1) * 5
				args := []string{"-quality", strconv.Itoa(quali), "./files_for_comparison/" + strconv.Itoa(i) + file_endung + "[0]", "-background", "white", "-alpha", "remove", "./Zwischenspeicher/" + strconv.Itoa(i) + "_converted.jpg"}
				start := time.Now()
				cmd := exec.Command("convert", args...)
				_, err := cmd.Output()
				elapsed := int(time.Since(start))
				if err != nil {
					fmt.Println(err)
				}
				benchmark_time_convertion_jpg = append(benchmark_time_convertion_jpg, [4]int{key, i, j, elapsed})
				err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_converted.jpg")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		fmt.Printf("done with the convertion tests of the %v files to jpg", file_endung)
	}

	//convert all to webp in different qualities
	for key, file_endung := range m {
		for i := 0; i < number_of_files; i++ {
			for j := 0; j < 20; j++ {
				quali := (j + 1) * 5
				args := []string{"-quality", strconv.Itoa(quali), "./files_for_comparison/" + strconv.Itoa(i) + file_endung + "[0]", "-background", "white", "-alpha", "remove", "./Zwischenspeicher/" + strconv.Itoa(i) + "_converted.webp"}
				start := time.Now()
				cmd := exec.Command("convert", args...)
				_, err := cmd.Output()
				elapsed := int(time.Since(start))
				if err != nil {
					fmt.Println(err)
				}
				benchmark_time_convertion_webp = append(benchmark_time_convertion_jpg, [4]int{key, i, j, elapsed})
				err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_converted.webp")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		fmt.Printf("done with the convertion test of the %v files to webp", file_endung)
	}

	//do all image manipulations with the JPGs and WEBPs in different qualities
	for key, file_endung := range m_2 {
		for i := 0; i < number_of_files; i++ {
			for j := 0; j < 20; j++ {
				//create file in new quality
				quali := (j + 1) * 5
				if file_endung == ".jpg" {
					args = []string{"-quality", strconv.Itoa(quali), "./files_for_comparison/" + strconv.Itoa(i) + file_endung + "[0]", "-background", "white", "-alpha", "remove", "./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung}
				} else {
					args = []string{"-quality", strconv.Itoa(quali), "./files_for_comparison/" + strconv.Itoa(i) + file_endung + "[0]", "./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung}
				}
				cmd := exec.Command("convert", args...)
				_, err := cmd.Output()
				if err != nil {
					fmt.Println(err)
				}
				//resize to new width x height
				test_number := 1
				width := 1920
				height := 1920
				start := time.Now()
				buffer, err := bimg.Read("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				newImage, err := bimg.NewImage(buffer).Resize(width, height)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				size, err := bimg.NewImage(newImage).Size()
				if size.Width != width || size.Height != height {
					fmt.Println("The image size is valid")
				}
				bimg.Write("./Zwischenspeicher/"+strconv.Itoa(i)+"_"+strconv.Itoa(quali)+"_temp"+file_endung, newImage)
				elapsed := int(time.Since(start))
				if key == 0 {
					benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_webp, [4]int{i, j, test_number, elapsed})
				}
				if key == 1 {
					benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, test_number, elapsed})
				}
				err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
				if err != nil {
					log.Fatal(err)
				}
				//rotate
				test_number = 2
				start = time.Now()
				buffer, err = bimg.Read("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				newImage, err = bimg.NewImage(buffer).Rotate(90)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				bimg.Write("./Zwischenspeicher/"+strconv.Itoa(i)+"_"+strconv.Itoa(quali)+"_temp"+file_endung, newImage)
				elapsed = int(time.Since(start))
				if key == 0 {
					benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_webp, [4]int{i, j, test_number, elapsed})
				}
				if key == 1 {
					benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, test_number, elapsed})
				}
				err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
				if err != nil {
					log.Fatal(err)
				}
				//force resize
				test_number = 3
				start = time.Now()
				width = 1000
				height = 500
				buffer, err = bimg.Read("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				newImage, err = bimg.NewImage(buffer).ForceResize(width, height)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				size_force, err := bimg.Size(newImage)
				if err != nil {
					log.Fatal(err)
				}
				if size_force.Width != width || size_force.Height != height {
					fmt.Fprintln(os.Stderr, "Incorrect image size")
				}
				bimg.Write("./Zwischenspeicher/"+strconv.Itoa(i)+"_"+strconv.Itoa(quali)+"_temp"+file_endung, newImage)
				elapsed = int(time.Since(start))
				if key == 0 {
					benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_webp, [4]int{i, j, test_number, elapsed})
				}
				if key == 1 {
					benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, test_number, elapsed})
				}
				err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
				if err != nil {
					log.Fatal(err)
				}
				//delete file
				err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		fmt.Printf("done with the manipulation tests of the %v files", file_endung)
	}

	//write alle results to files
	file, err := os.OpenFile("./results/benchmark_time_convertion_jpg"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter := bufio.NewWriter(file)
	for key := range m {
		_, _ = datawriter.WriteString(strconv.Itoa(key) + "\t" + strconv.Itoa(benchmark_time_convertion_jpg[key][0]) + "\t" + strconv.Itoa(benchmark_time_convertion_jpg[key][1]) + "\t" + strconv.Itoa(benchmark_time_convertion_jpg[key][2]) + "\t" + strconv.Itoa(benchmark_time_convertion_jpg[key][3]) + "\n")
	}
	datawriter.Flush()
	file.Close()
	file, err = os.OpenFile("./results/benchmark_time_convertion_webp"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter = bufio.NewWriter(file)
	for key := range m {
		_, _ = datawriter.WriteString(strconv.Itoa(key) + "\t" + strconv.Itoa(benchmark_time_convertion_webp[key][0]) + "\t" + strconv.Itoa(benchmark_time_convertion_webp[key][1]) + "\t" + strconv.Itoa(benchmark_time_convertion_webp[key][2]) + "\t" + strconv.Itoa(benchmark_time_convertion_webp[key][3]) + "\n")
	}
	datawriter.Flush()
	file.Close()
	file, err = os.OpenFile("./results/benchmark_time_manipulation_jpg"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter = bufio.NewWriter(file)
	for key := range m {
		_, _ = datawriter.WriteString(strconv.Itoa(key) + "\t" + strconv.Itoa(benchmark_time_manipulation_jpg[key][0]) + "\t" + strconv.Itoa(benchmark_time_manipulation_jpg[key][1]) + "\t" + strconv.Itoa(benchmark_time_manipulation_jpg[key][2]) + "\t" + strconv.Itoa(benchmark_time_manipulation_jpg[key][3]) + "\n")
	}
	datawriter.Flush()
	file.Close()
	file, err = os.OpenFile("./results/benchmark_time_manipulation_webp"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter = bufio.NewWriter(file)
	for key := range m {
		_, _ = datawriter.WriteString(strconv.Itoa(key) + "\t" + strconv.Itoa(benchmark_time_manipulation_webp[key][0]) + "\t" + strconv.Itoa(benchmark_time_manipulation_webp[key][1]) + "\t" + strconv.Itoa(benchmark_time_manipulation_webp[key][2]) + "\t" + strconv.Itoa(benchmark_time_manipulation_webp[key][3]) + "\n")
	}
	datawriter.Flush()
	file.Close()
	fmt.Printf("done")
}
