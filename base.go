package diyanet

const apiUrlPrefix = "https://awqatsalah.diyanet.gov.tr/"
const errorPrefix = "diyanet: "

type Result[T any] struct {
	Data      T      `json:"data"`
	IsSuccess bool   `json:"success"`
	Error     string `json:"message"`
}
