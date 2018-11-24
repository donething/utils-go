package file

import (
	"testing"
)

func TestGet(t *testing.T) {
	file, err := Get("/home/doneth/MyData")
	if err != nil {
		t.Fatal(err)
	}
	filesList, err := file.List(DIR)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v\n", filesList)
}

func TestDFile_ListPaths(t *testing.T) {
	file, err := Get("D:/Temp")
	if err != nil {
		t.Fatal(err)
	}
	pathsList, err := file.List(DIR)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v\n", pathsList[0])
}
