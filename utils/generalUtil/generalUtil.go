package generalUtil

import (
	"errors"
	"runtime"
	"strings"
	"time"

	"github.com/OSS-NovaCreate/echo-utils/utils/errorUtil/utilErrors"
)

// 呼び出し元関数名を取得する
func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// 現在時刻を取得する
func Now() time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(jst)
	return now
}

// SQL利用向けにフォーマットした文字列で現在時刻を取得する
func NowSQLFormat() string {
	now := Now()
	return now.Format("2006-01-02 15:04:05")
}

// 引数の時間をSQL向け文字列に変換する
func TimeSQLFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// メールアドレスの要件チェック
func CheckEmailRegCondition(email string) error {
	//アットマークの有無確認
	at := strings.Split(email, "@")
	if len(at) != 2 {
		return errors.New(utilErrors.NOT_EMAIL)
	}

	dot := strings.Split(at[1], ".")
	if len(dot) < 2 {
		return errors.New(utilErrors.NOT_EMAIL)
	}

	return nil
}

// タイムスタンプ取得する
func GetTimestamp() string {
	now := Now()
	return now.Format("20060102150405")
}
