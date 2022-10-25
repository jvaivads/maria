package service

type ApiError struct {
	Message string
	Status  int64
	Cause   Cause
}

type Cause struct {
	error error
	*Cause
}

func (c *Cause) String() string {
	if c == nil {
		return "[]"
	}
	return c.error.Error() + "; Cause:{:" + c.Cause.String() + "}"
}
