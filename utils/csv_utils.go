package utils

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadCsv() []string {
	f, err := os.Open("assets/nasdaq.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var tickers []string

	for i := 0; i < len(data); i++ {
		tickers = append(tickers, data[i][0])
	}

	return tickers
}

func CreateCsv() {
	csvFile, err := os.Create("buy_tickers.csv")

	if err != nil {
		panic(err)
	}

	wr := csv.NewWriter(csvFile)
	head := []string{"Ticker", "Variation", "Value"}

	errore := wr.Write(head)
	if errore != nil {
		panic(err)
	}
	wr.Flush()
	csvFile.Close()
}

func WriteCsv(row []string) {
	println("Writing ticker on csv")
	file, err := os.OpenFile("buy_tickers.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		println(err.Error())
	}

	wr := csv.NewWriter(file)

	werr := wr.Write(row)
	if werr != nil {
		println(werr.Error())
	}

	wr.Flush()
	file.Close()

}

func ChunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
