package repo

import (
	"time"
)

type status uint

const (
	draft status = iota
	deleted
	formed
	finished
	rejected
)

/*func (s *status) Scan(value any) error {
} */
/*func (s status) Value() (driver.Value, error) {
	return []string{"draft", "deleted", "formed", "finished", "rejected"}[s], nil
}*/

func (s status) String() string {
	strings := []string{"draft", "deleted", "formed", "finished", "rejected"}
	if int(s) >= len(strings) {
		return "unknown"
	}
	return strings[s]
}

type User struct {
	UserID      uint   `gorm:"primary_key"`
	Username    string `gorm:"type:varchar(30)"`
	Password    string `gorm:"type:varchar(30)"`
	IsModerator bool   `gorm:"type:bool"`
}

type DataService struct {
	DataID   uint   `gorm:"primarykey"`
	DataName string `gorm:"type:varchar(30)"`
	Encode   bool   `gorm:"type:bool"`
	Blob     string `gorm:"type:text"`
	Active   bool   `gorm:"type:bool"`
}

type EncryptDecryptRequest struct {
	RequestID    uint
	Status       status
	CreationDate time.Time
	FinishDate   time.Time
	FormDate     time.Time
	ModeratorID  uint
	CreatorID    uint
}

type EncryptDecryptToData struct {
	DataID    uint
	RequestID uint
	Result    string
}

func All() []any {
	return []any{
		&User{},
		&DataService{},
		&EncryptDecryptRequest{},
		&EncryptDecryptToData{},
	}
}
