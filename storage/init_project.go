package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/alameddinc/ysc/models"
	"os"
	"path"
	"sync"
	"time"
)

const filesPath = "./temp/"
const fileStorename = "files.txt"

var fileMap map[string]bool
var wgReadFiles sync.WaitGroup
var wgReadItems sync.WaitGroup
var bufferJsonValues chan models.RawValue
var bufferValueArr []models.Value

func init() {
	fileMap = make(map[string]bool)
	bufferJsonValues = make(chan models.RawValue)
	bufferValueArr = []models.Value{}
}

// ReadFileStorage reads filenames on files.txt
func ReadFileStorage() error {
	ctx, cancel := context.WithCancel(context.Background())
	go channelListener(ctx)
	files, err := fetchFileList()
	if err != nil {
		return err
	}
	for _, file := range files {
		fileMap[file] = true
		wgReadFiles.Add(1)
		go ReadFile(file)
	}
	wgReadFiles.Wait()
	cancel()
	wgReadItems.Wait()
	return nil
}

// ReadFile read file and fetch values
func ReadFile(filename string) error {
	defer wgReadFiles.Done()
	var file, err = os.OpenFile(path.Join(filesPath, filename), os.O_RDONLY, 0777)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m := scanner.Text()
		bufferJsonValues <- models.RawValue{m, filename}
	}
	return nil

}

// channelListener works on Goroutines for fetch all values
func channelListener(ctx context.Context) {
	wgReadItems.Add(1)
	defer wgReadItems.Done()
	for {
		select {
		case rawVal := <-bufferJsonValues:
			var tmpCore models.CoreValue
			if err := json.Unmarshal([]byte(rawVal.RawContent), &tmpCore); err != nil {
				continue
			}
			models.CachedValues[tmpCore.Key] = &models.Value{CoreValue: tmpCore, FilenameStamp: rawVal.Filename}
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// fetchFileList works for fetch all filenames
func fetchFileList() ([]string, error) {
	fileList := []string{}
	var file, err = os.OpenFile(path.Join(filesPath, fileStorename), os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileList = append(fileList, scanner.Text())
	}
	return fileList, scanner.Err()
}
