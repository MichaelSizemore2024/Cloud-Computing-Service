import requests
from bs4 import BeautifulSoup
import time

def scrape_quotes(start_page):
    while True:
        # Send a GET request to the website
        url = f'http://quotes.toscrape.com/page/{start_page}/'
        response = requests.get(url)

        # Check if the request was successful (status code 200)
        if response.status_code == 200:
            # Parse the HTML content
            soup = BeautifulSoup(response.text, 'html.parser')

            # Extract quotes and authors
            quotes = soup.select('.quote span.text')
            authors = soup.select('.quote small.author')

            # Print each quote with its author
            for quote, author in zip(quotes, authors):
                # Remove extra quotation marks
                cleaned_quote = quote.text.strip('“”')
                print(f'"{cleaned_quote}" - {author.text}')
                time.sleep(5)

            # Check if there is a "Next" button on the page
            next_button = soup.select('.next')
            if not next_button:
                break  # No more pages, exit the loop

            # Increment the page number for the next request
            start_page += 1

            # Add a delay before making the next request (2 seconds in this case)
            time.sleep(2)
        else:
            print(f"Error: Unable to fetch the webpage (Status Code: {response.status_code})")
            break  # Exit the loop on error

# Set the starting page
start_page = 1

# Run the scraper
scrape_quotes(start_page)