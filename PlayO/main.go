package main

import (
	"fmt"
	// "time"
	// "net/http"
)


type User struct {
	UserID string
	UserName string
	Email string
	Bookings []int
}

type Owner struct {
	OwnerID string
	User
	Properties []Property
}

type Location struct{
	Latitude int 
	Longitude int
}

type Property struct {
	GameType string
	PropertyID string
	Location
	Area float64
}
type Slot struct{
	SlotID int
	Property
	IsAvailable bool
	BookingID int

}

type UserManager struct{
	Users []User
}
func (um *UserManager) AddUser(u *User) User{
	User.UserID = uuid.New().string()
	return User
}
func(um *UserManager) Removeuser(u *User) int

type SlotManager struct{
	Slots []Slot
}

func (sm *SlotManager) GetSlot(p *Property) Slot
func (sm *SlotManager) BookSLot(s *Slot) int

type PropertyManager struct{
	Properties []Property
}


func (pm *PropertyManager)AddProperty(p *Property) int
func (pm *PropertyManager)RemoveProperty(p *Property) int
func (pm *PropertyManager)GetPropertyByLocation(l *Location) []Property

func main() {
	fmt.Println(" ")

	user1 := User{
		UserID: "naman",
		UserName: "nnnn",
		Email: "nmn@email",
		Bookings: []int{1,2,3},
		}

	fmt.Println(user1);


}
