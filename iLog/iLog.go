package iLog

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type log struct {
	addr         string
	locker       sync.Mutex
	file         *os.File
	fileOpenFlag bool
	isClosing    bool
	closed       bool
}

var a int

func New(addre string) *log {
	l := &log{addr: addre, locker: sync.Mutex{}}
	if l.open() {
		l.close()
		return l
	}
	fmt.Println("fail to create a log object")
	return &log{}
}
func (l *log) open() bool {
	var (
		f    *os.File
		err  error
		addr = l.addr
	)
	if addr == "" {
		return false
	}
	if isExist(addr) {
		f, err = os.OpenFile(addr, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("\n", err.Error())
			return false
		}
	} else {
		f, err = os.Create(addr)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	l.file = f
	l.fileOpenFlag = true
	return true
}
func (l log) close() {
	if !l.fileOpenFlag {
		return
	}
	if !l.closed {
		go func() {
			l.closed = true
			l.isClosing = true
			for {
				time.Sleep(10 * time.Second)
				if isClosing {
					l.isClosing = false
				} else {
					break
				}
			}
			l.file.Close()
			l.fileOpenFlag = false
			l.closed = false
		}()
	} else {
		l.isClosing = true
	}
}
func isExist(addr string) bool {
	_, err := os.Stat(addr)
	return err == nil || os.IsExist(err)
}
func (l log) Log(str string) bool {
	str = now() + str
	return l.write([]byte(str))
}
func (l log) LogError(err error) bool {
	return l.Log(err.Error())
}
func (l log) write(b []byte) bool {
	if !l.fileOpenFlag {
		if !l.open() {
			return false
		}
	}
	defer l.close()
	l.locker.Lock()
	defer l.locker.Unlock()
	l.file.Seek(0, 2)
	a++
	_, err := l.file.Write(b)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func now() string {
	return "\r\n##############" + time.Now().String() + "###############\r\n"
}
func (l log) Read() string {
	l.open()
	defer l.close()
	b := make([]byte, 1024)
	return string(l.read(b))
}
func (l log) read(b []byte) []byte {
	n, err := l.file.Read(b)
	if err != nil {
		fmt.Print(err.Error())
	}
	return b[:n]
}
func (l log) LogNum(v interface{}) bool {
	switch v.(type) {
	case string:
		return l.Log(v.(string))
	case int:
		return l.Log(strconv.Itoa(v.(int)))
	default:
		return false
	}
}
