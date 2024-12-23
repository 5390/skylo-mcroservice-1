package models

// SIMRecord represents a single SIM record from the input file.
type SIMRecord struct {
	IMSI      string `json:"imsi"`
	PIN1      string `json:"pin1"`
	PUK1      string `json:"puk1"`
	PIN2      string `json:"pin2"`
	PUK2      string `json:"puk2"`
	AAM1      string `json:"aam1"`
	KIUMTSEnc string `json:"ki_umts_enc"`
	ACC       string `json:"acc"`
}

// FileData represents the structure of the parsed file.
type FileData struct {
	Header struct {
		Quantity int    `json:"quantity"`
		VarIn    string `json:"var_in_list"` // Starting IMSI value
	} `json:"header"`

	Records           []SIMRecord `json:"records"` // SIM records from the file
	ValidationResults ValidationResult
}

type ValidationResult struct {
	InvalidImsiRangeCount int
	DuplicateImsiInFile   int
	AlreadyExistingImsi   int
}
