package repo

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status uint

const (
	Draft Status = iota
	Deleted
	Formed
	Finished
	Rejected
	UnknownStatus
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
}

type DataService struct {
	DataID    uint      `gorm:"primarykey" json:"data_id"`
	DataName  string    `gorm:"type:varchar(30)" json:"data_name"`
	Encode    bool      `gorm:"type:bool" json:"encode"`
	Blob      string    `gorm:"type:text" json:"blob"`
	Active    bool      `gorm:"type:bool" json:"active"`
	ImageUUID uuid.UUID `json:"image_uuid,omitempty"`
}

type DataServiceView struct {
	DataID   uint   `json:"data_id"`
	DataName string `json:"data_name"`
	Encode   bool   `json:"encode"`
	Blob     string `json:"blob"`
	Active   bool   `json:"active"`
	ImageURL string `json:"image_url,omitempty"`
}

type EncryptDecryptRequest struct {
	RequestID    uint `gorm:"primarykey"`
	Status       Status
	CreationDate time.Time `gorm:"default:NOW()"`
	FinishDate   *time.Time
	FormDate     *time.Time
	ModeratorID  *uint
	CreatorID    *uint
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
