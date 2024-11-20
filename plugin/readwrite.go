// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package plugin

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/benedictweis/tcpchat-server-go/application"
)

// handleRead is used to read from a reader and return the result on a channel.
func handleRead(ctx context.Context, reader io.Reader, messages chan<- application.MessageResult, sessionID string) {
	bufioReader := bufio.NewReader(reader)
	for {
		line, err := bufioReader.ReadString('\n')
		select {
		case <-ctx.Done():
			return
		case messages <- application.MessageResult{SessionID: sessionID, Message: line, Err: err}:
		}
	}
}

// handleWrite is used to write from a channel to a writer.
func handleWrite(ctx context.Context, writer io.Writer, messages <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messages:
			_, err := io.Copy(writer, bytes.NewBuffer([]byte(message)))
			if err != nil {
				slog.Warn("write error", "err", err)
			}
		}
	}
}
