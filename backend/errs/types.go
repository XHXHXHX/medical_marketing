package errs

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

type Err *Error

// NewError 生成一个完整的 Error
func NewError(code string, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

// NewSimpleError 生成 Error 时忽略 Details 字段
func NewSimpleError(code string, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func (e *Error) Reset() {}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) ProtoMessage() {}

func (e *Error) AsRaw() *Error {
	return (*Error)(e)
}

// Is 由 errors.Is(src, target) 调用，实现 Error 相等判断
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}
	if v, ok := target.(*Error); ok {
		return e.Code == v.Code
	}
	return false
}

// Wrap 方便函数, 等同于 errors.Wrap(err, msg)
func (e *Error) Wrap(msg string) error {
	return errors.Wrap(e, msg)
}

// Wrapf 方便函数, 等同于 errors.Wrapf(err, format, args...)
func (e *Error) Wrapf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// WithStack 方便函数,等同于 errors.WithStack(err)
func (e *Error) WithStack() error {
	return errors.WithStack(e)
}

func (e *Error) WithError(err error) error {
	return errors.Wrap(err, e.Message)
}

// Is 判断在 srcErr 的错误链中是否包含 targetErr
// 通过 Error.Code 来判断两个错误是否相同
func Is(srcErr, targetErr error) bool {
	return errors.Is(srcErr, targetErr)
}

// As 在错误链中找到第一个 *Error 类型的错误
func As(srcErr error) (*Error, bool) {
	if srcErr == nil {
		return nil, false
	}
	var err *Error
	if errors.As(srcErr, &err) {
		return NewError(err.Code, srcErr.Error()), true
	}
	return nil, false
}

// FromGRPCError 将 grpc 返回的 error 转换为 Error 类型
func FromGRPCError(err error) (*Error, bool) {
	sts, stsOk := status.FromError(err)
	if !stsOk || sts == nil {
		return nil, false
	}
	if details := sts.Details(); len(details) > 0 {
		if v, dtlOk := details[0].(*Error); dtlOk {
			return v, true
		}
	}
	return nil, false
}

