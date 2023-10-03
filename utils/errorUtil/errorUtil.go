package errorUtil

import (
	"errors"
	"os"
	"runtime"
	"strconv"

	"github.com/OSS-NovaCreate/echo-utils/utils/generalUtil"
)

func WriteErrorLog(log error) {
	file, _ := os.OpenFile("error.txt", os.O_WRONLY|os.O_APPEND, 0666)
	info, err := file.Stat()
	if err == nil {
		if info.Size() >= 100000000 {
			//ファイルサイズが100MB程度だった場合リネームする
			file.Close()
			os.Rename("error.txt", "error"+generalUtil.Now().Format("20060102150405")+".txt")
			file, _ = os.Create("error.txt")
		} else {
			//エラーを書き込む
			file.WriteString(log.Error())
			file.Close()
		}
	}
	data := []byte(log.Error())
	file.Write(data)
	file.Close()
}

// エラー発生個所の情報を取得する
func getErrorFuncInfo() string {
	pc, file, line, _ := runtime.Caller(2)
	return "file:" + file + " func:" + runtime.FuncForPC(pc).Name() + " line:" + strconv.Itoa(line) + " :: "
}

// エラー情報
func Error(log string) error {
	ei := getErrorFuncInfo()
	return errors.New(generalUtil.Now().Format("2006-01-02 15:04:05") + " " + ei + log + "\n")
}

// 関数名取得（エラーレスポンス用）
func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}
