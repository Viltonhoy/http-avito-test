package main

import (
	"fmt"
)

func main() {

	var name string
	fmt.Println("Как тебя зовут?")
	fmt.Scanf("%s\n", &name)

	var age int
	fmt.Println("Сколько тебе лет?")
	fmt.Scanf("%d\n", &age)

	fmt.Printf("Привет, %s, твой возраст - %d\n", name, age)

}
