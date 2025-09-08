package mapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	snakeName string
)

var CmdMapper = &cobra.Command{
	Use:   "mapper [subcommand]",
	Short: "xh mapper 相关代码生成",
	Long:  "usage: gtool xh mapper [subcommand] 用于xh-polaris的Mapper代码生成",
}

var CmdNewMapper = &cobra.Command{
	Use:   "new",
	Short: "生成mapper文件",
	Long:  "生成Mapper模板",
	Run: func(cmd *cobra.Command, args []string) {
		if snakeName == "" {
			fmt.Println("请使用 -n 参数指定名称（蛇形命名，如：user_info）")
			os.Exit(1)
		}
		if err := generate(); err != nil {
			fmt.Printf("生成失败: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("生成成功")
	},
}

func init() {
	CmdNewMapper.Flags().StringVarP(&snakeName, "name", "n", "", "Mapper名称（蛇形命名，如：user_info）")
	CmdMapper.AddCommand(CmdNewMapper)
}

func generate() error {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %w", err)
	}

	// 读取go.mod获取模块名
	goModule, err := getGoModule(wd)
	if err != nil {
		return fmt.Errorf("获取go module失败: %w", err)
	}

	// 转换蛇形命名为驼峰命名
	camelName := toCamelCase(snakeName)

	// 准备模板数据
	data := struct {
		PackageName string
		GoModule    string
		SnakeName   string
		Name        string
	}{
		PackageName: snakeName,
		GoModule:    goModule,
		SnakeName:   snakeName,
		Name:        camelName,
	}

	// 创建目录
	dirPath := filepath.Join(wd, "biz", "infra", "mapper", snakeName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	filePath := filepath.Join(dirPath, "mapper.go")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 定义模板
	tmpl := `package {{.PackageName}}

import (
	"github.com/xh-polaris/{{.GoModule}}/biz/infra/config"
	"github.com/zeromicro/go-zero/core/stores/monc"
)

var _ MongoMapper = (*mongoMapper)(nil)

const (
	collection     = "{{.SnakeName}}"
	cacheKeyPrefix = "cache:{{.SnakeName}}:"
)

type MongoMapper interface{}

type mongoMapper struct {
	conn *monc.Model
}

func New{{.Name}}MongoMapper(config *config.Config) MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, collection, config.CacheConf)
	return &mongoMapper{conn: conn}
}
`

	// 解析并执行模板
	t, err := template.New("mapper").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	return nil
}

// 从go.mod中获取模块名
func getGoModule(dir string) (string, error) {
	content, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("go.mod中没有找到module定义")
}

// 将蛇形命名转换为驼峰命名
func toCamelCase(s string) string {
	var sb strings.Builder
	nextUpper := true // 标记下一个字符是否需要大写（用于处理 _ 后的字母）
	for _, ch := range s {
		if ch == '_' {
			nextUpper = true // 遇到下划线，标记下一个字符需要大写
			continue         // 跳过当前下划线
		}
		if nextUpper {
			sb.WriteRune(unicode.ToUpper(ch)) // 大写当前字符
			nextUpper = false                 // 重置标记
		} else {
			sb.WriteRune(ch) // 否则直接写入字符
		}
	}
	return sb.String()
}
