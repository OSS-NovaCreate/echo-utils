package database

import (
	"os"

	"github.com/OSS-NovaCreate/echo-utils/utils/errorUtil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DatabaseInfo struct {
	host     string
	port     string
	user     string
	pass     string
	database string
}

var Database *sqlx.DB
var DbTx *sqlx.Tx

func costruct() (DatabaseInfo, error) {
	database := new(DatabaseInfo)
	database.host = os.Getenv("DB_HOST")
	database.port = os.Getenv("DB_PORT")
	database.user = os.Getenv("DB_USER")
	database.pass = os.Getenv("DB_PASS")
	database.database = os.Getenv("DB_DATABASE")

	return *database, nil
}

// DB接続
func Open() error {
	//環境変数から接続情報を取得する
	db, err := costruct()
	if err != nil {
		return err
	}

	//DBに接続する
	Database, err = sqlx.Open("mysql", db.user+":"+db.pass+"@tcp("+db.host+":"+db.port+")/"+db.database)
	if err != nil {
		//エラー処理
		e := errorUtil.Error(err.Error())
		errorUtil.WriteErrorLog(e)
		return e
	}

	return nil
}

// select文
func Select(sql string) (*sqlx.Rows, error) {
	//SQLを実行する
	rows, err := Database.Queryx(sql)
	if err != nil {
		//エラー処理
		e := errorUtil.Error(err.Error())
		return nil, e
	}

	return rows, nil
}

// insert文
func Insert(sql string) error {
	//SQLを実行する
	_, err := DbTx.Exec(sql)
	if err != nil {
		return errorUtil.Error(err.Error())
	}
	return nil
}

// update文
func Update(sql string) (int, error) {
	//SQLを実行する
	result, err := DbTx.Exec(sql)
	if err != nil {
		return 0, errorUtil.Error(err.Error() + "SQL:" + sql)
	}

	//更新した件数
	updateNum, err := result.RowsAffected()
	if err != nil {
		return 0, errorUtil.Error(err.Error())
	}
	return int(updateNum), nil
}

// トランザクション無しInsert文
func InsertNonTx(sql string) error {
	//SQLを実行する
	_, err := Database.Exec(sql)
	if err != nil {
		return errorUtil.Error(err.Error())
	}
	return nil
}

// IDを返すInsert文
func InsertReturnId(sql string) (int, error) {
	//SQLを実行する
	result, err := DbTx.Exec(sql)
	if err != nil {
		return 0, errorUtil.Error(err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errorUtil.Error(err.Error())
	}

	return int(id), nil
}
