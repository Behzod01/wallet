package wallet

import (
	"reflect"
	"testing"

	"github.com/Behzod01/wallet/pkg/types"
	"github.com/google/uuid"
)

func TestService_FindPaymentByID_success(t *testing.T) {
	//создаём сервис
	s := newTestService()

	//регистрируем там пользователя
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//пробуем найти платёж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
	//сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned=%v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	//создаём сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//пробуем найти несуществующий платёж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID(): must return error, returned nil")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned=%v", err)
		return
	}
}
func TestService_Reject_success(t *testing.T) {
	//создаём сервис
	s := newTestService()

	//регистрируем пользователья
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//пробуем отменить платёж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}
/*
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
}*/
