package main

import (
    "app/models/hrbc"
	"database/sql" // ①
    "gopkg.in/ini.v1"
    "os"
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

type ConfigList struct {
    Protocol  string
    Port      int
    DbName    string
    SQLDriver string
}

// ConfigList（構造体）を外部パッケージからも読み込めるようにパブリックで変数宣言
var Config ConfigList

const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
    // DataSourceName = "root:golang@tcp(common_hrbc_data_coordination_mysql_1:3310)/golang_db"
	DataSourceName = "root:golang@tcp(mysql)/golang_db" // mysql -u root -p"golang" golang_db
)
var usr = make(map[int]User)

func init() {
    cfg, err := ini.Load("config/local.ini")
    if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
    }

    Config = ConfigList{
        Protocol:  cfg.Section("server").Key("protocol").String(),
        Port:      cfg.Section("server").Key("port").MustInt(),
        DbName:    cfg.Section("db").Key("name").MustString("sample.sql"),
        SQLDriver: cfg.Section("db").Key("driver").String(),
    }
}

func main() {
    fmt.Printf("%T %v\n", Config.Protocol, Config.Protocol)
    fmt.Printf("%T %v\n", Config.Port, Config.Port)
    fmt.Printf("%T %v\n", Config.DbName, Config.DbName)
    fmt.Printf("%T %v\n", Config.SQLDriver, Config.SQLDriver)

    hrbc.SayHello()
	fmt.Println("Hello golang from docker! ")
    // database
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