package main


import "fmt" 
import "log" 
import "net/http" 
import "io/ioutil"
import "crypto/tls"
import "github.com/tidwall/gjson"
// import "strconv"
// import "os"

type project struct {
	name string
	status string
	uid string
}

func main() {

	site := "https://192.168.99.100:8443"
	tokenString := "CZQhzrOwjJ-l0ENg36AH9N6Z_T8LgaTBrO5shKGqVmU"
	projects := getProjects(site, tokenString)

	if projects == nil || len(projects) == 0 {
		fmt.Println("projects: empty!")
	} else {
	 	fmt.Println(projects[0])
	}

}

func getProjects(site string, tokenString string) []project {
	command := "projects"
	data := openshiftGet(site, tokenString, command)
	names := gjson.Get(string(data),`items.#.metadata.name`).Array()
	uids := gjson.Get(string(data),`items.#.metadata.uid`).Array()
	statuses := gjson.Get(string(data),`items.#.status.phase`).Array()

	if len(names) == 0 {
		log.Println("Result is empty!")
		return nil
	}

	projects := [] project{}

	for i := 0; i < len(names); i++ {
		proj := project{name: names[i].String(), uid: uids[i].String(), status: statuses[i].String()}
		projects = append(projects, proj)
	}

	return projects
}

/**
 * For openshift: 
 * to get token:
 * + oc login -u ... -p ...
 * + oc whoami -t
 * then copy the token string that can be used for 24 hours. 
 * 
 */
func openshiftGet(site string, tokenString string, command string) string {
	path := "/oapi/v1"
	
	url := fmt.Sprintf("%s%s/%s", site, path, command)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err);
		return ""
	}

	bearerToken := "Bearer " + tokenString
	req.Header.Add("Authorization", bearerToken)

	tr := &http.Transport {
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client {Transport: tr}

	resp, err := client.Do(req);
	if err != nil {
		log.Fatal("Do: ", err)
		return ""
	}

	data, _ := ioutil.ReadAll(resp.Body)

	return string(data)
}
