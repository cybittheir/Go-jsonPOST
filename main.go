package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

var (
	version string
	build   string
)

func main() {

	// Open our jsonFile

	jsonFile, err := os.Open("conf.json")

	// if we os.Open returns an error then handle it

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened conf.json")

	// defer the closing of our jsonFile so that we can parse it later on

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}

	var confResult map[string]map[string]string

	json.Unmarshal([]byte(byteValue), &result)
	json.Unmarshal([]byte(byteValue), &confResult)

	postprotocol := fmt.Sprintf("%s", confResult["Responder"]["protocol"])

	testport := "23"

	if postprotocol == "https" {

		testport = "443"

	} else if postprotocol == "http" {

		testport = "80"

	} else if postprotocol == "ftp" {

		testport = "21"
	}

	ResponderIPAddr := fmt.Sprintf("%s:%s", confResult["Responder"]["oldip4"], testport)
	httpposturl := fmt.Sprintf("%s://%s/%s", confResult["Responder"]["protocol"], confResult["Responder"]["oldip4"], confResult["Responder"]["confURL"])

	jsonPost, err := json.Marshal(result["newconfig"])

	fmt.Println("Trying connect to " + ResponderIPAddr)

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ResponderIPAddr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err == nil {

		conn.Close()

		fmt.Println("Connected successful to " + ResponderIPAddr)

		fmt.Println("Map of new config: ", result["newconfig"])

		fmt.Println("HTTP JSON POST URL:", httpposturl)

		fmt.Println("Sending POST request...")

		request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonPost))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		response, error := client.Do(request)

		if error != nil {

			panic(error)

		}

		defer response.Body.Close()

		fmt.Println("response Status:", response.Status)
		fmt.Println("response Headers:", response.Header)

		body, _ := ioutil.ReadAll(response.Body)

		fmt.Println("response Body:", string(body))

	} else {

		fmt.Println("No connection to " + ResponderIPAddr)

	}

}
