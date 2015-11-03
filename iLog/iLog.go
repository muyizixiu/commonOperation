package iLog

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var a int

type log struct {
	addr         string
	locker       sync.Mutex
	file         *os.File
	fileOpenFlag bool
	isClosing    bool
	closed       bool
}

const logDir = "/root/log/"

func New(addre string) *log {
	l := &log{addr: logDir + addre, locker: sync.Mutex{}}
	if l.open() {
		l.close()
		return l
	}
	fmt.Println("fail to create a log object")
	return &log{}
}

//打开log文件
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

//关闭log文件
func (l *log) close() {
	if !l.fileOpenFlag {
		return
	}
	if !l.closed {
		l.closed = true
		l.isClosing = true
		go func() {
			for {
				time.Sleep(10 * time.Second)
				if l.isClosing {
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

//文件是否存在
func isExist(addr string) bool {
	_, err := os.Stat(addr)
	return err == nil || os.IsExist(err)
}

//记录字符串
func (l log) Log(str string) bool {
	str = now() + str
	return l.write([]byte(str))
}

//错误记录
func (l log) LogError(err error) bool {
	return l.Log(err.Error())
}
func (l log) write(b []byte) bool {
	if !l.fileOpenFlag {
		if !l.open() {
			fmt.Println("fail")
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
	return "\r\n" + time.Now().String() + "   # "
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
