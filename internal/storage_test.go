// Author: yangzq80@gmail.com
// Date: 2021-02-23
//
package internal

import (
	"fmt"
	"strings"
	"testing"
)

func TestA(t *testing.T) {

}

func TestCmd(t *testing.T) {
	sts:=strings.Split("weight,desc",",")
	fmt.Println(sts[0],sts[1])
	//ExecCommand("vue init webpack demo")
	//bptRootIdxDir := "aa/b"
	//if ok := filesystem.PathIsExist(bptRootIdxDir); !ok {
	//	if err := os.MkdirAll(bptRootIdxDir, os.ModePerm); err != nil {
	//		fmt.Println(err)
	//		//return nil, err
	//	}
	//}
}
