package middleware

func IsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func IsClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func IsServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}
