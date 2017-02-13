package main

import (
	"fmt"
	"os/exec"

	"os"
	"sync"
	"io/ioutil"
	"strings"
	"path"
)

func main() {
	var dir string
	if len(os.Args) == 2 {
		dir = os.Args[1]
	} else {
		dir = "."
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	var filesToConvert []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.Contains(file.Name(), ".wav") {
			fmt.Printf("File '%s' is not a wav file. \n", file.Name())
			continue
		}

		filePath := path.Join(dir, file.Name())
		if _, err := os.Stat(filePath); err != nil {
			fmt.Printf("File '%s' does not exist. \n", file.Name())
			continue
		}

		if _, err := os.Stat(strings.Replace(filePath, ".wav", ".mp3", 1)); err == nil {
			fmt.Printf("There is already a .mp3 file for %s \n", file.Name())
			continue
		}

		filesToConvert = append(filesToConvert, filePath)
		fmt.Printf("Added '%s' to queue. \n", file.Name())
	}

	if len(filesToConvert) == 0 {
		fmt.Println("Nothing to do.")
		return
	}

	fmt.Println("")
	fmt.Println("###########################################################")
	fmt.Println("###########################################################")
	fmt.Println("###########################################################")
	fmt.Println("")

	var wg sync.WaitGroup
	wg.Add(len(filesToConvert))
	for _, p := range filesToConvert {
		go convert(p, &wg)
	}

	wg.Wait()

	fmt.Println("")
	fmt.Println("###########################################################")
	fmt.Println("")

	fmt.Println("All files converted.")
}

func convert(path string, wg *sync.WaitGroup) {
	cmd := exec.Command("lame", "--silent", "-h", "-V2", path)

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Build:", err)
		fmt.Println("Build:", string(output))
		fmt.Println("Aborting.")
		os.Exit(1)
	}

	fmt.Printf("Done converting '%s' \n", path)
	wg.Done()
}
