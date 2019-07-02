package jobservice

//Job contains the common set of fields needed for my job search
type Job struct {
	Company                 string
	Title                   string
	Description             string
	Qualifications          string
	PreferredQualifications string
	URL                     string
}
