package calgaryopendata

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/utilities/simplehttp"
)

func GetDevelopmentPermits() {

	url := "https://data.calgary.ca/resource/6933-unw5.json?$query=SELECT * WHERE applieddate > '2022-01-01T00:00:00.000'  AND latitude BETWEEN '51.022361' AND '51.038912' AND longitude BETWEEN '-114.117927' AND '-114.142638' ORDER BY applieddate DESC"
	response, err := simplehttp.SimpleGet(url, make(map[string]string))
	if err != nil {
		fmt.Println("error")
		panic(err)
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}
