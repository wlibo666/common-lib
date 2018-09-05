package webutils

const (
	ERR_FIELD_POSITION = "position"
	ERR_FIELD_ERR      = "err"
)

const (
	ERRNO_SUCCESS = 10000 + iota
	ERRNO_INVALID_API_TYPE
	ERRNO_INVALID_SIGN
	ERRNO_LOST_REQ_PARAM
	ERRNO_LOST_FORM_PARAM
	ERRNO_LOST_URL_PARAM
	ERRNO_INVALID_REQ_PARAM
	ERRNO_INVALID_FORM_PARAM
	ERRNO_INVALID_URL_PARAM
	ERRNO_CALL_CORE_API_FAILED
	ERRNO_QUERY_FROM_DB
	ERRNO_STORE_DB
	ERRNO_NOT_FOUND
	ERRNO_INTERNAL_ERROR
	ERRNO_USER_NOT_EXIST
	ERRNO_INCORRECT_PWD
	ERRNO_USED
)

var (
	ErrInfo map[int]string = make(map[int]string)
)

func init() {
	ErrInfo[ERRNO_INVALID_API_TYPE] = "Invalid api type"
	ErrInfo[ERRNO_INVALID_SIGN] = "Invalid or expired signature"
	ErrInfo[ERRNO_LOST_REQ_PARAM] = "Lost request param"
	ErrInfo[ERRNO_LOST_FORM_PARAM] = "Lost form param"
	ErrInfo[ERRNO_LOST_URL_PARAM] = "Lost URL param"
	ErrInfo[ERRNO_INVALID_REQ_PARAM] = "Invalid request param"
	ErrInfo[ERRNO_INVALID_FORM_PARAM] = "Invalid form param"
	ErrInfo[ERRNO_INVALID_URL_PARAM] = "Invalid URL param"
	ErrInfo[ERRNO_CALL_CORE_API_FAILED] = "Call core api failed"
	ErrInfo[ERRNO_QUERY_FROM_DB] = "Query from DB failed"
	ErrInfo[ERRNO_STORE_DB] = "Store to DB failed"
	ErrInfo[ERRNO_NOT_FOUND] = "Not found"
	ErrInfo[ERRNO_INTERNAL_ERROR] = "Internal error"
	ErrInfo[ERRNO_USER_NOT_EXIST] = "User is not exist"
	ErrInfo[ERRNO_INCORRECT_PWD] = "Incorrect password"
}

func AddErrno(errno int, err string) int {
	if errno >= ERRNO_SUCCESS && errno <= ERRNO_USED {
		return ERRNO_USED
	}
	ErrInfo[errno] = err
	return ERRNO_SUCCESS
}

func GetErrno(errno int) string {
	return ErrInfo[errno]
}
