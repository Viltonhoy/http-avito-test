package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

const (
	server   = "<localhost>"
	port     = 5432
	user     = "<postgres>"
	password = "<root>"
	database = "<localhost>"
)

func NewStore() (*Store, error) {
	//строка подключения
	var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, port, database)

	//создать пул соединений
	var db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
}

func hand() {
	//ReadClient
	err := ReadClient()
	if err != nil {
		log.Fatal("Error reading client", err.Error())
	}

	//UpdateClient
	err := UpdateClient(1, 1)
	if err != nil {
		log.Fatal("Error updating Employee: ", err.Error())
	}

}

func ReadClient() error {

	ctx := context.Background()

	//Check if DB is alive
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf("SELECT user_id, balance FROM bankacc;")

	//выполнить запрос
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return err
	}

	defer rows.Close()

	//итерация наборов результатов
	for rows.Next() {
		var user_id, balance int64

		//получить значение из строки
		err := rows.Scan(&user_id, &balance)
		if err != nil {
			return err
		}

		fmt.Printf("ID: %d, BALANCE: %d", user_id, balance)
	}

	return nil
}

func UpdateClient(user_id int64, balance int64) error {

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf(`INSERT INTO bankacc (user_id, balance) 
			VALUES ($ID, $BALANCE) ON CONFLICT (user_id) 
			DO UPDATE SET balance = (SELECT balance + $BALANCE FROM bankacc WHERE user_id = $ID) 
			WHERE bankacc.user_id = $ID`)

	_, err = db.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", user_id),
		sql.Named("BALANCE", balance))
	if err != nil {
		return err
	}
	return nil
}

func Transaction(user_id1, user_id2 int64, balance1, balance2 int64) error {

	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Printf(`...`)

	_, err = db.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", user_id),
		sql.Named("BALANCE", balance))
	if err != nil {
		return err
	}
	return nil
}
