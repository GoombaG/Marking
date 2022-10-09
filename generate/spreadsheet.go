package generate

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"unicode/utf8"

	"github.com/GoombaG/Marking/model"

	"github.com/xuri/excelize/v2"
)

// Generates a spreadsheet report for the entire class
func GenerateClassReport(s []model.Student, course model.Course) error {
	f := excelize.NewFile()
	Sheet1 := f.GetSheetName(0)

	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return err
	}

	// A slice of all the keys in generated to guarantee all future processing has the same order of strands
	keys := make([]string, len(course.Strands))
	i := 0
	for k := range course.Strands {
		keys[i] = k
		i++
	}

	// Set up column headings
	f.SetCellValue(Sheet1, "A1", "Student")
	for i = 0; i < len(keys); i++ {
		cellName, err := excelize.CoordinatesToCellName(i+2, 1)
		if err != nil {
			return err
		}
		f.SetCellValue(Sheet1, cellName, course.Strands[keys[i]].Name+" Term")
	}
	cellName, err := excelize.CoordinatesToCellName(i+2, 1)
	if err != nil {
		return err
	}
	f.SetCellValue(Sheet1, cellName, "Term Mark")
	cellName, err = excelize.CoordinatesToCellName(i+3, 1)
	if err != nil {
		return err
	}
	f.SetCellValue(Sheet1, cellName, "Masked Term Mark")
	cellName, err = excelize.CoordinatesToCellName(i+4, 1)
	if err != nil {
		return err
	}
	f.SetCellValue(Sheet1, cellName, "Summative Mark")
	cellName, err = excelize.CoordinatesToCellName(i+5, 1)
	if err != nil {
		return err
	}
	f.SetCellValue(Sheet1, cellName, "Exam Mark")
	cellName, err = excelize.CoordinatesToCellName(i+6, 1)
	if err != nil {
		return err
	}
	f.SetCellValue(Sheet1, cellName, "Final Mark")
	cellName, err = excelize.CoordinatesToCellName(i+7, 1)
	if err != nil {
		return err
	}
	f.SetCellValue(Sheet1, cellName, "Masked Final Mark")
	f.SetCellStyle("Sheet1", "A1", "Z1", boldStyle) // makes the top row bold

	// Add data under each heading
	for i = 0; i < len(s); i++ {
		// Add student's name
		cellName, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			return err
		}
		f.SetCellValue(Sheet1, cellName, s[i].Name)

		// Add strand marks
		j := 0
		for j = 0; j < len(keys); j++ {
			cellName, err = excelize.CoordinatesToCellName(j+2, i+2)
			if err != nil {
				return err
			}
			place(f, Sheet1, cellName, s[i].StrandMarks[keys[j]])
		}

		// Add unshadowed term mark
		cellName, err = excelize.CoordinatesToCellName(j+2, i+2)
		if err != nil {
			return err
		}
		place(f, Sheet1, cellName, s[i].TermMark)

		// Add shadowed term mark
		cellName, err = excelize.CoordinatesToCellName(j+3, i+2)
		if err != nil {
			return err
		}
		place(f, Sheet1, cellName, s[i].ShadowedTermMark)

		// Add summative mark
		cellName, err = excelize.CoordinatesToCellName(j+4, i+2)
		if err != nil {
			return err
		}
		place(f, Sheet1, cellName, s[i].SummativeMark)

		// Add exam mark
		cellName, err = excelize.CoordinatesToCellName(j+5, i+2)
		if err != nil {
			return err
		}
		place(f, Sheet1, cellName, s[i].ExamMark)

		// Add Unshadowed Final Mark
		cellName, err = excelize.CoordinatesToCellName(j+6, i+2)
		if err != nil {
			return err
		}
		place(f, Sheet1, cellName, s[i].UnshadowedFinalMark)

		// Add Shadowed Final Mark
		cellName, err = excelize.CoordinatesToCellName(j+7, i+2)
		if err != nil {
			return err
		}
		place(f, Sheet1, cellName, s[i].ShadowedFinalMark)
	}

	// Autofit all columns according to their text content
	// Modified code taken from https://github.com/qax-os/excelize/issues/92
	cols, err := f.GetCols(Sheet1)
	if err != nil {
		return err
	}
	i = 0
	for idx, col := range cols {
		if i > len(keys)+6 {
			break
		}
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return err
		}
		f.SetColWidth(Sheet1, name, name, float64(largestWidth))
		i++
	}
	// End code from https://github.com/qax-os/excelize/issues/92

	t := time.Now().Format("2006-01-02 Mon")
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	var reportPath string
	if runtime.GOOS == "windows" {
		reportPath = dir + "\\Downloads\\Class Report " + t + ".xlsx"
	} else {
		reportPath = dir + "/Downloads/Class Report " + t + ".xlsx"
	}

	//fmt.Println(reportPath)
	err = f.SaveAs(reportPath)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// Places a cell value but checks for a '-1', which indicates that there is No Mark
func place(f *excelize.File, sheet string, cellName string, mark float64) {
	if mark != -1 {
		f.SetCellValue(sheet, cellName, mark)
	} else {
		f.SetCellValue(sheet, cellName, "No Mark")
	}
}
