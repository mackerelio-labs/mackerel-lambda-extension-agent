// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-collections/go-datastructures/queue"
)

const defaultListenerPort = "4323"
const initialQueueSize = 5

// Used to listen to the Telemetry API
type TelemetryApiListener struct {
	httpServer *http.Server
	// LogEventsQueue is a synchronous queue and is used to put the received log events to be dispatched later
	LogEventsQueue *queue.Queue
	isSAMLocal     bool
}

func NewTelemetryApiListener(isSAMLocal bool) *TelemetryApiListener {
	return &TelemetryApiListener{
		httpServer:     nil,
		LogEventsQueue: queue.New(initialQueueSize),
		isSAMLocal:     isSAMLocal,
	}
}

func (s *TelemetryApiListener) listenOnAddress() string {
	var addr string
	if s.isSAMLocal {
		addr = ":" + defaultListenerPort
	} else {
		addr = "sandbox:" + defaultListenerPort
	}

	return addr
}

// Starts the server in a goroutine where the log events will be sent
func (s *TelemetryApiListener) Start() (string, error) {
	address := s.listenOnAddress()
	Logger.Info("Starting on address", address)
	s.httpServer = &http.Server{Addr: address}
	http.HandleFunc("/", s.http_handler)
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			Logger.Error("Unexpected stop on Http Server:", err)
			s.Shutdown()
		} else {
			Logger.Info("Http Server closed:", err)
		}
	}()
	return fmt.Sprintf("http://%s/", address), nil
}

// http_handler handles the requests coming from the Telemetry API.
// Everytime Telemetry API sends log events, this function will read them from the response body
// and put into a synchronous queue to be dispatched later.
// Logging or printing besides the error cases below is not recommended if you have subscribed to
// receive extension logs. Otherwise, logging here will cause Telemetry API to send new logs for
// the printed lines which may create an infinite loop.
func (s *TelemetryApiListener) http_handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Logger.Warning("Error reading body:", err)
		return
	}

	// Parse and put the log messages into the queue
	var slice []interface{}
	_ = json.Unmarshal(body, &slice)

	for _, el := range slice {
		s.LogEventsQueue.Put(el)
	}

	Logger.Info("logEvents received:", len(slice), " LogEventsQueue length:", s.LogEventsQueue.Len())
	slice = nil
}

// Terminates the HTTP server listening for logs
func (s *TelemetryApiListener) Shutdown() {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			Logger.Warning("Failed to shutdown http server gracefully:", err)
		} else {
			s.httpServer = nil
		}
	}
}
