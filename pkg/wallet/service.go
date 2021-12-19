package wallet

import (
	"errors"
	"github.com/MSHE97/wallet/pkg/types"
	"github.com/google/uuid"
)


type Error string

func (e Error) Error() string {
	return string(e)
}

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater then zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite payment not found")
var ErrLowBalance = errors.New("low balance")

type Service struct {
	NextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) FindAccountById(accountID int64) (*types.Account, error){
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			return acc, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentById(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID{
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) FindFavorite(favoriteID string) (*types.Favorite, error){
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID{
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

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
		return nil, ErrLowBalance
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

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentById(paymentID)
	if err == ErrPaymentNotFound {
		return err
	}
	payment.Status = types.PaymentStatusFail
	for _,acc := range s.accounts{
		if payment.AccountId == acc.ID{
			acc.Balance += payment.Amount
			return nil
		}
	}
	return ErrAccountNotFound
}


func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	// в начале найдём платёж по ID
	payment, err := s.FindPaymentById(paymentID)
	if err != nil {
		return nil, err
	}
	// отправляем повторный платёж
	repPay, err := s.Pay(payment.AccountId, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return repPay, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	// ищем платёж по ID
	payment, err := s.FindPaymentById(paymentID)
	if err != nil {
		return nil, err
	}
	// создаём избранный пратёж из найденного
	favorite := &types.Favorite{
		ID:	uuid.New().String(),
		AccountID: payment.AccountId,
		Name:	name,
		Amount:	payment.Amount,
		Category: payment.Category,
	}
	// добавляем его в слайс избранных
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error)  {
	// находим избранный платёж по ID
	favorite, err := s.FindFavorite(favoriteID)
	if err != nil {
		return nil, err
	}

	accountID := favorite.AccountID
	amount	  := favorite.Amount
	category  := favorite.Category

	// совершаем платёж по данным из избранного
	payment, err := s.Pay(accountID, amount, category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}