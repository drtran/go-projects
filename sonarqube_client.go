package main

import "fmt" 
import "log" 
import "net/http" 
import "io/ioutil"
import "os"
import "github.com/tidwall/gjson"
import "strconv"



func main() {

	if len(os.Args) == 1 {
		showUsage(os.Args[0])
		return
	}

	site := "http://localhost:9000"
	command := `get-coverage`

	if len(os.Args) > 1 {
		command, site = processArgs(os.Args[1:], `--site`)
	}
	
	dispatch(command, site)

	
}

func processArgs(args []string, key string) (command string, site string) {
	command = args[0]
	
	i := 1

	for i < len(args) {
		if args[i] == key {
			site = args[i+1]
			break
		}
		i++
	}
	return command, site
}

func showUsage(pgmName string) {
	println(fmt.Sprintf("%s command [options]", pgmName))
	println("\ncommands:") 
	println(`get-coverage: retrieve code coverage in percentage`)
	println("\noptions:")
	println("--site sitename: i.e. --site http://localhost:9000")
}

func dispatch(command string, site string) {
	if "get-coverage" == command {
		coverage := getCoverage(site)
		fmt.Printf("%d\n", coverage)
	} else if "get-complexity" == command {
		complexity := getComplexity(site)
		fmt.Printf("%d\n", complexity)
	}
}

func getCoverage(site string) int {
	data := callSonarQubeServer(site)

	coverageStr := gjson.Get(string(data),`component.measures.#[metric="line_coverage"].value`).String()
	f, _ := strconv.ParseFloat(coverageStr, 64)
	return int(f)
}

func getComplexity(site string) int {
	data := callSonarQubeServer(site)

	coverageStr := gjson.Get(string(data),`component.measures.#[metric="complexity"].value`).String()
	i, _ := strconv.Atoi(coverageStr)
	return i
}


func callSonarQubeServer(site string) string {
	path := "/api/measures/component"
	project := "org.springframework.samples:spring-petclinic"
	metrics := "ncloc,line_coverage,code_smells,complexity"
	additionalFields := "metrics,periods"
	
	url := fmt.Sprintf("%s%s?componentKey=%s&metricKeys=%s&additionalFields=%s", site, path, project, metrics, additionalFields)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err);
		return ""
	}

	client := &http.Client{}
	resp, err := client.Do(req);
	if err != nil {
		log.Fatal("Do: ", err)
		return ""
	}

	data, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("data: " + string(data))

	return string(data)
}

