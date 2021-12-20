package main

import (
	"database/sql" // ①
	"log"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // ②
)

// User db users
type User struct {
	ID       int
	Name     string
	Password string
}
const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
    // DataSourceName = "root:golang@tcp(common_hrbc_data_coordination_mysql_1:3310)/golang_db"
	DataSourceName = "root:golang@tcp(mysql)/golang_db"
)
var usr = make(map[int]User)

func main() {
    // database
	fmt.Println("Hello golang from docker! ")
    // DBにアクセスするためのオブジェクトを取得
	db, dbErr := sql.Open(DriverName, DataSourceName)
	if dbErr != nil {
		log.Print("error sql.Open:", dbErr)
	}
    defer db.Close()

    // DB接続を確認。接続されるまで試行（10回まで）
    for i := 0; i < 10; i++ {
		if err := db.Ping(); err != nil {
            fmt.Println(i)
            log.Print("PingError: ", err)
            continue;
        }
        fmt.Println("DB access success ")
        break
	}
    
    // usersテーブルの全てのレコードを取得するクエリの実行 ②
    rows, queryErr := db.Query("SELECT * FROM users")
    if queryErr != nil {
        log.Print("query error :", queryErr)
    }
    defer rows.Close()
    // ループを回してrowsからScanでデータを取得する。 ③
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Password); err != nil {
            log.Print(err)
        }
		
        usr[u.ID] = User{
            ID:       u.ID,
            Name:     u.Name,
            Password: u.Password,
        }

    }
    log.Print(usr)
}