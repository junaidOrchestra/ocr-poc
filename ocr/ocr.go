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

// func (e *Entity) String() string {
// 	var sb strings.Builder
// 	sb.WriteString(fmt.Sprintf("Type: %s, MentionText: %s", e.Type, e.MentionText))
// 	if len(e.Properties) > 0 {
// 		sb.WriteString(", Properties: [")
// 		for _, prop := range e.Properties {
// 			sb.WriteString(prop.String())
// 			sb.WriteString(", ")
// 		}
// 		sb.WriteString("]")
// 	}
// 	return sb.String()
// }

// func (o *OCRResponse) String() string {
// 	var sb strings.Builder
// 	sb.WriteString("Entities: [")
// 	for _, entity := range o.Entities {
// 		sb.WriteString(entity.String())
// 		sb.WriteString(", ")
// 	}
// 	sb.WriteString("]")
// 	return sb.String()
// }

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
	return "EUR" // Default value if currency is not found
}

type OCRClient interface {
	ProcessDocument(ctx context.Context, documentPath string) error
}
