package types

// Money - денежная суммы в минимальных еденицах (дирамы, копейки, центы и т.д.)
type Money int64

// PaymentCategory - представляет собой категорию, в которой был совершён платёж (cafe, auto, food, drugs, ...)
type PaymentCategory string

// PaymentStatus - представляет собой статус платежа
type PaymentStatus string

// Предопределённые статусы
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment представляет информацию о платеже
type Payment struct {
	ID        string
	AccountId int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

// Favorite представляет избранный платёж
type Favorite struct {
	ID			string
	AccountID	int64
	Name		string
	Amount		Money
	Category	PaymentCategory
}

type Phone string

type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Messenger interface {
	Send(message string) bool
	Receive() (message string, ok bool)
}

type Telegram struct {
}

func (t *Telegram) Send(message string) bool {
	return true
}

func (t *Telegram) Receive() (message string, ok bool) {
	return "", true
}

type error interface {
	Error() string
}
