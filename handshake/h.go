package h

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"regexp"
)

var (
	mask        = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

//for a websocket connection this function funcs as handshaking
//apply websocket protocal to establish a connection
func HandshakeOfWS(conn net.Conn) error {
	b := make([]byte, 1024)

	a := make([]byte, 1024)
	n, err := conn.Read(a)
	b = append(b, a[:n]...)
	if err != nil {
		if err == io.EOF {
		} else {
			return errors.New("error happen when read from connection!")
		}
	}

	re, err := regexp.Compile("Sec-WebSocket-Key: (.*)\r\n")
	if checkErr(err) {
		return errors.New("regexp error!")
	}
	result := re.FindStringSubmatch(string(b))
	if result == nil {
		return errors.New("not right format!")
	}
	key_tmp := append([]byte(result[1]), []byte(mask)...)
	key_tmp_tmp := sha1.Sum(key_tmp)
	encoder := base64.NewEncoding(base64Table)
	dst := make([]byte, 1024)
	encoder.Encode(dst, key_tmp_tmp[:])
	reply := "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + string(dst) + "\r\n\r\n"
	conn.Write([]byte(reply))
	return nil
}
func checkErr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
