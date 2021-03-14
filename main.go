package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChargersData struct {
	Chargers []struct {
		Array []int `json:"array"`
	} `json:"chargers"`
}

type RestaurantsData struct {
	Restaurants []struct {
		Name   string `json:"name"`
		Rating int    `json:"rating"`
	} `json:"restaurants"`
}

func readFile(path string, fileType string) interface{} {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Mobile; rv:15.0) Gecko/15.0 Firefox/15.0")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	if fileType == "charger" {
		var chargersdata ChargersData
		json.Unmarshal(body, &chargersdata)
		return chargersdata
	} else {
		var restaurantsdata RestaurantsData
		json.Unmarshal(body, &restaurantsdata)
		return restaurantsdata
	}
}

func isReachable(chargersdata ChargersData, initLevel int) []interface{} {

	results := []interface{}{}
	var reachable bool
	var level int

	for _, destinations := range chargersdata.Chargers {

		reachable = true
		level = initLevel

		destinationLength := len(chargersdata.Chargers)
		fmt.Println(destinationLength)

		for _, i := range destinations.Array {
			if level < 1 {
				results = append(results, 400)
				reachable = false
				break
			}
			// Every time the person travels 1 km, one charge drops and gains i charge
			level += i - 1
			// fmt.Println(i, level)
		}
		if reachable {
			results = append(results, 200)
		}
	}
	return results
}

func main() {
	r := gin.Default()
	chargersdata := readFile("https://s3-ap-southeast-1.amazonaws.com/he-public-data/chargers1e8f81f.json", "charger").(ChargersData)
	restaurantsdata := readFile("https://s3-ap-southeast-1.amazonaws.com/he-public-data/restaurants6010ade.json", "restaurants").(RestaurantsData)

	fmt.Println(restaurantsdata)

	fmt.Println(isReachable(chargersdata, 2))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	r.GET("/reachable/:dist", func(c *gin.Context) {
		dist := c.Param("dist")
		initLevel, err := strconv.Atoi(dist)

		if (err != nil) || (initLevel < 0) {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad parameter request"})
		} else {
			reachability := isReachable(chargersdata, initLevel)
			fmt.Println(reachability)
			c.JSON(http.StatusOK, gin.H{"data": reachability})
		}
	})

	r.Run()
}
