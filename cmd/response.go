package main

import (
	"sort"
	"strings"
)

type recipeDelivery struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
}

type recipeCountResponse struct {
	UniqueRecipeCount    int               `json:"unique_recipe_count"`
	CountPerRecipe       recipeCountList   `json:"count_per_recipe"`
	BusiestPostcode      postcodeCount     `json:"busiest_postcode"`
	CountPerPostcodeTime postcodeTimeCount `json:"count_per_postcode_and_time"`
	MatchByName          []string          `json:"match_by_name"`
}

type recipeCountList []recipeCount

func (l recipeCountList) filterByNames(names ...string) []string {
	list := make([]string, 0)

	for _, r := range l {
		recipeName := strings.ToLower(r.Recipe)
		for _, n := range names {
			searchName := strings.ToLower(n)
			if strings.Contains(recipeName, searchName) {
				list = append(list, r.Recipe)
				break
			}
		}
	}
	sort.Strings(list)
	return list
}

type recipeCount struct {
	Recipe        string `json:"recipe"`
	DeliveryCount int    `json:"count"`
}

type postcodeCount struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type postcodeTimeCount struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

func buildCountResponse(recipeCountSet recipeCountSet, postcodeCountSet postcodeCountSet, options recipeCountOptions) recipeCountResponse {
	sortedRecipeList := recipeCountSet.toSortedList()
	busiestPostcode := postcodeCountSet.findBusiestPostcode()

	return recipeCountResponse{
		UniqueRecipeCount: len(sortedRecipeList),
		CountPerRecipe:    sortedRecipeList,
		BusiestPostcode: postcodeCount{
			Postcode:      busiestPostcode,
			DeliveryCount: postcodeCountSet[busiestPostcode].deliveryCount,
		},
		CountPerPostcodeTime: postcodeTimeCount{
			Postcode:      options.postcode,
			From:          options.delivery.start.Format(timestampLayout),
			To:            options.delivery.end.Format(timestampLayout),
			DeliveryCount: postcodeCountSet[options.postcode].deliveryWithinTimeCount,
		},
		MatchByName: sortedRecipeList.filterByNames(options.recipes.names()...),
	}
}
