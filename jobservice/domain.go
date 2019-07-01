package jobservice

// Job represents a generic job to be returned by a service
type Job struct {
	Title                   string
	Description             string
	Qualifications          string
	PreferredQualifications string
	URL                     string
}
