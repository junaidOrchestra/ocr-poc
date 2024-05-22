// Sample vision-quickstart uses the Google Cloud Vision API to label an image.
package main

import (
	"context"
	"fmt"
	"log"
	mapper "ocr-poc/mapping"
	"ocr-poc/ocr"
	"ocr-poc/repository"
	"ocr-poc/service"
)

func main() {
	ctx := context.Background()

	supplierRepo := repository.NewInMemorySupplierRepository()
	supplierService := service.NewSupplierService(supplierRepo)
	services := service.NewServiceContainer(*supplierService)

	// Initialize the Google Document AI client
	googleClient, err := ocr.NewGoogleDocumentAIClient(ctx, "orchestra-ocr-poc-creds.json", "eu-documentai.googleapis.com:443")
	if err != nil {
		log.Fatalf("Failed to create Document AI client: %v", err)
	}

	// Process the document
	// documentPath := "invoice-2.png"
	documentPath := "invoice-3.pdf"
	ocrResp, err := googleClient.ProcessDocument(ctx, documentPath)
	if err != nil {
		log.Fatalf("Error processing document: %v", err)
	}

	// fmt.Println("response: ", resp.String())

	// Map the OCR response to an invoice
	invoice, err := mapper.MapToInvoice(ocrResp, services)
	if err != nil {
		log.Fatalf("Error mapping OCR response to invoice: %v", err)
	}

	fmt.Printf("Mapped Invoice: %+v\n", mapper.PrettyPrintInvoice(invoice))

}

// func main() {
// 	// Load the service account key file path from an environment variable
// 	//	serviceAccountKeyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

// 	// Create a new context
// 	ctx := context.Background()

// 	// Authenticate using the service account key
// 	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithCredentialsFile("orchestra-ocr-poc-creds.json"), option.WithEndpoint("eu-documentai.googleapis.com:443"))
// 	if err != nil {
// 		log.Fatalf("Failed to create Document AI client: %v", err)
// 	}

// 	// Replace "path/to/your/document.pdf" with the path to your document
// 	//	documentPath := "invoice.pdf"
// 	documentPath := "invoice-2.png"
// 	// Process the document
// 	if err := processDocument(ctx, client, documentPath); err != nil {
// 		log.Fatalf("Error processing document: %v", err)
// 	}
// }

// func processDocument(ctx context.Context, client *documentai.DocumentProcessorClient, documentPath string) error {
// 	// Read the document file
// 	//	documentBytes, err := os.ReadFile(documentPath)
// 	documentBytes, err := ioutil.ReadFile(documentPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read document file: %v", err)
// 	}

// 	// form parser
// 	req := &documentaipb.ProcessRequest{
// 		Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", "orchestra-ocr-poc", "eu", "f840c64af9f123f7"),
// 		Source: &documentaipb.ProcessRequest_RawDocument{
// 			RawDocument: &documentaipb.RawDocument{
// 				Content:  documentBytes,
// 				MimeType: "image/png", //"application/pdf",
// 			},
// 		},
// 	}

// 	// invoice parser
// 	// req := &documentaipb.ProcessRequest{
// 	// 	Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", "orchestra-ocr-poc", "eu", "513169612c58c534"),
// 	// 	Source: &documentaipb.ProcessRequest_RawDocument{
// 	// 		RawDocument: &documentaipb.RawDocument{
// 	// 			Content:  documentBytes,
// 	// 			MimeType: "image/png", // "application/pdf",
// 	// 		},
// 	// 	},
// 	// }

// 	// Process the document
// 	resp, err := client.ProcessDocument(ctx, req)
// 	if err != nil {
// 		return fmt.Errorf("failed to process document: %v", err)
// 	}

// 	fmt.Println("here it is ", len(resp.Document.Pages))
// 	// Extract information from the response
// 	// for _, page := range resp.Document.GetPages() {
// 	// 	for _, formField := range page.GetFormFields() {
// 	// 		fmt.Printf("Field: %s, \n Value: %s\n", formField.GetFieldName(), formField.GetFieldValue())
// 	// 	}
// 	// }

// 	fmt.Println("*** ENTITIES *** ")
// 	for _, entity := range resp.GetDocument().Entities {
// 		fmt.Printf("Field: %v, \n Value: %v\n", entity.GetType(), entity.MentionText)
// 	}

// 	fmt.Println("*** ENTITIES PROPERTIES *** ")
// 	for _, entity := range resp.GetDocument().Entities {
// 		fmt.Println("*** NEW ENTITY *** ")
// 		for _, prop := range entity.Properties {
// 			fmt.Printf("KEY: %v, \n Value: %v\n", prop.GetType(), prop.MentionText)
// 		}

// 	}

// 	//for _ entity := range resp.getDo

// 	///	fmt.Println("RESULT: ", resp.Document.Pages[0])

// 	return nil
// }

// func main() {
// 	ctx := context.Background()

// 	// Authenticate using the service account key
// 	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile("orchestra-ocr-poc-creds.json"))
// 	if err != nil {
// 		log.Fatalf("Failed to create Vision client: %v", err)
// 	}

// 	// // Creates a client.
// 	// client, err := vision.NewImageAnnotatorClient(ctx)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create client: %v", err)
// 	// }
// 	defer client.Close()

// 	// Sets the name of the image file to annotate.
// 	filename := "invoice.pdf"

// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatalf("Failed to read file: %v", err)
// 	}
// 	defer file.Close()
// 	image, err := vision.NewImageFromReader(file)
// 	if err != nil {
// 		log.Fatalf("Failed to create image: %v", err)
// 	}

// 	labels, err := client.DetectLabels(ctx, image, nil, 10)
// 	if err != nil {
// 		log.Fatalf("Failed to detect labels: %v", err)
// 	}

// 	fmt.Println("Labels:")
// 	for _, label := range labels {
// 		fmt.Println(label.Description)
// 	}
// }
