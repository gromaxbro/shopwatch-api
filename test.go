package main

import (
	// "encoding/json"
	// "encoding/json"
	"fmt"
	"io"
	// "strings"

	// "strings"
	// "io"
	// "log"
	"net/http"

	// "net/url"

	// "strings"
	// "github.com/PuerkitoBio/goquery"
)

func main(){
	request,error := http.NewRequest("GET","https://www.reliancedigital.in/product/hp-victus-15-15-fb3123ax-gaming-laptop-amd-ryzen-7-7445h16-gb-512-gb-ssd4-gb-nvidia-geforce-rtx-2050windows-11-homems-office-home-2024full-hd-3962-cm-156-inch-mica-silver-black-chrome-logo-mgj9z0-9533142?region_id=RDHD19&gad_campaignid=22588256564",nil)

	request.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36") 
	request.Header.Set("accept-language","en-GB,en;q=0.9")

	if error != nil{
        fmt.Println("error bro")
        return 
    }

    client := http.Client{}

    response,err := client.Do(request) 
    if err != nil {
    fmt.Println("error making request: ", err)
    	return 
	}
	defer response.Body.Close()
	dfwaf,_ := io.ReadAll(response.Body)
	fmt.Println(string(dfwaf))
}