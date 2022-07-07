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
	create_files := true
	//set true, if the convertion time tests shoufiles_for_comparisonld be performed (take the longest)
	convertion_test := true
	//set true, if the image manipulation time tests should be performed
	manipulation_test := true

	//settings for the original files
	original_file_endung := ".jpg"
	number_of_files_compared := 40
	number_of_files_max := 40

	//define file types
	m := make(map[int]string) //file types for convertion tests
	m[0] = ".jpg"
	m[1] = ".png"
	m[2] = ".webp"
	m[3] = ".tiff"
	m[4] = ".heic"
	m_2 := make(map[int]string) //file types for manipulation tests
	m_2[0] = ".jpg"
	m_2[1] = ".webp"

	//definition of other necessary variables
	quality_steps := 20
	var number_of_tests int //gets filled automatically
	var args []string
	benchmark_time_convertion_jpg := make([][4]int, 0) //4 entries: file_endung, file_number, quali, time
	benchmark_time_convertion_webp := make([][4]int, 0)
	benchmark_time_manipulation_jpg := make([][4]int, 0) //4 entries: file_number, quali, test_number, time
	benchmark_time_manipulation_webp := make([][4]int, 0)
	//formatted, takes the average over all pictures
	benchmark_time_convertion_jpg_formatted := make([][3]int, 0) //3 entries: file_endung, quali, time
	benchmark_time_convertion_webp_formatted := make([][3]int, 0)
	benchmark_time_manipulation_jpg_formatted := make([][3]int, 0) //3 entries: test_number, quali, time
	benchmark_time_manipulation_webp_formatted := make([][3]int, 0)

	//create all files necessary from the originals
	os.Mkdir("./Zwischenspeicher/", os.ModePerm)
	if create_files == true {
		for _, file_endung := range m {
			for i := 0; i < number_of_files_max; i++ {
				if file_endung == ".jpg" {
					args = []string{"./originals/" + strconv.Itoa(i) + original_file_endung + "[0]", "-background", "white", "-alpha", "remove", "./files_for_comparison/" + strconv.Itoa(i) + file_endung}
				} else {
					args = []string{"./originals/" + strconv.Itoa(i) + original_file_endung + "[0]", "./files_for_comparison/" + strconv.Itoa(i) + file_endung}
				}
				cmd := exec.Command("convert", args...)
				_, err := cmd.Output()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	//convertion time tests
	if convertion_test == true {
		//convert all to jpg in different qualities
		for key, file_endung := range m {
			for i := 0; i < number_of_files_compared; i++ {
				for j := 0; j < quality_steps; j++ {
					quali := (j + 1) * int(100/quality_steps)
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
			fmt.Printf("done with the convertion tests of the %v files to .jpg\n", file_endung)
		}
		//convert all to webp in different qualities
		for key, file_endung := range m {
			for i := 0; i < number_of_files_compared; i++ {
				for j := 0; j < quality_steps; j++ {
					quali := (j + 1) * int(100/quality_steps)
					args := []string{"-quality", strconv.Itoa(quali), "./files_for_comparison/" + strconv.Itoa(i) + file_endung + "[0]", "-background", "white", "-alpha", "remove", "./Zwischenspeicher/" + strconv.Itoa(i) + "_converted.webp"}
					start := time.Now()
					cmd := exec.Command("convert", args...)
					_, err := cmd.Output()
					elapsed := int(time.Since(start))
					if err != nil {
						fmt.Println(err)
					}
					benchmark_time_convertion_webp = append(benchmark_time_convertion_webp, [4]int{key, i, j, elapsed})
					err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_converted.webp")
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			fmt.Printf("done with the convertion test of the %v files to .webp\n", file_endung)
		}
		//format the results for benchmarking (takes the average over all images)
		for i := 0; i < len(m); i++ {
			for j := 0; j < quality_steps; j++ {
				formatted_entry := 0
				for k := 0; k < number_of_files_compared; k++ {
					line_nr := k*quality_steps + j + i*quality_steps*number_of_files_compared
					formatted_entry = formatted_entry + benchmark_time_convertion_jpg[line_nr][3]
				}
				formatted_entry = formatted_entry / number_of_files_compared
				benchmark_time_convertion_jpg_formatted = append(benchmark_time_convertion_jpg_formatted, [3]int{i, (j + 1) * int(100/quality_steps), formatted_entry})
			}
		}
		for i := 0; i < len(m); i++ {
			for j := 0; j < quality_steps; j++ {
				formatted_entry := 0
				for k := 0; k < number_of_files_compared; k++ {
					line_nr := k*quality_steps + j + i*quality_steps*number_of_files_compared
					formatted_entry = formatted_entry + benchmark_time_convertion_webp[line_nr][3]
				}
				formatted_entry = formatted_entry / number_of_files_compared
				benchmark_time_convertion_webp_formatted = append(benchmark_time_convertion_webp_formatted, [3]int{i, (j + 1) * int(100/quality_steps), formatted_entry})
			}
		}
		//write all results to files
		file, err := os.OpenFile("./results/benchmark_time_convertion_jpg"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter := bufio.NewWriter(file)
		for key := 0; key < len(benchmark_time_convertion_jpg_formatted); key++ {
			_, _ = datawriter.WriteString(strconv.Itoa(benchmark_time_convertion_jpg_formatted[key][0]) + "\t" + strconv.Itoa(benchmark_time_convertion_jpg_formatted[key][1]) + "\t" + strconv.Itoa(benchmark_time_convertion_jpg_formatted[key][2]) + "\n")
		}
		datawriter.Flush()
		file.Close()
		file, err = os.OpenFile("./results/benchmark_time_convertion_webp"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter = bufio.NewWriter(file)
		for key := 0; key < len(benchmark_time_convertion_webp_formatted); key++ {
			_, _ = datawriter.WriteString(strconv.Itoa(benchmark_time_convertion_webp_formatted[key][0]) + "\t" + strconv.Itoa(benchmark_time_convertion_webp_formatted[key][1]) + "\t" + strconv.Itoa(benchmark_time_convertion_webp_formatted[key][2]) + "\n")
		}
		datawriter.Flush()
		file.Close()
	}

	//do all image manipulations with jpg and webp in different qualities
	if manipulation_test == true {
		for key, file_endung := range m_2 {
			for i := 0; i < number_of_files_compared; i++ {
				for j := 0; j < quality_steps; j++ {
					number_of_tests = 0
					//create file in new quality
					quali := (j + 1) * int(100/quality_steps)
					if file_endung == ".jpg" {
						args = []string{"-quality", strconv.Itoa(quali), "./files_for_comparison/" + strconv.Itoa(i) + file_endung + "[0]", "-background", "white", "-alpha", "remove", "-resize", "1920x1920", "./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung}
					} else {
						args = []string{"-quality", strconv.Itoa(quali), "./files_for_comparison/" + strconv.Itoa(i) + file_endung + "[0]", "-resize", "1920x1920", "./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung}
					}
					cmd := exec.Command("convert", args...)
					_, err := cmd.Output()
					if err != nil {
						fmt.Println(err)
					}
					//do all tests
					//test 0 resize to new width x height
					width := 500
					height := 500
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
						benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_jpg, [4]int{i, j, number_of_tests, elapsed})
					}
					if key == 1 {
						benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, number_of_tests, elapsed})
					}
					os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
					number_of_tests++
					//test 1 resize to new width x height
					width = 860
					height = 860
					start = time.Now()
					buffer, err = bimg.Read("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
					newImage, err = bimg.NewImage(buffer).Resize(width, height)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
					size, err = bimg.NewImage(newImage).Size()
					if size.Width != width || size.Height != height {
						fmt.Println("The image size is valid")
					}
					bimg.Write("./Zwischenspeicher/"+strconv.Itoa(i)+"_"+strconv.Itoa(quali)+"_temp"+file_endung, newImage)
					elapsed = int(time.Since(start))
					if key == 0 {
						benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_jpg, [4]int{i, j, number_of_tests, elapsed})
					}
					if key == 1 {
						benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, number_of_tests, elapsed})
					}
					os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
					number_of_tests++
					//test 2 resize to new width x height
					width = 1200
					height = 1200
					start = time.Now()
					buffer, err = bimg.Read("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
					newImage, err = bimg.NewImage(buffer).Resize(width, height)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
					size, err = bimg.NewImage(newImage).Size()
					if size.Width != width || size.Height != height {
						fmt.Println("The image size is valid")
					}
					bimg.Write("./Zwischenspeicher/"+strconv.Itoa(i)+"_"+strconv.Itoa(quali)+"_temp"+file_endung, newImage)
					elapsed = int(time.Since(start))
					if key == 0 {
						benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_jpg, [4]int{i, j, number_of_tests, elapsed})
					}
					if key == 1 {
						benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, number_of_tests, elapsed})
					}
					os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
					number_of_tests++
					//test 3 rotate
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
						benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_jpg, [4]int{i, j, number_of_tests, elapsed})
					}
					if key == 1 {
						benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, number_of_tests, elapsed})
					}
					os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
					number_of_tests++
					//test 4 force resize
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
						benchmark_time_manipulation_jpg = append(benchmark_time_manipulation_jpg, [4]int{i, j, number_of_tests, elapsed})
					}
					if key == 1 {
						benchmark_time_manipulation_webp = append(benchmark_time_manipulation_webp, [4]int{i, j, number_of_tests, elapsed})
					}
					os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + "_temp" + file_endung)
					number_of_tests++

					//delete file after image manipulation tests
					err = os.Remove("./Zwischenspeicher/" + strconv.Itoa(i) + "_" + strconv.Itoa(quali) + file_endung)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			fmt.Printf("done with the manipulation tests of the %v files\n", file_endung)
		}
		//format the results for benchmarking (takes the average over all images)
		for i := 0; i < number_of_tests; i++ {
			for j := 0; j < quality_steps; j++ {
				formatted_entry := 0
				for k := 0; k < number_of_files_compared; k++ {
					line_nr := k*(number_of_tests*quality_steps) + j*number_of_tests + i
					formatted_entry = formatted_entry + benchmark_time_manipulation_jpg[line_nr][3]
				}
				formatted_entry = formatted_entry / number_of_files_compared
				benchmark_time_manipulation_jpg_formatted = append(benchmark_time_manipulation_jpg_formatted, [3]int{i, (j + 1) * int(100/quality_steps), formatted_entry})
			}
		}
		for i := 0; i < number_of_tests; i++ {
			for j := 0; j < quality_steps; j++ {
				formatted_entry := 0
				for k := 0; k < number_of_files_compared; k++ {
					line_nr := k*(number_of_tests*quality_steps) + j*number_of_tests + i
					formatted_entry = formatted_entry + benchmark_time_manipulation_webp[line_nr][3]
				}
				formatted_entry = formatted_entry / number_of_files_compared
				benchmark_time_manipulation_webp_formatted = append(benchmark_time_manipulation_webp_formatted, [3]int{i, (j + 1) * int(100/quality_steps), formatted_entry})
			}
		}
		//write all results to files
		file, err := os.OpenFile("./results/benchmark_time_manipulation_jpg"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter := bufio.NewWriter(file)
		for key := 0; key < len(benchmark_time_manipulation_jpg_formatted); key++ {
			_, _ = datawriter.WriteString(strconv.Itoa(benchmark_time_manipulation_jpg_formatted[key][0]) + "\t" + strconv.Itoa(benchmark_time_manipulation_jpg_formatted[key][1]) + "\t" + strconv.Itoa(benchmark_time_manipulation_jpg_formatted[key][2]) + "\n")
		}
		datawriter.Flush()
		file.Close()
		file, err = os.OpenFile("./results/benchmark_time_manipulation_webp"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter = bufio.NewWriter(file)
		for key := 0; key < len(benchmark_time_manipulation_webp_formatted); key++ {
			_, _ = datawriter.WriteString(strconv.Itoa(benchmark_time_manipulation_webp_formatted[key][0]) + "\t" + strconv.Itoa(benchmark_time_manipulation_webp_formatted[key][1]) + "\t" + strconv.Itoa(benchmark_time_manipulation_webp_formatted[key][2]) + "\n")
		}
		datawriter.Flush()
		file.Close()
	}
	fmt.Printf("done\n")
}
