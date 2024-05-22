package model

import "time"

type Supplier struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type IncomingInvoice struct {
	//Document                     DocumentRef `json:"document" validate:"required"`
	Supplier                    RelationReference  `json:"supplier" validate:"required"`
	ExternalNumber              string             `json:"externalNumber"`
	Date                        time.Time          `json:"invoiceDate" validate:"required"`
	DueDate                     *time.Time         `json:"dueDate"`
	Description                 string             `json:"description" validate:"omitempty"`
	Amount                      Amount             `json:"baseAmount" validate:"required"`
	VAT                         Amount             `json:"vatAmount" validate:"required"`
	AcceptDoubleExternalNumbers bool               `json:"acceptDuplicateExternalNumber"`
	Payment                     Payment            `json:"payment" validate:"required"`
	BookingItems                []BookingItem      `json:"bookingItems"`
	BookingPeriod               *BookingPeriodTemp `json:"bookingPeriod"`
	// ExclusiveAuthorizationGroups []models.Group    `json:"exclusiveAuthorizationGroups" validate:"omitempty"`
}

type RelationReference struct {
	ID   int    `json:"id"  validate:"required"`
	Name string `json:"name"`
}

type Amount struct {
	Currency string   `gorm:"column:Currency" json:",omitempty"`
	Value    *float64 `gorm:"column:Value" json:",omitempty"`
}

type BookingItem struct {
	CostCenter *CostCenter `json:"costCenter"`
	Account    AccountRef  `json:"account" validate:"required"`
	Amount     AmountType  `json:"amount" validate:"required"`
}

type BookingPeriodTemp struct {
	FirstPeriod BookingPeriod `json:"firstPeriod" validate:"required"`
	LastPeriod  BookingPeriod `json:"lastPeriod" validate:"required"`
}

// CostCenter defines the cost centers.
type CostCenter struct {
	ID          int    `json:"id" gorm:"column:ID;primary_key" validate:"required"`
	Description string `json:"description" gorm:"column:Description"`
	Status      int    `json:"status" gorm:"column:Status"`
}

type AccountRef struct {
	ID          int    `json:"id" gorm:"column:ID;primary_key" validate:"required"`
	Code        string `json:"code" gorm:"column:Code"`
	Description string `json:"description" gorm:"column:Description"`
	Status      int    `json:"-" gorm:"column:Status"`
}

type AmountType struct {
	Currency string  `json:"currency" gorm:"column:Currency" validate:"required"`
	Value    float64 `json:"value" gorm:"column:Value"`
}

type BookingPeriod struct {
	ID                 int        `json:"id" gorm:"column:ID;primary_key" validate:"required"`
	AccountingPeriodID int        `json:"accountingPeriodId" gorm:"column:AccountingPeriodID"`
	StartDate          time.Time  `json:"startDate" gorm:"column:StartDate"`
	EndDate            time.Time  `json:"endDate" gorm:"column:EndDate"`
	LockedTimestamp    *time.Time `json:"lockedTimestamp" gorm:"column:LockedTimestamp"`
}

type Payment struct {
	Reference          string        `json:"reference"`
	Description        string        `json:"description" validate:"required"`
	CreatePaymentOrder bool          `json:"createPaymentOrder"`
	PaymentOrder       *PaymentOrder `json:"paymentOrder" validate:"required_with=createPaymentOrder"`
}

type PaymentOrder struct {
	CounterAccount BankAccountRef `json:"counterAccount"`
	BankAccount    BankAccountRef `json:"bankAccount"  validate:"required"`
	Amount         Amount         `json:"amount"  validate:"required"`
}

type BankAccountRef struct {
	ID      int     `json:"id" gorm:"column:ID;primary_key"`
	Name    string  `json:"name" gorm:"column:Name"`
	Account string  `json:"account" gorm:"column:AccountNumber"`
	IBAN    *string `json:"iban" gorm:"column:IBAN"`
	Status  int     `json:"status" gorm:"column:Status"`
}
