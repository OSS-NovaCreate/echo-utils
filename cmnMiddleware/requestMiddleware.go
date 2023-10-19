package cmnMiddleware

import (
	"net/http"

	"github.com/OSS-NovaCreate/echo-utils/database"
	"github.com/OSS-NovaCreate/echo-utils/utils/errorUtil"
	"github.com/labstack/echo/v4"
)

// DBコネクションとトランザクションを開始する
func TransactionStart(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//トランザクション開始
		var err error
		database.DbTx, err = database.Database.Beginx()
		if err != nil {
			e := errorUtil.Error(err.Error())
			errorUtil.WriteErrorLog(e)
			return internalServerError(c, e.Error())
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

		return nil
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

// 意図しないエラーのレスポンス
func internalServerError(c echo.Context, log string) error {

	//エラーオブジェクト作成
	errObj := errorResponse{
		Error: "InternalServerError",
	}

	//既存のトランザクションをロールバック
	database.DbTx.Rollback()

	//エラーレスポンス
	return c.JSON(http.StatusInternalServerError, errObj)
}
