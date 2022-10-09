package process

import (
	"strings"

	c "github.com/GoombaG/Marking/data"
	"github.com/GoombaG/Marking/model"
	"github.com/GoombaG/Marking/read"
)

// Reads all the data from a provided 'sheet', and performs all calculations
// Returns complete structures of all students and their course
func ProcessSheet(sheet string) ([]model.Student, model.Course, string) {
	c.StoreEvaluationNames = false
	c.EvaluationsInOrder = nil
	students, course, errorString := read.ReadSheet(sheet)

	if errorString != "" {
		return students, course, errorString
	}

	for i := 0; i < len(students); i++ {
		students[i] = calculateStrandMarksTerm(course.Strands, students[i])
		students[i].StrandLetterMarks = calculateLetterStrandMarks(students[i])
		students[i].TermMark = calculateTotalMarkTerm(course.Strands, students[i].StrandMarks)
		students[i].ExamMark = calculateExamMark(course.Strands, students[i])
		students[i].SummativeMark = calculateSummativeMark(course.Strands, students[i])
		students[i].UnshadowedFinalMark = calculateFinalMark(students[i])

		strandMarksCopy := make(map[string]float64)
		for k, v := range students[i].StrandMarks {
			strandMarksCopy[k] = v
		}
		students[i].ShadowedTermMark, students[i].ShadowedFinalMark = calculateShadowedMarks(students[i], course)
		students[i].StrandMarks = strandMarksCopy

		names := strings.Split(students[i].Name, ",")
		students[i].LastName = names[0]
		students[i].FirstName = names[1]

		students[i] = purgeExtraLevels(students[i], course)
	}
	return students, course, ""
}
