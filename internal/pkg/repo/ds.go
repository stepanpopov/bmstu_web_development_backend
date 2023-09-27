package repo

const (
	Encode = iota
	Decode
)

type DataService struct {
	ID            uint `gorm:"primarykey"`
	Name          string
	Encode_Decode uint
	Blob          string
}
