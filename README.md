
# GoStock

GoStock is a simple service that were built to access the Stock Market and query
data about all Tickers listed on nasdaq.

## Requirements:

* Python 3+
* Go (Most recent version).

## Thecnical detail

The application has two Main steps:

### 1 - GoLang Scrapping

In this step, the application basically access Yahoo Finance API and get data about ALL Tickers on the market.
The application filters Tickers based on their performance.

* Growth Percentage - The user can say the minimum growth expected in a Date Range.
* Maximum Value of the Ticker - The user can say the maximum value that they would like to pay for a ticker.
* Market Movement - The user can say the minimum market movement they expect on this ticker.

After those three validations, The ticker that fits on this rules are appended to the `buy_tickers.csv`.

### 2 - Refinement with Python

In this step, the application gets the `buy_tickers.csv` and do another scrapping on Yahoo Finance.

On this time, other filters are applied.

* Buy Recommendations - We get all recommendations for this ticker in a date range and count how many `buy` recommendations we have.
* Sustainability - The user can say the minimum value of the Sustainability that they accept.

## Usage - Installation

First, you need to Install the yfinance, *a lib built in Python that will be used below on second step*:

```bash
pip install yfinance
```

So you can just run:

```bash
go run main.go
```

And you get all done. `main.go` runs the Go code and the Python code.

## Output

On the end of the main.go run, you will have 2 output files: `buy_tickers.csv` and `result.csv`.

`buy_tickers.csv`: Result of Go scrapping.
`result.csv`: Result of Python refinement.

The `result.csv` is the file that will show you the Best Tickers for your pocket.
