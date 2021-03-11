package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

const postcodeDefault string = "10120"
const deliveryTimeDefault string = "10AM-3PM"
const recipeNamesDefault string = "Potato,Veggie,Mushroom"
const timestampLayout string = "3PM"

func main() {
	// parses input option flags
	filePath := flag.String("file", "", "fixtures data file path (required)")
	postcode := flag.String("postcode", postcodeDefault, "postcode to search for")
	deliveryTime := flag.String("time", deliveryTimeDefault, "delivery time to search for")
	recipeNames := flag.String("recipes", recipeNamesDefault, "recipe(s) name(s) to search for, separated by commas")
	flag.Parse()
	options, err := parseCountOptions(*filePath, *postcode, *deliveryTime, *recipeNames)
	if err != nil {
		log.Fatal(err)
	}

	// reads input file content
	fileContent, err := ioutil.ReadFile(options.filePath)
	if err != nil {
		log.Fatal(err)
	}
	var recipeDeliveryInput []recipeDelivery
	err = json.Unmarshal([]byte(fileContent), &recipeDeliveryInput)
	if err != nil {
		log.Fatal(err)
	}

	// slices input for parallel processing
	numCPU := runtime.NumCPU()
	blocksize := len(recipeDeliveryInput) / numCPU
	c := make(chan partialCountSets)
	for i := 0; i < numCPU; i++ {
		start, end := i*blocksize, (i+1)*blocksize
		go partialCountRecipeDelivery(recipeDeliveryInput[start:end], options, c)
	}

	// merges partial counts into totals set
	recipeCountTotal := make(recipeCountSet, 0)
	postcodeCountTotal := make(postcodeCountSet, 0)
	partialCount := make([]partialCountSets, numCPU)
	for i := range partialCount {
		partialCount[i] = <-c
		recipeCountTotal.merge(partialCount[i].recipeCountSet)
		postcodeCountTotal.merge(partialCount[i].postcodeCountSet)
	}

	// outputs JSON response to stdout
	response := buildCountResponse(recipeCountTotal, postcodeCountTotal, options)
	printer := json.NewEncoder(os.Stdout)
	printer.Encode(response)
}

type partialCountSets struct {
	recipeCountSet   recipeCountSet
	postcodeCountSet postcodeCountSet
}

func partialCountRecipeDelivery(recipeDeliveryPart []recipeDelivery, options recipeCountOptions, c chan partialCountSets) {
	recipeCountSet, postcodeCountSet := countRecipeDelivery(recipeDeliveryPart, options)
	c <- partialCountSets{recipeCountSet, postcodeCountSet}
}
