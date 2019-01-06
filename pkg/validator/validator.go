//Package validator defines functions to handle validation on request
package validator

import (
	"../model"
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

//ValidateRequest defines set of rules to validate incoming request
//This function can be injected to handlers to perform validation
//Validates both mandatory fields as well as email format with regex.
func ValidateRequest(r *http.Request) (bool, string) {

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)

	var m = model.Metadata{}
	err := yaml.Unmarshal([]byte(bodyString), &m)

	if err != nil {
		return false, "Cannot unmarchall from request to object"
	}

	//Check mandatory fields and email format
	if result := checkFields(&m); len(result) > 0 {
		return false, strings.Join(result, "-")
	}
	return true, ""
}

//checkFields checks all mandatory fields inside request body and also validates email format
func checkFields(m *model.Metadata) []string {
	var emptyFields []string

	if m.Version == "" || len(m.Version) == 0 {
		emptyFields = append(emptyFields, "Version cannot be empty")
	}
	if m.Company == "" || len(m.Company) == 0 {
		emptyFields = append(emptyFields, "Company cannot be empty")
	}
	if m.Description == "" || len(m.Description) == 0 {
		emptyFields = append(emptyFields, "Description cannot be empty")
	}
	if m.License == "" || len(m.License) == 0 {
		emptyFields = append(emptyFields, "License cannot be empty")
	}
	if m.Source == "" || len(m.Source) == 0 {
		emptyFields = append(emptyFields, "Source cannot be empty")
	}
	if m.Title == "" || len(m.Title) == 0 {
		emptyFields = append(emptyFields, "Title cannot be empty")
	}
	if m.Website == "" || len(m.Website) == 0 {
		emptyFields = append(emptyFields, "Website cannot be empty")
	}

	if len(m.Maintainers) == 0 {
		emptyFields = append(emptyFields, "Maintainers cannot be empty")
	}
	if len(m.Maintainers) > 0 {
		for _, person := range m.Maintainers {
			if person.Name == "" {
				emptyFields = append(emptyFields, "Maintainer name cannot be empty")
			}
			if person.Email == "" {
				emptyFields = append(emptyFields, "Maintainer email cannot be empty")
			}
			if !isValidEmail(person.Email) {
				emptyFields = append(emptyFields, "Maintainer email address not correct")

			}
		}
	}
	return emptyFields
}

//isValidEmail validates email format
func isValidEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	return true
}
