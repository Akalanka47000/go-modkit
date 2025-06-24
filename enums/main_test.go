package enums

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringEnum(t *testing.T) {
	t.Parallel()

	Convey("should pick up values automatically", t, func() {
		type MealPreferenceType string

		type MealPreferences struct {
			Veg    MealPreferenceType
			NonVeg MealPreferenceType
		}

		MealPreference := New[MealPreferenceType](MealPreferences{})

		So(MealPreference.ListValues(), ShouldHaveLength, 2)

		So(MealPreference.Values.Veg, ShouldEqual, MealPreferenceType("Veg"))
		So(MealPreference.Values.NonVeg, ShouldEqual, MealPreferenceType("NonVeg"))

		Convey("should infer lowercase values when used with the lowercase option", func() {
			type MealPreferenceType string

			type MealPreferences struct {
				Veg    MealPreferenceType
				NonVeg MealPreferenceType
			}

			MealPreference := New[MealPreferenceType](MealPreferences{}, true)

			So(MealPreference.ListValues(), ShouldHaveLength, 2)

			So(MealPreference.Values.Veg, ShouldEqual, MealPreferenceType("veg"))
			So(MealPreference.Values.NonVeg, ShouldEqual, MealPreferenceType("nonveg"))
		})

		Convey("should check if a value is valid", func() {
			So(MealPreference.IsValid(MealPreference.Values.Veg), ShouldBeTrue)
			So(MealPreference.IsValid(MealPreference.Values.NonVeg), ShouldBeTrue)
			So(MealPreference.IsValid("Vegan"), ShouldBeFalse)
		})

		Convey("should check for equality and inequality", func() {
			So(MealPreference.Values.Veg == MealPreference.Values.NonVeg, ShouldBeFalse)
			So(MealPreference.Values.NonVeg != MealPreference.Values.Veg, ShouldBeTrue)
			So(MealPreference.Values.Veg == "Veg", ShouldBeTrue)
			So(MealPreference.Values.NonVeg == MealPreferenceType("NonVeg"), ShouldBeTrue)
			So(MealPreference.Values.Veg == "123", ShouldBeFalse)
		})
	})

	Convey("should handle custom values", t, func() {
		type MealPreferenceType string

		type MealPreferences struct {
			Veg    MealPreferenceType `json:"veg"`
			NonVeg MealPreferenceType `json:"non_veg"`
		}

		MealPreference := New[MealPreferenceType](MealPreferences{
			Veg:    "Vegetarian",
			NonVeg: "Non-Vegetarian",
		})

		So(MealPreference.ListValues(), ShouldHaveLength, 2)

		So(MealPreference.Values.Veg, ShouldEqual, MealPreferenceType("Vegetarian"))
		So(MealPreference.Values.NonVeg, ShouldEqual, MealPreferenceType("Non-Vegetarian"))
	})
}

func TestIntEnum(t *testing.T) {
	t.Parallel()

	Convey("should fail to pick up values automatically", t, func() {
		type StatusType int

		type Statuses struct {
			Pending StatusType
			Active  StatusType
			Closed  StatusType
		}

		Status := New[StatusType](Statuses{})

		So(Status.ListValues(), ShouldHaveLength, 3)

		So(Status.Values.Pending, ShouldEqual, 0)
		So(Status.Values.Active, ShouldEqual, 0)
		So(Status.Values.Closed, ShouldEqual, 0)
	})

	Convey("should handle custom values", t, func() {
		type StatusType int

		type Statuses struct {
			Pending StatusType `json:"pending"`
			Active  StatusType `json:"active"`
			Closed  StatusType `json:"closed"`
		}

		Status := New[StatusType](Statuses{
			Pending: 1,
			Active:  2,
			Closed:  3,
		})

		So(Status.ListValues(), ShouldHaveLength, 3)

		So(Status.Values.Pending, ShouldEqual, 1)
		So(Status.Values.Active, ShouldEqual, 2)
		So(Status.Values.Closed, ShouldEqual, 3)
	})
}
