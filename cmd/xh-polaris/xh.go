/*
Copyright © 2025 universero
*/

package xhpolaris

import (
	"github.com/spf13/cobra"
	"github.com/universero/gtool/cmd"
	"github.com/universero/gtool/cmd/xh-polaris/idl"
)

// XhCmd represents the xh-polaris command
var XhCmd = &cobra.Command{
	Use:   "xh [subcommand]",
	Short: "xh-polaris 常用的工具命令",
	Long:  `usage: gtool xh [subcommand] 用于xh-polaris中常见场景的代码生成`,
}

func init() {
	XhCmd.AddCommand(idl.CmdIdl)
	cmd.RootCmd.AddCommand(XhCmd)
}
