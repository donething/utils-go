package file

import (
	"testing"
)

func TestGet(t *testing.T) {
	file := Get("/home/doneth/MyData")
	filesList, err := file.List(DIR)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v\n", filesList)
}
