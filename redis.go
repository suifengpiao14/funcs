package funcs

import (
	"fmt"
	"strings"
)

func FormatRedisKey(key string) string {
	cacheKey := strings.ReplaceAll(key, " ", "_")
	if len(cacheKey) > 32 {
		prefix, suffix := cacheKey[:32], cacheKey[32:]
		cacheKey = fmt.Sprintf("%s_%s", prefix, Md5Lower(suffix))
	}
	return cacheKey
}
