package repo

import (
	"errors"
	"strings"
	"time"
)

type Status uint

const (
	Draft Status = iota
	Deleted
	Formed
	Finished
	Rejected
)

/*func (s *status) Scan(value any) error {
} */
/*func (s status) Value() (driver.Value, error) {
	return []string{"draft", "deleted", "formed", "finished", "rejected"}[s], nil
}*/

func convStrs() []string {
	return []string{"draft", "deleted", "formed", "finished", "rejected"}
}

func FromString(str string) (Status, error) {
	str = strings.ToLower(str)
	for i, v := range convStrs() {
		if v == str {
			return Status(i), nil
		}
	}
	return Status(0), errors.New("cant conv string to Status")
}

func (s Status) String() string {
	strings := convStrs()
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
	Avatar      string `gorm:"varchar(255)"`
}

type DataService struct {
	DataID   uint   `gorm:"primarykey"`
	DataName string `gorm:"type:varchar(30)"`
	Encode   bool   `gorm:"type:bool"`
	Blob     string `gorm:"type:text"`
	Active   bool   `gorm:"type:bool"`
}

type EncryptDecryptRequest struct {
	RequestID    uint `gorm:"primarykey"`
	Status       Status
	CreationDate time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
	FinishDate   time.Time
	FormDate     time.Time
	ModeratorID  uint
	CreatorID    uint
}

type EncryptDecryptToData struct {
	DataID    uint `gorm:"primarykey"`
	RequestID uint `gorm:"primarykey"`
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
