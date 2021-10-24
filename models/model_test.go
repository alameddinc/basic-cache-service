package models

import "testing"

// test file should be change to multiple when project will grow

func TestCreateValue(t *testing.T) {
	if v, err := CreateValue("testKey", "testVal"); err != nil || v.Content == "testVal" {
		if v2, err := CreateValue("testKey", "testVal2"); err != nil || v != v2 {
			t.Fatal("Not Updated Value")
		}
		t.Log("Success")
		return
	}
	t.Fatal("Not Created Value")
}

func TestFetchValue(t *testing.T) {
	if _, err := CreateValue("testKey2", "testVal"); err == nil {
		if v, err := FetchValue("testKey2"); err != nil || v.Content != "testVal" {
			t.Fatal("Not Fetch Value")
		}
		t.Log("Success")
		return
	}
	t.Fatal("Not Created Value")
}

func TestValue_Set(t *testing.T) {
	if v, err := CreateValue("testKey3", "testVal"); err == nil {
		if v.Set("testValNew"); err != nil || v.Content != "testValNew" {
			t.Fatal("Not Set Value")
		}
		if _, ok := CachedValues["testKey3"]; !ok {
			t.Fatal("Not Set Value")
		}
		if _, ok := ValueHistories[v.FilenameStamp]; !ok{
			t.Fatal("Not saved set event on history #1")
		}
		if history, ok := (ValueHistories[v.FilenameStamp])[v.Key]; !ok || history.Op != "w"{
			t.Log(history)
			t.Fatal("Not saved set event on history #2")
		}
		t.Log("Success")
		return
	}
	t.Fatal("Not Created Value")
}

func TestValue_Get(t *testing.T) {
	if v, err := CreateValue("testKey4", "testVal"); err == nil {
		if val := v.Get(); val != "testVal"{
			t.Fatal("Not Get Value")
		}
		t.Log("Success")
		return
	}
	t.Fatal("Not Created Value")
}

func TestValue_Delete(t *testing.T) {
	if v, err := CreateValue("testKey3", "testVal"); err == nil {
		if v.Delete(); err != nil {
			t.Fatal("Not Delete Value")
		}
		if _, ok := CachedValues["testKey3"]; ok {
			t.Fatal("Not Delete Value")
		}
		if _, ok := ValueHistories[v.FilenameStamp]; !ok{
			t.Fatal("Not saved delete event on history #1")
		}
		if history, ok := (ValueHistories[v.FilenameStamp])[v.Key]; !ok || history.Op != "d"{
			t.Fatal("Not saved delete event on history #2")
		}
		t.Log("Success")
		return
	}
	t.Fatal("Not Created Value")
}
