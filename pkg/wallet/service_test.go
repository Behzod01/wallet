package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Behzod01/wallet/pkg/types"
)

func TestService_FindPaymentByID_success(t *testing.T) {
	//создаём сервис
	s := &Service{ /*
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
			},*/
	}

	//регистрируем там пользователя
	phone := types.Phone("+992000000001")
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
	//пробуем найти платёж
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
	//сравниваем платежа
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := &Service{
		payments: []*types.Payment{
			{
				ID:     "1111",
				Amount: 5000_00,
			},
			{
				ID:     "4444",
				Amount: 1500_00,
			},
			{
				ID:     "2222",
				Amount: 50_000,
			},
			{
				ID:     "5555",
				Amount: 100_000,
			},
			{
				ID:     "3333",
				Amount: 160_000,
			},
		},
	}
	/*
		//регистрируем там пользователя
		phone := types.Phone("+992000000001")
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
		_, err := s.Pay(account.ID, 500_00, "phone")
		if err != nil {
			t.Errorf("Reject(): can not creat payment, error = %v", err)
			return
		}

		//пробуем найти несуществующий платёж
		_, err := s.FindPaymentByID(uuid.New().String())
		if err == nil {
			t.Errorf("FindPaymentByID(): must return error, returned nil")
			return
		}
		if err != ErrPaymentNotFound {
			t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned=%v",err)
			return
		}*/
	expected := ErrAccountNotFound

	result, err := s.FindPaymentByID("6666")
	if err != nil {
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
	//попробуем отменить платёж
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
}

func TestService_Reject_fail(t *testing.T) {
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
