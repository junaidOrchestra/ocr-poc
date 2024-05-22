package ocr

import (
	"context"
)

type OCRResponse struct {
	Entities []*Entity
}

type Entity struct {
	Type        string
	MentionText string
	Properties  []*Entity
}

func (o *OCRResponse) GetCurrency() string {
	for _, entity := range o.Entities {
		if entity.Type == "currency" {
			return entity.MentionText
		}
	}
	return "EUR-D" // Default value if currency is not found. -D is intentional to detect where currency is not detected.
}

func (o *OCRResponse) GetPaymentDescription() string {
	for _, entity := range o.Entities {
		if entity.Type == "payment_description" {
			return entity.MentionText
		}
	}
	return "EUR" // Default value if currency is not detected
}

type OCRClient interface {
	ProcessDocument(ctx context.Context, documentPath string) error
}
