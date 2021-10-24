package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alameddinc/ysc/models"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

var (
	SYNC_LIFECYCLE = 10 * time.Second
)

var wgWriteAll sync.WaitGroup

func Sync(debug bool) {
	for true {
		filesPath := path.Join(filesPath, fileStorename)
		files, err := os.OpenFile(filesPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			return
		}
		writer, closeWriter := createNewWriter(files)
		defer closeWriter()
		models.NextTimestamp = time.Now().Unix()
		log.Println("Running...")
		time.Sleep(SYNC_LIFECYCLE)
		scanner := bufio.NewScanner(files)

		currentFiles := map[string]bool{}
		for scanner.Scan() {
			line := scanner.Text()
			if _, ok := currentFiles[line]; !ok {
				currentFiles[line] = true
			}
		}
		for fn, works := range models.ValueHistories {
			if _, ok := currentFiles[fn]; !ok {
				writer.WriteString(fn + "\n")
			}
			wgWriteAll.Add(1)
			go SyncToFile(works, fn)
			delete(models.ValueHistories, fn)
		}
		wgWriteAll.Wait()
		closeWriter()
		if debug{
			return
		}
	}
}

func SyncToFile(works map[string]*models.ValueHistory, fn string) {
	lineCount := 0
	defer wgWriteAll.Done()
	currentFilepath := path.Join(filesPath, fn)
	newFilePath := path.Join(filesPath, fmt.Sprintf("%s.backup", fn))
	newFile, err := os.OpenFile(newFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	writer, closeWriter := createNewWriter(newFile)
	defer closeWriter()
	// new values
	affectedRowCount := syncNewValues(works, writer)
	file, err := os.OpenFile(currentFilepath, os.O_RDONLY, 0777)
	// When new file
	if err != nil {
		closeWriter()
		os.Rename(newFilePath, currentFilepath)
		return
	}
	affectedRowCount += aviableRows(works, file, writer)
	closeWriter()
	file.Close()
	removeAndMove(currentFilepath, newFilePath, lineCount)
	return

}

func aviableRows(works map[string]*models.ValueHistory, file *os.File, writer *bufio.Writer) int {
	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var tmpCore models.CoreValue
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &tmpCore); err != nil {
			continue
		}

		if _, ok := works[tmpCore.Key]; ok {
			continue
		}
		if _, err := writer.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			log.Printf("%s can not writed on file", tmpCore.Key)
			continue
		}
		lineCount++
	}
	return lineCount
}

func createNewWriter(file *os.File) (*bufio.Writer, func()) {
	writer := bufio.NewWriter(file)
	return writer, func() {
		writer.Flush()
		file.Close()
	}
}

func syncNewValues(works map[string]*models.ValueHistory, writer *bufio.Writer) int {
	lineCount := 0
	for _, v := range works {
		if v.Op == "d" {
			continue
		}
		tmpVal, err := json.Marshal(v.Value.CoreValue)
		if err != nil {
			log.Printf("%s can not converted", v.Value)
			continue
		}
		if _, err := writer.WriteString(fmt.Sprintf("%s\n", string(tmpVal))); err != nil {
			log.Printf("%s can not writed on file", v.Value)
			continue
		}
		lineCount++
	}
	return lineCount
}

func removeAndMove(c string, n string, lineCount int) {
	os.Remove(c)
	if lineCount == 0 {
		os.Remove(n)
		return
	}
	os.Rename(n, c)
}
