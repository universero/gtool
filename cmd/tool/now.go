/*
Copyright © 2025 universero
*/

package tool

import (
	"fmt"
	"os"
	"time"

	"github.com/universero/gtool/cmd"

	"github.com/spf13/cobra"
)

// nowCmd represents the now command
var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Display current timestamp with timezone support",
	Long: `Display current time and Unix timestamp in specified timezone.

Examples:
  gtool now                 # Show local time
  gtool now -z UTC          # Show time in UTC
  gtool now --tz=Asia/Tokyo # Show time in specific zone

Supported timezone formats:
  - IANA Time Zone names (e.g. Asia/Shanghai)
  - 'Local' for system default
  - 'UTC' or 'GMT'`,
	Run: func(cmd *cobra.Command, args []string) {
		tzName, _ := cmd.Flags().GetString("tz")

		loc, err := time.LoadLocation(tzName)
		if err != nil {
			fmt.Printf("Error: invalid timezone %s\n", tzName)
			os.Exit(1)
		}

		now := time.Now().In(loc)

		fmt.Printf("Current Time: %s\n", now.Format("2006-01-02 15:04:05 MST"))
		fmt.Printf("Timestamp   : %d\n", now.Unix())
	},
}

func init() {
	cmd.RootCmd.AddCommand(nowCmd)
	nowCmd.Flags().StringP("tz", "z", "Local", `Timezone specification (e.g. Asia/Shanghai). 
Valid options: IANA names, 'Local', 'UTC'`)
}
