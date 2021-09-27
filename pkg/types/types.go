package types

//Money  представляет собой денежную сумму в минимальных единицах(центы, копейки, и т.д)
type Money int64

//Category  представляет собой категорию в которой был совершен платёж
type PaymetCategory string

//Category  представляет собой статус платежа
type PaymetStatus string

//Предопределлённые статусы платежей
const (
	PaymetStatusOk         PaymetStatus = "OK"
	PaymetStatusFail       PaymetStatus = "FAIL"
	PaymetStatusInProgress PaymetStatus = "INPROGRESS"
)

//Payment представляет  информацию о платеже
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymetCategory
	Status    PaymetStatus
}

type Phone string

//Account представляет информацию о счёте пользователя
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
