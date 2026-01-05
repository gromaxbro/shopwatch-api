package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type data_obj struct{
	Title string
	Price string
	Seller string
	Image string
	Reviews string
}

type AGRating struct {
	RatingValue float64 `json:"ratingValue"`
	ReviewCount int     `json:"reviewCount"`
}

type Offer struct {
	Price    float64 `json:"price"`
	Currency string  `json:"priceCurrency"`
}

type ProductLD struct {
	Type   string   `json:"@type"`
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	Offers Offer    `json:"offers"`
	Rating AGRating `json:"aggregateRating"`
}

func flipkart(url string) data_obj{
	request,error := http.NewRequest("GET",url,nil)

	request.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") 
	request.Header.Set("accept-language","en-GB,en;q=0.9")

	if error != nil{
        fmt.Println("error bro")
        return data_obj{}
    }

    client := http.Client{}

    response,err := client.Do(request) 
    if err != nil {
    fmt.Println("error making request: ", err)
    	return  data_obj{}
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Failed to parse the HTML document", err)
    }

	scriptContent := doc.Find("script#jsonLD").Text()
	

	if strings.Contains(scriptContent, "][") {
	    scriptContent = strings.Split(scriptContent, "][")[0] + "]"
	}
	// fmt.Println(scriptContent)
	var ldData []ProductLD
	var finalPrice int
	var productName string
	var Image string
	var Rating float64
	seller := doc.Find("#sellerName").Text()
	errr := json.Unmarshal([]byte(scriptContent), &ldData)
	if errr != nil {
		fmt.Println("Error parsing JSON-LD:", err)
		return data_obj{}
	}

		for _, item := range ldData {
	    // Only extract data if the type is "Product"
	    if item.Type == "Product" {
	        finalPrice = int(item.Offers.Price)
	        productName = item.Name
	        Image = item.Image
	        Rating =  item.Rating.RatingValue
	        break // We found it, so we can stop loopng
	    }
	}
	

	str := strconv.FormatFloat(Rating, 'f', 2, 64)
	fg := strconv.Itoa(finalPrice)
	data_ret := data_obj{productName,fg,seller,Image,str}
	return data_ret
}

func amazon_scrape(url string) data_obj {
	request,error := http.NewRequest("GET",url,nil)

	request.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") 
	request.Header.Set("accept-language","en-GB,en;q=0.9")

	if error != nil{
        fmt.Println("error bro")
        return data_obj{}
    }

    client := http.Client{}

    response,err := client.Do(request) 
    if err != nil {
    fmt.Println("error making request: ", err)
    	return data_obj{}
	}
	defer response.Body.Close()

	fmt.Println(response.StatusCode)

    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Failed to parse the HTML document", err)
    }
    title := strings.TrimSpace(doc.Find("#productTitle").First().Text())  
    price := doc.Find("span.a-price-whole").First().Text()
    seller := doc.Find("#sellerProfileTriggerId").First().Text()
    image,boob := doc.Find("#landingImage").First().Attr("src")
    reviews := doc.Find("#averageCustomerReviews_feature_div  span.a-icon-alt").Text()

    // fmt.Println(title)
    // fmt.Println((price))
    // fmt.Println(seller)
    fmt.Println(boob)
    // fmt.Println(reviews)

    result := data_obj{title,price,seller,image,reviews}
    return result
}

func amazon_search(search string) []string {
	safeSearch := url.QueryEscape(search)
	baseUrl := "https://www.amazon.in/s/?url=search-alias%%3Daps&field-keywords=%s"
	finalUrl := fmt.Sprintf(baseUrl, safeSearch)
	request,error := http.NewRequest("GET",finalUrl,nil)

	request.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") 
	request.Header.Set("accept-language","en-GB,en;q=0.9")

	if error != nil{
        fmt.Println("error bro")
        return []string{}
    }

    client := http.Client{}

    response,err := client.Do(request) 
    if err != nil {
    fmt.Println("error making request: ", err)
    	return  []string{}
	}
	defer response.Body.Close()

	// df,_ := io.ReadAll(response.Body)
	// fmt.Println(string(df))
	doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Failed to parse the HTML document", err)
    }
    base, _ := url.Parse("https://www.amazon.in") 

    return_data := []string{}
	doc.Find(`div[data-component-type="s-search-result"]`).Each(func(i int, item *goquery.Selection) { // [web:26]
		a := item.Find("a.a-link-normal.s-no-outline").First()
		href, ok := a.Attr("href")
		if !ok || href == "" {
			return
		}

		ref, err := url.Parse(href)
		if err != nil {
			return
		}

		abs := base.ResolveReference(ref) // makes absolute URL from relative [web:30]
		fmt.Println(abs.String())
		return_data = append(return_data, abs.String())
	})
	return return_data
}


func Mainbody(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("wassap niggers\n"))

	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	msg := strings.TrimSpace(string(b))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	amazon := amazon_scrape(msg)

	json.NewEncoder(w).Encode(amazon) 
}

func Searchbody(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("wassap niggers\n"))

	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	msg := strings.TrimSpace(string(b))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	amazon := amazon_search(msg)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // This stops the \u0026 conversion
	encoder.Encode(amazon)
}

func FlipBody(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("wassap niggers\n"))

	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	msg := strings.TrimSpace(string(b))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	amazon := flipkart(msg)

	json.NewEncoder(w).Encode(amazon) 
}

func main(){
	fmt.Println("....")

	mux := http.NewServeMux()
	mux.HandleFunc("/amazon", Mainbody)
	mux.HandleFunc("/search", Searchbody)
	mux.HandleFunc("/flipkart", FlipBody)

	fmt.Println("Listening to 1234")
	http.ListenAndServe(":1234", mux)

	// amazon := amazon_scrape("https://www.amazon.in/sspa/click?ie=UTF8&spc=MTo2ODI5NjQzMjE2MjY3NjA0OjE3NjcyODEzMTQ6c3BfYXRmOjIwMDkzOTEwMzczODk4OjowOjo&url=%2FBodyband-Strengthener-Adjustable-Equipment-Black-Orange%2Fdp%2FB0BKZK9JGB%2Fref%3Dsr_1_1_sspa%3Fcrid%3D2VPKWSXP5XW9I%26dib%3DeyJ2IjoiMSJ9.3MeWd_ttU_OKDywIniqi_unGrLHpsKAGm6Ky3MNN795-G48a7lqODGMMK3mebp_OHkcDg7D0cyw4hRDCCOsCcOuX0a0J-_JpqMCEbzu6eKLP_sOW_RVilB-MLPDisUIqyVDOoDUqlAeERGSA4oax7ATY4hSazNuSENRz5B9R_KyGEhXY1a5U9PnFUHIvlzT9UkfKuo0lXPyZJE-imH_hhWcOWytbYbYQ31BiiS_mmsHE0lzvKd3qmKkj4aQqcfacIBLwSLhY5dOiVIuby2XqelhLFQNJGjsWEC1y0fZd3RQ.sbRooHxtElxvr6BpPKb2yI6kv5GW0gAotOZJSQgyOg0%26dib_tag%3Dse%26keywords%3Dgripper%26qid%3D1767281314%26sprefix%3Dgripper%252Caps%252C373%26sr%3D8-1-spons%26aref%3DjOxRHrvPO8%26sp_csd%3Dd2lkZ2V0TmFtZT1zcF9hdGY%26psc%3D1&aref=jOxRHrvPO8&sp_cr=ZAZ")
	// fmt.Println(amazon)
	
}