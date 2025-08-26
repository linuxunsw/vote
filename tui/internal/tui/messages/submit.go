package messages

// Sends data to root model to trigger submit
type SubmitFormMsg struct {
	Name    string
	Email   string
	Discord string

	Roles     []string
	Statement string
	Url       string
}

// Tells the UI result of submission - if an error is
// present the submission was unsuccessful
/* type FormSubmissionResultMsg struct {
	err error
}*/
