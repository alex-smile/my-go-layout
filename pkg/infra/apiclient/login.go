package apiclient

import "net/http"

func GetUser(username string) error {

	return NewAPIError(http.StatusInternalServerError, "get user error")
}
