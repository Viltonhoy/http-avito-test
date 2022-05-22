package generatetable

import (
	"fmt"
	"math/rand"
	"time"
)

type Posting struct {
	accountID        int64
	CBjournal        operationType
	accountingPeriod int
	amount           int64
	date             time.Time
	addressee        *int64
}

type operationType string

const (
	operationTypeDeposit    operationType = "deposit"
	operationTypeWithdrawal operationType = "withdrawal"
	operationTypeTransfer   operationType = "transfer"
)

const cacheBookAccountID = int64(0)

func GenerateTableData(userCount, totalRecordCount int) []Posting {
	postingTable := make([]Posting, 0, totalRecordCount)
	userTotalBalances := make(map[int64]int64, userCount)

	for i := 1; i <= userCount; i++ {
		year := time.Now().Year()
		amountValue := rand.Int63n(10000000) * 100

		postingTable = append(
			postingTable,
			Posting{
				accountID:        int64(i),
				CBjournal:        operationTypeDeposit,
				accountingPeriod: year,
				amount:           amountValue,
				date:             time.Now(),
			},
			Posting{
				accountID:        cacheBookAccountID,
				CBjournal:        operationTypeDeposit,
				accountingPeriod: year,
				amount:           -1 * amountValue,
				date:             time.Now(),
			},
		)

		userTotalBalances[int64(i)] = amountValue
	}

	casheBookOperation := []string{string(operationTypeDeposit), string(operationTypeTransfer), string(operationTypeWithdrawal)}

	var i int = 1
	n := (totalRecordCount - 2*userCount) / 2
	for i <= n {

		i++
		year := time.Now().Year()

		switch casheBookOperation[rand.Intn(len(casheBookOperation))] {
		case string(operationTypeDeposit):
			accountID := rand.Intn(userCount-1) + 1
			amount := rand.Int63n(10000000) * 100

			postingTable = append(
				postingTable,
				Posting{
					accountID:        int64(accountID),
					CBjournal:        operationTypeDeposit,
					accountingPeriod: year,
					amount:           amount,
					date:             time.Now(),
				},
				Posting{
					accountID:        cacheBookAccountID,
					CBjournal:        operationTypeDeposit,
					accountingPeriod: year,
					amount:           -1 * amount,
					date:             time.Now(),
				},
			)

			userTotalBalances[int64(accountID)] += amount
		case string(operationTypeWithdrawal):
			accountID := rand.Intn(userCount-1) + 1
			if userTotalBalances[int64(accountID)] == 1 {
				i--
				break
			}
			amount := rand.Int63n(userTotalBalances[int64(accountID)]-1) + 1

			postingTable = append(
				postingTable,
				Posting{
					accountID:        int64(accountID),
					CBjournal:        operationTypeWithdrawal,
					accountingPeriod: year,
					amount:           amount * -1,
					date:             time.Now(),
				},
				Posting{
					accountID:        cacheBookAccountID,
					CBjournal:        operationTypeWithdrawal,
					accountingPeriod: year,
					amount:           amount,
					date:             time.Now(),
				},
			)

			userTotalBalances[int64(accountID)] -= amount
		case string(operationTypeTransfer):
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
					accountID:        senderID,
					CBjournal:        operationTypeTransfer,
					accountingPeriod: year,
					amount:           amount * -1,
					date:             time.Now(),
					addressee:        &recipientID,
				},
				Posting{
					accountID:        recipientID,
					CBjournal:        operationTypeTransfer,
					accountingPeriod: year,
					amount:           amount,
					date:             time.Now(),
					addressee:        &senderID,
				},
			)
			userTotalBalances[recipientID] += amount

			userTotalBalances[senderID] = oldAmount - amount
		default:
			fmt.Println("WTF")
		}
	}
	return postingTable
}
