/*
Copyright © 2025 universero
*/

package tool

import (
	"fmt"
	"strconv"
	"time"

	"github.com/universero/gtool/cmd"

	"github.com/spf13/cobra"
)

var timestampCmd = &cobra.Command{
	Use:   "timestamp <unix_timestamp>",
	Short: "Convert a Unix timestamp to formatted datetime",
	Long: `Convert seconds since Unix epoch (1970-01-01 UTC) to human-readable format with timezone support.

Examples:
  gtool timestamp 1620000000               # Use local timezone
  gtool timestamp 1620000000 -z UTC         # UTC timezone
  gtool timestamp 1620000000 --tz=Asia/Shanghai # Specific timezone

Supported timezone formats:
  - IANA Time Zone names (e.g. Asia/Shanghai, America/New_York)
  - 'Local' for system default timezone
  - 'UTC' or 'GMT' for universal time`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 解析时间戳
		ts, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println("错误：无效的时间戳格式")
			return
		}

		// 获取时区参数
		tz, _ := cmd.Flags().GetString("tz")

		// 加载时区
		loc, err := time.LoadLocation(tz)
		if err != nil {
			fmt.Printf("error：invalid timezone %s\n", tz)
			return
		}

		// 转换时间
		t := time.Unix(ts, 0).In(loc)

		// 格式化输出
		fmt.Printf("Timestamp: %d\n", ts)
		fmt.Printf("Formated : %s\n", t.Format("2006-01-02 15:04:05 MST"))
		fmt.Printf("UTC  Time: %s\n", t.UTC().Format("2006-01-02 15:04:05 MST"))
	},
}

func init() {
	cmd.RootCmd.AddCommand(timestampCmd)

	// 添加时区参数
	timestampCmd.Flags().StringP("tz", "z", "Local", `Timezone specification:
Use IANA timezone name (e.g. Asia/Shanghai), 
'Local' for system timezone, 
or 'UTC' for universal time`)
}
