package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	if len(str) > 0 {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(str))
		cipherStr := md5Ctx.Sum(nil)
		return hex.EncodeToString(cipherStr)
	}
	return ""
}
