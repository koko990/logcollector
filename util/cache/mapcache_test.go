package cache_test

import (
	"fmt"
	"git/inspursoft/board/src/collector/util/cache"
	"testing"
)

func TestRegister(t *testing.T) {
	temp := cache.Register("myRepo")
	temp.Put("asd", "asddd")
	isOk := temp.MatchQueryKey("as")
	fmt.Println(temp.Get("asd"), isOk)
}
