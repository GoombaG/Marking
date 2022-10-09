package model

type Course struct {
	Strands map[string]Strand
	Name    string
	Teacher string
	MarkTypes map[string]Print
}

type Strand struct {
	Name         string
	Weight       float64
	Expectations map[string]Expectation
}

type Expectation struct {
	Name        string
	Evaluations []Evaluation
}

type Evaluation struct {
	Name   string
	Weight float64
	Type   string
	ID     int
	Colour string
}

type Student struct {
	Name                string
	FirstName           string
	LastName            string
	StrandMarks         map[string]float64 // map key is strand name
	StrandLetterMarks   map[string]string  // map key is strand name
	TermMark            float64
	ExamMark            float64
	SummativeMark       float64
	ShadowedTermMark    float64
	UnshadowedFinalMark float64
	ShadowedFinalMark   float64
	Marks               map[int]string // map key is evaluation ID
	ShowSummary         bool
}

type Class struct {
	Students           []Student
	Course             Course
	MarkTypes          []string
	MarkColWidth       float64
	LandscapeOption    bool
	Legend             bool
	EvaluationsInOrder []Print
}

type Print struct {
	Name   string
	Colour string
}