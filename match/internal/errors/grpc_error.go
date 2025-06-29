package errors

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewGRPCError 建立自訂 gRPC 錯誤，code 可自訂
func NewGRPCError(code codes.Code, msg string) error {
	st := status.New(code, msg)
	detail := &errdetails.ErrorInfo{
		Reason:   code.String(), // 用 code 字串表示 Reason
		Domain:   "global",
		Metadata: map[string]string{"message": msg},
	}
	stWithDetails, err := st.WithDetails(detail)
	if err != nil {
		// 如果附加 detail 失敗，回傳原本 status.Err()
		return st.Err()
	}
	return stWithDetails.Err()
}
