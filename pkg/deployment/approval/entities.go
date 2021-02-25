package approval

// ApproverError is an error related to approvers
type ApproverError struct {
	m string
}

func (e ApproverError) Error() string {
	return e.m
}
