package arc

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type ChunkedReader struct {
	rd              io.Reader
	chunkSize       int
	buf             bytes.Buffer
	mu              sync.Mutex
	ticker          *time.Ticker
	closed          chan struct{}
	chunk_available chan bool
	err             error
}

func NewChunkedReader(reader io.Reader, interval time.Duration, chunkSize int) *ChunkedReader {
	r := ChunkedReader{
		rd:        reader,
		chunkSize: chunkSize,
		ticker:    time.NewTicker(interval),
		closed:    make(chan struct{}),
	}

	go r.fill()

	return &r
}

func (c *ChunkedReader) Read() (chunk []byte, err error) {
	//immediately return if there is a full chunk available
	c.mu.Lock()
	if c.buf.Len() > c.chunkSize {
		c.mu.Unlock()
		return c.chunk(), nil
	}
	c.mu.Unlock()

	for {
		select {
		case <-c.ticker.C:
			if b := c.chunk(); b != nil {
				return b, nil
			}
		case <-c.closed:
			b := c.chunk()
			if c.buf.Len() == 0 {
				//return the error when there is nothing left
				return b, c.err
			}
			return b, nil
		case <-c.chunk_available:
			return c.chunk(), nil
		}
	}

}

func (c *ChunkedReader) chunk() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.buf.Len() > 0 {
		n := c.chunkSize
		if c.buf.Len() < n {
			n = c.buf.Len()
		} else {
			//if we have a full chunk look in the last
			//10% of the buffer for a newline and split there
			//this gives us nice newline terminated chunks
			//in most cases
			buf := c.buf.Bytes()[:c.chunkSize]
			lookBack := len(buf) - 1 - c.chunkSize/10
			for i := len(buf) - 1; i > lookBack; i-- {
				if buf[i] == '\n' {
					n = i + 1
					break
				}
			}
		}
		chunk := make([]byte, n)
		c.buf.Read(chunk)
		return chunk

	}
	return nil
}

func (c *ChunkedReader) fill() {
	defer close(c.closed)
	tmpBuf := make([]byte, 512)
	for {
		n, err := c.rd.Read(tmpBuf)
		if n > 0 {
			c.mu.Lock()
			c.buf.Write(tmpBuf[:n])
			c.mu.Unlock()
		}
		if err != nil {
			c.err = err
			return
		}
	}

}
