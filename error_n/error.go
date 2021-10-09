package error_n

type Error struct {
	ErrStr string
}

func (e Error) Error() string {
	return e.ErrStr
}

//生成一个新error
func NewError(errStr string) Error {
	return Error{ErrStr: errStr}
}
