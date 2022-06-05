package main

import (
	"go.uber.org/zap"
)

type Person struct {
	Name string
}

func rename(person *Person) {
	person.Name = "Alice"
}

func main() {

	// a := 5355
	// b := -2

	// c := -0.32

	// val := decimal.NewFromInt(int64(a)).Div(decimal.NewFromInt(int64(b)))

	// val2 := decimal.NewFromFloat(c).Exponent()

	// val3 := decimal.New(int64(a), int32(b))
	// fmt.Println(val, val2, val3)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	user_id := 3

	logger.With(zap.Int("user_id", user_id))

	logger.Info(`sdfsdf`)

	logger.Info(`qwer`)

}
