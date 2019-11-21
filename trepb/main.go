package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	name := os.Getenv("NAME")
	cpf := os.Getenv("CPF")
	outputFolder := os.Getenv("OUTPUT_FOLDER")
	flag.Parse()
	if *month == 0 || *year == 0 {
		log.Fatalf("Month or year not provided. Please provide those to continue. --mes={} --ano={}\n")
	}
	if outputFolder == "" {
		outputFolder = "./output"
	}

	if err := os.Mkdir(outputFolder, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatalf("Error creating output folder(%s): %q", outputFolder, err)
	}

	filePath := filePath(outputFolder, *month, *year)
	if err := crawl(filePath, name, cpf, *month, *year); err != nil {
		log.Fatalf("Crawler error(%02d-%04d): %q", *month, *year, err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file (%s): %q", filePath, err)
	}
	defer f.Close()

	table, err := loadTable(f)
	if err != nil {
		log.Fatalf("error while loading data table from %s: %q", filePath, err)
	}

	records, err := employeeRecords(table)
	if err != nil {
		log.Fatalf("error while parsing data from table (%s): %q", filePath, err)
	}

	employees, err := json.MarshalIndent(records, "\n", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling error: %q", err)
	}
	fmt.Printf("%s", employees)
}
