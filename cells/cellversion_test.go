package cells

import (
	"bytes"
	"testing"
)
import "encoding/hex"

func TestCanReadAndWriteCellVersion(t *testing.T) {
	versionsBytes, _ := hex.DecodeString("000300040005")
	readBuffer := bytes.NewBuffer(versionsBytes)

	cellVersion, deserializeErr := ReadCellVersion(readBuffer)
	if deserializeErr != nil {
		t.Errorf("cell version deserialize failed: %s", deserializeErr)
	}

	writeBuffer := new(bytes.Buffer)
	serializeErr := cellVersion.Write(writeBuffer)
	if serializeErr != nil {
		t.Errorf("cell version serialize failed: %s", serializeErr)
	}

	if !bytes.Equal(versionsBytes, writeBuffer.Bytes()) {
		t.Errorf("input and output byte arrays are not equal")
	}
}
