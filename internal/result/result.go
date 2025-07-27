package result

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func SuccessResult(data any) *Result {
	return NewResult(0, "success", data)
}
func ListResult(data any, page, pageSize, total int64) *Result {
	return NewResult(0, "success", map[string]any{"list": data, "page": page, "pageSize": pageSize, "total": total})
}
func SimpleResult(msg string) *Result {
	return NewResult(0, msg, nil)
}
func DataResult(msg string, data any) *Result {
	return NewResult(0, msg, data)
}
func NewResult(code int, msg string, data any) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func (r Result) AddError(err error) *Result {
	if err != nil {
		r.Msg += "：" + err.Error()
	}
	return &r
}
