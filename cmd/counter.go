package main

import (
	"sort"
)

type recipeCountSet map[string]int

func (s recipeCountSet) add(recipe string) {
	s[recipe]++
}

func (s recipeCountSet) merge(o recipeCountSet) {
	for k, v := range o {
		s[k] += v
	}
}

func (s recipeCountSet) toSortedList() recipeCountList {
	list := make(recipeCountList, 0)

	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		list = append(list, recipeCount{
			Recipe:        k,
			DeliveryCount: s[k],
		})
	}

	return list
}

type postcodeCountSet map[string]*postcodeMatches
type postcodeMatches struct {
	deliveryCount           int
	deliveryWithinTimeCount int
}

func (s postcodeCountSet) add(postcode string, isWithinTime bool) {
	if !s.exists(postcode) {
		s[postcode] = &postcodeMatches{}
	}

	s[postcode].deliveryCount++
	if isWithinTime {
		s[postcode].deliveryWithinTimeCount++
	}
}

func (s postcodeCountSet) merge(o postcodeCountSet) {
	for k, v := range o {
		if !s.exists(k) {
			s[k] = &postcodeMatches{}
		}

		s[k].deliveryCount += v.deliveryCount
		s[k].deliveryWithinTimeCount += v.deliveryWithinTimeCount
	}
}

func (s postcodeCountSet) findBusiestPostcode() string {
	maxKey := ""
	maxVal := 0

	for postcode, matches := range s {
		if matches.deliveryCount > maxVal {
			maxKey = postcode
			maxVal = matches.deliveryCount
		}
	}

	return maxKey
}

func (s postcodeCountSet) exists(postcode string) bool {
	return s[postcode] != nil
}

func countRecipeDelivery(recipeDeliveryInput []recipeDelivery, options recipeCountOptions) (recipeCountSet, postcodeCountSet) {
	recipeCountSet := make(recipeCountSet, 0)
	postcodeCountSet := make(postcodeCountSet, 0)

	for _, r := range recipeDeliveryInput {
		recipeCountSet.add(r.Recipe)

		matchesSearchPostcode := r.Postcode == options.postcode
		if !matchesSearchPostcode {
			postcodeCountSet.add(r.Postcode, false)
			continue
		}
		deliveryPeriod, _ := parseDeliveryPeriod(r.Delivery)
		isWithinTime := options.delivery.includes(deliveryPeriod)
		postcodeCountSet.add(r.Postcode, isWithinTime)
	}

	return recipeCountSet, postcodeCountSet
}
