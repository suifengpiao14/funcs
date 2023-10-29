package funcs

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

//Md5Lower md5 小写
func Md5Lower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//RandSimple 简单随机数
func RandSimple(maxNumber int) (randInt int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randInt = r.Intn(maxNumber)
	return randInt
}
