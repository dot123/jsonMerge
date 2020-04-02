package main

import (
	"encoding/json"
	"flag"
	log "github.com/sirupsen/logrus"
	_ "howett.net/plist"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	dataDict = make(map[string]interface{})
	input    = "./"
	output   = "./output.json"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})

	startTime := time.Now().UnixNano()

	flag.StringVar(&input, "input", "./", "输入目录")
	flag.StringVar(&output, "output", "./output.json", "输出文件路径")
	flag.Parse()
	filepath.Walk(input, walkFunc)

	writeJSON(output, dataDict)

	endTime := time.Now().UnixNano()
	log.Infof("总耗时:%v毫秒\n", (endTime-startTime)/1000000)
	time.Sleep(time.Millisecond * 1000)
}

func walkFunc(files string, info os.FileInfo, err error) error {
	if err != nil {
		log.Error(err)
		return err
	}
	_, fileName := filepath.Split(files)
	if path.Ext(files) == ".json" {
		data, err := ioutil.ReadFile(filepath.FromSlash(input + "/" + fileName))
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		v := make(map[string]interface{})
		err = json.Unmarshal(data, &v)

		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		dataDict[fileName] = v
	}
	return nil
}

//字典转字符串
func map2Str(dataDict map[string]interface{}) string {
	b, err := json.Marshal(dataDict)
	if err != nil {
		log.Errorln(err)
		return ""
	}
	return string(b)
}

//写JSON文件
func writeJSON(fileName string, dataDict map[string]interface{}) {
	file, err := os.OpenFile(filepath.FromSlash(fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) //不存在创建清空内容覆写
	if err != nil {
		log.Errorln("open file failed.", err.Error())
		return
	}

	defer file.Close()
	//字典转字符串
	file.WriteString(map2Str(dataDict))
}
