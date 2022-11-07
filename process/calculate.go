package process

import (
	"fmt"
	"math"

	c "github.com/GoombaG/Marking/data"
	"github.com/GoombaG/Marking/model"
)

// Calculates and returns single strand mark for the term
func calculateStrandMarkTerm(strand model.Strand, student model.Student) float64 {
	strandMark := 0.0
	totalWeight := 0.0
	for _, s := range strand.Expectations {
		for i := 0; i < len(s.Evaluations); i++ {
			if s.Evaluations[i].Type == "E" || s.Evaluations[i].Type == "S" {
				continue
			}
			numericMark := 0.0
			numericMark = c.LetterToNumeric[student.Marks[s.Evaluations[i].ID]]
			if numericMark == math.MaxFloat64 {
				continue
			}
			numericMark = numericMark * s.Evaluations[i].Weight
			strandMark += numericMark
			totalWeight += s.Evaluations[i].Weight
		}
	}
	if totalWeight == 0 {
		return -1
	}
	return strandMark / totalWeight
}

// Calculates all strand marks for the term, and returns an updated student
func calculateStrandMarksTerm(strands map[string]model.Strand, student model.Student) model.Student {
	student.StrandMarks = make(map[string]float64)
	for k := range strands {
		student.StrandMarks[k] = calculateStrandMarkTerm(strands[k], student)
	}
	return student
}

// Calculates and returns the overall mark for the term
func calculateTotalMarkTerm(strands map[string]model.Strand, strandMarks map[string]float64) float64 {
	termMark := 0.0
	totalWeight := 0.0
	for k := range strands {
		if strandMarks[k] == -1 {
			continue
		}
		totalWeight += strands[k].Weight
		termMark += strandMarks[k] * strands[k].Weight
	}
	if totalWeight == 0 {
		return -1
	}
	return termMark / totalWeight
}

// Calculates and returns the mark for a specific type of evaluation
// This can be used to calculate the summative or exam mark
func calculateTypeMark(strands map[string]model.Strand, student model.Student, chr string) float64 {
	totalMark := 0.0
	totalWeight := 0.0
	for k, strand := range strands {
		_ = k
		for l, expectation := range strand.Expectations {
			_ = l
			for m, evaluation := range expectation.Evaluations {
				_ = m
				if evaluation.Type != chr {
					continue
				}
				numericMark := c.LetterToNumeric[student.Marks[evaluation.ID]]
				if numericMark == math.MaxFloat64 {
					continue
				}
				numericMark = numericMark * evaluation.Weight
				totalMark += numericMark
				totalWeight += evaluation.Weight
			}
		}
	}
	if totalWeight == 0 {
		return -1
	}
	return totalMark / totalWeight
}

// Calculates and returns the exam mark
func calculateExamMark(strands map[string]model.Strand, student model.Student) float64 {
	return calculateTypeMark(strands, student, "E")
}

// Calculates and returns the summative mark
func calculateSummativeMark(strands map[string]model.Strand, student model.Student) float64 {
	return calculateTypeMark(strands, student, "S")
}

// Calculates and returns an overall mark
func calculateFinalMark(student model.Student) float64 {
	// If no marks
	if student.TermMark == -1 && student.ExamMark == -1 && student.SummativeMark == -1 {
		return -1.0
	}

	if c.SummaryWeightKey {
		totalMark := 0.0
		totalWeight := 0.0
		if student.TermMark != -1 {
			totalMark += student.TermMark * c.SummaryWeights["Term"]
			totalWeight += c.SummaryWeights["Term"]
		}
		if student.SummativeMark != -1 {
			totalMark += student.SummativeMark * c.SummaryWeights["Summative"]
			totalWeight += c.SummaryWeights["Summative"]
		}
		if student.ExamMark != -1 {
			totalMark += student.ExamMark * c.SummaryWeights["Exam"]
			totalWeight += c.SummaryWeights["Exam"]
		}

		if totalWeight == 0.0 {
			return -1
		}
		return totalMark / totalWeight
	} else {
		if student.ExamMark == -1 && student.SummativeMark == -1 {
			return student.TermMark
		} else if student.ExamMark == -1 {
			return student.TermMark*0.7 + student.SummativeMark*0.3
		} else if student.SummativeMark == -1 {
			return student.TermMark*0.7 + student.ExamMark*0.3
		} else {
			return student.TermMark*+student.ExamMark*0.2 + student.SummativeMark*0.1
		}
	}
}

// Shadows overall term marks using exam marks
func shadowTermMarks(student model.Student, course model.Course) model.Student {
	for key, strand := range course.Strands {
		examMark := 0.0
		examMarkWeight := 0.0
		for _, s := range strand.Expectations {
			for i := 0; i < len(s.Evaluations); i++ {
				if s.Evaluations[i].Type != "E" {
					continue
				}
				numericMark := 0.0
				numericMark = c.LetterToNumeric[student.Marks[s.Evaluations[i].ID]]
				if numericMark == math.MaxFloat64 {
					continue
				}
				numericMark = numericMark * s.Evaluations[i].Weight
				examMark += numericMark
				examMarkWeight += s.Evaluations[i].Weight
			}
		}
		if examMarkWeight == 0.0 {
			continue
		}
		examMark /= examMarkWeight
		if student.StrandMarks[key] == -1 {
			continue
		}

		student.StrandMarks[key] = math.Max(student.StrandMarks[key], examMark)
	}
	return student
}

// Calculates and returns both the shadowed term mark and shadowed final mark
func calculateShadowedMarks(student model.Student, course model.Course) (float64, float64) {
	student = shadowTermMarks(student, course)
	student.TermMark = calculateTotalMarkTerm(course.Strands, student.StrandMarks)
	return student.TermMark, calculateFinalMark(student)
}

// Removes ignored levels so they aren't displayed on reports
func purgeExtraLevels(student model.Student, course model.Course) model.Student {
	for _, strand := range course.Strands {

		for _, expectation := range strand.Expectations {
			for i := 0; i < len(expectation.Evaluations); i++ {
				toPurge := true
				for j := 0; j < len(c.MarksInOrder); j++ {
					if student.Marks[expectation.Evaluations[i].ID] == c.MarksInOrder[j] ||
						c.LetterToNumeric[student.Marks[expectation.Evaluations[i].ID]] == math.MaxFloat64 {
						toPurge = false
						break
					}
				}

				if toPurge {
					numericMark := c.LetterToNumeric[student.Marks[expectation.Evaluations[i].ID]]
					newMark := purgedLetterMark(numericMark)
					student.Marks[expectation.Evaluations[i].ID] = newMark
				}
			}
		}
	}

	return student
}

func purgedLetterMark(mark float64) string {
	diff := 1000.0
	letterMark := "ERROR"
	for letterMark == "ERROR" {
		for key, element := range c.LetterToNumeric {
			toPurge := true
			for j := 0; j < len(c.MarksInOrder); j++ {
				if key == c.MarksInOrder[j] {
					toPurge = false
					break
				}
			}

			if toPurge || ((math.Abs(element-mark)) >= 1 && (element < mark)) || math.Abs(element-mark) > diff {
				continue
			}
			diff = math.Abs(element - mark)
			letterMark = key
		}
		mark--
	}
	return letterMark
}

// Calculates corresponding letter grades for each of a given student's strand marks
// Returns the strand marks
func calculateLetterStrandMarks(student model.Student) map[string]string {
	student.StrandLetterMarks = make(map[string]string)
	for key, mark := range student.StrandMarks {
		if mark == -1 {
			student.StrandLetterMarks[key] = ""
		} else {
			student.StrandLetterMarks[key] = purgedLetterMark(mark)
		}
	}
	fmt.Println(student.StrandLetterMarks)
	return student.StrandLetterMarks
}
