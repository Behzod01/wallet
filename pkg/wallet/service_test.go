package wallet

import (
	"reflect"
	"testing"

	"github.com/Behzod01/wallet/pkg/types"
	"github.com/google/uuid"
)

func TestService_RegisterAccount_success(t *testing.T) {
	s := newTestService()
	account, err := s.RegisterAccount(defaultTestAccount.phone)
	if err != nil {
		t.Error(err)
	}
	phone := types.Phone("+992000000001")
	_, err = s.RegisterAccount(phone)
	expected := ErrPhoneRegistered
	if !reflect.DeepEqual(expected, err) {
		t.Errorf("want alredy registered, now:%v",err)
		return
	}
	err = s.Deposit(account.ID, -500)
	if err == nil {
		t.Error(err)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	//создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//пробуем найти платёж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error=%v", err)
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

	//сравниваем платежи
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

func TestService_Reject_notfound(t *testing.T) {
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
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	got, err := s.FindPaymentByID(payments[0].ID)
	if err != nil {
		t.Errorf("Repeat(): error=%v", err)
		return
	}
	//пробуем повторить платёж
	payment, err := s.Repeat(payments[0].ID)
	if err != nil {
		t.Errorf("Repeat(): can't repeat error=%v", err)
		return
	}

	if reflect.DeepEqual(got, payment) {
		t.Errorf("wrong repeat of payment, error=%v", err)
	}
}
