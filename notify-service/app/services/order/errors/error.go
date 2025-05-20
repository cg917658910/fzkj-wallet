package errors

import (
	"bytes"
	"fmt"
	"strings"
)

// ErrorType 代表错误类型。
type ErrorType string

// 错误类型常量。
const (
	// ERROR_TYPE_SCHEDULER 代表调度器错误。
	ERROR_TYPE_SCHEDULER ErrorType = "scheduler error"
)

// NotifyError 代表错误的接口类型。
type NotifyError interface {
	// Type 用于获得错误的类型。
	Type() ErrorType
	// Error 用于获得错误提示信息。
	Error() string
}

// myNotifyError 代表错误的实现类型。
type myNotifyError struct {
	// errType 代表错误的类型。
	errType ErrorType
	// errMsg 代表错误的提示信息。
	errMsg string
	// fullErrMsg 代表完整的错误提示信息。
	fullErrMsg string
}

// NewNotifyError 用于创建一个新的错误值。
func NewNotifyError(errType ErrorType, errMsg string) NotifyError {
	return &myNotifyError{
		errType: errType,
		errMsg:  strings.TrimSpace(errMsg),
	}
}

// NewNotifyErrorBy 用于根据给定的错误值创建一个新的错误值。
func NewNotifyErrorBy(errType ErrorType, err error) NotifyError {
	return NewNotifyError(errType, err.Error())
}

func (ce *myNotifyError) Type() ErrorType {
	return ce.errType
}

func (ce *myNotifyError) Error() string {
	if ce.fullErrMsg == "" {
		ce.genFullErrMsg()
	}
	return ce.fullErrMsg
}

// genFullErrMsg 用于生成错误提示信息，并给相应的字段赋值。
func (ce *myNotifyError) genFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("cg notify error: ")
	if ce.errType != "" {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s", buffer.String())
	return
}

// IllegalParameterError 代表非法的参数的错误类型。
type IllegalParameterError struct {
	msg string
}

// NewIllegalParameterError 会创建一个IllegalParameterError类型的实例。
func NewIllegalParameterError(errMsg string) IllegalParameterError {
	return IllegalParameterError{
		msg: fmt.Sprintf("illegal parameter: %s",
			strings.TrimSpace(errMsg)),
	}
}

func (ipe IllegalParameterError) Error() string {
	return ipe.msg
}
