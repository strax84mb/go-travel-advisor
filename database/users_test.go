package database

/*
import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShow(t *testing.T) {
	Convey("Do some test", t, func() {
		x := 1
		Convey("When the integer is incremented", func() {
			x++
			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
		Convey("Should be 1", func() {
			So(x, ShouldEqual, 1)
		})
		Convey("When value is decreased", func() {
			x--
			Convey("Value should be zero", func() {
				So(x, ShouldEqual, 0)
			})
		})
	})
}

func setupTests(db *dbWrapperMock) {
	db.setZeros()
	gdb = db
}

func TestSaveUser(t *testing.T) {
	Convey("Having called a saveUser", t, func() {
		Convey("While username is taken", func() {
			user := User{
				Username: "admin",
				Password: "some_other_password",
				Role:     UserRoleUser,
			}
			setupTests(&dbWrapperMock{
				Value:         []interface{}{&user},
				NoRecordFound: []bool{false},
				Errors:        []error{nil},
			})
			err := saveUser(user)
			Convey("UsernameTakenError is returned", func() {
				So(err, ShouldHaveSameTypeAs, &UsernameTakenError{})
			})
		})
		Convey("While username is not taken", func() {
			user := User{
				Username: "strahinjamalabosna",
				Password: "some_password",
				Role:     UserRoleUser,
			}
			Convey("While saving is successful", func() {
				setupTests(&dbWrapperMock{
					Value:         []interface{}{&User{}},
					Errors:        []error{nil},
					NoRecordFound: []bool{true},
				})
				err := saveUser(user)
				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
			})
			Convey("While saving results in error", func() {
				user.Username = "settofail"
				setupTests(&dbWrapperMock{
					Value:         []interface{}{&User{}},
					Errors:        []error{errors.New("insert error")},
					NoRecordFound: []bool{true},
				})
				err := saveUser(user)
				Convey("StatementError is returned", func() {
					So(err, ShouldHaveSameTypeAs, &StatementError{})
					So("Error while saving user", ShouldEqual, err.Error())
				})
			})
		})
	})
}
*/
