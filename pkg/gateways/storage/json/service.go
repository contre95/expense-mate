package json

type JSONStorage struct {
	Path string
}

func NewStorage(path string) *JSONStorage {
	return &JSONStorage{path}
}
