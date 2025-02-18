package models

type File struct {
	Name        string
	Data        []byte
	Bucket      string
	ContentType string
}

type MinioErr struct {
	Error    error
	Bucket   string
	FileName string
	ObjectId string
}

type MinioResponse struct {
	Url      string
	ObjectId string
}
