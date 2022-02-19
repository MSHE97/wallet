package wallet

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"wallet/pkg/types"
)

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

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID{
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) FindFavoriteById(favoriteID string) (*types.Favorite, error){
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
	var targetPayment *types.Payment
	var targetAccount *types.Account
	targetPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return ErrPaymentNotFound
	}
	
	targetAccount, err = s.FindAccountById(targetPayment.AccountId)
	if err != nil {
		return ErrAccountNotFound
	}
	targetPayment.Status = types.PaymentStatusFail
	targetAccount.Balance += targetPayment.Amount
	return nil
}


func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	// в начале найдём платёж по ID
	payment, err := s.FindPaymentByID(paymentID)
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
	payment, err := s.FindPaymentByID(paymentID)
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
	favorite, err := s.FindFavoriteById(favoriteID)
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

func (s *Service) ExportToFile(path string) error {
	// создаём или перезаписываем файл
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err := file.Close()
		if err !=nil {
			log.Print(err)
		}
	}()

	// записываем аккаунты в файл
	var line string
	for _, acc := range s.accounts {
		line += strconv.FormatInt( acc.ID, 10) + ";" + string(acc.Phone) +
				";" + strconv.FormatInt(int64(acc.Balance), 10) + "|"
	}
	_, err = file.Write([]byte( line[:len(line)-1] ))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (s *Service) ImportFromFile(path string) error {
	// открываем файл
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err := file.Close()
		if err !=nil {
			log.Print(err)
		}
	}()

	// считываем файл
	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF { // файл закончился
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := string(content)
	//log.Print(data)
	// парсим данные
	rows := strings.Split(data, "|")
	var id int64
	var phone types.Phone
	var amount int64
	for _, row := range rows{
		metaData := strings.Split(row, ";")
		id, err = strconv.ParseInt(metaData[0], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}
		phone = types.Phone(metaData[1])
		amount, err = strconv.ParseInt(metaData[2], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}
		account := &types.Account{
			ID:     id,
			Phone:   phone,
			Balance: types.Money(amount),
		}

		s.accounts = append(s.accounts, account)
	}
	return nil
}

func (s *Service) Export(path string) error {
	return nil
}

func (s *Service) Import(path string) error {
	return nil
}