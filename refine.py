import pandas as pd
from IPython.display import display, HTML
import threading
import concurrent.futures
import yfinance as yf
import sys

lock = threading.Lock()

# buy_tickers.csv is the output of main.go.
# main.go create this file and run refine.py automatically, so please run main.go.
tickers = pd.read_csv('buy_tickers.csv', dtype={'Ticker':'string'})

result_csv = {'ticker': [], 'recommendation': [], 'sustainability':[]}

def generate_chunks(lst, n) -> list:
    for i in range(0, len(lst), n):
        yield lst[i:i + n]


def do_stock_query(tickers: list, lock):
  
    for ticker in tickers:    
        print(ticker)
        try:
            tk = yf.Ticker(ticker)

            recommendations_df = pd.DataFrame(tk.recommendations)
            recommendations_df.reset_index(inplace=True)

            sustainability_df = pd.DataFrame(tk.sustainability);

            governace = sustainability_df.filter(regex='governanceScore', axis=0)
            governace.reset_index(inplace=True)


            if not governace.empty:
                sustainability_value = governace['Value']
            else:
                continue

            sustainability_ok =  sustainability_value <= 9


            if sustainability_ok.values[0]:
                
                if not recommendations_df.empty:
                    filtered_recommendations_df = recommendations_df[(recommendations_df['Date'] > '2022-04-01')]


                if not filtered_recommendations_df.empty:
                    try:
                        # This opperation can throws a exception
                        buy_recommendations_amount = filtered_recommendations_df['To Grade'].value_counts().Buy
                        
                        if buy_recommendations_amount >= 1:
                            # Thread Safe appending values to result_csv
                            with lock:
                                result_csv['ticker'].append(ticker)
                                result_csv['recommendation'].append(filtered_recommendations_df['To Grade'].value_counts().to_string())
                                result_csv['sustainability'].append(sustainability_value.to_string())

                        else:
                            print(f"Buy Recommendations Below exptectations, current: {buy_recommendations_amount}")
                    except:
                        print("Cannot Count Buy recommendations.")

            else:
                print(f"Sustainability below expectations on Ticker {ticker}. Current Value: {sustainability_value}")

        
        
        except Exception as e:
            print("General Error")
            print(sys.exc_info()[0])
 
# 
# Here are the most important point about the analysis.
# The first argument is the list of Tickers, the second is the size of the 
# chunk that will be generated.
# Let's Say that the List have a size of 20 and you are building chunks of 5.
# On the end, you will have 4 chunks. 
# Lower we get each chunk and create a executor (A Thread) for each chunk.
#
# I Really recommend to keep your chunk size in a small number, so theoretically you will have more
# processing power.
#
chunks = generate_chunks(tickers['Ticker'].str.strip(), 5)

# Parallel Search on stock market
with concurrent.futures.ThreadPoolExecutor() as executor:
    futures = []
    
    # Read comments between lines 75-85 to more understanding.
    for chunk in chunks:
        futures.append(executor.submit(do_stock_query, tickers=chunk, lock=lock))
    
    
    executor.map(do_stock_query, chunks)
        
        
result = pd.DataFrame(data=result_csv)
result.to_csv('result.csv')

# If you are on Jupyter Notebook, this is a good way to visualize the results.
# display(HTML(df.to_html()))
     
