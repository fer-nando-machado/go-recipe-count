package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecipeSearchSet(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should add recipe": func(t *testing.T) {
			// given
			set := make(recipeSearchSet)

			// when
			set.add("Chocolate")

			// then
			assert.Equal(t, 1, len(set))
			assert.True(t, set["Chocolate"])
		},
		"should add recipes in bulk": func(t *testing.T) {
			// given
			set := make(recipeSearchSet)

			// when
			set.addBulk("Apple,Cake,Lemonade", ",")

			// then
			assert.Equal(t, 3, len(set))
			assert.True(t, set["Apple"])
			assert.True(t, set["Cake"])
			assert.True(t, set["Lemonade"])
		},
		"should list recipes as names array": func(t *testing.T) {
			// given
			set := make(recipeSearchSet)
			set.addBulk("Apple,Cake,Lemonade", ",")

			// when
			names := set.names()

			// then
			assert.Equal(t, 3, len(names))
			assert.ElementsMatch(t, [...]string{"Apple", "Cake", "Lemonade"}, names)
		},
		"should check if recipe exists": func(t *testing.T) {
			// given
			set := make(recipeSearchSet)
			set.add("Banana")

			// then
			assert.True(t, set.exists("Banana"))
			assert.False(t, set.exists("Orange"))
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestParseDeliveryPeriod(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should parse delivery period": func(t *testing.T) {
			// given
			timestamp := "12AM-12PM"

			// when
			deliveryPeriod, err := parseDeliveryPeriod(timestamp)

			// then
			assert.Equal(t, time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC), deliveryPeriod.start)
			assert.Equal(t, time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC), deliveryPeriod.end)
			assert.NoError(t, err)
		},
		"should parse delivery period (with spaces)": func(t *testing.T) {
			// given
			timestamp := "1AM - 1PM"

			// when
			deliveryPeriod, err := parseDeliveryPeriod(timestamp)

			// then
			assert.Equal(t, time.Date(0, 1, 1, 1, 0, 0, 0, time.UTC), deliveryPeriod.start)
			assert.Equal(t, time.Date(0, 1, 1, 13, 0, 0, 0, time.UTC), deliveryPeriod.end)
			assert.NoError(t, err)
		},
		"should parse delivery period (with weekday)": func(t *testing.T) {
			// given
			timestamp := "Tuesday 11AM - 11PM"

			// when
			deliveryPeriod, err := parseDeliveryPeriod(timestamp)

			// then
			assert.Equal(t, time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), deliveryPeriod.start)
			assert.Equal(t, time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC), deliveryPeriod.end)
			assert.NoError(t, err)
		},
		"should not parse delivery period (missing dash)": func(t *testing.T) {
			// given
			timestamp := "12AM12PM"

			// when
			_, err := parseDeliveryPeriod(timestamp)

			// then
			assert.Error(t, err)
		},
		"should not parse delivery period (missing AM/PM)": func(t *testing.T) {
			// given
			timestamp := "12MM-12MM"

			// when
			_, err := parseDeliveryPeriod(timestamp)

			// then
			assert.Error(t, err)
		},
		"should not parse delivery period (invalid numbers)": func(t *testing.T) {
			// given
			timestamp := "12AM-13PM"

			// when
			_, err := parseDeliveryPeriod(timestamp)

			// then
			assert.Error(t, err)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestIncludesDeliveryPeriod(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should include delivery period": func(t *testing.T) {
			// given
			deliveryWindow, _ := parseDeliveryPeriod("10AM-3PM")
			deliveryPeriod, _ := parseDeliveryPeriod("10AM-2PM")

			// when
			includes := deliveryWindow.includes(deliveryPeriod)

			// then
			assert.True(t, includes)
		},
		"should not include delivery period": func(t *testing.T) {
			// given
			deliveryWindow, _ := parseDeliveryPeriod("10AM-3PM")
			deliveryPeriod, _ := parseDeliveryPeriod("9AM-2PM")

			// when
			includes := deliveryWindow.includes(deliveryPeriod)

			// then
			assert.False(t, includes)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestParseCountOptions(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should parse count options": func(t *testing.T) {
			// given
			filePath := "path/to/file.json"
			postcode := "99999"
			deliveryTime := "12AM-12PM"
			recipeNames := "Potato,Pie"

			// when
			options, err := parseCountOptions(filePath, postcode, deliveryTime, recipeNames)

			// then
			assert.Equal(t, "path/to/file.json", options.filePath)
			assert.Equal(t, "99999", options.postcode)
			assert.Equal(t, time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC), options.delivery.start)
			assert.Equal(t, time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC), options.delivery.end)
			assert.Equal(t, 2, len(options.recipes))
			assert.NoError(t, err)
		},
		"should parse count options and fill in default fields": func(t *testing.T) {
			// given
			filePath := "path/to/file.json"

			// when
			options, err := parseCountOptions(filePath, "", "", "")

			// then
			assert.Equal(t, "path/to/file.json", options.filePath)
			assert.Equal(t, "10120", options.postcode)
			assert.Equal(t, time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), options.delivery.start)
			assert.Equal(t, time.Date(0, 1, 1, 15, 0, 0, 0, time.UTC), options.delivery.end)
			assert.Equal(t, 3, len(options.recipes))
			assert.NoError(t, err)
		},
		"should not parse count options when missing required fields": func(t *testing.T) {
			// given
			filePath := ""
			postcode := "99999"
			deliveryTime := "12AM-12PM"
			recipeNames := "Potato,Pie"

			// when
			_, err := parseCountOptions(filePath, postcode, deliveryTime, recipeNames)

			// then
			assert.Error(t, err)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
