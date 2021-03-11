package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeCountSet(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should add multiple recipes": func(t *testing.T) {
			// given
			set := make(recipeCountSet)

			// when
			set.add("Maple")
			set.add("Maple")
			set.add("Maple")
			set.add("Syrup")
			set.add("Syrup")
			set.add("Jam")

			// then
			assert.Equal(t, 3, len(set))
			assert.Equal(t, 3, set["Maple"])
			assert.Equal(t, 2, set["Syrup"])
			assert.Equal(t, 1, set["Jam"])
		},
		"should merge sets": func(t *testing.T) {
			// given
			set := make(recipeCountSet)
			set.add("Maple")
			set.add("Maple")
			set.add("Maple")
			set.add("Syrup")
			set.add("Syrup")
			set.add("Jam")
			setOther := make(recipeCountSet)
			setOther.add("Maple")
			setOther.add("Tangerine")
			setOther.add("Jam")
			setOther.add("Tangerine")
			setOther.add("Ham")

			// when
			set.merge(setOther)

			// then
			assert.Equal(t, 5, len(set))
			assert.Equal(t, 4, set["Maple"])
			assert.Equal(t, 2, set["Syrup"])
			assert.Equal(t, 2, set["Jam"])
			assert.Equal(t, 2, set["Tangerine"])
			assert.Equal(t, 1, set["Ham"])
		},
		"should return alphabetically sorted list": func(t *testing.T) {
			// given
			set := make(recipeCountSet)
			set.add("Maple")
			set.add("Maple")
			set.add("Maple")
			set.add("Syrup")
			set.add("Syrup")
			set.add("Jam")

			// when
			list := set.toSortedList()

			// then
			expected := make(recipeCountList, 0)
			expected = append(expected,
				recipeCount{Recipe: "Jam", DeliveryCount: 1},
				recipeCount{Recipe: "Maple", DeliveryCount: 3},
				recipeCount{Recipe: "Syrup", DeliveryCount: 2},
			)
			assert.Equal(t, 3, len(list))
			assert.Equal(t, expected, list)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestPostcodeCountSet(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should add multiple postcodes": func(t *testing.T) {
			// given
			set := make(postcodeCountSet)

			// when
			set.add("30000", false)
			set.add("30000", true)
			set.add("30000", true)
			set.add("20000", false)
			set.add("20000", true)
			set.add("10000", false)

			// then
			assert.Equal(t, 3, len(set))
			assert.Equal(t, &postcodeMatches{deliveryCount: 3, deliveryWithinTimeCount: 2}, set["30000"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 2, deliveryWithinTimeCount: 1}, set["20000"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 1, deliveryWithinTimeCount: 0}, set["10000"])
		},
		"should merge sets": func(t *testing.T) {
			// given
			set := make(postcodeCountSet)
			set.add("30000", false)
			set.add("30000", true)
			set.add("30000", true)
			set.add("20000", false)
			set.add("20000", true)
			set.add("10000", false)
			setOther := make(postcodeCountSet)
			setOther.add("30000", true)
			setOther.add("40000", true)
			setOther.add("10000", false)
			setOther.add("40000", false)
			setOther.add("50000", true)

			// when
			set.merge(setOther)

			// then
			assert.Equal(t, 5, len(set))
			assert.Equal(t, &postcodeMatches{deliveryCount: 1, deliveryWithinTimeCount: 1}, set["50000"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 2, deliveryWithinTimeCount: 1}, set["40000"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 4, deliveryWithinTimeCount: 3}, set["30000"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 2, deliveryWithinTimeCount: 1}, set["20000"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 2, deliveryWithinTimeCount: 0}, set["10000"])
		},
		"should return busiest postcode": func(t *testing.T) {
			// given
			set := make(postcodeCountSet)
			set.add("30000", false)
			set.add("30000", true)
			set.add("30000", true)
			set.add("20000", false)
			set.add("20000", true)
			set.add("10000", false)

			// when
			busiest := set.findBusiestPostcode()

			// then
			assert.Equal(t, "30000", busiest)
		},
		"should check if postcode exists": func(t *testing.T) {
			// given
			set := make(postcodeCountSet)
			set.add("10000", true)

			// then
			assert.True(t, set.exists("10000"))
			assert.False(t, set.exists("11111"))
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestCountRecipeDelivery(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should s": func(t *testing.T) {
			// given
			var recipeDeliveryInput []recipeDelivery
			recipeDeliveryInput = append(recipeDeliveryInput,
				recipeDelivery{
					Postcode: "10120",
					Recipe:   "Cherry Balsamic Pork Chops",
					Delivery: "Wednesday 10AM - 3PM",
				},
				recipeDelivery{
					Postcode: "10208",
					Recipe:   "Creamy Dill Chicken",
					Delivery: "Thursday 11AM - 2PM",
				},
				recipeDelivery{
					Postcode: "10120",
					Recipe:   "Cherry Balsamic Pork Chops",
					Delivery: "Thursday 9AM - 3PM",
				},
				recipeDelivery{
					Postcode: "10186",
					Recipe:   "Cherry Balsamic Pork Chops",
					Delivery: "Saturday 1AM - 8PM",
				},
				recipeDelivery{
					Postcode: "10120",
					Recipe:   "Hot Honey Barbecue Chicken Legs",
					Delivery: "Wednesday 10AM - 4PM",
				},
				recipeDelivery{
					Postcode: "10208",
					Recipe:   "Hot Honey Barbecue Chicken Legs",
					Delivery: "Wednesday 1AM - 12PM",
				})

			deliveryWindow, _ := parseDeliveryPeriod("10AM-3PM")
			options := recipeCountOptions{
				filePath: "path/to/file.json",
				postcode: "10120",
				delivery: deliveryWindow,
			}

			// when
			recipeCountSet, postcodeCountSet := countRecipeDelivery(recipeDeliveryInput, options)

			// then
			assert.Equal(t, 3, len(recipeCountSet))
			assert.Equal(t, 3, recipeCountSet["Cherry Balsamic Pork Chops"])
			assert.Equal(t, 1, recipeCountSet["Creamy Dill Chicken"])
			assert.Equal(t, 2, recipeCountSet["Hot Honey Barbecue Chicken Legs"])

			assert.Equal(t, 3, len(postcodeCountSet))
			assert.Equal(t, &postcodeMatches{deliveryCount: 3, deliveryWithinTimeCount: 1}, postcodeCountSet["10120"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 1, deliveryWithinTimeCount: 0}, postcodeCountSet["10186"])
			assert.Equal(t, &postcodeMatches{deliveryCount: 2, deliveryWithinTimeCount: 0}, postcodeCountSet["10208"])
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
