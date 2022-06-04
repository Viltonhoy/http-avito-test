package main

import (
	"fmt"
)

type Person struct {
	Name string
}

func rename(person *Person) {
	person.Name = "Alice"
}

func main() {

	person := &Person{
		Name: "Bob",
	}

	fmt.Println(person)
	rename(person)
	fmt.Println(person)
	// a := 5355
	// b := -2

	// c := -0.32

	// val := decimal.NewFromInt(int64(a)).Div(decimal.NewFromInt(int64(b)))

	// val2 := decimal.NewFromFloat(c).Exponent()

	// val3 := decimal.New(int64(a), int32(b))
	// fmt.Println(val, val2, val3)

}
