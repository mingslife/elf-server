package utils

func Offset(limit, page int) int {
	if page != -1 && limit != -1 {
		return (page - 1) * limit
	}
	return -1
}
