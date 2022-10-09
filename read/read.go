/*
This file contains all the functions that read from the xlsx file
*/

package read

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	c "github.com/GoombaG/Marking/data"
	"github.com/GoombaG/Marking/model"

	"github.com/xuri/excelize/v2"
)

// Reads in all neccesary data from the spreadsheet
func ReadSheet(sheet string) ([]model.Student, model.Course, string) {

	var newCourse model.Course
	f, err := excelize.OpenFile(sheet)
	if err != nil {
		errorString := fmt.Sprintf("Error: %v, trouble opening file: \"%s\"", err, sheet)
		return nil, newCourse, errorString
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	errorString := readMarks(f)
	if errorString != "" {
		return nil, newCourse, errorString
	}

	newCourse.Strands, errorString = readStrands(f)
	if errorString != "" {
		return nil, newCourse, errorString
	}

	newCourse.Strands, errorString = readExpectations(f, newCourse.Strands)
	if errorString != "" {
		return nil, newCourse, errorString
	}

	students, errorString := readStudents(f, newCourse)
	if errorString != "" {
		return students, newCourse, errorString
	}

	students, newCourse, errorString = readEvaluations(f, students, newCourse)
	if errorString != "" {
		return students, newCourse, errorString
	}

	errorString = readSummary(f)
	if errorString != "" {
		return students, newCourse, errorString
	}

	errorString, newCourse = readCourse(f, newCourse)
	if errorString != "" {
		return students, newCourse, errorString
	}

	return students, newCourse, ""
}

// Reads in all strand data from the spreadsheet
func readStrands(f *excelize.File) (map[string]model.Strand, string) {
	strand := make(map[string]model.Strand)

	column, row, errorString := findCoordinates(f, "Strand Code")
	if errorString != "" {
		return strand, errorString
	}
	row++

	for {
		var newStrand model.Strand

		newName, errorString := readCellFromCoordinates(f, column+2, row)
		if errorString != "" {
			return strand, errorString
		}
		if newName == "" {
			return strand, ""
		}
		newStrand.Name = newName

		newCode, errorString := readCellFromCoordinates(f, column, row)
		if newCode == "" {
			errorString = fmt.Sprintf("Error: strand \"%s\" is missing a code", newName)
		}
		if errorString != "" {
			return strand, errorString
		}

		newWeight, errorString := readCellFromCoordinates(f, column+1, row)
		if newWeight == "" {
			errorString = fmt.Sprintf("Error: strand \"%s\" is missing a weight", newName)
		}
		if errorString != "" {
			return strand, errorString
		}

		var err error
		newStrand.Weight, err = strconv.ParseFloat(newWeight, 64)
		if err != nil {
			errorString := fmt.Sprintf("Error: %v, trouble parsing Strand Weight", err)
			return strand, errorString
		}

		newStrand.Expectations = make(map[string]model.Expectation)

		row++
		strand[newCode] = newStrand
	}
}

// Reads in all expectation data from the spreadsheet
func readExpectations(f *excelize.File, strands map[string]model.Strand) (map[string]model.Strand, string) {

	column, row, errorString := findCoordinates(f, "Expectation Strand")
	if errorString != "" {
		return strands, errorString
	}
	row++

	for {
		var newExpectation model.Expectation

		newName, errorString := readCellFromCoordinates(f, column+1, row)
		if errorString != "" {
			return strands, errorString
		}
		if newName == "" {
			return strands, ""
		}
		newExpectation.Name = newName

		newCode, errorString := readCellFromCoordinates(f, column+2, row)
		if newCode == "" {
			errorString = fmt.Sprintf("Error: strand \"%s\" is missing a code", newName)
		}
		if errorString != "" {
			return strands, errorString
		}

		newStrandCode, errorString := readCellFromCoordinates(f, column, row)
		if newStrandCode == "" {
			errorString = fmt.Sprintf("Error: strand \"%s\" is missing a strand code", newName)
		}
		if errorString != "" {
			return strands, errorString
		}

		row++

		// verify that the strand code exists
		validStrandCode := false
		for key := range strands {
			if key == newStrandCode {
				validStrandCode = true
			}
		}
		if validStrandCode {
			strands[newStrandCode].Expectations[newCode] = newExpectation
		} else {
			errorString = fmt.Sprintf("Error: strand \"%s\" has an invalid strand code \"%s\"", newName, newStrandCode)
			return strands, errorString
		}
	}
}

func readMarkTypes(f *excelize.File, course model.Course) (model.Course, string) {
	column, row, errorString := findCoordinates(f, "Evaluation Type")
	if errorString != "" {
		return course, errorString
	}

	for {
		var newMarkType model.Print
		newType, errorString := readCellFromCoordinates(f, column, row)
		if errorString != "" {
			return course, errorString
		}
		if newType == "" {
			return course, ""
		}
		newMarkType.Name = newType

		newColour, errorString := cell.GetStyle
	}
}

// Reads in some evaluation data from the spreadsheet
// DOESN'T READ the evaluation key/name block
func readEvaluations(f *excelize.File, students []model.Student, course model.Course) ([]model.Student, model.Course, string) {
	column, row, errorString := findCoordinates(f, "Evaluation")
	if errorString != "" {
		return students, course, errorString
	}
	column++

	IDNumber := 1

	for {
		var newEvaluation model.Evaluation
		newKey, errorString := readCellFromCoordinates(f, column, row)
		if errorString != "" {
			return students, course, errorString
		}
		if newKey == "" {
			return students, course, ""
		}
		newEvaluation.Name = newKey

		newExpectationCode, errorString := readCellFromCoordinates(f, column, row+1)
		if newExpectationCode == "" {
			errorString = fmt.Sprintf("Error: evaluation \"%s\" is missing an expectation code", newKey)
		}
		if errorString != "" {
			return students, course, errorString
		}

		newWeight, errorString := readCellFromCoordinates(f, column, row+2)
		if newWeight == "" {
			errorString = fmt.Sprintf("Error: evaluation \"%s\" is missing a weight", newKey)
		}
		if errorString != "" {
			return students, course, errorString
		}

		var err error
		newEvaluation.Weight, err = strconv.ParseFloat(newWeight, 64)
		if err != nil {
			errorString := fmt.Sprintf("Error: %v, trouble parsing evaluation weight", err)
			return students, course, errorString
		}

		newType, errorString := readCellFromCoordinates(f, column, row+3)
		if newType == "" {
			errorString = fmt.Sprintf("Error: evaluation \"%s\" is missing a type", newKey)
		}
		if errorString != "" {
			return students, course, errorString
		}
		newEvaluation.Type = newType

		switch newEvaluation.Type {
		case "Q":
			newEvaluation.Colour = "#8200BE" // PURPLE
			break
		case "E":
			newEvaluation.Colour = "#0000FF" // BLUE
			break
		case "T":
			newEvaluation.Colour = "#FF0000" // RED
			break
		case "S":
			newEvaluation.Colour = "#000000" // BLACK
			break
		case "P":
			newEvaluation.Colour = "#FF55ED" // PINK
			break
		default:
			errorString := fmt.Sprintf("Error: evaluation \"%s\" has an invalid type \"%s\"", newKey, newEvaluation.Type)
			return students, course, errorString
		}

		newEvaluation.ID = IDNumber
		IDNumber++

		for i := 0; i < len(students); i++ {
			newMark, errorString := readCellFromCoordinates(f, column, row+5+i)
			if newMark == "" {
				errorString = fmt.Sprintf("Error: %s is missing a mark for %s", students[i].Name, newEvaluation.Name)
			}
			if errorString != "" {
				return students, course, errorString
			}

			validMark := false
			for key := range c.LetterToNumeric {
				if newMark == key {
					validMark = true
					break
				}
			}

			if !validMark {
				errorString = fmt.Sprintf("Error: %s has an invalid mark for %s", students[i].Name, newEvaluation.Name)
				return students, course, errorString
			}

			students[i].Marks[newEvaluation.ID] = newMark
		}

		validExpectationCode := false
		for k := range course.Strands {
			tempStrand := course.Strands[k]
			for s := range course.Strands[k].Expectations {
				tempExpectation := course.Strands[k].Expectations[s]
				if s == newExpectationCode {
					tempExpectation.Evaluations = append(tempExpectation.Evaluations, newEvaluation)
					tempStrand.Expectations[s] = tempExpectation
					course.Strands[k] = tempStrand
					validExpectationCode = true
					break
				}
			}
			if validExpectationCode {
				break
			}
		}

		if !validExpectationCode {
			errorString = fmt.Sprintf("Error: evaluation \"%s\" has an invalid expectation code \"%s\"", newKey, newExpectationCode)
			return students, course, errorString
		}

		newName, errorString := readCellFromCoordinates(f, column, row+4)
		if errorString != "" {
			return students, course, errorString
		}

		if newName != "" {

			var newPrintEvaluation model.Print
			newPrintEvaluation.Colour = newEvaluation.Colour
			newPrintEvaluation.Name = newKey + ": " + newName

			repeat := false
			for i := 0; i < len(c.EvaluationsInOrder); i++ {
				if newPrintEvaluation.Name == c.EvaluationsInOrder[i].Name {
					repeat = true
				}
			}

			if !repeat {
				c.EvaluationsInOrder = append(c.EvaluationsInOrder, newPrintEvaluation)
			}
		}
		column++
	}
}

// Reads in all student data from the spreadsheet
func readStudents(f *excelize.File, course model.Course) ([]model.Student, string) {
	var students []model.Student

	column, row, errorString := findCoordinates(f, "Evaluation")
	if errorString != "" {
		return students, errorString
	}
	row += 5
	startColumn := column

	for {
		column = startColumn
		var newStudent model.Student
		newName, errorString := readCellFromCoordinates(f, column, row)
		if errorString != "" {
			return students, errorString
		}
		if newName == "" {
			break
		}
		newStudent.Name = newName
		newStudent.Marks = make(map[int]string)
		students = append(students, newStudent)
		row++
	}
	return students, ""
}

func readMarks(f *excelize.File) string {
	column, row, errorString := findCoordinates(f, "Level")
	if errorString != "" {
		return errorString
	}

	c.LetterToNumeric = make(map[string]float64)
	c.MarksInOrder = nil

	lock, errorString := readCellFromCoordinates(f, column+3, row+1)
	if errorString != "" {
		return errorString
	}

	if strings.ToUpper(lock) != "Y" {
		c.LetterToNumeric = map[string]float64{
			"A": math.MaxFloat64, "B": 0, "R-": 25, "R": 35, "R+": 44, "1-": 52, "1": 55, "1+": 58, "2-": 62, "2": 65, "2+": 68, "3-": 72, "3": 75, "3+": 78, "4-": 86, "4": 94, "4+": 100}
		c.MarksInOrder = append(c.MarksInOrder, "A", "B", "R-", "R", "R+", "1-", "1", "1+", "2-", "2", "2+", "3-", "3", "3+", "4-", "4", "4+")
	} else {
		for {
			row++
			newLetterMark, errorString := readCellFromCoordinates(f, column, row)
			if errorString != "" {
				return errorString
			}
			if newLetterMark == "" {
				break
			}

			newNumericMark, errorString := readCellFromCoordinates(f, column+1, row)
			if newNumericMark == "" {
				errorString = fmt.Sprintf("Error: level %s is missing a percent", newLetterMark)
			}
			if errorString != "" {
				return errorString
			}

			newIgnore, errorString := readCellFromCoordinates(f, column+2, row)
			if errorString != "" {
				return errorString
			}

			if strings.ToUpper(newIgnore) == "" {
				c.MarksInOrder = append(c.MarksInOrder, newLetterMark)
			}

			if strings.ToUpper(newNumericMark) == "NM" || strings.ToUpper(newNumericMark) == "NO MARK" {
				c.LetterToNumeric[newLetterMark] = math.MaxFloat64
			} else {
				newFloatMark, err := strconv.ParseFloat(newNumericMark, 64)
				if err != nil {
					errorString := fmt.Sprintf("Error: %v, trouble parsing level percent", err)
					return errorString
				}
				c.LetterToNumeric[newLetterMark] = newFloatMark
			}
		}
	}
	return ""
}

func readSummary(f *excelize.File) string {
	column, row, errorString := findCoordinates(f, "Summary")
	if errorString != "" {
		return errorString
	}

	c.SummaryWeights = make(map[string]float64)
	c.SummaryWeightKey = false

	lock, errorString := readCellFromCoordinates(f, column+2, row+1)
	if errorString != "" {
		return errorString
	}

	if strings.ToUpper(lock) != "Y" {
		return ""
	} else {
		c.SummaryWeightKey = true
		for {
			row++
			newType, errorString := readCellFromCoordinates(f, column, row)
			if errorString != "" {
				return errorString
			}
			if newType == "" {
				break
			}
			if newType != "Exam" && newType != "Term" && newType != "Summative" {
				errorString = fmt.Sprintf("Error: summary type \"%s\" is invalid", newType)
				return errorString
			}

			newWeight, errorString := readCellFromCoordinates(f, column+1, row)
			if newWeight == "" {
				errorString = fmt.Sprintf("Error: missing weight for %s", newType)
			}
			if errorString != "" {
				return errorString
			}

			weight, err := strconv.ParseFloat(newWeight, 64)
			if err != nil {
				errorString := fmt.Sprintf("Error: %v, trouble parsing summary weight", err)
				return errorString
			}
			c.SummaryWeights[newType] = weight
		}
	}
	return ""
}

func readCourse(f *excelize.File, course model.Course) (string, model.Course) {
	column, row, errorString := findCoordinates(f, "Course")
	if errorString != "" {
		return errorString, course
	}

	courseName, errorString := readCellFromCoordinates(f, column+1, row)
	if errorString != "" {
		return errorString, course
	}
	course.Name = courseName

	teacher, errorString := readCellFromCoordinates(f, column+1, row+1)
	if errorString != "" {
		return errorString, course
	}
	course.Teacher = teacher

	return "", course
}

// Takes in the spreadsheet file and a column and a row number.
// Returns the contents of the cell at the provided column/row
func readCellFromCoordinates(f *excelize.File, column, row int) (string, string) {
	var cellValue string = ""
	currentLocation, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		errorString := fmt.Sprintf("Error: %v, trouble converting coordinates to cell name", err)
		return "", errorString
	}

	cellValue, err = f.GetCellValue(f.GetSheetName(0), currentLocation)
	if err != nil {
		errorString := fmt.Sprintf("Error: %v, trouble getting cell value", err)
		return "", errorString
	}

	return cellValue, ""
}

// Takes in the spreadsheet file and a string, and returns the
// coordinates of the first cell found that contains that text
func findCoordinates(f *excelize.File, keyword string) (int, int, string) {
	startLocation, err := f.SearchSheet(f.GetSheetName(0), keyword)
	if err != nil {
		errorString := fmt.Sprintf("Error: %v, trouble finding \"%s\"", err, keyword)
		return -1, -1, errorString
	}

	if len(startLocation) == 0 {
		errorString := fmt.Sprintf("Error: trouble finding keyword \"%s\"", keyword)
		return -1, -1, errorString
	}

	f.GetCellValue(f.GetSheetName(0), startLocation[0])
	column, row, err := excelize.CellNameToCoordinates(startLocation[0])
	if err != nil {
		errorString := fmt.Sprintf("Error: %v, trouble converting cell name to coordinates", err)
		return -1, -1, errorString
	}

	return column, row, ""
}
