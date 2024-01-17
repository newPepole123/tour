package main

import (
	"log"

	"github.com/gitdownload/tour/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err: %v", err)
	}
}

//发送

// func SendMail(mailTo []string, subject string, body string, filePath string, fileName string) error {
// 	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码

// 	mailConn := map[string]string{
// 		"user": "niwota0418@163.com",
// 		"pass": "KOCEHJOGINWJEJOU",
// 		"host": "smtp.163.com",
// 		"port": "465",
// 	}

// 	// mailConn := map[string]string{
// 	// 	"user": "1257371603@qq.com",
// 	// 	"pass": "zbajuosixfcxidaf",
// 	// 	"host": "smtp.exmail.qq.com",
// 	// 	"port": "465",
// 	// }

// 	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

// 	m := gomail.NewMessage()

// 	m.SetHeader("From", m.FormatAddress(mailConn["user"], "中数科技"))
// 	// 这种方式可以添加别名，即“XX官方”
// 	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
// 	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
// 	//m.SetHeader("From", mailConn["user"])

// 	m.SetHeader("To", mailTo...)    //发送给多个用户
// 	m.SetHeader("Subject", subject) //设置邮件主题
// 	m.SetBody("text/html", body)    //设置邮件正文

// 	if len(filePath) > 0 {
// 		m.Attach(filePath,
// 			gomail.Rename(fileName),
// 			gomail.SetHeader(map[string][]string{
// 				"Content-Disposition": {
// 					fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", fileName)),
// 				},
// 			}),
// 		)
// 	}

// 	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

// 	err := d.DialAndSend(m)
// 	return err

// }
// func main() {
// 	//定义收件人
// 	mailTo := []string{
// 		"1257371603@qq.com",
// 	}
// 	//邮件主题为"Hello"
// 	subject := "蚁言每天报表"
// 	// 邮件正文
// 	body := "''20211011年报表''"

// 	err := SendMail(mailTo, subject, body, "D:\\20210602\\20211011.xlsx", "20211011.xlsx")
// 	if err != nil {
// 		log.Println(err)
// 		fmt.Println("send fail")
// 		return
// 	}

// 	fmt.Println("send successfully")

// }

// func main() {
// 	// file := xlsx.NewFile()
// 	// sheet, err := file.AddSheet("2021-10-10")
// 	// if err != nil {
// 	// 	fmt.Printf("创建sheet有误：%v \n", err)
// 	// 	return
// 	// }
// 	// row := sheet.AddRow()
// 	// row.SetHeightCM(1) //设置每行的高度
// 	// cell := row.AddCell()
// 	// cell.Value = "haha"
// 	// cell = row.AddCell()
// 	// cell.Value = "1234567"
// 	// err = file.Save("file.xlsx")
// 	// if err != nil {
// 	// 	fmt.Printf("保存文件有误：%v \n", err)
// 	// 	return
// 	// }
// 	ps := []model.OrderReport{}
// 	p := model.OrderReport{}
// 	p.Date = "2021-12-02"
// 	p.FirstLevelAmount = 15.5
// 	p.SecondLevelAmount = 10.5
// 	p.OrderAmount = 50
// 	ps = append(ps, p)
// 	GeneratePeopleExcel(ps)
// }

// func GeneratePeopleExcel(peo []model.OrderReport) (error, bool) {

// 	t := make([]string, 0)
// 	t = append(t, "日期")
// 	t = append(t, "一级分佣金额")
// 	t = append(t, "二级分佣金额")
// 	t = append(t, "订单金额")

// 	file := xlsx.NewFile()
// 	sheet, err := file.AddSheet("2021年11月10号报表")
// 	if err != nil {
// 		fmt.Printf("创建sheet有误：%v \n", err)
// 		return err, false
// 	}
// 	titleRow := sheet.AddRow()
// 	xlsRow := util.NewRow(titleRow, t)
// 	err = xlsRow.SetRowTitle()
// 	if err != nil {
// 		fmt.Printf("设置标题有误：%v \n", err)
// 		return err, false
// 	}
// 	for _, v := range peo {
// 		ostr := v.ToOrderReportStr()
// 		currentRow := sheet.AddRow()
// 		tmp := make([]string, 0)
// 		tmp = append(tmp, ostr.Date)
// 		tmp = append(tmp, ostr.FirstLevelAmount)
// 		tmp = append(tmp, ostr.SecondLevelAmount)
// 		tmp = append(tmp, ostr.OrderAmount)
// 		// tmp = append(tmp, v.Marry)
// 		// tmp = append(tmp, v.Address)

// 		xlsRow := util.NewRow(currentRow, tmp)
// 		err := xlsRow.GenerateRow()
// 		if err != nil {
// 			fmt.Printf("生成行有误：%v \n", err)
// 			return err, false
// 		}
// 	}
// 	err = file.Save("2021-12-02.xls")
// 	if err != nil {
// 		fmt.Printf("保存xlsx有误：%v \n", err)
// 		return err, false
// 	}

// 	// time.Sleep(time.Duration(5) * time.Second)

// 	// os.Remove("人员信息.xls")
// 	return nil, true
// }
