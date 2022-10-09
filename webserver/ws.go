package ws

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/GoombaG/Marking/data"
	"github.com/GoombaG/Marking/generate"
	"github.com/GoombaG/Marking/model"
	"github.com/GoombaG/Marking/process"
)

//go:embed requestPageTemplate.html
var requestPageTemplate string

//go:embed fileSelect.html
var fileSelect string

var students []model.Student
var course model.Course
var err error
var errorString = ""

func Start() {
	http.HandleFunc("/", requestFile)
	http.HandleFunc("/mf", markFile)
	http.HandleFunc("/downloadSheet", downloadSheet)
	http.HandleFunc("/singlePDF", createHTML)
	http.HandleFunc("/manyPDF", createPDFs)

	http.ListenAndServe("localhost:8080", nil)
}

func requestFile(w http.ResponseWriter, req *http.Request) {
	if errorString != "" {
		fmt.Fprintf(w, errorString+"\n\n")
		errorString = ""
	} else {
		fmt.Fprintf(w, fileSelect)
	}
}

func markFile(w http.ResponseWriter, req *http.Request) {
	if errorString != "" {
		http.Redirect(w, req, "/", 302)
	}

	if err := req.ParseForm(); err != nil {
		errorString = fmt.Sprintf("ParseForm() err: %v", err)
		http.Redirect(w, req, "/", 302)
	}
	filePath := req.FormValue("markFile")

	students, course, errorString = process.ProcessSheet(filePath)
	if errorString != "" {
		http.Redirect(w, req, "/", 302)
	}

	t, err := template.New("Request Page").Parse(requestPageTemplate)
	if err != nil {
		errorString = fmt.Sprintf("Error parsing request page template %v", err)
		http.Redirect(w, req, "/", 302)
	}

	f, err := os.OpenFile("requestPage.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		errorString = fmt.Sprintf("Error creating output file: %v", err)
		http.Redirect(w, req, "/", 302)
	}

	defer f.Close()

	var c model.Class
	c.Students = students
	c.Course = course
	c.LandscapeOption = false
	if len(data.EvaluationsInOrder) > 0 {
		c.LandscapeOption = true
	}

	t.Execute(f, c)

	http.ServeFile(w, req, "requestPage.html")
}

func downloadSheet(w http.ResponseWriter, req *http.Request) {
	err = generate.GenerateClassReport(students, course)
	if err != nil {
		errorString = fmt.Sprintf("Error creating class report %v", err)
	}
}

func createHTML(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var selectedStudents []model.Student
	for _, s := range students {
		if _, ok := req.Form["SELECTED"+s.Name]; ok {
			if _, ok := req.Form["SUMMARY"+s.Name]; ok {
				s.ShowSummary = true
			} else {
				s.ShowSummary = false
			}
			selectedStudents = append(selectedStudents, s)
		}
	}
	if _, ok := req.Form["LANDSCAPE"]; ok {
		generate.GenerateHTMLReport(selectedStudents, course, true)
	} else {
		generate.GenerateHTMLReport(selectedStudents, course, false)
	}
	http.ServeFile(w, req, "report.html")
}

func createPDFs(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	for _, s := range students {
		var selectedStudents []model.Student
		if _, ok := req.Form["SELECTED"+s.Name]; ok {
			if _, ok := req.Form["SUMMARY"+s.Name]; ok {
				s.ShowSummary = true
			} else {
				s.ShowSummary = false
			}
			selectedStudents = append(selectedStudents, s)
			if _, ok := req.Form["LANDSCAPE"]; ok {
				generate.GenerateHTMLReport(selectedStudents, course, true)
			} else {
				generate.GenerateHTMLReport(selectedStudents, course, false)
			}
		} else {
			continue
		}

		var pdfName string
		t := time.Now().Format("2006-01-02 Mon")
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error obtaining home directory: ", err)
		}

		if runtime.GOOS == "windows" {
			pdfName = dir + "\\Downloads\\" + s.LastName + "_" + s.FirstName + " " + t + ".pdf"
			os.Remove(pdfName)

			var out []byte
			var err error

			if _, ok := req.Form["LANDSCAPE"]; ok {
				out, err = exec.Command(dir+"\\..\\..\\Program Files\\wkhtmltopdf\\bin\\wkhtmltopdf", "--disable-smart-shrinking", "--footer-right", "Page [page] out of [topage]", "-O", "landscape", "--dpi", "300", "--page-size", "Letter",
					"report.html", pdfName).Output()
			} else {
				out, err = exec.Command(dir+"\\..\\..\\Program Files\\wkhtmltopdf\\bin\\wkhtmltopdf", "--disable-smart-shrinking", "--footer-right", "Page [page] out of [topage]", "--dpi", "300", "--page-size", "Letter",
					"report.html", pdfName).Output()
			}

			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Printf(string(out))
			}

		} else {
			pdfName = dir + "/Downloads/" + s.LastName + "_" + s.FirstName + " " + t + ".pdf"
			os.Remove(pdfName)

			var out []byte
			var err error

			if _, ok := req.Form["LANDSCAPE"]; ok {
				out, err = exec.Command("wkhtmltopdf", "--disable-smart-shrinking", "--footer-right", "Page [page] out of [topage]", "-O", "landscape", "--page-size", "Letter",
					"report.html", pdfName).Output()
			} else {
				out, err = exec.Command("wkhtmltopdf", "--disable-smart-shrinking", "--footer-right", "Page [page] out of [topage]", "--page-size", "Letter",
					"report.html", pdfName).Output()
			}

			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Printf(string(out))
			}
		}
	}
}
