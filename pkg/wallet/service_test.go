package wallet

import (
	"testing"
	"github.com/MSHE97/wallet/pkg/types"
)

type testService struct{
	*Service
}

func newTestService()	*testService{
	return &testAccount{Service: &Service{}}
}

type testAccount struct{
	phone 	types.Phone
	balance types.Money
	payments []struct{
		amount 		types.Money
		category	types.PaymentCategory
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

func (s *testService) addAccount(data testAccount) (types.Account, []*types.Payment, error){
	// регистрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	// пополняем его счёт
	err = s.Deposite(account.ID, data.Balance)
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
		got, err := s.FindPayById(payment.ID)
		if err != nil {
			t.Errorf("FindPaymentById(): error = %v", err)
			return
		}

		// сравниваем платежи
		if !reflect.DeepEqual(payment, got) {
			t.Errorf("FindPaymentById(): wrong payment returned = %v", err)
			return
		}
}

func TestService_FindPaymentById_fail(t *testing.T ) {
	// создадим экземпляр сервиса
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}
	
	// пробуем найти не существующий платёж
	_, err := s.FindPaymentById(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentById(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentById(): must return ErrPaymentNotFound, returned %v", err)
		return
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
	if err := nil {
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
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error()
		return
	}
	// пробуем отменить не существующий платёж
	err := s.Reject(uuid.New().String())
	if err != ErrPaymentNotFound{
		t.Errorf("Reject(): payment musn't be found")
	}
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
	if err == nil {
		t.Errorf("got: nil error")
	}
	if account != nil {
		t.Errorf("want: nil, got accaunt: %v", account.ID)
	}
}