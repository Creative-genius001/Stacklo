package utils

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

type ErrorJSON struct {
	Message    string                 `json:"message"`
	StatusCode int                    `json:"status_code"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
	StackTrace string                 `json:"stack_trace,omitempty"`
}

func NewError(status int, message string) *ErrorJSON {
	return &ErrorJSON{
		Message:    message,
		StatusCode: status,
		// StackTrace: captureStackTrace(3),
	}
}

func captureStackTrace(skip int) string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])

	trace := ""
	for {
		frame, more := frames.Next()
		trace += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
	return trace
}

func JSONErrorResp(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"status_code": code,
		"message":     message,
	})
}
