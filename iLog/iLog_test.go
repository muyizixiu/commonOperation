package iLog

import (
	"testing"
)

var (
	test_addr []string = []string{"1.txt", "2.txt", "shirasd/dfad"}
)

func Test_Log(t *testing.T) {
	test_log := make([]*log, len(test_addr))
	for i, v := range test_addr {
		test_log[i] = New(v)
	}
	for _, v := range test_log {
		for i := 0; i < 100000; {
			v.Log("hello world")
			i++
		}
	}
}
