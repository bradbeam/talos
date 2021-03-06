package log

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	filechunker "github.com/autonomy/talos/src/initramfs/pkg/chunker/file"
)

var instance = map[string]*Log{}
var mu = &sync.Mutex{}

// Log represents the log of a service. It supports streaming of the contents of
// the log file by way of implementing the chunker.Chunker interface.
type Log struct {
	Name   string
	Path   string
	source filechunker.Source
}

// New initializes and registers a log for a service.
func New(name string) (*Log, error) {
	if l, ok := instance[name]; ok {
		return l, nil
	}
	logpath := FormatLogPath(name)
	w, err := os.Create(logpath)
	if err != nil {
		return nil, fmt.Errorf("create log file: %s", err.Error())
	}

	l := &Log{
		Name:   name,
		Path:   logpath,
		source: w,
	}

	mu.Lock()
	instance[name] = l
	mu.Unlock()

	return l, nil
}

// Write implements io.WriteCloser.
func (l *Log) Write(p []byte) (n int, err error) {
	return l.source.Write(p)
}

// Close implements io.WriteCloser.
func (l *Log) Close() error {
	return l.source.Close()
}

// Read implements chunker.Chunker.
func (l *Log) Read(ctx context.Context) <-chan []byte {
	c := filechunker.NewChunker(l.source)
	return c.Read(ctx)
}

// FormatLogPath formats the path the log file.
func FormatLogPath(p string) string {
	return path.Join("/var/log", p+".log")
}
