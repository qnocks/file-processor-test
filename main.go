package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	CSVExt  = ".csv"
	JSONExt = ".json"
)

type Product struct {
	Name   string `json:"product"`
	Price  int    `json:"price"`
	Rating int    `json:"rating"`
}

func readCSV(filename string, ch chan<- Product) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %s\n", err)
		return
	}

	for _, record := range records {
		price, err := strconv.Atoi(record[1])
		if err != nil {
			continue
		}

		rating, err := strconv.Atoi(record[2])
		if err != nil {
			continue
		}

		ch <- Product{
			Name:   record[0],
			Price:  price,
			Rating: rating,
		}
	}

	close(ch)
}

func readJSON(filename string, ch chan<- []Product) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading JSON file: %s\n", err)
		return
	}

	var products []Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %s\n", err)
		return
	}

	ch <- products
	close(ch)
}

func findMostExpensiveProduct(products []Product, ch chan<- Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price > products[j].Price
	})

	ch <- products[0]
	close(ch)
}

func findHighestRatedProduct(products []Product, ch chan<- Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].Rating > products[j].Rating
	})

	ch <- products[0]
	close(ch)
}

func main() {
	var products []Product
	mostExpensiveCh := make(chan Product)
	highestRatedCh := make(chan Product)

	if len(os.Args) < 2 {
		log.Fatalf("File to proccess not found. Specify as command line argument")
	}
	filename := os.Args[1]
	ext := filepath.Ext(filename)

	switch ext {
	case CSVExt:
		ch := make(chan Product)
		go readCSV(filename, ch)

		for p := range ch {
			products = append(products, p)
		}
	case JSONExt:
		ch := make(chan []Product)
		go readJSON(filename, ch)

		for p := range ch {
			products = append(products, p...)
		}
	default:
		log.Fatalf("Unsupported file format: %s\n", ext)
	}

	if len(products) == 0 {
		log.Fatal("No products found in the file")
	}

	productsCopy := make([]Product, len(products))
	copy(productsCopy, products)

	go findMostExpensiveProduct(products, mostExpensiveCh)
	go findHighestRatedProduct(productsCopy, highestRatedCh)

	mostExpensiveProduct := <-mostExpensiveCh
	highestRatedProduct := <-highestRatedCh

	fmt.Printf("Самый дорогой продукт - %s\n", mostExpensiveProduct.Name)
	fmt.Printf("С самым высоким рейтингом - %s\n", highestRatedProduct.Name)
}
