package wallet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
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
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}
	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("can't find payment, error=%v", err)
	}
	payment, err := s.Pay(pay.AccountID, pay.Amount, pay.Category)
	if err != nil {
		return nil, fmt.Errorf("can't create payment again, error=%v", err)
	}

	return payment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	findpay, err := s.FindFavoriteByID(favoriteID)

	if err != nil {
		return nil, ErrFavoriteNotFound
	}
	pay, err := s.Pay(findpay.AccountID, findpay.Amount, findpay.Category)

	if err != nil {
		return nil, ErrPaymentNotFound
	}

	return pay, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {

	var favorite *types.Favorite

	for _, fav := range s.favorites {

		if fav.ID == favoriteID {
			favorite = fav
			break

		}
	}
	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}
	return favorite, nil
}

func (s *Service) ExportToFile(path string) error {
	str := ""

	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	for _, account := range s.accounts {
		str += strconv.Itoa(int(account.ID)) + ";"
		str += string(account.Phone) + ";"
		str += strconv.Itoa(int(account.Balance)) + "|"
	}
	_, err = file.Write([]byte(str))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}

	content := make([]byte, 0)
	buf := make([]byte, 4)

	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := string(content)

	accounts := strings.Split(data, "|")
	accounts = accounts[:len(accounts)-1]

	for _, account := range accounts {

		splits := strings.Split(account, ";")

		id, err := strconv.Atoi(splits[0])
		if err != nil {
			log.Print(err)
			return err
		}

		phone := splits[1]

		balance, err := strconv.Atoi(splits[2])
		if err != nil {
			log.Print(err)
			return err
		}

		s.accounts = append(s.accounts, &types.Account{
			ID:      int64(id),
			Phone:   types.Phone(phone),
			Balance: types.Money(balance),
		})
	}
	return nil
}

func (s *Service) Export(dir string) error {

	if len(s.accounts) > 0 {
		file, err := os.OpenFile(dir+"/accounts.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				if err == nil {
					cerr = err
				}
			}
		}()
		accstr := ""
		for _, account := range s.accounts {
			accstr += strconv.Itoa(int(account.ID)) + ";"
			accstr += string(account.Phone) + ";"
			accstr += strconv.Itoa(int(account.Balance)) + "\n"
		}
		file.WriteString(accstr)
	}

	if len(s.payments) > 0 {
		fil, err := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if cerr := fil.Close(); cerr != nil {
				if err == nil {
					cerr = err
				}
			}
		}()

		paystr := ""
		for _, payment := range s.payments {
			paystr += string(payment.ID) + ";"
			paystr += strconv.Itoa(int(payment.AccountID)) + ";"
			paystr += strconv.Itoa(int(payment.Amount)) + ";"
			paystr += string(payment.Category) + ";"
			paystr += string(payment.Status) + "\n"
		}
		fil.WriteString(paystr)
	}

	if len(s.favorites) > 0 {
		files, err := os.OpenFile(dir+"/favorites.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if cerr := files.Close(); cerr != nil {
				if err == nil {
					cerr = err
				}
			}
		}()

		favstr := ""
		for _, favorite := range s.favorites {
			favstr += favorite.ID + ";"
			favstr += strconv.Itoa(int(favorite.AccountID)) + ";"
			favstr += favorite.Name + ";"
			favstr += strconv.Itoa(int(favorite.Amount)) + ";"
			favstr += string(favorite.Category) + "\n"
		}
		files.WriteString(favstr)
	}
	return nil
}

func (s *Service) Import(dir string) error {

	_, err := os.Stat(dir + "/accounts.dump")
	if err == nil {
		file, err := os.ReadFile(dir + "/accounts.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		accstr := string(file)
		accounts := strings.Split(accstr, "\n")

		if len(accounts) > 0 {
			accounts = accounts[:len(accounts)-1]
		}
		for _, account := range accounts {
			splits := strings.Split(account, ";")
			id, err := strconv.Atoi(splits[0])
			if err != nil {
				log.Print(err)
				return err
			}
			phone := splits[1]
			balance, err := strconv.Atoi(splits[2])
			if err != nil {
				log.Print(err)
				return err
			}
			s.accounts = append(s.accounts, &types.Account{
				ID:      int64(id),
				Phone:   types.Phone(phone),
				Balance: types.Money(balance),
			})
		}
	}
	// Payments==========================================

	_, err = os.Stat(dir + "/payments.dump")

	if err == nil {
		file, err := os.ReadFile(dir + "/payments.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		paystr := string(file)
		payments := strings.Split(paystr, "\n")

		if len(payments) > 0 {
			payments = payments[:len(payments)-1]
		}
		for _, payment := range payments {
			splits := strings.Split(payment, ";")
			id := splits[0]
			accountid, err := strconv.Atoi(splits[1])
			if err != nil {
				log.Print(err)
				return err
			}
			amount, err := strconv.Atoi(splits[2])
			if err != nil {
				log.Print(err)
				return err
			}
			category := splits[3]
			status := splits[4]
			s.payments = append(s.payments, &types.Payment{
				ID:        id,
				AccountID: int64(accountid),
				Amount:    types.Money(amount),
				Category:  types.PaymentCategory(category),
				Status:    types.PaymentStatus(status),
			})

		}

	}
	//Favorites =======================================================
	_, err = os.Stat(dir + "/favorites.dump")

	if err == nil {
		file, err := os.ReadFile(dir + "/favorites.dump")
		if err != nil {
			log.Print(err)
			return err
		}
		favstr := string(file)
		favorites := strings.Split(favstr, "\n")
		if len(favorites) > 0 {
			favorites = favorites[:len(favorites)-1]
		}
		for _, favorite := range favorites {
			splits := strings.Split(favorite, ";")
			id := splits[0]
			accountid, err := strconv.Atoi(splits[1])
			if err != nil {
				log.Print(err)
				return err
			}
			name := splits[2]
			amount, err := strconv.Atoi(splits[3])
			if err != nil {
				log.Print(err)
				return err
			}
			category := types.PaymentCategory(splits[4])
			s.favorites = append(s.favorites, &types.Favorite{
				ID:        id,
				AccountID: int64(accountid),
				Name:      name,
				Amount:    types.Money(amount),
				Category:  types.PaymentCategory(category),
			})
		}
	}
	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	payments := make([]types.Payment, 0)
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			payments = append(payments, *payment)
		}
	}
	return payments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	var file *os.File
	var err error
	if len(payments) == 0 {
		return nil
	}
	if len(payments) <= records {
		file, err = os.Create(dir + "/payments.dump")
		if err != nil {
			return err
		}
	} else {
		file, err = os.Create(dir + "/payments1.dump")
		if err != nil {
			return err
		}
	}
	x := 1
	i := 1
	for _, payment := range payments {
		log.Println(strconv.Itoa(i) + " " + strconv.Itoa(x) + " " + strconv.Itoa(records))
		if i%records == 1 && i != 1 {
			x++
			file, err = os.Create(dir + "/payments" + strconv.Itoa(x) + ".dump")
			if err != nil {
				return err
			}
		}
		_, err := file.Write([]byte(payment.ID + ";"))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(strconv.FormatInt(payment.AccountID, 10) + ";"))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(strconv.FormatInt(int64(payment.Amount), 10) + ";"))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(payment.Category + ";"))
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = file.Write([]byte(payment.Status + "\n"))
		if err != nil {
			log.Print(err)
			return err
		}
		i++
	}
	return nil
}
