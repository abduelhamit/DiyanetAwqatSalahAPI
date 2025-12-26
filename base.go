package diyanet

const apiURLPrefix = "https://awqatsalah.diyanet.gov.tr/"
const errorPrefix = "diyanet: "

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
