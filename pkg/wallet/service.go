package wallet

import (
	"errors"
	"wallet/pkg/types"

	"github.com/google/uuid"
)

// "github.com/MSHE97/wallet/package/types"

type Error string

func (e Error) Error() string {
	return string(e)
}

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater thwn zero")
var ErrAccountNotFound = errors.New("account not found")

type Service struct {
	NextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

// func RegisterAccount(service *Service, phone types.Phone) {
// 	for _, account := range service.accounts {
// 		if account.Phone == phone {
// 			return
// 		}
// 	}

// 	service.NextAccountID++
// 	service.accounts = append(service.accounts, &types.Account{
// 		ID:      service.NextAccountID,
// 		Phone:   phone,
// 		Balance: 0,
// 	})
// }

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.NextAccountID++
	account := &types.Account{
		ID:      s.NextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount < 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrAccountNotFound
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountId: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

// type UUID [16]byte

// func (uuid UUID) String() string {
// 	var buf [36]byte
// 	encodeHex(buf[:], uuid)
// 	return string(buf[:])
// }
