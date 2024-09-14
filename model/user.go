package model

type Car struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Brand  string `json:"brand" gorm:"not null"`
	Model  string `json:"model" gorm:"not null"`
	Year   int    `json:"year" gorm:"not null"`
	UserId string `json:"userId" gorm:"not null"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

type User struct {
	ID     string  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string  `json:"name" gorm:"not null"`
	Age    int     `json:"age" gorm:"not null"`
	Cars   []Car   `json:"cars"`
	Groups []Group `json:"groups" gorm:"many2many:group_users"`
}

type Group struct {
	ID    string `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"not null"`
	Users []User `json:"users" gorm:"many2many:group_users"`
}
