package storage

import (
	"encoding/json"
	"github.com/alameddinc/ysc/models"
	"os"
	"path"
	"testing"
	"time"
)

func TestReadStorageFile(t *testing.T) {
	testResults := map[string]string{
		"test1": "1",
		"test2": "2",
		"test3": "3",
		"test4": "4",
		"test5": "5",
		"test6": "6",
	}
	wgReadFiles.Add(1)
	go func() {
		for {
			if len(testResults) == 0 {
				break
			}
			select {
			case rawVal := <-bufferJsonValues:
				var tmpCore models.CoreValue
				if err := json.Unmarshal([]byte(rawVal.RawContent), &tmpCore); err != nil {
					continue
				}
				if v, ok := testResults[tmpCore.Key]; !ok || v != tmpCore.Content {
					t.Fatal("Can not Read")
				}
				delete(testResults, tmpCore.Key)
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
	ReadStorageFile("mock.txt")
	t.Log("Success")
	return
}

func TestSync(t *testing.T) {
	SYNC_LIFECYCLE = time.Millisecond
	testSyncFilename := "testsync.txt"
	testContent := time.Now().String()
	testValueHistory := models.ValueHistory{&models.Value{FilenameStamp: testSyncFilename, CoreValue: models.CoreValue{
		Key:     "testCase",
		Content: testContent,
	}}, "w"}
	models.ValueHistories[testSyncFilename] = map[string]*models.ValueHistory{}
	models.ValueHistories[testSyncFilename][testValueHistory.Value.Key] = &testValueHistory
	Sync(true)
	if _, err := os.Stat(path.Join(filesPath, testSyncFilename)); os.IsNotExist(err){
		t.Fatal("Can not Created file")
	}
	if err := os.Remove(path.Join(filesPath, testSyncFilename)); err != nil{
		t.Fatal("Can not Deleted File")
	}
	t.Log("Success")
}
