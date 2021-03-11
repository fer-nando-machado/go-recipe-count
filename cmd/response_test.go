package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipesCountList(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should filter by recipes names (case-insensitive)": func(t *testing.T) {
			// given
			recipeCountList := make(recipeCountList, 0)
			recipeCountList = append(recipeCountList,
				recipeCount{Recipe: "Starfish and coffee"},
				recipeCount{Recipe: "Maple syrup and jam"},
				recipeCount{Recipe: "Butterscotch clouds"},
				recipeCount{Recipe: "Tangerine"},
				recipeCount{Recipe: "Side order of ham"},
			)

			// when
			recipes := recipeCountList.filterByNames("Coffee", "tangerine")

			// then
			assert.ElementsMatch(t, [...]string{"Starfish and coffee", "Tangerine"}, recipes)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

//buildResponse(recipeCountSet recipeCountSet, postcodeCountSet postcodeCountSet, options recipeCountOptions)
func TestBuildResponse(t *testing.T) {
	tests := map[string]func(*testing.T){
		"should build a response": func(t *testing.T) {
			// given
			recipeSet := make(recipeCountSet)
			recipeSet.add("Starfish and coffee")
			recipeSet.add("Starfish and coffee")
			recipeSet.add("Starfish and coffee")
			recipeSet.add("Maple syrup and jam")
			recipeSet.add("Maple syrup and jam")
			recipeSet.add("Butterscotch clouds")

			postcodeSet := make(postcodeCountSet)
			postcodeSet.add("30000", false)
			postcodeSet.add("30000", false)
			postcodeSet.add("30000", false)
			postcodeSet.add("20000", false)
			postcodeSet.add("20000", true)
			postcodeSet.add("10000", false)

			deliveryWindow, _ := parseDeliveryPeriod("10AM-3PM")
			recipeSearch := make(recipeSearchSet)
			recipeSearch.addBulk("Coffee,jam", ",")
			options := recipeCountOptions{
				postcode: "20000",
				delivery: deliveryWindow,
				recipes:  recipeSearch,
			}

			// when
			response := buildCountResponse(recipeSet, postcodeSet, options)

			// then
			expectedCountPerRecipe := make(recipeCountList, 0)
			expectedCountPerRecipe = append(expectedCountPerRecipe,
				recipeCount{Recipe: "Butterscotch clouds", DeliveryCount: 1},
				recipeCount{Recipe: "Maple syrup and jam", DeliveryCount: 2},
				recipeCount{Recipe: "Starfish and coffee", DeliveryCount: 3},
			)
			expectedMatchByName := make([]string, 0)
			expectedMatchByName = append(expectedMatchByName,
				"Maple syrup and jam", "Starfish and coffee",
			)
			expected := recipeCountResponse{
				UniqueRecipeCount: 3,
				CountPerRecipe:    expectedCountPerRecipe,
				BusiestPostcode: postcodeCount{
					Postcode:      "30000",
					DeliveryCount: 3,
				},
				CountPerPostcodeTime: postcodeTimeCount{
					Postcode:      "20000",
					From:          options.delivery.start.Format(timestampLayout),
					To:            options.delivery.end.Format(timestampLayout),
					DeliveryCount: 1,
				},
				MatchByName: expectedMatchByName,
			}

			assert.Equal(t, expected, response)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
