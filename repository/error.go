package repository

type DataNotFound struct{}

func (e *DataNotFound) Error() string {
	return "data not found"
}
