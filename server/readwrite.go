package server

import (
	"bufio"
	"bytes"
	"context"
	"io"
)

// handleRead is used to read from a reader and return the result on a channel
func handleRead(ctx context.Context, reader io.Reader, messages chan<- MessageResult, sessionId string) {
	bufioReader := bufio.NewReader(reader)
	for {
		line, err := bufioReader.ReadString('\n')
		select {
		case <-ctx.Done():
			return
		case messages <- MessageResult{sessionId, line, err}:
		}
	}
}

// handleWrite is used to write from a channel to a writer
func handleWrite(ctx context.Context, writer io.Writer, messages <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messages:
			io.Copy(writer, bytes.NewBuffer([]byte(message)))
		}
	}
}
