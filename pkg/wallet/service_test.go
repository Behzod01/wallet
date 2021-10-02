package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Behzod01/wallet/pkg/types"
	"github.com/google/uuid"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
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
	phone:   "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *Service) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	//регистрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error=%v", err)
	}

	//пополняем его счёт
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error=%v", err)
	}

	//выполняем платежи
	//можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		//тогда здесь работаем просто через index, а не append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error=%v", err)
		}
	}
	return account, payments, nil
}

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
		t.Errorf("want alredy registered, now:%v", err)
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

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "mobile connection")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	payment := uuid.New().String()
	_, err = s.FavoritePayment(payment, "mobile connection")
	if err == ErrFavoriteNotFound {
		t.Errorf("%v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "mobile connection")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	_, err = s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error("PayFromFavorite(): must return error, returned nil")
		return
	}
	fmt.Println(payment.ID)
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	payment := payments[0]

	favorite := uuid.New().String()
	_, err = s.PayFromFavorite(favorite)
	if err == nil {
		t.Error("PayFromFavorite(): must return error, returned nil")
		return
	}
	fmt.Println(payment.ID)
}

func  TestExportToFile(t *testing.T){

	s := newTestService()

	s.RegisterAccount("748362578456")
	s.RegisterAccount("73826283758")
	s.RegisterAccount("7835628753")

err:=s.ExportToFile("../data/export.txt")

if err==nil{

	t.Errorf("cannot export  = %v",err)
	return
}
}

func  TestImportFromFile(t *testing.T){

	s := newTestService()

	s.RegisterAccount("748362578456")
	s.RegisterAccount("73826283758")
	s.RegisterAccount("7835628753")


err := s.ImportFromFile("../data/export.txt")

if err==nil{

	t.Errorf("cannot export  = %v",err)
	return
}
}

func TestExport(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
	  t.Error(err)
	  return
	}
  
	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "new")
	if err != nil {
	  t.Errorf("FavoritePayment(): error = %v", err)
	  return
	}
  
	err = s.Export("../data")
  
	if err == nil {
	  t.Error(err)
	}
}
  
func TestImport(t *testing.T) {
	s := newTestService()
	account, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
	  t.Error(err)
	  return
	}
  
	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "new")
	if err != nil {
	  t.Errorf("FavoritePayment(): error = %v", err)
	  return
	}
  
	_ = s.Export("data")
  
	err = s.Import("data")
  
	if !reflect.DeepEqual(account, s.accounts[0]) {
	  t.Errorf(("ImportF(): wrong account returned = %v"), err)
	  return
	}
}

func TestService_Export_success(t *testing.T) {
  svc := Service{}

  acc, err := svc.RegisterAccount("+992000000001")

  if err != nil {
    t.Errorf("method RegisterAccount returned not nil error, account => %v", acc)
  }

  err = svc.Deposit(acc.ID, 100_00)
  if err != nil {
    t.Errorf("method Deposit returned not nil error, error => %v", err)
  }

  _, err = svc.Pay(acc.ID, 1, "Cafe")
  payN, err := svc.Pay(acc.ID, 2, "Auto")
  _, err = svc.Pay(acc.ID, 3, "MarketShop")
  if err != nil {
    t.Errorf("method Pay returned not nil error, err => %v", err)
  }

  _, err = svc.FavoritePayment(payN.ID, "love")
  if err != nil {
    t.Errorf("method Pay returned not nil error, err => %v", err)
  }
  
  err = svc.Export("../../data")
  if err != nil {
    t.Errorf("method Export returned not nil error, err => %v", err)
  }
}

func BenchmarkSumPayments(b *testing.B) {

	s := Service{}

	account, err := s.RegisterAccount("909796600")
	if err != nil {
		b.Error(err)
	}

	err = s.Deposit(account.ID, 10_000_00)
	err = s.Deposit(account.ID, 20_000_00)
	err = s.Deposit(account.ID, 30_000_00)

	if err != nil {
		b.Errorf("can not deposit account, error = %v", err)
		return
	}
	//осуществляем платёж на его счёт

	_, err = s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 2000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 3000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}
	want := types.Money(6000_00)

	for i := 0; i < b.N; i++ {
		result := s.SumPayments(10)
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}

}


/*
func BenchmarkSumPaymentsWithProgress(b *testing.B) {

	s:=Service{}

	account,err:=s.RegisterAccount("909796600")
	if err!= nil{
		b.Error(err)
	}

	err =s.Deposit(account.ID, 10_000_00)
	err= s.Deposit(account.ID, 20_000_00)
	err= s.Deposit(account.ID, 30_000_00)

	if err != nil {
		b.Errorf("can not deposit account, error = %v", err)
		return
	}
	//осуществляем платёж на его счёт

	_, err = s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 2000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 3000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 2000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 3000_00, "auto")
	if err != nil {
		b.Errorf(" can not creat payment, error = %v", err)
		return
	}

	want := types.Money(11000_00)

	for i := 0; i < b.N; i++ {
		result:=s.SumPaymentsWithProgress()

		prog:= <-result
		if prog.Result != want {
			b.Fatalf("invalid result, got %v", prog.Result)
		}
	}

}

func TestSumPaymentsWithProgress(t *testing.T) {

	s:=Service{}

	account,err:=s.RegisterAccount("909796600")
	if err!= nil{
		t.Error(err)
	}

	err =s.Deposit(account.ID, 10_000_00)
	err= s.Deposit(account.ID, 20_000_00)
	err= s.Deposit(account.ID, 30_000_00)

	if err != nil {
		t.Errorf("can not deposit account, error = %v", err)
		return
	}
	//осуществляем платёж на его счёт

	_, err = s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 2000_00, "auto")
	if err != nil {
		t.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 3000_00, "auto")
	if err != nil {
		t.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 2000_00, "auto")
	if err != nil {
		t.Errorf(" can not creat payment, error = %v", err)
		return
	}
	_, err = s.Pay(account.ID, 3000_00, "auto")
	if err != nil {
		t.Errorf(" can not creat payment, error = %v", err)
		return
	}

	want := types.Money(11000_00)


		result:=s.SumPaymentsWithProgress()

		prog:= <-result
		if prog.Result != want {
			t.Fatalf("invalid result, got %v", prog.Result)
		}
	}
*/
