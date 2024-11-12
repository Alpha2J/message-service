package constant

type BusinessCode struct {
	Code    int
	Message string
}

var (
	INVALID_REQUEST   = BusinessCode{-5, "invalid request"}
	SYSTEM_BUSY       = BusinessCode{-4, "system busy"}
	SESSION_EXPIRED   = BusinessCode{-3, "session expired"}
	UNKNOWN_EXCEPTION = BusinessCode{-2, "unknown exception"}
	COMMON_BIZ_ERROR  = BusinessCode{-1, "common biz error"}
	SUCCESS           = BusinessCode{0, "success"}
)
