package wallet

import (
"fmt"
"github.com/MSHE97/wallet/pkg/types"
"github.com/google/uuid"
"reflect"
"testing"
)

type testService struct{
	*Service
}

func newTestService() *testService {
	return &testService{ Service: &Service{}}
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone: 		"+99200000001",
	balance: 	10_000_00,
	payments:	[]struct{
		amount		types.Money
		category	types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}


func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error){
	// регистрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	// пополняем его счёт
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	// выполним платежи
	// можем создать слайс нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments{
		// тогда можно работать здесь просто через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
		}
	}
	return account, payments, nil
}

func TestFindAccountById_success(t *testing.T) {
	svc := &Service{}
	_, err := svc.RegisterAccount("000000001")
	if err != nil {
		t.Errorf("Can't register new account")
		return
	}

	var ID int64 = 1
	account, err := svc.FindAccountById(ID)
	if err != nil {
		t.Errorf("Error founding acount ID: %v ", ID)
	}
	if ID != account.ID{
		t.Errorf("want: %v, got: %v", ID, account.ID)
	}
}

func TestFindAccountById_fail(t *testing.T ) {
	svc := &Service{}
	var ID int64 = 1
	account, err := svc.FindAccountById(ID)
	if err != ErrAccountNotFound {
		t.Errorf("got incorrect error")
	}
	if account != nil {
		t.Errorf("want: nil, got accaunt: %v", account.ID)
	}
}

func TestService_FindPaymentById_success(t *testing.T) {
	// создадим экземпляр сервиса
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}

	// пробуем найти платёж
	payment  := payments[0]
	got, err := s.FindPaymentById(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentById(): error = %v", err)
		return
	}

	// сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentById(): wrong payment returned = %v", err)
	}
}



func TestService_FindPaymentById_fail(t *testing.T ) {
	// создадим экземпляр сервиса
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}

	// пробуем найти не существующий платёж
	_, err = s.FindPaymentById(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentById(): must return error, returned nil")
	}
}

func TestService_Reject_success(t *testing.T){
	// создадим экземпляр сервиса
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}

	// пробуем отменить платёж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentById(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by ID, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountById(payment.AccountId)
	if err != nil {
		t.Errorf("Reject(): can't find account by ID, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance{
		t.Errorf("Reject(): balance didn't change, account = %v", savedAccount)
		return
	}
}

func TestReject_fail(t *testing.T){
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}
	// пробуем отменить не существующий платёж
	err = s.Reject(uuid.New().String())
	if err != ErrPaymentNotFound{
		t.Errorf("Reject(): payment musn't be found")
	}
}

func TestService_Repeat_success(t *testing.T)  {
	// создаём сервис
	s := newTestService()
	// создаём аккаунт с платежом
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}

	//проведём опрерацию повторного платежа Repeat()
	repeatPayID := payments[0].ID
	payment, err := s.Repeat(repeatPayID)
	if err != nil {
		t.Errorf("Repeat(): payment making error = %v", err)
	}
	if payment.Status == types.PaymentStatusFail {
		t.Errorf("Repeat(): failed making repeat payment")
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	// создаём аккаунт с платежом
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}

	// добавим платёж в избранные
	paymentID := payments[0].ID
	name := "One of my favorites!)"
	favorite := &types.Favorite{
		ID: 		"",
		AccountID:  payments[0].AccountId,
		Name: 		name,
		Amount: 	payments[0].Amount,
		Category: 	payments[0].Category,
	}

	got, err := s.FavoritePayment(paymentID, name)
	if err != nil {
		t.Errorf("FavoritePayment(): fail including to favorites, error = %v", err)
	}
	favorite.ID = got.ID
	if !reflect.DeepEqual(got, favorite) {
		t.Errorf("FavoritePayment(): included wrong parameter, got = \n%v\n want = \n%v\n", got, favorite)
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	// создаём аккаунт с платежом
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}

	// добавим платёж в избранные
	var payment types.Payment = *payments[0]
	name := "One of my favorites!)"
	favorite, err := s.FavoritePayment(payment.ID, name)
	if err != nil {
		t.Errorf("PayFromFavorite(): fail creating favorite payment, error = %v", err)
	}

	// выполним платёж из избранного
	favoritePay, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): payment haven't done, err = %v", err)
	}
	var want *types.Payment = &payment
	want.ID = favoritePay.ID
	if !reflect.DeepEqual(favoritePay, want) {
		t.Errorf("PayFromFavorite(): wrong payment parameter, got = \n%v\n want = \n%v\n", favoritePay, want)
	}
}