package diyanet

const apiURLPrefix = "https://awqatsalah.diyanet.gov.tr/"
const errorPrefix = "diyanet: "

// Config holds the configuration parameters for the Diyanet Awqat Salah service.
type Config struct {
	// Email is the user's email address used for authentication.
	Email string

	// Password is the user's password used for authentication.
	Password string
}

// Result is a generic response envelope returned by Diyanet Awqat Salah APIs.
// It wraps the actual payload, a success indicator, and any server-provided message.
type Result[T any] struct {
	// Data contains the response payload when the request is successful.
	Data T `json:"data"`
	// Ok indicates whether the API call succeeded.
	Ok bool `json:"success"`
	// Error carries a message when the call fails; it is empty on success.
	Error string `json:"message"`
}
