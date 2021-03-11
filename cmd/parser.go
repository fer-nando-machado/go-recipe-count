package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type recipeCountOptions struct {
	filePath string
	postcode string
	delivery deliveryPeriod
	recipes  recipeSearchSet
}

type deliveryPeriod struct {
	start time.Time
	end   time.Time
}

func (p deliveryPeriod) includes(o deliveryPeriod) bool {
	return !o.start.Before(p.start) && !o.end.After(p.end)
}

type recipeSearchSet map[string]bool

func (s recipeSearchSet) add(recipe string) {
	s[recipe] = true
}

func (s recipeSearchSet) addBulk(recipes string, separator string) {
	for _, r := range strings.Split(recipes, separator) {
		s.add(r)
	}
}

func (s recipeSearchSet) names() []string {
	names := make([]string, 0, len(s))
	for k := range s {
		names = append(names, k)
	}
	return names
}

func (s recipeSearchSet) exists(recipe string) bool {
	return s[recipe]
}

func parseDeliveryPeriod(deliveryTime string) (deliveryPeriod, error) {
	rexp := regexp.MustCompile(`(1[012]|[1-9])(\\s)?(AM|PM)-(1[012]|[1-9])(\\s)?(AM|PM)`)
	timestamp := rexp.FindString(strings.ReplaceAll(deliveryTime, " ", ""))
	if timestamp == "" {
		return deliveryPeriod{}, errors.New("badly formatted delivery time string")
	}

	deliveryTimes := strings.Split(timestamp, "-")
	deliveryStart, _ := time.Parse(timestampLayout, deliveryTimes[0])
	deliveryEnd, _ := time.Parse(timestampLayout, deliveryTimes[1])
	return deliveryPeriod{
		start: deliveryStart,
		end:   deliveryEnd,
	}, nil
}

func parseCountOptions(filePath string, postcode string, deliveryTime string, recipeNames string) (recipeCountOptions, error) {
	if len(filePath) == 0 {
		return recipeCountOptions{}, errors.New("file is a required argument")
	}
	if len(postcode) == 0 {
		postcode = postcodeDefault
	}
	if len(deliveryTime) == 0 {
		deliveryTime = deliveryTimeDefault
	}
	if len(recipeNames) == 0 {
		recipeNames = recipeNamesDefault
	}

	deliveryPeriod, _ := parseDeliveryPeriod(deliveryTime)

	recipeSearchSet := make(recipeSearchSet)
	recipeSearchSet.addBulk(recipeNames, ",")

	return recipeCountOptions{
		filePath,
		postcode,
		deliveryPeriod,
		recipeSearchSet,
	}, nil
}
