package middleware

import (
	"net/http"

	"github.com/OSS-NovaCreate/echo-utils/database"
	"github.com/OSS-NovaCreate/echo-utils/utils/errorUtil"
	"github.com/OSS-NovaCreate/echo-utils/utils/requestUtil"
	"github.com/labstack/echo/v4"
)

func Request(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestUtil.SetContext(c)
		requestUtil.SetUserAgent(c.Request().UserAgent())
		next(c)
		return nil
	}
}

// DBコネクションとトランザクションを開始する
func ConnectionStart(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//DBとコネクションする
		err := database.Open()
		if err != nil {
			return internalServerError(errorUtil.Error(err.Error()).Error())
		}

		//トランザクション開始
		database.DbTx, err = database.Database.Beginx()
		if err != nil {
			e := errorUtil.Error(err.Error())
			errorUtil.WriteErrorLog(e)
			return internalServerError(e.Error())
		}

		next(c)
		return nil
	}
}

// コミットしてDBのコネクションを終了する
func ConnectionEnd(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		next(c)
		database.DbTx.Commit()
		database.Database.Close()

		return nil
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

// 意図しないエラーのレスポンス
func internalServerError(log string) error {

	//エラーオブジェクト作成
	errObj := errorResponse{
		Error: "InternalServerError",
	}

	//既存のトランザクションをロールバック
	database.DbTx.Rollback()

	//コンテキストを取得
	c := requestUtil.GetContext()

	//エラーレスポンス
	return c.JSON(http.StatusInternalServerError, errObj)
}
