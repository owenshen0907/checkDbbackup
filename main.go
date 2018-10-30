// checkDbbackup project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	//	"path/filepath"
	"flag"
	"strconv"
	"strings"
	"time"

	"github.com/larspensjo/config"
)

var TOPIC = make(map[string]string)
var ti = time.Now().Format("20060102")

func main() {
	readLogin()
	contentFileName := TOPIC["BodyPrefix"] + ti + TOPIC["BodyStuffix"]

	contentFile, _ := os.OpenFile(contentFileName, os.O_RDWR|os.O_CREATE, 0666)
	SEEK_END, _ := contentFile.Seek(0, os.SEEK_END)

	defer contentFile.Close()
	_, _ = contentFile.WriteAt([]byte("\r\n"), SEEK_END)
	curDate := timeMinusDay(0)
	beforDay := timeMinusDay(1)

	checkCDD(contentFile, curDate, beforDay)
	contentFile.WriteString("-----------------------\r\n")
	checkFL(contentFile, curDate, beforDay)
	b, _ := ioutil.ReadFile(TOPIC["signature"])
	contentFile.WriteString(string(b))

}
func checkCDD(contentFile *os.File, curDate string, beforDay string) {
	TheOnebuckupCDD := "YLCHCDD_" + beforDay + "_23"
	TheTwobuckupCDD := "YLCHCDD_" + curDate + "_01"
	CDDattechFile := "YLCHCDD_" + beforDay

	oneCDDDataFlag := "未检测到" + beforDay + ":车抵贷23点数据库备份文件，请检查问题！"
	twoCDDDataFlag := "未检测到" + curDate + ":车抵贷1点数据库备份文件，请检查问题！"

	//CDDAPPFlag := "车抵贷当前备份应用为" + beforDay + "升级程序。"
	CDDFileFlag := "车抵贷未检测到" + beforDay + "当天附件备份。请检查问题！"

	files1, _ := ListDir(TOPIC["sourceDataFile1"], ".sql.zip")
	//fmt.Println(len(files1))
	for _, v := range files1 {
		fmt.Println(v)
		if strings.Contains(v, TheOnebuckupCDD) {
			oneCDDDataFlag = beforDay + ":车抵贷23点数据库已正常备份。"
		}
		if strings.Contains(v, TheTwobuckupCDD) {
			twoCDDDataFlag = curDate + ":车抵贷1点数据库已正常备份。"

		}
	}
	attechFileCount := checkAttech(TOPIC["linkSqlCDD"], TOPIC["sqlCheckAttech"])
	if attechFileCount == 0 {
		fmt.Println("当日生成附件为空！！")
		CDDFileFlag = beforDay + ":车抵贷当天没有合同录入未生成附件。"
	} else {
		filesAttech1, _ := ListDir(TOPIC["sourceAttachFile1"], ".zip")
		for _, v := range filesAttech1 {
			fmt.Println(v)
			if strings.Contains(v, CDDattechFile) {
				CDDFileFlag = beforDay + ":车抵贷当天生成" + strconv.Itoa(attechFileCount) + "个附件，附件已正常备份。\r\n" + "附件压缩路径：" + v
			}
		}
	}

	fmt.Println(oneCDDDataFlag)
	fmt.Println(twoCDDDataFlag)
	fmt.Println(CDDFileFlag)
	//fmt.Println(CDDAPPFlag)
	contentFile.WriteString(oneCDDDataFlag + "\r\n")
	contentFile.WriteString(twoCDDDataFlag + "\r\n")
	contentFile.WriteString(CDDFileFlag + "\r\n")
	//contentFile.WriteString(CDDFileFlag + "\r\n")
}
func checkFL(contentFile *os.File, curDate string, beforDay string) {

	TheOnebuckupFL := "YLCHFL_" + beforDay + "_23"
	TheTwobuckupFL := "YLCHFL_" + curDate + "_01"
	FLattechFile := "YLCHFL_" + beforDay

	oneFLDataFlag := "未检测到" + beforDay + ":融资租赁23点数据库备份文件，请检查问题！"
	twoFLDataFlag := "未检测到" + curDate + ":融资租赁1点数据库备份文件，请检查问题！"

	//FLAPPFlag := "融资租赁当前备份应用为" + beforDay + "升级程序。"
	FLFileFlag := "融资租赁未检测到" + beforDay + "当天附件备份。请检查问题！"

	files2, _ := ListDir(TOPIC["sourceDataFile2"], ".sql.zip")
	//fmt.Println(len(files2))
	for _, v := range files2 {
		fmt.Println(v)
		if strings.Contains(v, TheOnebuckupFL) {
			oneFLDataFlag = beforDay + ":融资租赁23点数据库已正常备份。"
		}
		if strings.Contains(v, TheTwobuckupFL) {
			twoFLDataFlag = curDate + ":融资租赁1点数据库已正常备份。"
		}
	}
	attechFileCount := checkAttech(TOPIC["linkSqlRZZL"], TOPIC["sqlCheckAttech"])
	if attechFileCount == 0 {
		fmt.Println("当日生成附件为空！！")
		FLFileFlag = beforDay + ":融资租赁当天没有合同录入未生成附件。"
	} else {
		filesAttech2, _ := ListDir(TOPIC["sourceAttachFile2"], ".zip")
		for _, v := range filesAttech2 {
			fmt.Println(v)
			if strings.Contains(v, FLattechFile) {
				FLFileFlag = beforDay + ":融资租赁当天生成" + strconv.Itoa(attechFileCount) + "个附件，附件已正常备份。\r\n" + "附件压缩路径：" + v
			}
		}
	}

	//	filesAttech2, _ := ListDir(TOPIC["sourceAttachFile2"], ".sql.zip")

	fmt.Println()
	fmt.Println(oneFLDataFlag)
	fmt.Println(twoFLDataFlag)
	fmt.Println(FLFileFlag)
	//fmt.Println(FLAPPFlag)

	contentFile.WriteString(oneFLDataFlag + "\r\n")
	contentFile.WriteString(twoFLDataFlag + "\r\n")
	contentFile.WriteString(FLFileFlag + "\r\n")
	//contentFile.WriteString(FLFileFlag + "\r\n")

}

func checkAttech(linkInfo, script string) int {
	linkSql(linkInfo)
	attechFileArr := ReadData(script)
	fmt.Println(len(attechFileArr))
	return len(attechFileArr)
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10) //初始化file切片，预留十个位置
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

func timeMinusDay(day int64) string {
	timestamp := time.Now().Unix()
	timestamp = timestamp - 60*60*24*day
	ti := time.Unix(timestamp, 0)
	return ti.Format("20060102")
}
func readLogin() {
	var (
		configFile = flag.String("configfile", "config.ini", "General configuration file")
	)
	flag.Parse()
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		fmt.Println("read ini error")
		return
	}
	if cfg.HasSection("exe") {
		section, err := cfg.SectionOptions("exe")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("exe", v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
}
