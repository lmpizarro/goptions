package rfx

import (
	"fmt"

	"net/http"
	"os"

	"encoding/csv"
)

type credentials struct {
	User     string
	Password string
	Account  string
}

func ReadCredentials(file string, v bool) credentials {

	csvFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	if v {
		fmt.Println("Successfully Opened CSV file")
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	line := csvLines[0]
	emp := credentials{
		User:     line[0],
		Password: line[1],
		Account:  line[2],
	}
	return emp
}

func Token(user, password string) string {
	r, err := http.NewRequest("POST", Url_Auth, nil)

	if err != nil {
		panic(err)
	}
	r.Header.Add("X-Password", password)
	r.Header.Add("X-Username", user)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.StatusCode)
	}

	token := res.Header["X-Auth-Token"]
	return token[0]
}

func Login() string{
	cred := ReadCredentials("./env.csv", false)
	user := cred.User
	password := cred.Password
	return Token(user, password)
}