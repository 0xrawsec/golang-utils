package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/0xrawsec/golang-utils/log"
)

type OptionalHeader struct {
	This float32
	Is   float64
	Just uint8
	A    uint32
	Test int64
}

type Header struct {
	Magic          []byte
	Array          [8]byte
	TableLen       uint64
	OptionalHeader OptionalHeader
}

func (h Header) String() string {
	return fmt.Sprintf("Magic: %s Array: %s TableLen: %d", h.Magic, h.Array, h.TableLen)
}

type OffsetTable []uint64

func init() {
	log.InitLogger(log.LDebug)
}

func TestPack(t *testing.T) {
	b := [8]byte{}
	copy(b[:], "Fooobaar")
	h := Header{Magic: []byte("Foobar"), Array: b, TableLen: 1337}
	writter := new(bytes.Buffer)
	if err := binary.Write(writter, binary.LittleEndian, h.TableLen); err != nil {
		t.Error(err)
	}
	t.Logf("Successfully encoded: %q", writter.Bytes())
}

func TestMarshal(t *testing.T) {
	oh := OptionalHeader{3, 4, 1, 0, 10}
	b := [8]byte{}
	copy(b[:], "Fooobaar")
	h := Header{Magic: []byte("Foobar"), TableLen: 1337, Array: b, OptionalHeader: oh}
	enc, err := Marshal(&h, binary.LittleEndian)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Successfully encoded: %q", enc)
}

func TestUnmarshal(t *testing.T) {
	oh := OptionalHeader{3, 4, 1, 0, 10}
	b := [8]byte{}
	copy(b[:], "Fooobaar")
	h := Header{Magic: []byte("Foobar"), TableLen: 1337, Array: b, OptionalHeader: oh}
	enc, err := Marshal(&h, binary.LittleEndian)
	if err != nil {
		t.Error(err)
	}
	reader := bytes.NewReader(enc)
	nh := Header{Magic: []byte("aaa"), TableLen: 42}
	err = Unmarshal(reader, &nh, binary.LittleEndian)
	if err != nil {
		t.Error(err)
	}
	t.Log(nh)
	if reflect.DeepEqual(nh, h) {
		t.Logf("Successfully decoded: %v", nh)
	}
}

func TestUnmarshalInitSlice(t *testing.T) {
	array := []byte{41, 42, 43}
	slice := make([]byte, 3)
	data, err := Marshal(&array, binary.LittleEndian)
	if err != nil {
		t.Error(err)
	}
	t.Logf("data: %q\n", data)
	reader := bytes.NewReader(data)
	len := int64(0)
	err = Unmarshal(reader, &len, binary.LittleEndian)
	if err != nil {
		t.Error(err)
	}
	err = UnmarshaInitSlice(reader, &slice, binary.LittleEndian)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(array[:], slice) {
		t.Error("Test failed")
	}
}
