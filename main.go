/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/universero/gtool/cmd"
	_ "github.com/universero/gtool/cmd/tool"
	_ "github.com/universero/gtool/cmd/xh-polaris" // 确保导入xh-polaris包
)

func main() {
	cmd.Execute()
}
