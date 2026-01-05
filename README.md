# shopwatch-api
An api for scraping product information from amazon,flipcart,croma,etc


| Endpoint | Method | Body (Raw Text) | Description |
| :--- | :--- | :--- | :--- |
| `/amazon` | `POST` | `Product URL` | Scrapes full details from an Amazon product page. |
| `/flipkart` | `POST` | `Product URL` | Scrapes details from a Flipkart product page (via JSON-LD). |
| `/search` | `POST` | `Search Query` | Returns an array of Amazon product URLs for that query. |


### example

**Scrape Flipkart Product:**
```
curl -X POST http://localhost:1234/flipkart -d "https://www.flipkart.com/product-url-here"
```

output
```js
{
  "Title": "Product Name",
  "Price": "60990",
  "Seller": "Retailer Name",
  "Image": "https://image-url.jpg",
  "Reviews": "4.5"
}
```

# ______________IN DEVELOPMENT______________