package statuspage

import "net/http"

// StatusExplanation returns a short, human-readable description of the given HTTP
// status code.
func StatusExplanation(statusCode int) string {
	switch statusCode {
	// 4xx
	case http.StatusBadRequest:
		return "Your browser has sent a malformed request."
	case http.StatusUnauthorized:
		return "You must be authenticated to use this service."
	case http.StatusForbidden:
		return "You do not have access to this service."
	case http.StatusNotFound:
		return "The page you've requested could not be found."
	case http.StatusNotAcceptable:
		return "The content of this page is not accepted by your browser."
	case http.StatusProxyAuthRequired:
		return "You must be authenticated with the proxy server to use this service."
	case http.StatusRequestTimeout:
		return "Your browser did not send a request in a timely manner."
	case http.StatusRequestEntityTooLarge:
		return "Your browser has sent a request that's too large to process."
	case http.StatusRequestURITooLong:
		return "Your browser has sent a request with a URI that's too large to process."
	case http.StatusUpgradeRequired:
		return "Maybe you're trying to access a WebSocket server?"
	case http.StatusTooManyRequests:
		return "Your request has been rate-limited, please descrease the number of requests."
	case http.StatusRequestHeaderFieldsTooLarge:
		return "Your browser has sent a request header that is too large to process."
	case http.StatusUnavailableForLegalReasons:
		return "Your request has been denied for legal reasons."

	// 5xx
	case http.StatusNotImplemented:
		return "The feature you've requested is not supported."
	case http.StatusBadGateway:
		return "The service you've requested could not be contacted, please try again."
	case http.StatusServiceUnavailable:
		return "The service you've requested is temporarily unavailable, please try again."
	case http.StatusGatewayTimeout:
		return "The service you've requested did not respond in a timely manner, please try again."
	case http.StatusHTTPVersionNotSupported:
		return "Your browser's HTTP version is not supported."
	}

	if 400 <= statusCode && statusCode <= 599 {
		return "We're sorry, something went wrong!"
	}

	return "That's all we know."
}
