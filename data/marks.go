package data

// Allows for conversion of string grades to numeric grades
var LetterToNumeric = map[string]float64{}

// This is used to preserve the order of the string grades
var MarksInOrder []string

// This is used to allow teachers to alter weighting of exam/summative/term
var SummaryWeights = map[string]float64{}
var SummaryWeightKey bool
