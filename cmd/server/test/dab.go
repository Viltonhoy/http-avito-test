package main

import (
	"fmt"
	"math/rand"
	"time"
)

type deposit struct {
	user_id int64
	journal string
	period  string
	cashe   int64
	dat     time.Time
}

type withdrawal struct {
	user_id int64
	journal string
	period  string
	cashe   int64
	dat     time.Time
}

type transfer struct {
	user_id   int64
	journal   string
	period    string
	cashe     int64
	dat       time.Time
	addressee string
}

func main() {
	fmt.Println(depos())
	fmt.Println(withdr())
}

var d deposit
var w withdrawal

func depos() (deposit, deposit) {
	var cb deposit

	rand.Seed(time.Now().Unix())

	d.user_id = rand.Int63n(10)
	d.journal = "deposit"
	d.period = time.Now().Format("2006")
	d.cashe = rand.Int63n(10000) * 100
	d.dat = time.Now()

	cb.user_id = 0
	cb.journal = d.journal
	cb.period = d.period
	cb.cashe = d.cashe * -1
	cb.dat = d.dat

	return d, cb
}

func withdr() (withdrawal, withdrawal) {
	var cb withdrawal

	rand.Seed(time.Now().Unix())

	w.user_id = rand.Int63n(10)
	w.journal = "withdrawal"
	w.period = time.Now().Format("2006")
	w.cashe = rand.Int63n(10000) * -100
	w.dat = time.Now()

	cb.user_id = 0
	cb.journal = w.journal
	cb.period = w.period
	cb.cashe = w.cashe * -1
	cb.dat = w.dat

	return w, cb
}

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// type posting struct {
// 	user_id int64
// 	journal string
// 	period  string
// 	cashe   int64
// 	dat     time.Time
// }

// type postingTable []posting

// func main() {

// 	p := posting{
// 		account_id(),
// 		cb_journal(),
// 		accounting_period(),
// 		amount(),
// 		date(),
// 	}

// 	cb := posting{
// 		0,
// 		cb_journal(),
// 		accounting_period(),
// 		amount(),
// 		date(),
// 	}
// 	var tp postingTable
// 	for i := 0; i < 3; i++ {

// 		switch {
// 		case p.journal == "withdrawal":
// 			tp = append(tp, p, cb)
// 		}
// 	}
// 	fmt.Println(tp)
// }

// func account_id() int64 {
// 	rand.Seed(time.Now().Unix())

// 	return rand.Int63n(10)
// }

// func cb_journal() string {
// 	journal := make([]string, 0)

// 	journal = append(journal, "deposit", "withdrawal", "transfer")

// 	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
// 	message := fmt.Sprint(journal[rand.Intn(len(journal))])
// 	return message
// }

// func accounting_period() string {
// 	date := time.Now()
// 	return date.Format("2006")
// }

// func amount() int64 {
// 	rand.Seed(time.Now().Unix())

// 	return rand.Int63n(10000)
// }

// func date() time.Time {
// 	date := time.Now()
// 	return date
// }
