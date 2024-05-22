// mapper/mapper.go
package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"ocr-poc/model"
	"ocr-poc/ocr"
	"ocr-poc/service"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func MapToInvoice(ocrResp *ocr.OCRResponse, services *service.ServiceContainer) (*model.IncomingInvoice, error) {
	var invoice model.IncomingInvoice

	currency := ocrResp.GetCurrency()
	
	invoice.Supplier = model.RelationReference{}

	for _, entity := range ocrResp.Entities {
		switch entity.Type {
		case "supplier_registration":
			id, _ := strconv.ParseInt(entity.MentionText, 10, 64)
			supplier, _ := services.SupplierService.GetByID(id)
			invoice.Supplier = model.RelationReference(*supplier)
		case "invoice_id":
			invoice.ExternalNumber = entity.MentionText
		case "invoice_date":
			date, err := parseDate(entity.MentionText)
			if err != nil {
				return nil, fmt.Errorf("failed to parse date: %v", err)
			}
			invoice.Date = date
		case "due_date":
			date, err := parseDate(entity.MentionText)
			if err != nil {
				return nil, fmt.Errorf("failed to parse due date: %v", err)
			}
			invoice.DueDate = &date
		case "description":
			invoice.Description = entity.MentionText

		case "total_amount":
			amount, err := parseAmount(entity.MentionText)
			if err != nil {
				return &model.IncomingInvoice{}, fmt.Errorf("failed to parse total amount: %v", err)
			}
			invoice.Amount = model.Amount{Currency: currency, Value: amount}
		case "total_tax_amount":
			vatAmount, err := parseAmount(entity.MentionText)
			if err != nil {
				return &model.IncomingInvoice{}, fmt.Errorf("failed to parse VAT amount: %v", err)
			}
			invoice.VAT = model.Amount{Currency: currency, Value: vatAmount}
		case "accept_double_external_numbers":
			invoice.AcceptDoubleExternalNumbers = parseBool(entity.MentionText)
		case "payment":
			invoice.Payment = parsePayment(entity, currency)
		// case "line_item":
		// 	bookingItem := parseBookingItem(entity, currency)
		// 	invoice.BookingItems = append(invoice.BookingItems, bookingItem)
		case "booking_period":
			invoice.BookingPeriod = parseBookingPeriod(entity)
		}
	}

	return &invoice, nil
}

// Non-English to English month names mapping
var monthNames = map[string]string{
	"januari":   "January",
	"februari":  "February",
	"maart":     "March",
	"april":     "April",
	"mei":       "May",
	"juni":      "June",
	"juli":      "July",
	"augustus":  "August",
	"september": "September",
	"oktober":   "October",
	"november":  "November",
	"december":  "December",
}

// translateMonthNames replaces non-English month names with English ones
func translateMonthNames(dateStr string) string {
	for nonEng, eng := range monthNames {
		dateStr = strings.ReplaceAll(dateStr, nonEng, eng)
	}
	return dateStr
}

// parseDate tries to parse a date from various possible formats
func parseDate(dateStr string) (time.Time, error) {

	dateStr = translateMonthNames(dateStr)

	// Try to parse the date using the dateparse library
	parsedDate, err := dateparse.ParseAny(dateStr)
	if err == nil {
		return parsedDate, nil
	}

	dateFormats := []string{
		"02 Jan 2006", "02/01/2006", "January 2, 2006", "02-01-2006", "2006-01-02",
		"02-01-2006", "02/01/2006", "02.01.2006", "2006.01.02", "02 January 2006",
	}

	// Try to parse the date using custom formats
	for _, format := range dateFormats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return parsedDate, fmt.Errorf("failed to parse date: %s", dateStr)
}

func parseAmount(amountStr string) (*float64, error) {
	// Remove spaces for easier processing
	amountStr = strings.ReplaceAll(amountStr, " ", "")
	var normalizedStr string
	if strings.Count(amountStr, ",") == 1 && strings.Count(amountStr, ".") > 0 {
		normalizedStr = strings.Replace(amountStr, ".", "", -1)
		normalizedStr = strings.Replace(normalizedStr, ",", ".", 1)
	} else {
		normalizedStr = strings.Replace(amountStr, ",", ".", 1)
	}

	value, err := strconv.ParseFloat(normalizedStr, 64)
	if err != nil {
		return nil, errors.New("failed to parse amount")
	}

	return &value, nil
}

func parseFloat(valueStr string) *float64 {
	var value float64
	var err error
	normalizedStr := strings.ReplaceAll(valueStr, ",", "")
	value, err = strconv.ParseFloat(normalizedStr, 64)
	if err != nil {
		return &value
	}
	return &value
}

func parseBool(valueStr string) bool {
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false
	}
	return value
}

func parsePayment(entity *ocr.Entity, currency string) model.Payment {
	var payment model.Payment
	for _, prop := range entity.Properties {
		switch prop.Type {
		case "reference":
			payment.Reference = prop.MentionText
		case "description":
			payment.Description = prop.MentionText
		case "create_payment_order":
			payment.CreatePaymentOrder = parseBool(prop.MentionText)
		case "payment_order":
			payment.PaymentOrder = parsePaymentOrder(prop, currency)
		}
	}
	return payment
}

func parsePaymentOrder(entity *ocr.Entity, currency string) *model.PaymentOrder {
	var paymentOrder model.PaymentOrder
	for _, prop := range entity.Properties {
		switch prop.Type {
		case "counter_account":
			paymentOrder.CounterAccount = parseBankAccountRef(prop)
		case "bank_account":
			paymentOrder.BankAccount = parseBankAccountRef(prop)
		case "amount":
			paymentAmount, err := parseAmount(prop.MentionText)
			if err != nil {
				fmt.Println("failed to parse payment order amount: %v", err)
				return &model.PaymentOrder{}
			}
			paymentOrder.Amount = model.Amount{
				Currency: currency,
				Value:    paymentAmount,
			}
		}
	}
	return &paymentOrder
}

func parseBankAccountRef(entity *ocr.Entity) model.BankAccountRef {
	var bankAccountRef model.BankAccountRef
	for _, prop := range entity.Properties {
		switch prop.Type {
		case "id":
			bankAccountRef.ID = parseInt(prop.MentionText)
		case "name":
			bankAccountRef.Name = prop.MentionText
		case "account":
			bankAccountRef.Account = prop.MentionText
		case "iban":
			bankAccountRef.IBAN = &prop.MentionText
		}
	}
	return bankAccountRef
}

func parseBookingItem(entity *ocr.Entity, currency string) model.BookingItem {
	var bookingItem model.BookingItem
	for _, prop := range entity.Properties {
		switch prop.Type {
		case "line_item/quantity":
			quantity, _ := strconv.ParseFloat(strings.ReplaceAll(prop.MentionText, ",", ""), 64)
			bookingItem.Amount = model.AmountType{
				Currency: currency,
				Value:    quantity,
			}
		case "line_item/unit":
			// Handle unit if needed
		case "line_item/unit_price":
			// Handle unit price if needed
		case "line_item/amount":

			bookingItem.Amount = model.AmountType{
				Currency: currency,
				Value:    *parseFloat(prop.MentionText),
			}
		case "line_item/description":
			// Assuming description is a separate field and not part of bookingItem
			bookingItem.CostCenter = &model.CostCenter{Description: prop.MentionText}
		}
	}
	return bookingItem
}

func parseBookingPeriod(entity *ocr.Entity) *model.BookingPeriodTemp {
	var bookingPeriod model.BookingPeriodTemp
	for _, prop := range entity.Properties {
		switch prop.Type {
		case "first_period":
			bookingPeriod.FirstPeriod = parseBookingPeriodDetails(prop)
		case "last_period":
			bookingPeriod.LastPeriod = parseBookingPeriodDetails(prop)
		}
	}
	return &bookingPeriod
}

func parseBookingPeriodDetails(entity *ocr.Entity) model.BookingPeriod {
	var bookingPeriod model.BookingPeriod
	for _, prop := range entity.Properties {
		switch prop.Type {
		case "id":
			bookingPeriod.ID = parseInt(prop.MentionText)
		case "accounting_period_id":
			bookingPeriod.AccountingPeriodID = parseInt(prop.MentionText)
		case "start_date":
			bookingPeriod.StartDate, _ = parseDate(prop.MentionText)
		case "end_date":
			bookingPeriod.EndDate, _ = parseDate(prop.MentionText)
		case "locked_timestamp":
			timestamp, _ := parseDate(prop.MentionText)
			bookingPeriod.LockedTimestamp = &timestamp
		}
	}
	return bookingPeriod
}

func parseInt(valueStr string) int {
	value, err := strconv.Atoi(strings.ReplaceAll(valueStr, ",", ""))
	if err != nil {
		return 0
	}
	return value
}

// PrettyPrintInvoice converts the IncomingInvoice to a beautified JSON string
func PrettyPrintInvoice(invoice *model.IncomingInvoice) string {
	jsonBytes, err := json.MarshalIndent(invoice, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal invoice: %v", err)
		return ""
	}
	return string(jsonBytes)
}
