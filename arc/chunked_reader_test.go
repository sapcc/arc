// +build !integration

package arc

import (
	"io"
	"testing"
	"time"
)

type mockReader struct {
	Data  []string
	Delay time.Duration
}

func (r *mockReader) Read(p []byte) (int, error) {
	time.Sleep(r.Delay)
	if len(r.Data) == 0 {
		return 0, io.EOF
	}
	chunk := []byte(r.Data[0])
	if len(p) < len(chunk) {
		panic("given byte slice to small for mockReader")
	}
	copy(p, chunk)
	r.Data = r.Data[1:]
	return len(chunk), nil

}

func TestChunkedReader(t *testing.T) {

	rd := mockReader{Data: []string{"read1", "read2"}}

	chunker := NewChunkedReader(&rd, 1*time.Second, 256)

	chunk, err := chunker.Read()

	if err != io.EOF {
		t.Error("Expected io.EOF error")
	}

	expected_chunk := "read1read2"
	received_chunk := string(chunk)

	if received_chunk != expected_chunk {
		t.Errorf("Read was not properly chunked. Expected: %s, Got: %s", expected_chunk, received_chunk)
	}

}

func TestChunkSize(t *testing.T) {
	rd := mockReader{Data: []string{"read1", "read2"}}

	chunker := NewChunkedReader(&rd, 1*time.Second, 4)

	chunk, err := chunker.Read()

	if err != nil {
		t.Error("Unexpected err")
	}
	if string(chunk) != "read" {
		t.Error("Unexpected chunk")
	}

	chunk, err = chunker.Read()
	if err != nil {
		t.Error("Unexpected err")
	}
	if string(chunk) != "1rea" {
		t.Error("Unexpected chunk")
	}

	chunk, err = chunker.Read()
	if err == nil {
		t.Error("Expected err")
	}
	if string(chunk) != "d2" {
		t.Error("Unexpected chunk")
	}

}

func TestTimedChunk(t *testing.T) {
	rd := mockReader{Data: []string{"read1", "read2"}, Delay: 50 * time.Millisecond}

	chunker := NewChunkedReader(&rd, 10*time.Millisecond, 256)

	chunk, err := chunker.Read()

	if err != nil {
		t.Error("Unexpected err")
	}
	if string(chunk) != "read1" {
		t.Error("Unexpected chunk")
	}

}

func TestChunkAtNewLine(t *testing.T) {
	rd := mockReader{Data: []string{"123456789012345678\n90"}}
	chunker := NewChunkedReader(&rd, 10*time.Millisecond, 20)

	chunk, _ := chunker.Read()
	if string(chunk) != "123456789012345678\n" {
		t.Errorf("Unexpected chunk [%s]", string(chunk))
	}

}
