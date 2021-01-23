package utils

import (
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

func RandAuthToken() string {
	buf := make([]byte, 32)
	_, err := crand.Read(buf)
	if err != nil {
		return RandomString(64)
	}
	return fmt.Sprintf("%x", buf)
}

// 生成长度为length的随机字符串
func RandomString(n int64) string {
	var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	letterLength := len(letters)
	for i := range result {
		result[i] = letters[rand.Intn(letterLength)]
	}
	return string(result)
}

// 生成32位MD5摘要
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// 生成0-max之间随机数
func RandNumber(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

// convert GBK to UTF-8
func GBKToUTF8(s string) ([]byte, error) {
	I := bytes.NewReader([]byte(s))
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// convert UTF-8 to GBK
func UTF8ToGBK(s string) ([]byte, error) {
	I := bytes.NewReader([]byte(s))
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// convert BIG5 to UTF-8
func Big5ToUTF8(s string) ([]byte, error) {
	I := bytes.NewReader([]byte(s))
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// convert UTF-8 to BIG5
func UTF8ToBig5(s string) ([]byte, error) {
	I := bytes.NewReader([]byte(s))
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 批量替换字符串
func ReplaceStrings(s string, old []string, replace []string) string {
	if s == "" {
		return s
	}
	if len(old) != len(replace) {
		return s
	}
	for i, v := range old {
		s = strings.Replace(s, v, replace[i], 1000)
	}
	return s
}

func InStringSlice(slice []string, element string) bool {
	element = strings.TrimSpace(element)
	for _, v := range slice {
		if strings.TrimSpace(v) == element {
			return true
		}
	}
	return false
}

// 转义json特殊字符
func EscapeJson(s string) string {
	specialChars := []string{"\\", "\b", "\f", "\n", "\r", "\t", "\""}
	replaceChars := []string{"\\\\", "\\b", "\\f", "\\n", "\\r", "\\t", "\\\""}
	return ReplaceStrings(s, specialChars, replaceChars)
}

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(PanicTrace(e))
		}
	}()
	f()
	return
}

// PanicTrace panic调用链跟踪
func PanicTrace(err interface{}) string {
	stackBuf := make([]byte, 4096)
	n := runtime.Stack(stackBuf, false)
	return fmt.Sprintf("panic: %v %s", err, stackBuf[:n])
}
