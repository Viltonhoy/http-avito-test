package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func NewStore() error { //*Store, error
	conf := config.New()

	//строка подключения
	var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		conf.Server, conf.User, conf.Port, conf.Database)

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
	return nil
}

func (s *Store) Close() {
	s.db.Close()
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

func ReadClient(s *Store) error {

	ctx := context.Background()

	//Check if DB is alive
	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf("SELECT * FROM bankacc;")

	//выполнить запрос
	rows, err := s.db.QueryContext(ctx, tsql)
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

func UpdateClient(s *Store, user_id int64, balance int64) error {

	ctx := context.Background()
	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf(`INSERT INTO bankacc (user_id, balance) 
			VALUES ($ID, $BALANCE) ON CONFLICT (user_id) 
			DO UPDATE SET balance = (SELECT balance + $BALANCE FROM bankacc WHERE user_id = $ID) 
			WHERE bankacc.user_id = $ID`)

	_, err = s.db.ExecContext(
		ctx,
		tsql,
		sql.Named("ID", user_id),
		sql.Named("BALANCE", balance))
	if err != nil {
		return err
	}
	return nil
}

func Transaction(s *Store, user_id1, user_id2 int64, balance1, balance2 int64) error {

	ctx := context.Background()
	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	ftsql := fmt.Sprintf(`UPDATE bankacc SET balance = (SELECT balance + -$BALANCE FROM bankacc WHERE user_id = $ID);`)

	_, err = s.db.ExecContext(
		ctx,
		ftsql,
		sql.Named("ID", user_id1),
		sql.Named("BALANCE", balance1))
	if err != nil {
		tx.Rollback()
		return err
	}

	stsql := fmt.Sprintf(`INSERT INTO bankacc (user_id, balance) VALUES ($ID, $VALUE) ON CONFLICT (user_id) 
	do UPDATE SET balance = (SELECT balance + $BALANCE FROM bankacc WHERE user_id = $ID);`)

	_, err = s.db.ExecContext(
		ctx,
		stsql,
		sql.Named("ID", user_id2),
		sql.Named("BALANCE", balance2))
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}
