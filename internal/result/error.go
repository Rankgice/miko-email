package result

var (
	ErrorReqParam     = ErrorResult(10000, "请求参数错误")
	ErrorBindingParam = ErrorResult(10001, "绑定参数错误")
	ErrorAdd          = ErrorResult(10100, "添加失败")
	ErrorUpdate       = ErrorResult(10101, "更新失败")
	ErrorDelete       = ErrorResult(10102, "删除失败")
	ErrorSelect       = ErrorResult(10103, "查询失败")
	ErrorCopy         = ErrorResult(10106, "复制失败")
	ErrorNotFound     = ErrorResult(10110, "未查询到数据")
	ErrorTxCommit     = ErrorResult(10113, "事务提交失败")
	ErrorDataVerify   = ErrorResult(10114, "数据校验失败")
)

func ErrorResult(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
	}
}

func ErrorSimpleResult(msg string) *Result {
	return &Result{
		Code: 20000,
		Msg:  msg,
	}
}
