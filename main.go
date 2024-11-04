package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// Ticket struct
type Ticket struct {
	LoginUser   string    `json:"loginUser"`
	Terminal    string    `json:"terminal"`
	Date        string    `json:"date"`
	Time        string    `json:"time"`
	PaymentDate string    `json:"payment_date"`
	PaymentTime string    `json:"payment_time"`
	PaymentType string    `json:"payment_type"`
	Tag         Tag       `json:"tag"`
	Payments    []Payment `json:"payments"`
	Orders      []Order   `json:"orders"`
}

// Tag struct. Example: {"Pax":"100", "PaxTime":"2020/10/10"}
type Tag struct {
	Pax     int       `json:"Pax"`
	PaxTime time.Time `json:"PaxTime"`
}

// PaymentInformation struct. Example: {"RefNo":"100", "RefTime":"2020/10/10"}
type PaymentInformation struct {
	RefNo   int       `json:"RefNo"`
	RefTime time.Time `json:"RefTime"`
}

// Payment struct
type Payment struct {
	Name               string             `json:"name"`
	Tendered           string             `json:"tendered"`
	PaymentInformation PaymentInformation `json:"payment_information"`
}

// Order struct
type Order struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

func main() {
	var templatePath string
	var outputPath string
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Print to POS",
		Long:  "Print to a POS printer interpolating a template with an object",
		Run: func(cmd *cobra.Command, args []string) {
			// Resolve template path

			ticket := createTicket()

			parsedTemplate := loadTemplate(templatePath)
			err := parsedTemplate.Execute(os.Stdout, ticket)
			if err != nil {
				fmt.Println("Error rendering template:", err)
			}
		},
	}

	// Add flags to the root command
	rootCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to the template file (required)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to the output file")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Mark template flag as required
	// rootCmd.MarkFlagRequired("template")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadTemplate(templatePath string) *template.Template {
	if templatePath == "" {
		templatePath = defaultTemplatePath()
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		os.Exit(1)
	}

	return tmpl
}

func defaultTemplatePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	return filepath.Join(homeDir, "default_template")
}

func createTicket() Ticket {
	return Ticket{
		LoginUser:   "betasve",
		Terminal:    "pos terminal",
		Date:        time.Date(2024, time.October, 31, 0, 0, 0, 00, time.UTC).Format("2006-01-02"),
		Time:        time.Now().Format("15:04:05"),
		PaymentDate: time.Date(2024, time.October, 31, 17, 35, 24, 00, time.FixedZone("GMT", 2)).Format("2006-01-02"),
		PaymentTime: time.Date(2024, time.October, 31, 17, 35, 24, 00, time.FixedZone("GMT", 2)).Format("15:04:05"),
		PaymentType: "credit_card",
		Tag: Tag{
			Pax:     100,
			PaxTime: time.Date(2024, time.October, 31, 17, 35, 24, 00, time.FixedZone("GMT", 2)),
		},
		Payments: []Payment{
			{
				Name:     "Payment 1",
				Tendered: "is tendered",
				PaymentInformation: PaymentInformation{
					RefNo:   101,
					RefTime: time.Date(2024, time.October, 31, 17, 35, 24, 00, time.FixedZone("GMT", 2)),
				},
			},
		},
		Orders: []Order{
			{
				Name:     "Fried beans",
				Quantity: "3",
				Price:    "1.99",
			},
		},
	}
}
