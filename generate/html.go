package generate

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"

	c "github.com/GoombaG/Marking/data"
	"github.com/GoombaG/Marking/model"
)

//go:embed reportTemplate.html
var reportTemplate string

// Generates an HTML page with reports for some students
func GenerateHTMLReport(students []model.Student, course model.Course, legend bool) {
	t, err := template.New("Report Page").Parse(reportTemplate)
	if err != nil {
		fmt.Println("Error parsing template")
	}

	reportFileName := "report.html"

	// Creates a new file or overwrites the current one
	f, err := os.OpenFile(reportFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error creating output file")
	}

	defer f.Close()

	var class model.Class
	class.Students = students
	class.Course = course
	// These are all the mark types that are displayed on the page
	class.MarkTypes = c.MarksInOrder
	class.MarkColWidth = 62.0 / float64(len(c.MarksInOrder))
	if legend {
		class.MarkColWidth = 50.0 / float64(len(c.MarksInOrder))
	}
	class.EvaluationsInOrder = c.EvaluationsInOrder
	class.Legend = legend

	t.Execute(f, class)
}
