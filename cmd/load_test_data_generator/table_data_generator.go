package main

import (
	"fmt"
	"http-avito-test/internal/storage"
	"math/rand"
	"time"
)

type Posting struct {
	AccountID        int64
	CBjournal        storage.OperationType
	AccountingPeriod string
	Amount           int64
	Date             time.Time
	Addressee        *int64
}

const cacheBookAccountID = int64(0)

//  GenerateTableData generates a slice of user data values ​​that will be added to the table for performance tests.
//  Function takes userCount and totalRecordCount int values.
func GenerateTableData(userCount, totalRecordCount int) []Posting {
	postingTable := make([]Posting, 0, totalRecordCount)
	userTotalBalances := make(map[int64]int64, userCount)

	// Generating first values ​​for each user id with big amount values for
	//
	for i := 1; i <= userCount; i++ {
		year := fmt.Sprintf(`Period: %d`, time.Now().Year())
		amountValue := rand.Int63n(10000000) * 100

		postingTable = append(
			postingTable,
			Posting{
				AccountID:        int64(i),
				CBjournal:        storage.OperationTypeDeposit,
				AccountingPeriod: year,
				Amount:           amountValue,
				Date:             time.Now(),
			},
			Posting{
				AccountID:        cacheBookAccountID,
				CBjournal:        storage.OperationTypeDeposit,
				AccountingPeriod: year,
				Amount:           -1 * amountValue,
				Date:             time.Now(),
			},
		)

		userTotalBalances[int64(i)] = amountValue
	}

	casheBookOperation := []string{string(storage.OperationTypeDeposit), string(storage.OperationTypeTransfer), string(storage.OperationTypeWithdrawal)}

	var i int = 1
	n := (totalRecordCount - 2*userCount) / 2
	for i <= n {

		i++
		year := fmt.Sprintf(`Period: %d`, time.Now().Year())

		switch casheBookOperation[rand.Intn(len(casheBookOperation))] {
		case string(storage.OperationTypeDeposit):
			accountID := rand.Intn(userCount-1) + 1
			amount := rand.Int63n(10000000) * 100

			postingTable = append(
				postingTable,
				Posting{
					AccountID:        int64(accountID),
					CBjournal:        storage.OperationTypeDeposit,
					AccountingPeriod: year,
					Amount:           amount,
					Date:             time.Now(),
				},
				Posting{
					AccountID:        cacheBookAccountID,
					CBjournal:        storage.OperationTypeDeposit,
					AccountingPeriod: year,
					Amount:           -1 * amount,
					Date:             time.Now(),
				},
			)

			userTotalBalances[int64(accountID)] += amount
		case string(storage.OperationTypeWithdrawal):
			accountID := rand.Intn(userCount-1) + 1
			if userTotalBalances[int64(accountID)] == 1 {
				i--
				break
			}
			amount := rand.Int63n(userTotalBalances[int64(accountID)]-1) + 1

			postingTable = append(
				postingTable,
				Posting{
					AccountID:        int64(accountID),
					CBjournal:        storage.OperationTypeWithdrawal,
					AccountingPeriod: year,
					Amount:           amount * -1,
					Date:             time.Now(),
				},
				Posting{
					AccountID:        cacheBookAccountID,
					CBjournal:        storage.OperationTypeWithdrawal,
					AccountingPeriod: year,
					Amount:           amount,
					Date:             time.Now(),
				},
			)

			userTotalBalances[int64(accountID)] -= amount
		case string(storage.OperationTypeTransfer):
			var senderID, oldAmount, amount int64

			for k, v := range userTotalBalances {
				senderID = k
				oldAmount = v
				break
			}

			if userTotalBalances[senderID] == 1 {
				i--
				break
			}

			amount = rand.Int63n(oldAmount-1) + 1

			delete(userTotalBalances, senderID)

			var recipientID int64
			for k := range userTotalBalances {
				recipientID = k
				break
			}

			postingTable = append(
				postingTable,
				Posting{
					AccountID:        senderID,
					CBjournal:        storage.OperationTypeTransfer,
					AccountingPeriod: year,
					Amount:           amount * -1,
					Date:             time.Now(),
					Addressee:        &recipientID,
				},
				Posting{
					AccountID:        recipientID,
					CBjournal:        storage.OperationTypeTransfer,
					AccountingPeriod: year,
					Amount:           amount,
					Date:             time.Now(),
					Addressee:        &senderID,
				},
			)
			userTotalBalances[recipientID] += amount

			userTotalBalances[senderID] = oldAmount - amount
		}
	}
	return postingTable
}

// import (
// 	"fmt"
// 	"http-avito-test/internal/storage"
// 	"math/rand"
// 	"time"

// 	"github.com/shopspring/decimal"
// )

// const cacheBookAccountID = int64(0)

// //reuse
// type posting struct {
// 	storage.ReadUserHistoryResult
// 	AccountingPeriod string
// }

// //  GenerateTableData generates a slice of user data values ​​that will be added to the table for performance tests.
// //  Function takes userCount and totalRecordCount int values.
// func GenerateTableData(userCount, totalRecordCount int) []posting {
// 	postingTable := make([]posting, 0, totalRecordCount)
// 	userTotalBalances := make(map[int64]decimal.Decimal, userCount)

// 	// Generating first values ​​for each user id with big amount values for
// 	//
// 	for i := 1; i <= userCount; i++ {
// 		year := fmt.Sprintf(`Period: %d`, time.Now().Year())
// 		amountValue := decimal.NewFromInt(rand.Int63n(10000000) * 100)

// 		postingTable = append(
// 			postingTable,
// 			posting{
// 				storage.ReadUserHistoryResult{
// 					AccountID: int64(i),
// 					CBjournal: storage.OperationTypeDeposit,
// 					Amount:    amountValue,
// 					Date:      time.Now(),
// 				},
// 				year,
// 			},
// 			posting{
// 				storage.ReadUserHistoryResult{
// 					AccountID: cacheBookAccountID,
// 					CBjournal: storage.OperationTypeDeposit,
// 					Amount:    amountValue.Mul(decimal.NewFromInt(-1)),
// 					Date:      time.Now(),
// 				},
// 				year,
// 			},
// 		)

// 		userTotalBalances[int64(i)] = amountValue
// 	}

// 	casheBookOperation := []string{string(storage.OperationTypeDeposit), string(storage.OperationTypeTransfer), string(storage.OperationTypeWithdrawal)}

// 	var i int = 1
// 	n := (totalRecordCount - 2*userCount) / 2
// 	for i <= n {

// 		i++
// 		year := fmt.Sprintf(`Period: %d`, time.Now().Year())

// 		switch casheBookOperation[rand.Intn(len(casheBookOperation))] {
// 		case string(storage.OperationTypeDeposit):
// 			accountID := rand.Intn(userCount-1) + 1
// 			amount := decimal.NewFromInt(rand.Int63n(10000000) * 100)

// 			postingTable = append(
// 				postingTable,
// 				posting{
// 					storage.ReadUserHistoryResult{
// 						AccountID: int64(accountID),
// 						CBjournal: storage.OperationTypeDeposit,
// 						Amount:    amount,
// 						Date:      time.Now(),
// 					},
// 					year,
// 				},
// 				posting{
// 					storage.ReadUserHistoryResult{
// 						AccountID: cacheBookAccountID,
// 						CBjournal: storage.OperationTypeDeposit,
// 						Amount:    amount.Mul(decimal.NewFromInt(-1)),
// 						Date:      time.Now(),
// 					},
// 					year,
// 				},
// 			)
// 			decimal.
// 				userTotalBalances[int64(accountID)] += amount
// 		case string(storage.OperationTypeWithdrawal):
// 			accountID := rand.Intn(userCount-1) + 1
// 			if userTotalBalances[int64(accountID)] == decimal.NewFromInt(1) {
// 				i--
// 				break
// 			}
// 			amount := decimal.NewFromInt(rand.Int63n(userTotalBalances[int64(accountID)]-decimal.NewFromInt(-1)) + 1)

// 			postingTable = append(
// 				postingTable,
// 				posting{
// 					storage.ReadUserHistoryResult{
// 						AccountID: int64(accountID),
// 						CBjournal: storage.OperationTypeWithdrawal,
// 						Amount:    amount.Mul(decimal.NewFromInt(-1)),
// 						Date:      time.Now(),
// 					},
// 					year,
// 				},
// 				posting{
// 					storage.ReadUserHistoryResult{
// 						AccountID: cacheBookAccountID,
// 						CBjournal: storage.OperationTypeWithdrawal,
// 						Amount:    amount,
// 						Date:      time.Now(),
// 					},
// 					year,
// 				},
// 			)

// 			userTotalBalances[int64(accountID)] -= amount
// 		case string(storage.OperationTypeTransfer):
// 			var senderID int64
// 			var oldAmount, amount decimal.Decimal
// 			for k, v := range userTotalBalances {
// 				senderID = k
// 				oldAmount = v
// 				break
// 			}

// 			if userTotalBalances[senderID] == 1 {
// 				i--
// 				break
// 			}

// 			amount = decimal.NewFromInt(rand.Int63n(oldAmount.Add(decimal.NewFromInt(-1))) + 1)

// 			delete(userTotalBalances, senderID)

// 			var recipientID int64
// 			for k := range userTotalBalances {
// 				recipientID = k
// 				break
// 			}

// 			postingTable = append(
// 				postingTable,
// 				posting{
// 					storage.ReadUserHistoryResult{
// 						AccountID: senderID,
// 						CBjournal: storage.OperationTypeTransfer,
// 						Amount:    amount.Mul(decimal.NewFromInt(-1)),
// 						Date:      time.Now(),
// 						Addressee: &recipientID,
// 					},
// 					year,
// 				}, posting{
// 					storage.ReadUserHistoryResult{
// 						AccountID: recipientID,
// 						CBjournal: storage.OperationTypeTransfer,
// 						Amount:    amount,
// 						Date:      time.Now(),
// 						Addressee: &senderID,
// 					},
// 					year},
// 			)
// 			userTotalBalances[recipientID] += amount

// 			userTotalBalances[senderID] = oldAmount.Add(amount.Mul(decimal.NewFromInt(-1)))
// 		}
// 	}
// 	return postingTable
// }
