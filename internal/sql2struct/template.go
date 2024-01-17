package sql2struct

import (
	"fmt"
	"os"
	"text/template"

	"github.com/gitdownload/tour/internal/word"
)

//单表生成
const structTpl = `package model
	type {{.TableName | ToCamelCase}} struct {
	{{range .Columns}} {{$length := len .Comment}} {{if gt $length 0}}//{{.Comment}} {{else}}// {{.Name}} {{end}}
		{{$typeLen := len .Type}} {{if gt $typeLen 0}} {{.Name | ToCamelCase}} {{.Type}} {{.Tag}}{{else}}{{.Name}}{{end}}
	{{end}}
	}

	func (model {{.TableName | ToCamelCase}}) TableName() string {
		return "{{.TableName}}"
	}


	func (obj *{{.TableName | ToCamelCase}}) GetFields() []interface{} {
		v := reflect.ValueOf(obj).Elem()
		var list []interface{}
		for i := 0; i < v.NumField(); i++ {
			vv := v.Field(i).Addr().Interface()
			list = append(list, vv)
		}
		return list
	}

	type {{.TableName | ToCamelCase}}Vo struct {
		{{range .Columns}} {{$length := len .Comment}} {{if gt $length 0}}//{{.Comment}} {{else}}// {{.Name}} {{end}}
			{{$typeLen := len .Type}} {{if gt $typeLen 0}} {{.Name | ToCamelCase}} {{.Type}} {{.Tag}}{{else}}{{.Name}}{{end}}
		{{end}}
		}
	
	func (vo {{.TableName | ToCamelCase}}Vo) To{{.TableName | ToCamelCase}}() *{{.TableName | ToCamelCase}} {
		a := {{.TableName | ToCamelCase}}{}
		{{range .Columns}}
		a.{{.Name | ToCamelCase}} = vo.{{.Name | ToCamelCase}}
		{{end}}
		return &a
		//if len(vo.ImageUrl) > 0 {
			
			//a.ImageUrl = helper.SliceToString(vo.ImageUrl)
		//}
	}

	func (a {{.TableName | ToCamelCase}}) To{{.TableName | ToCamelCase}}Vo() *{{.TableName | ToCamelCase}}Vo {
		vo := {{.TableName | ToCamelCase}}Vo{}
		{{range .Columns}}
		vo.{{.Name | ToCamelCase}} = a.{{.Name | ToCamelCase}}
		{{end}}
		return &vo

		//if len(a.ImageUrl) > 0 {
			//vo.ImageUrl = strings.Split(a.ImageUrl, ",")
		//} else {
			//vo.ImageUrl = []string{}
		//}
	}

	type {{.TableName | ToCamelCase}}Res struct {
		
	}
	
	

	`

type StructTemplalte struct {
	structTpl string
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type StructTemplalteDB struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplalte {
	return &StructTemplalte{structTpl: structTpl}
}

func (t *StructTemplalte) AssemblyColumns(tbColumns []*Column) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {

		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", word.UnderscoreToLowerCamelCase(column.ColumnName))
		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DBTypeToStructType[column.DataType],
			Tag:     tag,
			Comment: column.ColumnComment,
		})
	}

	return tplColumns
}

func (t *StructTemplalte) Generate(tableName string, tplColumns []*StructColumn) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
	}).Parse(t.structTpl))

	tplDB := StructTemplalteDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	fileName := "./internal/sql2struct/model/" + tableName + ".go"
	f, err2 := os.Create(fileName)
	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

// {{$typeLen := len .Type}} {{if gt $typeLen 0}} {{.Name | ToCamelCase}} {{.Type}} {{.Tag}}{{else}}{{.Name}}{{end}}

//多表生成
const structTplRes = `package repository

const (
	_{{.TableName | ToLowerCamelCase}}TableName = "{{.TableName}}"
	_get{{.TableName | ToCamelCase}}By{{.TableName | ToCamelCase}}{{range .ContactTable}}And{{. | ToCamelCase}}{{end}}Fields = "{{$ColumnLen := len .Columns}}{{range $index,$value := .Columns}}{{$value.TableName}}.{{$value.ColumnName}}{{if le $index ($ColumnLen |DescTow)}},{{end}}{{end}}"
	
	_{{.TableName | ToLowerCamelCase}}Select    = "select " + _get{{.TableName | ToCamelCase}}By{{.TableName | ToCamelCase}}{{range .ContactTable}}And{{. | ToCamelCase}}{{end}}Fields + " from " + _{{.TableName | ToLowerCamelCase}}TableName + {{.ContactCondition}}
	_{{.TableName | ToLowerCamelCase}}Count     = "select count({{.TableName}}.id) from " + _{{.TableName | ToLowerCamelCase}}TableName
	
)


type {{.TableName | ToCamelCase}}Repository struct {
}

func (r *{{.TableName | ToCamelCase}}Repository) Select{{.TableName | ToCamelCase}}By{{.TableName | ToCamelCase}}{{range .ContactTable}}And{{. | ToCamelCase}}{{end}}(tool *model.SelectTool) ([]*model.{{.TableName | ToCamelCase}}, error){

	var sqlx strings.Builder
	sqlx.WriteString(_{{.TableName | ToLowerCamelCase}}Select)
	tool.CompletionSelectWhereClause(&sqlx)
	fmt.Printf("完整sql语句是：%v \n", sqlx.String())
	rows, err := query(sqlx.String(), tool.GetSelectArgs()...)
	if err != nil {
		fmt.Printf("query err：%v \n", err)
		return nil, err
	}
	defer rows.Close()
	list := []*model.{{.TableName | ToCamelCase}}{}
	for rows.Next() {
		item := model.{{.TableName | ToCamelCase}}{}
		err := rows.Scan(item.GetFields()...)
		if err != nil {
			fmt.Println("rows.Scan err:", err)
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil

}
`

type StructTemplalteRes struct {
	structTplRes string
}

type StructTemplalteDBMany struct {
	TableName        string
	Columns          []*Column
	ContactTable     []string
	ContactCondition string
}

func NewStructTemplateRes() *StructTemplalteRes {
	return &StructTemplalteRes{structTplRes: structTplRes}
}

// func (t *StructTemplalte) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
// 	tplColumns := make([]*StructColumn, 0, len(tbColumns))
// 	for _, column := range tbColumns {

// 		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", word.UnderscoreToLowerCamelCase(column.ColumnName))
// 		tplColumns = append(tplColumns, &StructColumn{
// 			Name:    column.ColumnName,
// 			Type:    DBTypeToStructType[column.DataType],
// 			Tag:     tag,
// 			Comment: column.ColumnComment,
// 		})
// 	}

// 	return tplColumns
// }

func (t *StructTemplalteRes) GenerateRes(tableName string, tplColumns []*Column, contactTable []string, contactCondition string) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"DescTow":          word.DescTow,
		"DescFour":         word.DescFour,
		"DescThree":        word.DescThree,
	}).Parse(t.structTplRes))

	tplDB := StructTemplalteDBMany{
		TableName:        tableName,
		Columns:          tplColumns,
		ContactTable:     contactTable,
		ContactCondition: contactCondition,
	}

	fileName := "./internal/sql2struct/repository/" + tableName + ".go"
	f, err2 := os.Create(fileName)
	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

const structTplMd = `
|          字段           |                                               含义                                                 |   类型   |
|: ---------------------: | :-----------------------------------------------------------------------------------------------: | :------: |
	{{range .Columns}}{{.ColumnName | CenterStr25}}{{.ColumnComment|CenterStr90}}{{.DataType | CenterStr18}}
	{{end}}
	`

//生成文档
type StructTemplalteMd struct {
	structTplMd string
}

type StructTemplaltMdDB struct {
	Columns []*Column
}

func NewStructTemplateMd() *StructTemplalteMd {
	return &StructTemplalteMd{structTplMd: structTplMd}
}

func (t *StructTemplalteMd) GenerateMd(tplColumns []*Column) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"CenterStr25":      word.CenterStr25,
		"CenterStr30":      word.CenterStr40,
		"CenterStr18":      word.CenterStr18,
		"CenterStr90":      word.CenterStr90,
	}).Parse(t.structTplMd))

	tplDB := StructTemplaltMdDB{

		Columns: tplColumns,
	}

	fileName := "./internal/sql2struct/md/" + "test" + ".md"
	f, err2 := os.Create(fileName)
	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

//单表生成

const structTplOneRes = `package repository

const (
	_{{.TableName | ToLowerCamelCase}}TableName = "{{.TableName}}"
	_get{{.TableName | ToCamelCase}}Fields = "{{$ColumnLen := len .Columns }}{{range $index,$value := .Columns}}{{$value.ColumnName}}{{if le $index ($ColumnLen |DescTow)}},{{end}}{{end}}"
	
	_{{.TableName | ToLowerCamelCase}}Select    = "select " + _get{{.TableName | ToCamelCase}}Fields + " from " + _{{.TableName | ToLowerCamelCase}}TableName
	_{{.TableName | ToLowerCamelCase}}Count     = "select count({{.TableName}}.id) from " + _{{.TableName | ToLowerCamelCase}}TableName
	
)

var (
	_{{.TableName | ToLowerCamelCase}}ID             uint64
	_{{.TableName | ToLowerCamelCase}}Repository     contract.{{.TableName | ToCamelCase}}Repository
	_once{{.TableName | ToCamelCase}}Repository sync.Once
)


func Create{{.TableName | ToCamelCase}}Repository() contract.{{.TableName | ToCamelCase}}Repository {

	_once{{.TableName | ToCamelCase}}Repository.Do(func() {
		_{{.TableName | ToLowerCamelCase}}Repository = &{{.TableName | ToCamelCase}}Repository{}

		if _{{.TableName | ToLowerCamelCase}}ID == 0 {
			_{{.TableName | ToLowerCamelCase}}ID, _ = max("{{.TableName}}", "id")

			if _{{.TableName | ToLowerCamelCase}}ID == 0 {
				_{{.TableName | ToLowerCamelCase}}ID = WebConfig().App.AppID - WebConfig().App.AppNum
			}
		}
	})

	return _{{.TableName | ToLowerCamelCase}}Repository
}



type {{.TableName | ToCamelCase}}Repository struct {
}


func (r *{{.TableName | ToCamelCase}}Repository) Create{{.TableName | ToCamelCase}}ID() uint64 {
	return atomic.AddUint64(&_{{.TableName | ToLowerCamelCase}}ID, WebConfig().App.AppNum)
}


func (r *{{.TableName | ToCamelCase}}Repository) Select{{.TableName | ToCamelCase}}ByTool(tool *model.SelectTool) ([]*model.{{.TableName | ToCamelCase}}, error){

	var sqlx strings.Builder
	sqlx.WriteString(_{{.TableName | ToLowerCamelCase}}Select)
	tool.CompletionSelectWhereClause(&sqlx)
	fmt.Printf("完整sql语句是：%v \n", sqlx.String())
	rows, err := query(sqlx.String(), tool.GetSelectArgs()...)
	if err != nil {
		fmt.Printf("query err：%v \n", err)
		return nil, err
	}
	defer rows.Close()
	list := []*model.{{.TableName | ToCamelCase}}{}
	for rows.Next() {
		item := model.{{.TableName | ToCamelCase}}{}
		err := rows.Scan(item.GetFields()...)
		if err != nil {
			fmt.Println("rows.Scan err:", err)
			return nil, err
		}
		list = append(list, &item)
	}
	return list, nil

}



func (r *{{.TableName | ToCamelCase}}Repository) Create{{.TableName | ToCamelCase}}({{.TableName | ToLowerCamelCase}} *model.{{.TableName | ToCamelCase}}) (uint64, error) {
	// terror := helper.NewTError()

	insertSQL := "insert into {{.TableName}}({{$ColumnLen := len .Columns }}{{range $index,$value := .Columns}}{{if le $index ($ColumnLen | DescThree)}}{{$value.ColumnName}}{{end}}{{if le $index ($ColumnLen | DescFour)}},{{end}}{{end}}) values({{$ColumnLen := len .Columns }}{{range $index,$value := .Columns}}{{if le $index ($ColumnLen |DescThree)}}?{{end}}{{if le $index ($ColumnLen |DescFour)}},{{end}}{{end}})"

	{{.TableName | ToLowerCamelCase}}.Id = r.Create{{.TableName | ToCamelCase}}ID()
	{{.TableName | ToLowerCamelCase}}.CreatedAt = now()

	_, err := exec(insertSQL, {{$ColumnLen := len .Columns }}{{range $index,$value := .Columns}}{{if le $index ($ColumnLen | DescThree)}}{{.TableName | ToLowerCamelCase}}.{{$value.ColumnName | ToCamelCase}}{{end}}{{if le $index ($ColumnLen | DescFour)}},{{end}}{{end}})
	if err != nil {
		fmt.Printf("创建有误：%v \n", err)
		return 0, err
	}

	return {{.TableName | ToLowerCamelCase}}.Id, nil
}


func (r *{{.TableName | ToCamelCase}}Repository) Update{{.TableName | ToCamelCase}}ByTool(tool *model.UpdateTool, pip *sql.Tx) (int64, error) {
	var sqlx strings.Builder
	sqlx.WriteString("update {{.TableName}} ")
	tool.CompletionUpdateSentence(&sqlx)
	//  // fmt.Println("完整sql语句是：", sqlx.String())
	var result sql.Result
	var err error
	if pip != nil {
		result, err = txExec(pip, sqlx.String(), tool.GetUpdateArgs()...)
	} else {
		result, err = exec(sqlx.String(), tool.GetUpdateArgs()...)
	}
	if err != nil {
		fmt.Println("exec err:", err)
		return 0, err
	}
	return result.RowsAffected()
}


func (r *{{.TableName | ToCamelCase}}Repository) Count{{.TableName | ToCamelCase}}ByTool(tool *model.SelectTool) (int, error) {

	ssql := "select count(id) from {{.TableName}}"

	var sqlx strings.Builder
	sqlx.WriteString(ssql)
	tool.CompletionSelectWhereClause(&sqlx)

	// // fmt.Printf("完整sql语句是：%v \n", sqlx.String())
	row := queryRow(sqlx.String(), tool.GetSelectArgs()...)
	// rows, err := query(sqlx.String(),tool.GetSelectArgs()...)

	if row.Err() != nil {
		fmt.Printf("queryrow err：%v \n", row.Err())
		return 0, row.Err()
	}

	var count int
	err := row.Scan(&count)
	if err != nil {
		fmt.Printf("queryRow err：%v \n", err)
		return 0, err

	}

	return count, nil

}


`

type StructTemplalteOneRes struct {
	structTplOneRes string
}

type StructTemplalteOneDB struct {
	TableName string
	Columns   []*Column
}

func NewStructTemplateOneRes() *StructTemplalteOneRes {
	return &StructTemplalteOneRes{structTplOneRes: structTplOneRes}
}

func (t *StructTemplalteOneRes) GenerateOneRes(tableName string, tplColumns []*Column) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"DescTow":          word.DescTow,
		"DescFour":         word.DescFour,
		"DescThree":        word.DescThree,
	}).Parse(t.structTplOneRes))

	tplDB := StructTemplalteOneDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	fileName := "./internal/sql2struct/repository/" + tableName + ".go"
	f, err2 := os.Create(fileName)

	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

//单表contract生成

const structTplOneContract = `package contract

type {{.TableName | ToCamelCase}}Repository interface {

	Select{{.TableName | ToCamelCase}}ByTool(tool *model.SelectTool) ([]*model.{{.TableName | ToCamelCase}}, error)
	Create{{.TableName | ToCamelCase}}({{.TableName | ToLowerCamelCase}} *model.{{.TableName | ToCamelCase}}) (uint64, error) 
	Update{{.TableName | ToCamelCase}}ByTool(tool *model.UpdateTool, pip *sql.Tx) (int64, error) 
	Count{{.TableName | ToCamelCase}}ByTool(tool *model.SelectTool) (int, error)


}
`

type StructTemplalteOneContract struct {
	structTplOneContract string
}

func NewStructTemplateOneContract() *StructTemplalteOneContract {
	return &StructTemplalteOneContract{structTplOneContract: structTplOneContract}
}

func (t *StructTemplalteOneContract) GenerateOneContract(tableName string, tplColumns []*Column) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"DescTow":          word.DescTow,
		"DescFour":         word.DescFour,
		"DescThree":        word.DescThree,
	}).Parse(t.structTplOneContract))

	tplDB := StructTemplalteOneDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	fileName := "./internal/sql2struct/contract/" + tableName + ".go"
	f, err2 := os.Create(fileName)

	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

const structTplOneController = `package controller


var (
	_{{.TableName | ToLowerCamelCase}}Controller     *{{.TableName | ToCamelCase}}Controller
	_once{{.TableName | ToCamelCase}}Controller sync.Once
)

func Create{{.TableName | ToCamelCase}}Controller() *{{.TableName | ToCamelCase}}Controller {
	_once{{.TableName | ToCamelCase}}Controller.Do(func() {
		_{{.TableName | ToLowerCamelCase}}Controller = new({{.TableName | ToCamelCase}}Controller)
	})
	return _{{.TableName | ToLowerCamelCase}}Controller
}

type {{.TableName | ToCamelCase}}Controller struct {
}

func (this *{{.TableName | ToCamelCase}}Controller) Create{{.TableName | ToCamelCase}}(ctx *web.Context) (web.Data, error) {

	body := model.{{.TableName | ToCamelCase}}{}
	err := ctx.TryParseBody(&body)
	res := model.CreateResultRes()
	if err != nil {
		fmt.Printf("解析参数有误：%v \n", err)
		return res.SetErrorData(-1, "解析参数有误"), nil
	}

	errC := validator.Validate{{.TableName | ToCamelCase}}(&body)
	if errC != nil {
		return res.SetErrorData(errC.ErrCode, errC.ErrMsg), nil
	}
	result, errC := proxy.Create{{.TableName | ToCamelCase}}(&body)
	if errC != nil {
		return res.SetErrorData(errC.ErrCode, errC.ErrMsg), nil
	}
	return res.SetSuccessData(result), nil
	
}


func (this *{{.TableName | ToCamelCase}}Controller) Select{{.TableName | ToCamelCase}}By(ctx *web.Context) (web.Data, error) {
	var (
		name     string
		page     int
		pageSize int
	)

	res := model.CreateResultRes()
	ctx.TryParseQuery("page", &page)
	//ctx.TryParseQuery("name", &name)
	ctx.TryParseQuery("pageSize", &pageSize)

	result, errC := proxy.Select{{.TableName | ToCamelCase}}By(page, pageSize)
	if errC != nil {
		return res.SetErrorData(errC.ErrCode, errC.ErrMsg), nil
	}
	return res.SetSuccessData(result), nil
}

func (this *{{.TableName | ToCamelCase}}Controller) Update{{.TableName | ToCamelCase}}(ctx *web.Context) (web.Data, error) {

	body := model.{{.TableName | ToCamelCase}}{}

	res := model.CreateResultRes()

	err := ctx.TryParseBody(&body)
	if err != nil {
		fmt.Printf("解析参数有误：%v \n", err)
		return res.SetErrorData(-1, "解析参数有误："+err.Error()), nil
	}

	if body.Id <= 0 {
		fmt.Println("id值有误")
		return res.SetErrorData(-1, "ID值有误"), nil
	}

	result, errC := proxy.Update{{.TableName | ToCamelCase}}(&body)
	if errC != nil {
		fmt.Println(errC)
		return res.SetErrorData(errC.ErrCode, errC.ErrMsg), nil
	}

	return res.SetSuccessData(result), nil
}

func (this *{{.TableName | ToCamelCase}}Controller) Del{{.TableName | ToCamelCase}}(ctx *web.Context) (web.Data, error) {

	var ID uint64
	ctx.TryParseQuery("ID", &ID)

	res := model.CreateResultRes()

	if ID <= 0 {
		fmt.Println("id值有误")
		return res.SetErrorData(-1, "ID值有误"), nil
	}

	result, errC := proxy.Del{{.TableName | ToCamelCase}}ByID(ID)
	if errC != nil {
		fmt.Println(errC)
		return res.SetErrorData(errC.ErrCode, errC.ErrMsg), nil
	}

	return res.SetSuccessData(result), nil
}

`

type StructTemplalteOneController struct {
	structTplOneController string
}

func NewStructTemplateOneController() *StructTemplalteOneController {
	return &StructTemplalteOneController{structTplOneController: structTplOneController}
}

func (t *StructTemplalteOneController) GenerateOneController(tableName string, tplColumns []*Column) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"DescTow":          word.DescTow,
		"DescFour":         word.DescFour,
		"DescThree":        word.DescThree,
	}).Parse(t.structTplOneController))

	tplDB := StructTemplalteOneDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	fileName := "./internal/sql2struct/controller/" + tableName + ".go"
	f, err2 := os.Create(fileName)

	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

const structTplOneRouter = `package route



func {{.TableName | ToLowerCamelCase}}Route(app *web.Application, prefix string) {

	controller := controller.Create{{.TableName | ToCamelCase}}Controller()
	//-----------后台管理在用⬇-----------------
	 app.Post(prefix+"/create/{{.TableName}}/", middleware.Chain(controller.Create{{.TableName | ToCamelCase}}))
	 app.Get(prefix+"/get/{{.TableName}}/", middleware.Chain(controller.Select{{.TableName | ToCamelCase}}By))
	 app.Put(prefix+"/update/{{.TableName}}/", middleware.Chain(controller.Update{{.TableName | ToCamelCase}}))
	 //app.Delete(prefix+"/del/{{.TableName}}/", middleware.Chain(controller.Del{{.TableName | ToCamelCase}}))

}


`

type StructTemplalteOneRouter struct {
	structTplOneRouter string
}

func NewStructTemplateOneRouter() *StructTemplalteOneRouter {
	return &StructTemplalteOneRouter{structTplOneRouter: structTplOneRouter}
}

func (t *StructTemplalteOneRouter) GenerateOneRouter(tableName string, tplColumns []*Column) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"DescTow":          word.DescTow,
		"DescFour":         word.DescFour,
		"DescThree":        word.DescThree,
	}).Parse(t.structTplOneRouter))

	tplDB := StructTemplalteOneDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	fileName := "./internal/sql2struct/route/" + tableName + ".go"
	f, err2 := os.Create(fileName)

	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

const structTplOneProxy = `package proxy

func Create{{.TableName | ToCamelCase}}({{.TableName | ToLowerCamelCase}} *model.{{.TableName | ToCamelCase}}) (uint64, *model.CommonError) {
	repo := repository.Create{{.TableName | ToCamelCase}}Repository()
	result, err := repo.Create{{.TableName | ToCamelCase}}({{.TableName | ToLowerCamelCase}})
	if err != nil {
		fmt.Printf("create fail：%v \n", err)
		return 0, model.CreateCommonError(-5, "create fail")
	}

	return result, nil
}


func GetEvenDel{{.TableName | ToCamelCase}}ByID(id uint64) (*model.{{.TableName | ToCamelCase}}, *model.CommonError) {
	repo := repository.Create{{.TableName | ToCamelCase}}Repository()
	tool := model.CreateSelectTool()
	tool.AddWhereCase("id", "=", id)
	result, err := repo.Select{{.TableName | ToCamelCase}}ByTool(tool)
	if err != nil {
		fmt.Printf("select by id fail：%v \n", err)
		return nil, model.CreateCommonError(-5, "select {{.TableName | ToLowerCamelCase}} fail by id")
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}



func GetExist{{.TableName | ToCamelCase}}ByID(id uint64) (*model.{{.TableName | ToCamelCase}}, *model.CommonError) {
	repo := repository.Create{{.TableName | ToCamelCase}}Repository()
	tool := model.CreateSelectTool()
	tool.AddWhereCase("id", "=", id)
	tool.AddWhereCase("isnull(deleted_at)", "", nil)
	result, errC := repo.Select{{.TableName | ToCamelCase}}ByTool(tool)
	if errC != nil {
		fmt.Printf("select {{.TableName | ToLowerCamelCase}} by id fail：%v \n", errC)
		return nil, model.CreateCommonError(-5, "select {{.TableName | ToLowerCamelCase}} fail")
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}


func Del{{.TableName | ToCamelCase}}ByID(id uint64) (int64, *model.CommonError) {
	repo := repository.Create{{.TableName | ToCamelCase}}Repository()
	tool := model.CreateUpdateTool()
	tool.AddSetCase("deleted_at", helper.Now())
	tool.AddWhereCase("id", "=", id)

	result, err := repo.Update{{.TableName | ToCamelCase}}ByTool(tool, nil)
	if err != nil {
		fmt.Printf("delete {{.TableName | ToLowerCamelCase}} fail：%v \n", err)
		return 0, model.CreateCommonError(-5, "delete {{.TableName | ToLowerCamelCase}} fail"+err.Error())
	}

	return result, nil
}

func Select{{.TableName | ToCamelCase}}By(excludeID int,page int, pagesize int) (*model.{{.TableName | ToCamelCase}}Res, *model.CommonError) {
	result := model.{{.TableName | ToCamelCase}}Res{}
	repo := repository.Create{{.TableName | ToCamelCase}}Repository()
	tool := model.CreateSelectTool()
	if excludeID != -1 {
		tool.AddWhereCase("id", "!=", excludeID)
	}

	
	

	total, err := repo.Count{{.TableName | ToCamelCase}}ByTool(tool)
	if err != nil {
		fmt.Printf("select {{.TableName | ToLowerCamelCase}} error ：%v \n", err)
		return nil, model.CreateCommonError(-5, "select {{.TableName | ToLowerCamelCase}} error："+err.Error())
	}

	result.Total = total

	orderBy := " order by created_at desc "
	var sortPage string

	if page != 0 || pagesize != 0 {
		offset := (page - 1) * pagesize
		sortPage = " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(pagesize)
	}
	last := orderBy + sortPage
	tool.Last = last

	{{.TableName | ToLowerCamelCase}}s, err := repo.Select{{.TableName | ToCamelCase}}ByTool(tool)
	if err != nil {
		fmt.Printf("select {{.TableName | ToLowerCamelCase}} error：%v \n", err)
		return nil, model.CreateCommonError(-5, "select {{.TableName | ToLowerCamelCase}} error："+err.Error())
	}

	vos := []*model.{{.TableName | ToCamelCase}}Vo{}

	for i := 0; i < len({{.TableName | ToLowerCamelCase}}s); i++ {
		vo := {{.TableName | ToLowerCamelCase}}[i].To{{.TableName | ToCamelCase}}Vo()
		vos = append(vos, vo)
	}

	result.{{.TableName | ToCamelCase}}s = vos

	return &result, nil
}


func Update{{.TableName | ToCamelCase}}(body *model.{{.TableName | ToCamelCase}}) (int64, *model.CommonError) {
	repo := repository.Create{{.TableName | ToCamelCase}}Repository()

	tool := model.CreateUpdateTool()
	if body.Name != "" {
		tool.AddSetCase("name", body.Name)
	}

	tool.AddSetCase("updated_at", time.Now())
	tool.AddWhereCase("id", "=", body.Id)
	result, err := repo.Update{{.TableName | ToCamelCase}}ByTool(tool, nil)
	if err != nil {
		fmt.Printf("更新标签有误：%v \n", err)
		return 0, model.CreateCommonError(-5, "更新标签有误：")
	}
	return result, nil

}

`

type StructTemplalteOneProxy struct {
	structTplOneProxy string
}

func NewStructTemplateOneProxy() *StructTemplalteOneProxy {
	return &StructTemplalteOneProxy{structTplOneProxy: structTplOneProxy}
}

func (t *StructTemplalteOneProxy) GenerateOneProxy(tableName string, tplColumns []*Column) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"DescTow":          word.DescTow,
		"DescFour":         word.DescFour,
		"DescThree":        word.DescThree,
	}).Parse(t.structTplOneProxy))

	tplDB := StructTemplalteOneDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	fileName := "./internal/sql2struct/proxy/" + tableName + ".go"
	f, err2 := os.Create(fileName)

	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}

const structTplOneValidator = `package validator

func Validate{{.TableName | ToCamelCase}}(a *model.{{.TableName | ToCamelCase}}) *model.CommonError {


	{{$ColumnLen := len .Columns}}{{range $index,$value := .Columns}}{{if le $index ($ColumnLen |DescFour)}}
	if a.{{$value.ColumnName | ToCamelCase}} == ""{
		return model.CreateCommonError(-1,"{{.ColumnComment}}不能为空")
	}
	{{end}}{{end}}
	

	return nil

}


func Validate{{.TableName | ToCamelCase}}Vo(a *model.{{.TableName | ToCamelCase}}Vo) *model.CommonError {

	{{$ColumnLen := len .Columns}}{{range $index,$value := .Columns}}{{if le $index ($ColumnLen |DescFour)}}
	if a.{{$value.ColumnName | ToCamelCase}} == ""{
		return model.CreateCommonError(-1,"{{.ColumnComment}}不能为空")
	}
	{{end}}{{end}}
	


	

	return nil

}



`

type StructTemplalteOneValidator struct {
	structTplOneValidator string
}

func NewStructTemplateOneValidator() *StructTemplalteOneValidator {
	return &StructTemplalteOneValidator{structTplOneValidator: structTplOneValidator}
}

func (t *StructTemplalteOneValidator) AssemblyColumns(tbColumns []*Column) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {

		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", word.UnderscoreToLowerCamelCase(column.ColumnName))
		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DBTypeToStructType[column.DataType],
			Tag:     tag,
			Comment: column.ColumnComment,
		})
	}

	return tplColumns
}

func (t *StructTemplalteOneValidator) GenerateOneValidator(tableName string, tplColumns []*Column) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":      word.UnderscoreToUpperCamelCase,
		"ToLowerCamelCase": word.UnderscoreToLowerCamelCase,
		"DescTow":          word.DescTow,
		"DescFour":         word.DescFour,
		"DescThree":        word.DescThree,
	}).Parse(t.structTplOneValidator))

	tplDB := StructTemplalteOneDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	fileName := "./internal/sql2struct/validator/" + tableName + ".go"
	f, err2 := os.Create(fileName)

	if err2 != nil {
		panic(err2)

	}

	err := tpl.Execute(f, tplDB)
	if err != nil {
		return err

	}
	return nil
}
