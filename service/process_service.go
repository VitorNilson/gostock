package service

import (
	"fmt"
	"strconv"
	"sync"
	"vitornilson1998/gostock/utils"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
	"github.com/shopspring/decimal"
)

func ProcessChunk(pipe <-chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for element := range pipe {
		for _, ticker := range element {

			params := buildParams(ticker, 1, 1, 2022, 14, 10, 2022)

			iter := chart.Get(params)

			var totalPositiveVariation float64
			var totalNegativeVariation float64
			
			for iter.Next() {

				var OneHundred float64 = 100.00

				if iter.Bar().Open.Equal(decimal.Zero) {
					println("Open zero on ticker ", ticker)
					continue
				}

				div := periodVariation(iter)

				percentage := (div * OneHundred) - OneHundred

				var zero float64

				if percentage >= zero {
					totalPositiveVariation += percentage
				} else {
					totalNegativeVariation += percentage
				}

			}

			q, err := quote.Get(ticker)
			if err != nil {
				println("Error while querying ticker ", ticker)
			}
		
			if q == nil {
				println("The ticker quote is nil ", ticker)
				continue
			}

			variation := (totalPositiveVariation + totalNegativeVariation)

			expectedGrowth := 30.00
			maxValueThatIWantToPay := 35.00
			marketMovement := 1_000_000

			isTickerInsideParameters := validateTicker(variation, expectedGrowth, q.Bid, maxValueThatIWantToPay, q.RegularMarketVolume, marketMovement)

			if isTickerInsideParameters {

				println("Buy ", ticker, " per ", float64(q.Bid))
				variation := totalPositiveVariation + totalNegativeVariation

				row := []string{ticker, strconv.FormatFloat(variation, 'f', 6, 64), strconv.FormatFloat(q.Bid, 'f', 6, 64)}
				print(row)
				utils.WriteCsv(row)
			}

			if err := iter.Err(); err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}

func buildParams(ticker string, startDay int, startMonth int, startYear,
	endDay int, endMonth int, endYear int) *chart.Params {

	return &chart.Params{
		Symbol:   ticker,
		Interval: datetime.OneMonth,
		Start:    &datetime.Datetime{Day: startDay, Month: startMonth, Year: startYear},
		End:      &datetime.Datetime{Day: endDay, Month: endMonth, Year: endYear},
	}

}

// Divides the close by the open value and then returns the percentage.
func periodVariation(iter *chart.Iter) float64 {
	value, _ := iter.Bar().Close.Div(iter.Bar().Open).BigFloat().Float64()
	return value
}

// variation: Current % of variation in a period.
// growthGoal: Minimum % of growth expected in the period.
// actualBid: Current quotation of the stock.
// expectedBid: Maximum value that I would like to pay for.
// regularMarketVolume: Current market volume (Average of transactions)
// regularVolumeGoal: Minimum market volume that I would like to have on my stock (liquidity)
func validateTicker(variation float64, growthGoal float64, actualBid float64, expectedBid float64, regularMarketVolume int, regularVolumeGoal int) bool {
	return variation >= growthGoal && actualBid <= expectedBid && regularMarketVolume >= regularVolumeGoal
}
