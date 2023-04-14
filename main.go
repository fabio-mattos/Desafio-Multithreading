package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ApiCepDto struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCepDto struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	url_api_cep := "https://cdn.apicep.com/file/apicep/"
	url_via_cep := "http://viacep.com.br/ws/"

	channelApiCep := make(chan ApiCepDto)
	channelViaCep := make(chan ViaCepDto)

	cep := "88036-610"

	go func() {
		url := url_api_cep + cep + ".json"

		response := request(url)
		if response != nil {
			var apiCepDto ApiCepDto

			if err := json.Unmarshal(response, &apiCepDto); err != nil {
				return
			}

			channelApiCep <- apiCepDto
		}
	}()

	go func() {
		url := url_via_cep + cep + "/json/"

		response := request(url)
		if response != nil {
			var viaCepDto ViaCepDto

			if err := json.Unmarshal(response, &viaCepDto); err != nil {
				return
			}

			channelViaCep <- viaCepDto
		}
	}()

	select {
	case responseDaApiCep := <-channelApiCep:
		fmt.Printf("A resposta da API CEP foi mais rápida: %+v\n", responseDaApiCep)

	case responseDaViaCep := <-channelViaCep:
		fmt.Printf("A resposta da Via CEP foi mais rápida: %+v \n", responseDaViaCep)

	case <-time.After(time.Second):
		fmt.Println("Request timeout")
	}
}

func request(url string) []byte {
	response, err := http.Get(url)

	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	return body
}
