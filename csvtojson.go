package main



import (
	"os"
	"flag"
	"errors"
	"fmt"
	"path/filepath"
	"encoding/csv"
	"io"
)
type inputFile struct {
	filepath string
	separator string
	pretty bool
}

func getFileData() (inputFile, error) {
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("A filepath argument is required")
	}
	separator := flag.String("separator", "comma", "Column separator")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")

	flag.Parse()

	fileLocation := flag.Arg(0)

	if !(*separator == "comma" || * separator == "semicolon") {
		return inputFile{}, errors.New("Only comma or semicolon separators are allowed")
	}

	return inputFile{fileLocation, *separator, *pretty}, nil


}

func checkIfValidFile(fileName string) (bool, error) {
	if fileExtension := filepath.Ext(fileName); fileExtension != ".csv" {
		return false, fmt.Errorf("File %s is not a CSV fiel", fileName)
	}
	if _, err := os.Stat(fileName); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("File %s does not exist", fileName)
	}

	return true, nil


}

func processCsvFile(fileData inputFile, writerChannel chan<- map[string]string) {

	file, err := os.Open(fileData.filepath)
	check(err)
	defer file.Close()

	
	var headers, line []string
	 
	reader := csv.NewReader(file)
  	
	if fileData.separator == "semicolon" {
		reader.Comma = ';'
	}
  	
	headers, err = reader.Read()
	check(err) // Again, error checking
  	
	for {
		
		line, err = reader.Read()
		
		if err == io.EOF {
			close(writerChannel)
			break
		} else if err != nil {
			exitGracefully(err) 
		}
		
		record, err := processLine(headers, line)

		if err != nil { 
			fmt.Printf("Line: %sError: %s\n", line, err)
			continue
		}
		
		writerChannel <- record
	}
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
 }

 func check(e error) {
	if e != nil {
	   exitGracefully(e)
	}
 }


 func processLine(headers []string, dataList []string) (map[string]string, error) {
	
	if len(dataList) != len(headers) {
		return nil, errors.New("Line doesn't match headers format. Skipping")
	}
	
	recordMap := make(map[string]string)
	
	for i, name := range headers {
		recordMap[name] = dataList[i]
	}
	
	return recordMap, nil
}



func main(){
	fileData, err := getFileData()
	if(err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(fileData)

}