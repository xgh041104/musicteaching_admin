package md5

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	// 测试用例
	sum := Md5("administrator!q@w#e$r%t") // 调用函数

	fmt.Println(sum)

}
