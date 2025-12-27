package logger

import (
	"fmt"
	"sync"
)

// Logger handles broadcasting logs to subscribers (frontend) and console
type Logger struct {
	mu          sync.Mutex
	subscribers []chan string
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		subscribers: make([]chan string, 0),
	}
}

// Printf formats and logs a message to console and all subscribers
func (l *Logger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	
	// 1. Print to standard console (for debugging)
	fmt.Println(msg)

	// 2. Broadcast to all connected frontend clients
	l.mu.Lock()
	defer l.mu.Unlock()
	
	for _, ch := range l.subscribers {
		select {
		case ch <- msg:
		default:
			// If client is too slow, drop message to prevent blocking automation
		}
	}
}

// Subscribe returns a channel that receives log messages
func (l *Logger) Subscribe() chan string {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	ch := make(chan string, 100)
	l.subscribers = append(l.subscribers, ch)
	return ch
}

// Unsubscribe removes a subscriber (simplified for this use case)
// In a full prod app, we'd track IDs to remove specific channels
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, ch := range l.subscribers {
		close(ch)
	}
	l.subscribers = nil
}
