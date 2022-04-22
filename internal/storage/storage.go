package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

var server = "<localhost>"
var port = 5432
var user = "<postgres>"
var password = "<root>"
var database = "<localhost>"

	//строка подключения
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, port, database)

	var err error
	//создать пул соединений
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	//ReadClient
	count, err:=ReadClient()
	if err !=nil {
		log.Fatal("Error reading client", err.Error())
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)

	func ReadClient() (int, error) {
		ctx:= context.Background()

		//Check if DB is alive
		err:=db.PingContext(ctx)
		if err != nil{
			return -1, err
		}

		tsql := fmt.Sprintf("SELECT user_id, bankacc FROM users;")

		//выполнить запрос
		rows, err := db.QueryContext(ctx, tsql)
		if err !=nil{
			return -1, err
		}

		defer rows.Close()

		var count int

		//итерация наборов результатов
		for rows.Next() {
			var user_id, bankacc int64

			//получить значение из строки
			err:= rows.Scan(&user_id, &bankacc)
			if err != nil{
				return -1, err
			}

			fmt.Printf("ID: %d, BANK: %d", user_id, bankacc)
			count++
		}

		return count, nil
	}

	func UpdateClient(id int64, bankacc int64)(int64, error){
		ctx:=context.Background()

		err:=db.PingContext(ctx)
		if err !=nil{
			return -1, err
		}

		tsql:=fmt.Sprintf("insert into users (user_id, bankacc) values (1, 5) on conflict (user_id) do update set bankacc = (select bankacc + 5 from users where user_id = 1) where userss.user_id = 1")
		
		result, err :=db.ExecContext(
			ctx,
			tsql,
			sql.Named("ID", user_id),
			sql.Named("Bank", bankacc))
		if err != nil{
			return -1, err
		}
		return result.RowsAffected()	
	}