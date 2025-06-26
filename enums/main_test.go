package enums

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringEnum(t *testing.T) {
	t.Parallel()

	Convey("should pick up values automatically", t, func() {
		type MealPreferences struct {
			String
			Veg    string
			NonVeg string
		}

		MealPreference := New(MealPreferences{})

		So(MealPreference.Values(), ShouldHaveLength, 2)

		So(MealPreference.Veg, ShouldEqual, "Veg")
		So(MealPreference.NonVeg, ShouldEqual, "NonVeg")

		Convey("should infer lowercase values when used with the lowercase option", func() {
			type MealPreferences struct {
				String
				Veg    string
				NonVeg string
			}

			MealPreference := New(MealPreferences{}, Lowercase())

			So(MealPreference.Values(), ShouldHaveLength, 2)

			So(MealPreference.Veg, ShouldEqual, "veg")
			So(MealPreference.NonVeg, ShouldEqual, "nonveg")
		})

		Convey("should infer lowercase values when used with the uppercase option", func() {
			type MealPreferences struct {
				String
				Veg    string
				NonVeg string
			}

			MealPreference := New(MealPreferences{}, Uppercase())

			So(MealPreference.Values(), ShouldHaveLength, 2)

			So(MealPreference.Veg, ShouldEqual, "VEG")
			So(MealPreference.NonVeg, ShouldEqual, "NONVEG")
		})

		Convey("should infer values when used with a pure struct definition", func() {
			MealPreference := New(struct {
				String
				Veg    string
				NonVeg string
			}{}, Lowercase())

			So(MealPreference.Values(), ShouldHaveLength, 2)

			So(MealPreference.Veg, ShouldEqual, "veg")
			So(MealPreference.NonVeg, ShouldEqual, "nonveg")
		})

		Convey("should check if a value is valid", func() {
			So(MealPreference.IsValid(MealPreference.Veg), ShouldBeTrue)
			So(MealPreference.IsValid(MealPreference.NonVeg), ShouldBeTrue)
			So(MealPreference.IsValid("Vegan"), ShouldBeFalse)
		})

		Convey("should validate a value", func() {
			So(MealPreference.Validate(MealPreference.Veg), ShouldBeNil)
			So(MealPreference.Validate(MealPreference.NonVeg), ShouldBeNil)
			So(MealPreference.Validate("Vegan"), ShouldEqual, fmt.Errorf("invalid value for type MealPreference: Vegan. Valid values include: [Veg NonVeg]"))
		})
	})

	Convey("should handle custom values", t, func() {
		type MealPreferences struct {
			String
			Veg    string
			NonVeg string
		}

		MealPreference := New(MealPreferences{
			Veg:    "Vegetarian",
			NonVeg: "Non-Vegetarian",
		})

		So(MealPreference.Values(), ShouldHaveLength, 2)

		So(MealPreference.Veg, ShouldEqual, "Vegetarian")
		So(MealPreference.NonVeg, ShouldEqual, "Non-Vegetarian")
	})
}

func TestIntEnum(t *testing.T) {
	t.Parallel()

	Convey("should fail to pick up values automatically", t, func() {
		type Statuses struct {
			Int
			Pending int
			Active  int
			Closed  int
		}

		Status := New(Statuses{})

		So(Status.Values(), ShouldHaveLength, 3)

		So(Status.Pending, ShouldEqual, 0)
		So(Status.Active, ShouldEqual, 0)
		So(Status.Closed, ShouldEqual, 0)
	})

	Convey("should handle custom values", t, func() {
		type Statuses struct {
			Int
			Pending int
			Active  int
			Closed  int
		}

		Status := New(Statuses{
			Pending: 1,
			Active:  2,
			Closed:  3,
		})

		So(Status.Values(), ShouldHaveLength, 3)

		So(Status.Pending, ShouldEqual, 1)
		So(Status.Active, ShouldEqual, 2)
		So(Status.Closed, ShouldEqual, 3)
	})
}
