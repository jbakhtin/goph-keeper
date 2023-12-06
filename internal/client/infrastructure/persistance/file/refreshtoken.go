package objectstorage

type refreshtoken interface {
	Read() (int, error)
	Write() (int, error)
	Create() (bool, error)
}
