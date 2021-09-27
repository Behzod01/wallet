package wallet

import (
	"errors"

	"github.com/Behzod01/wallet/pkg/types"
	"github.com/google/uuid"
)

type errorString struct {
	s string
}

func New(text string) error {
	return &errorString{text}
}

func (e *errorString) Error() string {
	return e.s
}

var ErrAccountNotFound = errors.New("account not found")
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("not enough balance")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments       []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
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

	//zachislenie sredstv poka ne rasmatrivaem kak platezh
	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymetCategory)(*types.Payment, error) {
	if amount <= 0 {
		return nil,  ErrAmountMustBePositive
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
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymetStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
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
	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment

	for _, p := range s.payments {
		if p.ID == paymentID {
			payment = p
			break
		}
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)

	if payment.ID == paymentID {
		payment.Status = types.PaymetStatusFail
	} else if payment.ID == err.Error() {
		return ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payment.AccountID)

	if payment.AccountID == account.ID {
		account.Balance += payment.Amount 		
	}
	return nil
	
}