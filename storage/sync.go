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
		filesFullPath := path.Join(filesPath, fileStorename)
		files, err := os.OpenFile(filesFullPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
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
			_, err := os.Stat(path.Join(filesPath, line))
			if _, ok := currentFiles[line]; !ok && err == nil {
				currentFiles[line] = true
			}
		}
		for fn, works := range models.ValueHistories {
			fmt.Println(models.ValueHistories)
			if _, ok := currentFiles[fn]; !ok {
				currentFiles[fn] = true
			}
			wgWriteAll.Add(1)
			go SyncToFile(works, fn)
			delete(models.ValueHistories, fn)
		}
		os.Truncate(path.Join(filesPath, fileStorename), 0)
		for fn, _ := range currentFiles {
			writer.WriteString(fn + "\n")
		}
		wgWriteAll.Wait()
		closeWriter()
		if debug {
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

// syncNewValues works for availrable values on ValueHistories
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

// syncNewValues works for new values on ValueHistories
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

// removeAndMove delete newFile with currentFile when lineCount is zero
func removeAndMove(c string, n string, lineCount int) {
	os.Remove(c)
	if lineCount == 0 {
		os.Remove(n)
		return
	}
	os.Rename(n, c)
}
