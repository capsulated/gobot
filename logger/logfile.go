package logger

import (
	"os"
	"sync"
)

type LogFile struct {
	path  string
	mutex *sync.Mutex
	file  *os.File
}

func NewLogFile(file string) (*LogFile, error) {
	lf := new(LogFile)
	lf.path = file
	lf.mutex = &sync.Mutex{}
	lf.Open()
	return lf, nil
}

func (l *LogFile) Open() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	var err error
	l.file, err = os.OpenFile(l.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	return err
}

func (l *LogFile) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.Close()
}

func (l *LogFile) Sync() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.Sync()
}

func (l *LogFile) Reopen() error {
	var err error
	err = l.Sync()
	if err != nil {
		return err
	}
	err = l.Close()
	if err != nil {
		return err
	}
	err = l.Open()
	if err != nil {
		return err
	}
	return nil
}

func (l *LogFile) Write(bytes []byte) (n int, err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.file.Write(bytes)
}
