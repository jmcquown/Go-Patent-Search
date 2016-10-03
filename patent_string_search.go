package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"log"
	"io"
	
)

func main() {
	//Read input from the terminal
	reader := bufio.NewReader(os.Stdin)

	//Get patent #
	//Use 20110250858 for testing
	fmt.Print("Enter publication #: ")
	pubNum, _ := reader.ReadString('\n')
	//Remove \n from the pub #
	pubNum = strings.TrimRight(pubNum, "\n")

	//Get paragraph and line #
	fmt.Print("Enter paragraph # (ex: [0001]): ")
	paragraphNum, _ := reader.ReadString('\n')
	paragraphNum = strings.TrimRight(paragraphNum, "\n")

	// fmt.Print("Enter line # (ex: 1-3): ")
	// lineNum, _ := reader.ReadString('\n')
	// lineNum = strings.TrimRight(lineNum, "\n")

	// fmt.Print("Publication #: " + pubNum)
	// fmt.Print("Paragraph: " + paragraphNum)
	// fmt.Print("Line(s): " + lineNum)

	fmt.Println(pubNum)

	//Create the URL for the patent by appending the publication number at the appropriate locations
	patentURL := bytes.NewBuffer(nil)
	patentURL.WriteString("http://appft.uspto.gov/netacgi/nph-Parser?Sect1=PTO1&Sect2=HITOFF&d=PG01&p=1&u=%2Fnetahtml%2FPTO%2Fsrchnum.html&r=1&f=G&l=50&s1=%22")
	patentURL.WriteString(pubNum)
	patentURL.WriteString("%22.PGNR.&OS=DN/")
	patentURL.WriteString(pubNum)
	patentURL.WriteString("&RS=DN/")
	patentURL.WriteString(pubNum)

	//Print out the url
	fmt.Println(patentURL.String() + "\n")

	//Get request for the webpage
	response, _ := http.Get(patentURL.String())
	
	body := response.Body
	defer body.Close()

	//Convert body (io.Reader) to a string
	tempBuffer := new(bytes.Buffer)
	tempBuffer.ReadFrom(body)
	htmlBody := tempBuffer.String()

	paragraphIndex := strings.Index(htmlBody, paragraphNum)	
	fmt.Println("Index of paragraph in text %d\n", paragraphIndex)

	//Seek to that index in the body string
	bodyReader := strings.NewReader(htmlBody)
	//Cast paragraphIndex to int64 b/c fuck you and seek from the start
	if _, err := bodyReader.Seek(int64(paragraphIndex), io.SeekStart); err != nil {
		log.Fatal(err)
	}

	//Create a string to hold the paragraph we want
	paragraphData := ""
	//Variable to account for the fact that the first char read is a [
	//When secondBracket is true, then we know that the next [ that appears is the start of a new paragraph
	secondBracket := false

	//Read the string character by character until we reach a [ (91 in ASCII)
	for {
		if char, _, err := bodyReader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		//If the char is [ and it is the second bracket || the char is a < (only happens after the last paragraph)
		} else if (char == 91 && secondBracket == true) || char == 60 {
				secondBracket = true
				break
		} else {	//Append the character to the string
			paragraphData = paragraphData + string(char)
		}
	}
	//Print the paragraph
	fmt.Println(paragraphData)
}
