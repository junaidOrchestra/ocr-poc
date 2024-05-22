// ocr/google_document_ai.go
package ocr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	documentai "cloud.google.com/go/documentai/apiv1"
	"google.golang.org/api/option"

	"cloud.google.com/go/documentai/apiv1/documentaipb"
)

type GoogleDocumentAIClient struct {
	client *documentai.DocumentProcessorClient
}

func NewGoogleDocumentAIClient(ctx context.Context, credentialsFile, endpoint string) (*GoogleDocumentAIClient, error) {
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithCredentialsFile(credentialsFile), option.WithEndpoint(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to create Document AI client: %v", err)
	}
	return &GoogleDocumentAIClient{client: client}, nil
}

func (g *GoogleDocumentAIClient) ProcessDocument(ctx context.Context, documentPath string) (*OCRResponse, error) {
	documentBytes, err := os.ReadFile(documentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read document file: %v", err)
	}

	req := &documentaipb.ProcessRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", "orchestra-ocr-poc", "eu", "513169612c58c534"),
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  documentBytes,
				MimeType: getMimeType(documentPath),
			},
		},
	}

	resp, err := g.client.ProcessDocument(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %v", err)
	}

	// Extract the text
	documentText := resp.GetDocument().GetText()

	// Define a regex pattern to capture payment terms
	paymentTermsPattern := regexp.MustCompile(`(?i)\b(remit to|make checks payable to|Betaling gaarne binnen|payment terms:|due within|send payment to|payment instructions)\b.*`)

	paymentTerms := ""
	// Search for payment terms in the extracted text
	match := paymentTermsPattern.FindString(documentText)
	if match != "" {
		paymentTerms = match
	}

	ocrResponse := &OCRResponse{}
	for _, entity := range resp.GetDocument().Entities {
		ocrResponse.Entities = append(ocrResponse.Entities, convertEntity(entity))
	}
	ocrResponse.Entities = append(ocrResponse.Entities, &Entity{Type: "payment_description", MentionText: paymentTerms})

	return ocrResponse, nil
}

func convertEntity(entity *documentaipb.Document_Entity) *Entity {
	e := &Entity{
		Type:        entity.GetType(),
		MentionText: entity.GetMentionText(),
	}
	fmt.Printf("%v \t %v", e.Type, e.MentionText)
	fmt.Println()

	for _, prop := range entity.Properties {
		e.Properties = append(e.Properties, convertEntity(prop))
	}
	return e
}

// Determine MIME type based on file extension
func getMimeType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".pdf":
		return "application/pdf"
	default:
		return ""
	}
}
