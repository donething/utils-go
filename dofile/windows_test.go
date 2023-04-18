package dofile

import (
	"fmt"
	"github.com/donething/utils-go/dotext"
	"testing"
)

func TestGetDriverInfo(t *testing.T) {
	free, total, avail, err := GetDriverInfo("")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(dotext.BytesHumanReadable(free),
		dotext.BytesHumanReadable(total), dotext.BytesHumanReadable(avail))
}
