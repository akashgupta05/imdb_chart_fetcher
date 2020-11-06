package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Input struct {
	ChartURL   string
	ItemsCount int
}

func validateInput() (*Input, error) {
	if len(os.Args) != 3 {
		return nil, errors.New("incorrect no. of parameters")
	}

	chartURL := os.Args[1]
	_, err := url.Parse(chartURL)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return nil, errors.New("invalid count")
	}
	if count <= 0 {
		return nil, errors.New("count should be greater than 0")
	}

	return &Input{ChartURL: chartURL, ItemsCount: count}, nil

}

func main() {
	input, err := validateInput()
	if err != nil {
		fmt.Printf("input validation failed : %v", err.Error())
		return
	}

	imdbCharts := NewImdbCharts()
	imdbCharts.fetch(input.ChartURL, input.ItemsCount)

	result, err := json.Marshal(imdbCharts.Charts)
	if err != nil {
		fmt.Printf("json marshaling failed : %v", err.Error())
		return
	}

	fmt.Print(string(result))
}
