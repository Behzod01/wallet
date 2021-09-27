package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Behzod01/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	service := &Service{
		accounts: []*types.Account{
			{
				ID:      1111,
				Phone:   "90-000-00-01",
				Balance: 10_000,
			},
			{
				ID:      4444,
				Phone:   "90-444-00-01",
				Balance: 1500_00,
			},
			{
				ID:      2222,
				Phone:   "93-000-33-33",
				Balance: 50_000,
			},
			{
				ID:      5555,
				Phone:   "77-700-07-77",
				Balance: 100_000,
			},
			{
				ID:      3333,
				Phone:   "91-111-11-01",
				Balance: 160_000,
			},
		},
	}
	expected := types.Account{
		ID:      5555,
		Phone:   "77-700-07-77",
		Balance: 100_000,
	}

	result, err := service.FindAccountByID(5555)

	if err == nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	service := &Service{
		accounts: []*types.Account{
			{
				ID:      1111,
				Phone:   "90-000-00-01",
				Balance: 10_000,
			},
			{
				ID:      4444,
				Phone:   "90-444-00-01",
				Balance: 1500_00,
			},
			{
				ID:      2222,
				Phone:   "93-000-33-33",
				Balance: 50_000,
			},
			{
				ID:      5555,
				Phone:   "77-700-07-77",
				Balance: 100_000,
			},
			{
				ID:      3333,
				Phone:   "91-111-11-01",
				Balance: 160_000,
			},
		},
	}
	expected := ErrAccountNotFound

	result, err := service.FindAccountByID(55555)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	service := &Service{
		payments: []*types.Payment{
			{
				ID: "1111",
				Amount: 5000_00,				
			},
			{
				ID:      "4444",
				Amount: 1500_00,
			},
			{
				ID:      "2222",
				Amount: 50_000,
			},
			{
				ID:      "5555",
				Amount: 100_000,
			},
			{
				ID:      "3333",
				Amount: 160_000,
			},
		},
	}
	expected := types.Payment{
		ID:      "3333",
		Amount: 160_000,
	}

	result, err := service.FindPaymentByID("3333")

	if err == nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_FindPaymentByID_notfound(t *testing.T) {
	service := &Service{
		payments: []*types.Payment{
			{
				ID: "1111",
				Amount: 5000_00,				
			},
			{
				ID:      "4444",
				Amount: 1500_00,
			},
			{
				ID:      "2222",
				Amount: 50_000,
			},
			{
				ID:      "5555",
				Amount: 100_000,
			},
			{
				ID:      "3333",
				Amount: 160_000,
			},
		},
	}
	expected := ErrAccountNotFound

	result, err := service.FindPaymentByID("1111")

	if err == nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_Reject_success(t *testing.T) {
	//создаём сервис
	s := Service{}

	//регистрируем пользователья
	phone := types.Phone("+992777000777")
	account, err := s.RegisterAccount(phone)

	if err != nil {
		t.Errorf("Reject(): can not register account, error = %v", err)
		return
	}
	//пополняем его счёт
	err = s.Deposit(account.ID, 1000_00)

	if err != nil {
		t.Errorf("Reject(): can not deposit account, error = %v", err)
		return
	}
	//осуществляем платёж на его счёт
	payment, err := s.Pay(account.ID, 500_00, "phone")
	if err != nil {
		t.Errorf("Reject(): can not creat payment, error = %v", err)
		return
	}
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
}

func TestService_Reject_notfound(t *testing.T) {
	//создаём сервис
	s := Service{}

	//регистрируем пользователья
	phone := types.Phone("+992777000777")
	account, err := s.RegisterAccount(phone)

	if err != nil {
		t.Errorf("Reject(): can not register account, error = %v", err)
		return
	}
	//пополняем его счёт
	err = s.Deposit(account.ID, 1000_00)

	if err != nil {
		t.Errorf("Reject(): can not deposit account, error = %v", err)
		return
	}
	//осуществляем платёж на его счёт
	payment, err := s.Pay(account.ID, 500_00, "phone")
	if err != nil {
		t.Errorf("Reject(): can not creat payment, error = %v", err)
		return
	}
	err = s.Reject(payment.ID)
	if err == ErrPaymentNotFound {
		t.Errorf("Reject(): error = %v", err)
		return
	}
}
