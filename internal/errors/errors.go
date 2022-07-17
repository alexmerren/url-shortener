package errors

type baseError struct {
	err error
	msg string
}

func (b *baseError) Error() string {
	return b.msg
}

type internalServerError struct {
	baseError
}

func NewInternalServerError(err error) *internalServerError {
	return &internalServerError{
		baseError{
			err: err,
			msg: err.Error(),
		},
	}
}

type notFoundError struct {
	baseError
}

func NewNotFoundError(err error) *notFoundError {
	return &notFoundError{
		baseError{
			err: err,
			msg: err.Error(),
		},
	}
}

type badRequestError struct {
	baseError
}

func NewBadRequestError(err error) *badRequestError {
	return &badRequestError{
		baseError{
			err: err,
			msg: err.Error(),
		},
	}
}
