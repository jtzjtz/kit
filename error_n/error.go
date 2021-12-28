package error_n

type ErrorN struct {
	ErrStr string
}

func (e ErrorN) Error() string {
	return e.ErrStr
}

//生成一个新error
func Error(errStr string) ErrorN {
	return ErrorN{ErrStr: errStr}
}
