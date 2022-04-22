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

func main(){
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



	//UpdateClient
	update, err:= UpdateClient("...")
	if err != nil{
		log.Fatal("Error updating Employee: ", err.Error())
	}
	fmt.Printf("Updated %d row(s) successfully.\n", update)




	func ReadClient() (int, error) {
		ctx:= context.Background()

		//Check if DB is alive
		err:=db.PingContext(ctx)
		if err != nil{
			return -1, err
		}

		tsql := fmt.Sprintf("SELECT user_id, balance FROM bankacc;")

		//выполнить запрос
		rows, err := db.QueryContext(ctx, tsql)
		if err !=nil{
			return -1, err
		}

		defer rows.Close()

		var count int

		//итерация наборов результатов
		for rows.Next() {
			var user_id, balance int64

			//получить значение из строки
			err:= rows.Scan(&user_id, &balance)
			if err != nil{
				return -1, err
			}

			fmt.Printf("ID: %d, BANK: %d", user_id, balance)
			count++
		}

		return count, nil
	}



	func UpdateClient(id int64, balance int64)(int64, error){
		ctx:=context.Background()

		err:=db.PingContext(ctx)
		if err !=nil{
			return -1, err
		}

		tsql:=fmt.Sprintf(`INSERT INTO bankacc (user_id, balance) 
			VALUES (@ID, @Bank) ON CONFLICT (user_id) 
			DO UPDATE SET BALANCE = (SELECT BALANCE + @Bank FROM bankacc WHERE user_id = @ID) 
			WHERE bankacc.user_id = @ID`)
		
		result, err :=db.ExecContext(
			ctx,
			tsql,
			sql.Named("ID", user_id),
			sql.Named("Bank", balance))
		if err != nil{
			return -1, err
		}
		return result.RowsAffected()	
	}

