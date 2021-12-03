package wallet

import (
	"testing"
	//"github.com/MSHE97/wallet/pkg/types"
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
	if err == nil {
		t.Errorf("got: nil error")
	}
	if account != nil {
		t.Errorf("want: nil, got accaunt: %v", account.ID)
	}
}