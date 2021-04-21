package knative

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var kubeconfigFile string

func init() {
	kubeconfigFile = os.Getenv("KUBECONFIG_FILE")
}

func Create(id, imageName string) (string, error) {
	out, err := exec.Command("kn", "service", "create", id, "--image", imageName, "--kubeconfig", kubeconfigFile,
		"--env", fmt.Sprintf("TARGET=\"Func ID: %s\"", id), "--user", "1000").Output()
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("exec: %s", err)
	}

	fmt.Printf("Output %s\n", out)
	myString := string(out[:])

	re := regexp.MustCompile("http://.*")
	url := re.FindString(myString)
	fmt.Println(strings.TrimSuffix(url, "\n"))

	return url, nil
}
