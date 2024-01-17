package cmd

import (
	"fmt"
	"strings"

	"github.com/gitdownload/tour/internal/sql2struct"
	"github.com/spf13/cobra"
)

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string
var contact string
var selectTableName string
var mdTableNames string

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql转换和处理",
	Long:  "sql转换和处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var sql2structCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql转换",
	Long:  "sql转换",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo := &sql2struct.DBInfo{
			DBType:   dbType,
			Host:     host,
			UserName: username,
			Password: password,
			Charset:  charset,
		}

		fmt.Printf("mdTableNames的值为：%v \n", mdTableNames)
		//生成文档
		columnses := []*sql2struct.Column{}

		fmt.Printf("md是否为空：%v \n", mdTableNames == "")

		if mdTableNames != "" {
			fmt.Println("开始生成文档")
			mdStr := strings.Split(mdTableNames, ",")
			for i := 0; i < len(mdStr); i++ {

				dbModel := sql2struct.NewDBModel(dbInfo)
				err := dbModel.Connect()
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				columns, err := dbModel.GetColumns(dbName, mdStr[i])
				if err != nil {
					fmt.Printf("err=%v", err)
				}
				columnses = append(columnses, columns...)

			}

			template := sql2struct.NewStructTemplateMd()

			err := template.GenerateMd(columnses)
			if err != nil {
				fmt.Printf("err=%v", err)

			}

		} else {
			//生成model和repository
			fmt.Println("开始生成model和repository")
			if contact != "" {

				//先生成model

				dbModel := sql2struct.NewDBModel(dbInfo)
				err := dbModel.Connect()
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				columns, err := dbModel.GetColumns(dbName, tableName)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				fmt.Printf("columns的长度为：%v \n", len(columns))

				template := sql2struct.NewStructTemplate()
				templateColumns := template.AssemblyColumns(columns)
				err = template.Generate(tableName, templateColumns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				//生成repository
				condition := "\""

				conStr := strings.Split(contact, ",")
				selectTableStr := strings.Split(selectTableName, ",")
				for i := 0; i < len(conStr); i++ {
					condition += " left join " + selectTableStr[i] + " on " + conStr[i]

				}
				condition += "\""

				// ct := sql2struct.TableColumnAndTableName{}
				dbModel = sql2struct.NewDBModel(dbInfo)
				err = dbModel.Connect()
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				columns, err = dbModel.GetColumns(dbName, tableName)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				templateRes := sql2struct.NewStructTemplateRes()
				// templateColumns := template.AssemblyColumns(columns)
				err = templateRes.GenerateRes(tableName, columns, selectTableStr, condition)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

			} else {

				fmt.Println("==================生成单表====================")
				// 先生成model
				ct := sql2struct.TableColumnAndTableName{}
				dbModel := sql2struct.NewDBModel(dbInfo)
				err := dbModel.Connect()
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				columns, err := dbModel.GetColumns(dbName, tableName)
				if err != nil {
					fmt.Printf("err=%v", err)

				}
				fmt.Printf("columns的长度为：%v \n", len(columns))
				ct.TableColumns = columns
				ct.TableName = tableName

				//创建一个模板
				template := sql2struct.NewStructTemplate()
				templateColumns := template.AssemblyColumns(columns)
				err = template.Generate(tableName, templateColumns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				//生成repository
				// ct := sql2struct.TableColumnAndTableName{}
				dbModel = sql2struct.NewDBModel(dbInfo)
				err = dbModel.Connect()
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				columns, err = dbModel.GetColumns(dbName, tableName)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				templateRes := sql2struct.NewStructTemplateOneRes()
				// templateColumns := template.AssemblyColumns(columns)
				err = templateRes.GenerateOneRes(tableName, columns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				templateContract := sql2struct.NewStructTemplateOneContract()
				err = templateContract.GenerateOneContract(tableName, columns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				templateController := sql2struct.NewStructTemplateOneController()
				err = templateController.GenerateOneController(tableName, columns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				templateRouter := sql2struct.NewStructTemplateOneRouter()
				err = templateRouter.GenerateOneRouter(tableName, columns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				templateProxy := sql2struct.NewStructTemplateOneProxy()
				err = templateProxy.GenerateOneProxy(tableName, columns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

				templateValidator := sql2struct.NewStructTemplateOneValidator()
				err = templateValidator.GenerateOneValidator(tableName, columns)
				if err != nil {
					fmt.Printf("err=%v", err)

				}

			}

		}

	},
}

func init() {
	sqlCmd.AddCommand(sql2structCmd)
	sql2structCmd.Flags().StringVarP(&username, "username", "", "", "请输入数据库的用户")
	sql2structCmd.Flags().StringVarP(&password, "password", "", "", "请输入数据库的密码")
	sql2structCmd.Flags().StringVarP(&host, "host", "", "8.129.17.80:3306", "请输入数据库的Host")
	// sql2structCmd.Flags().StringVarP(&host, "host", "", "localhost", "请输入数据库的Host")
	sql2structCmd.Flags().StringVarP(&charset, "charset", "", "utf8mb4", "请输入数据库编码")
	sql2structCmd.Flags().StringVarP(&dbType, "type", "", "mysql", "请输入数据库实例类型")
	sql2structCmd.Flags().StringVarP(&dbName, "db", "", "", "请输入数据库的账号")
	sql2structCmd.Flags().StringVarP(&tableName, "table", "", "", "请输入表名称")
	sql2structCmd.Flags().StringVarP(&contact, "contact", "", "", "请输入联系")
	sql2structCmd.Flags().StringVarP(&selectTableName, "selectTableName", "", "", "请输入查询关联表")

	sql2structCmd.Flags().StringVarP(&mdTableNames, "mdTableNames", "", "", "输入文档表")
	// mdTableNames

	rootCmd.AddCommand(sqlCmd)
}
