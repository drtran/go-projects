package main


import "fmt" 
import "log" 
import "net/http" 
import "os"
import "io/ioutil"
import "strings"
import "time"
import "strconv"

type project struct {
	name string
	status string
	uid string
}

func main() {

	site := "http://pet-clinic-nsd-poc.192.168.99.101.nip.io/petclinic/"
	grepMatchString := "Version:"
	site = os.Args[1]
	grepMatchString = os.Args[2]
	sleepTimeInSeconds, _ := strconv.Atoi(os.Args[3])
	previousLine := ""
	trimmedLine := ""

	var count int = 0

	for {
		docLines := strings.Split(curl(site), "\n");

		for _, docLine := range docLines {
			if strings.Contains(docLine, grepMatchString) {
				trimmedLine = strings.Trim(docLine, " \t")
				time.Sleep(time.Duration(sleepTimeInSeconds) * time.Second)
				if previousLine != trimmedLine {
					previousLine = trimmedLine
					fmt.Println("")
					count = 1
				}
				fmt.Println(fmt.Sprintf ("%d:%s", count, trimmedLine))
				count++
			}
		}
	}
}

func curl(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err);
		return ""
	}
	client := &http.Client {}

	resp, err := client.Do(req);
	if err != nil {
		log.Fatal("Do: ", err)
		return ""
	}

	data, _ := ioutil.ReadAll(resp.Body)

	return string(data)
}
