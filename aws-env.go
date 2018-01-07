package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
	"os"
	"strings"
	"bufio"
)

func main() {
	if os.Getenv("AWS_ENV_PATH") == "" {
		log.Println("aws-env running locally, without AWS_ENV_PATH")
		return
	}


	ExportVariables(os.Getenv("AWS_ENV_PATH"), "")
}

func CreateClient() *ssm.SSM {
	session := session.Must(session.NewSession())
	return ssm.New(session)
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ExportVariables(path string, nextToken string) {

	client := CreateClient()

	input := &ssm.GetParametersByPathInput{
		Path:           &path,
		WithDecryption: aws.Bool(true),
	}

	if nextToken != "" {
		input.SetNextToken(nextToken)
	}

	output, err := client.GetParametersByPath(input)

	if err != nil {
		log.Panic(err)
	}

	f, err := os.Create(".env")
    	check(err)
	w := bufio.NewWriter(f)    

	for _, element := range output.Parameters {
		
		name := *element.Name
		value := *element.Value

		env := strings.Trim(name[len(path):], "/")
		value = strings.Replace(value, "\n", "\\n", -1)
		
		line := fmt.Sprintf("%s='%s'\n", env, value)

		w.WriteString(line)
	}

	w.Flush()

	if output.NextToken != nil {
		ExportVariables(path, *output.NextToken)
	}
}

