package idl

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	commonProto     string // -c 参数：客户端proto文件
	serverProtoFile string // -s 参数：服务端proto文件
	serviceName     string // -n 参数：服务名称
)

var CmdIdl = &cobra.Command{
	Use:   "idl [subcommand]",
	Short: "xh IDL 相关代码生成",
	Long:  `usage: gtool xh idl [subcommand] 用于xh-polaris的IDL中的代码生成`,
}

var GenServiceCmd = &cobra.Command{
	Use:   "gen-svc",
	Short: "生成服务定义",
	Long:  `根据指定的proto文件生成服务定义`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateService(); err != nil {
			fmt.Printf("生成服务失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("服务生成完成")
	},
}

func init() {
	// 添加flags
	GenServiceCmd.Flags().StringVarP(&commonProto, "common-proto", "c", "", "message文件路径")
	GenServiceCmd.Flags().StringVarP(&serverProtoFile, "service-proto", "s", "", "service文件路径")
	GenServiceCmd.Flags().StringVarP(&serviceName, "name", "n", "", "服务名称")

	// 标记必需参数
	if err := GenServiceCmd.MarkFlagRequired("common-proto"); err != nil {
		return
	}
	if err := GenServiceCmd.MarkFlagRequired("service-proto"); err != nil {
		return
	}
	if err := GenServiceCmd.MarkFlagRequired("name"); err != nil {
		return
	}

	CmdIdl.AddCommand(GenServiceCmd)
}

func generateService() error {
	// 1. 读取并解析common proto文件获取message
	messages, exits, err := parseProtoMessages(commonProto)
	if err != nil {
		return err
	}

	// 2. 生成service内容
	serviceContent, err := generateServiceContent(messages, exits)
	if err != nil {
		return err
	}

	// 3. 写入server proto文件
	return writeServerProto(serviceContent)
}

// parseProtoMessages 解析proto文件提取Message定义
func parseProtoMessages(filePath string) ([]Message, map[Message]bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	var messages []Message
	exits := map[Message]bool{}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "message") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := strings.TrimSuffix(parts[1], "{")
				messages = append(messages, Message{Name: name})
				exits[Message{Name: name}] = true
			}
		}
	}
	return messages, exits, nil
}

// generateServiceContent 生成服务定义内容
func generateServiceContent(messages []Message, exits map[Message]bool) (string, error) {
	const templateStr = `syntax = "proto3";

package core_api;

option go_package = "{{.GoPackage}}";

import "basic/http.proto";
import "basic/re.proto";
// 请手动完成依赖导入

service {{.ServiceName}} {
{{- range .Methods}}
    rpc {{.Name}}({{.Request}}) returns ({{.Response}}) {
        option(http.) = "";
    }
{{- end}}
}
`

	data := struct {
		GoPackage   string
		ServiceName string
		Methods     []Method
	}{
		GoPackage:   toSnakeCase(serviceName),
		ServiceName: toCamelCase(serviceName),
		Methods:     generateMethods(messages, exits),
	}

	tpl, err := template.New("service").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tpl.Execute(&buf, data)
	return buf.String(), err
}

// writeServerProto 写入服务定义到server proto文件
func writeServerProto(content string) error {
	return os.WriteFile(serverProtoFile, []byte(content), 0644)
}

// 辅助函数：生成方法列表
func generateMethods(messages []Message, exits map[Message]bool) []Method {
	var methods []Method
	for _, msg := range messages {
		if strings.HasSuffix(msg.Name, "Req") {
			baseName := strings.TrimSuffix(msg.Name, "Req")
			method := Method{
				Name:    baseName,
				Request: msg.Name,
			}
			if exits[Message{baseName + "Resp"}] {
				method.Response = baseName + "Resp"
			} else {
				method.Response = "basic.Response"
			}
			methods = append(methods, method)
		}
	}
	return methods
}

// 辅助函数：转换为蛇形命名
func toSnakeCase(s string) string {
	// 实现蛇形命名转换逻辑
	return strings.ToLower(s)
}

// 辅助函数：转换为驼峰命名
func toCamelCase(s string) string {
	// 使用 cases.Title 来处理，这里指定简体中文作为语言标签
	c := cases.Title(language.SimplifiedChinese)
	return c.String(s)
}

type Message struct {
	Name string
}

type Method struct {
	Name     string
	Request  string
	Response string
}
