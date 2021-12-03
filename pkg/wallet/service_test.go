package wallet

import (
	"testing"
	"github.com/MSHE97/wallet/pkg/types"
)

func TestFindAccountById_exists(t *testing.T) {
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

func TestFindAccountById_notExists(t *testing.T ) {
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

func TestReject_existsPay(t *testing.T){
	svc := &Service{}
	acc, err0 := svc.RegisterAccount("000000001")
	if err0 != nil {
		switch err0 {
		case ErrPhoneRegistered:
			t.Errorf("Can't register new account")			
		}	
	}

	err1 := svc.Deposit(acc.ID, 2500)
	if err1 != nil {
		switch err1 {
		case ErrAmountMustBePositive:
			t.Errorf("Deposite must be positive")
		case ErrAccountNotFound:
			t.Errorf("Account not found")
		}
		return
	}

	payment, err2 := svc.Pay(acc.ID, 2500, "cafe")
	if err2 != nil {
		switch err2 {
		case ErrAccountNotFound:
			t.Errorf("Account not found")
		case ErrAmountMustBePositive:
			t.Errorf("Amount must be greater then zero")
		case ErrLowBalance:
			t.Errorf("Low balance")
		}
	}

	err := svc.Reject(payment.ID)
	if err != nil{
		switch err {
		case ErrPaymentNotFound:
			t.Errorf("Payment not found")
		}
	}

	if payment.Status != types.PaymentStatusFail{
		t.Errorf("Reject wasn't happen")
	}
}

func TestReject_notExists(t *testing.T){
	svc := &Service{}
	err := svc.Reject("Not existing payment ID")
	if err != ErrPaymentNotFound{
		t.Errorf("Payment found")
	}
}

func TestFindPaymentById_exists(t *testing.T) {
	svc := &Service{}
	acc, err0 := svc.RegisterAccount("000000001")
	if err0 != nil {
		switch err0 {
		case ErrPhoneRegistered:
			t.Errorf("Can't register new account")			
		}	
	}

	err1 := svc.Deposit(acc.ID, 2500)
	if err1 != nil {
		switch err1 {
		case ErrAmountMustBePositive:
			t.Errorf("Deposite must be positive")
		case ErrAccountNotFound:
			t.Errorf("Account not found")
		}
	}

	payment, err2 := svc.Pay(acc.ID, 2500, "cafe")
	if err2 != nil {
		switch err2 {
		case ErrAccountNotFound:
			t.Errorf("Account not found")
		case ErrAmountMustBePositive:
			t.Errorf("Amount must be greater then zero")
		case ErrLowBalance:
			t.Errorf("Low balance")
		}
	}

	foundPay, err := svc.FindPaymentById(payment.ID)
	if err != nil {
		switch err {
		case ErrPaymentNotFound:
			t.Errorf("Payment not found")	
		}			
	}
	if foundPay != payment {
		t.Errorf("Found incorrect payment")
	}
}

func TestFindPaymentById_notExists(t *testing.T) {
	svc := &Service{}
	_, err := svc.FindPaymentById("Not existing payment ID")
	if err != ErrPaymentNotFound {
		t.Errorf("Payment found")
	}
}