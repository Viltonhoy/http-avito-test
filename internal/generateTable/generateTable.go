package generatetable

import (
	"fmt"
	"math/rand"
	"time"
)

type Posting struct {
	AccountID        int64
	CBjournal        operationType
	AccountingPeriod string
	Amount           int64
	Date             time.Time
	Addressee        *int64
}

type operationType string

const (
	operationTypeDeposit    operationType = "deposit"
	operationTypeWithdrawal operationType = "withdrawal"
	operationTypeTransfer   operationType = "transfer"
)

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
				CBjournal:        operationTypeDeposit,
				AccountingPeriod: year,
				Amount:           amountValue,
				Date:             time.Now(),
			},
			Posting{
				AccountID:        cacheBookAccountID,
				CBjournal:        operationTypeDeposit,
				AccountingPeriod: year,
				Amount:           -1 * amountValue,
				Date:             time.Now(),
			},
		)

		userTotalBalances[int64(i)] = amountValue
	}

	casheBookOperation := []string{string(operationTypeDeposit), string(operationTypeTransfer), string(operationTypeWithdrawal)}

	var i int = 1
	n := (totalRecordCount - 2*userCount) / 2
	for i <= n {

		i++
		year := fmt.Sprintf(`Period: %d`, time.Now().Year())

		switch casheBookOperation[rand.Intn(len(casheBookOperation))] {
		case string(operationTypeDeposit):
			accountID := rand.Intn(userCount-1) + 1
			amount := rand.Int63n(10000000) * 100

			postingTable = append(
				postingTable,
				Posting{
					AccountID:        int64(accountID),
					CBjournal:        operationTypeDeposit,
					AccountingPeriod: year,
					Amount:           amount,
					Date:             time.Now(),
				},
				Posting{
					AccountID:        cacheBookAccountID,
					CBjournal:        operationTypeDeposit,
					AccountingPeriod: year,
					Amount:           -1 * amount,
					Date:             time.Now(),
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
					AccountID:        int64(accountID),
					CBjournal:        operationTypeWithdrawal,
					AccountingPeriod: year,
					Amount:           amount * -1,
					Date:             time.Now(),
				},
				Posting{
					AccountID:        cacheBookAccountID,
					CBjournal:        operationTypeWithdrawal,
					AccountingPeriod: year,
					Amount:           amount,
					Date:             time.Now(),
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
					AccountID:        senderID,
					CBjournal:        operationTypeTransfer,
					AccountingPeriod: year,
					Amount:           amount * -1,
					Date:             time.Now(),
					Addressee:        &recipientID,
				},
				Posting{
					AccountID:        recipientID,
					CBjournal:        operationTypeTransfer,
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
