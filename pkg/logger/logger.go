package logger

import (
	"fmt"
	"sync"
)

// handles broadcasting logs
type Logger struct {
	mu          sync.Mutex
	subscribers []chan string
}

// creates a new Logger instance
func New() *Logger {
	return &Logger{
		subscribers: make([]chan string, 0),
	}
}

// logs a message to console and all subscribers
func (l *Logger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)

	fmt.Println(msg)

	l.mu.Lock()
	defer l.mu.Unlock()

	for _, ch := range l.subscribers {
		select {
		case ch <- msg:
		default:
		}
	}
}

// Subscribe - a channel that receives log messages
func (l *Logger) Subscribe() chan string {
	l.mu.Lock()
	defer l.mu.Unlock()

	ch := make(chan string, 100)
	l.subscribers = append(l.subscribers, ch)
	return ch
}

// Unsubscribe removes a channel from the subscribers list
func (l *Logger) Unsubscribe(ch chan string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, sub := range l.subscribers {
		if sub == ch {
			l.subscribers = append(l.subscribers[:i], l.subscribers[i+1:]...)
			close(ch)
			break
		}
	}
}

// Close closes all subscriber channels
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, ch := range l.subscribers {
		close(ch)
	}
	l.subscribers = nil
}
