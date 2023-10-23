package usecase

import (
	"bytes"
	"fmt"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"os"
)

func ExampleNewPDFGenerator() {

	file, err2 := os.Open("./Komon - Report Dev/index.html")
	if err2 != nil {
		fmt.Println(err2)
	}
	defer file.Close()
	b1 := make([]byte, 21000)
	_, err2 = file.Read(b1)

	pdfGenerator, err := pdf.NewPDFGenerator()
	pdfGenerator.Dpi.Set(200)
	pdfGenerator.PageSize.Set(pdf.PageSizeA4)

	pdfGenerator.Orientation.Set(pdf.OrientationPortrait)
	page := pdf.NewPageReader(bytes.NewReader(b1))

	page.EnableLocalFileAccess.Set(true)

	page.Allow.Set("C:\\Users\\softw\\go-clean-arch-mongo\\Komon - Report Dev")

	pdfGenerator.AddPage(page)

	err = pdfGenerator.Create()
	if err != nil {
		fmt.Println("pdfGenerator.Create error: ", err)
	}

	err = pdfGenerator.WriteFile("./test.pdf")

	// Output: Done
}
