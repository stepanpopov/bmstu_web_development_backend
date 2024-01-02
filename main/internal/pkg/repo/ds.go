package repo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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
	fmt.Println(str)
	str = strings.ToLower(str)
	fmt.Println(str)
	for i, v := range convStrs() {
		fmt.Println(v, str)
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

type JWTClaims struct {
	jwt.StandardClaims
	UserUUID    uuid.UUID `json:"user_uuid"`
	IsModerator bool      `json:"scopes"`
}

type User struct {
	UserID      uuid.UUID `gorm:"primary_key"`
	Username    string    `gorm:"type:varchar(30)"`
	Password    string    `gorm:"type:varchar(30)"`
	IsModerator bool      `gorm:"type:boolean"`
}

type DataService struct {
	DataID    uint      `gorm:"primarykey" json:"data_id"`
	DataName  string    `gorm:"type:varchar(30)" json:"data_name"`
	Encode    bool      `gorm:"type:bool" json:"encode"`
	Blob      string    `gorm:"type:text" json:"blob"`
	Active    bool      `gorm:"type:bool" json:"active"`
	ImageUUID uuid.UUID `json:"image_uuid,omitempty" gorm:"type:uuid"`
}

type EncryptDecryptRequest struct {
	RequestID     uint `gorm:"primarykey"`
	Status        Status
	CreationDate  time.Time `gorm:"default:NOW()"`
	FinishDate    *time.Time
	FormDate      *time.Time
	ModeratorID   *uuid.UUID `gorm:"type:uuid"`
	CreatorID     *uuid.UUID `gorm:"type:uuid"`
	ResultCounter uint
}

type EncryptDecryptRequestView struct {
	RequestID    uint `gorm:"primarykey"`
	Status       Status
	CreationDate time.Time `gorm:"default:NOW()"`
	FinishDate   *time.Time
	FormDate     *time.Time
	Moderator    *string `gorm:"column:username"`
	Creator      *string `gorm:"column:username"`
}

type EncryptDecryptRequestViewWithCount struct {
	RequestID     uint `gorm:"primarykey"`
	Status        Status
	CreationDate  time.Time `gorm:"default:NOW()"`
	FinishDate    *time.Time
	FormDate      *time.Time
	Moderator     *string `gorm:"column:username"`
	Creator       *string `gorm:"column:username"`
	ResultCounter uint
}

type Calculated struct {
	ID      uint
	Result  string
	Success bool
}

type EncryptDecryptToData struct {
	DataID    uint `gorm:"primarykey"`
	RequestID uint `gorm:"primarykey"`
	Result    *string
	Success   *bool
}

/*type DataServiceResult struct {
	Result  string
	Success bool
}*/

type DataServiceWithOptResult struct {
	DataService
	Result  *string `gorm:"column:result"`
	Success *bool   `gorm:"column:success"`
}

func All() []any {
	return []any{
		&User{},
		&DataService{},
		&EncryptDecryptRequest{},
		&EncryptDecryptToData{},
	}
}
