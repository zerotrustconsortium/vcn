package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

func firstFile(dir string) (io.Reader, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		return os.Open(dir + "/" + f.Name())
	}
	return nil, fmt.Errorf("empty directory: %s", dir)
}

func contains(xs []string, x string) bool {
	for _, a := range xs {
		if a == x {
			return true
		}
	}
	return false
}

func readPassword(msg string) (string, error) {
	fmt.Print(msg)
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println(".")
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func createDirectoryInfrastructure() {
	err := os.MkdirAll(WalletDirectory(), os.FileMode(0700))
	if err != nil {
		log.Fatal(err)
	}
}

// PrintErrorURLCustom takes custom domain and status code
func PrintErrorURLCustom(domain string, code int) {

	fmt.Print("Get help for this error at:\n")

	color.Set(StyleError())
	fmt.Print(formatErrorURLCustom(domain, code))
	color.Unset()

	fmt.Println()
	return

}
func formatErrorURLCustom(domain string, status int) string {

	errorPage := ErrorWikiURL()

	return fmt.Sprintf("%s%s-%d", errorPage, domain, status)

}

// PrintErrorURLByEndpoint takes API errors and creates github wiki links
func PrintErrorURLByEndpoint(resource string, verb string, status int) {

	fmt.Print("Get help for this error at:\n")

	color.Set(StyleError())
	fmt.Print(formatErrorURLByEndpoint(resource, verb, status))
	color.Unset()

	fmt.Println()
	return

}
func formatErrorURLByEndpoint(resource string, verb string, status int) string {

	errorPage := ErrorWikiURL()

	// get last part of endpoint
	x := strings.Split(resource, "/")
	resource = x[len(x)-1]

	return fmt.Sprintf("%s%s-%s-%d", errorPage, resource, strings.ToLower(verb), status)

}

func checkJQExists() {
	_, err := exec.LookPath("jq")
	if err != nil {
		fmt.Printf("<jq> is not installed.")
		PrintErrorURLCustom("jq", 404)
		os.Exit(1)
	}
}

func getDockerHash(param string) (hash string) {

	checkJQExists()

	dockerID := strings.Replace(param, "docker:", "", 1)

	// TODO: sanitize even further
	// so far, let's check dockerID is a string without whitespaces
	dockerID = strings.Replace(dockerID, " ", "", -1)

	// docker inspect <image> | jq -c ".[0].Id"
	cmd := fmt.Sprintf(`docker inspect %s | jq -c ".[0].Id"`, dockerID)
	output, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
		PrintErrorURLCustom("docker", 500)
		os.Exit(1)
	}
	// ugly hack: better read in the docker information
	// with a proper os.exec below
	if strings.TrimSpace(string(output)) == "null" {
		fmt.Printf("Docker image cannot be found: <%s>", dockerID)
		PrintErrorURLCustom("docker", 404)
		os.Exit(1)
	}

	hash = string(output)
	hash = strings.Replace(hash, `"`, ``, -1)
	hash = strings.Replace(hash, "sha256:", "", 1)

	/*
		c1 := exec.Command("docker", "inspect", dockerID)
		c2 := exec.Command("jq", "-c", ".[0].Id")

		c2.Stdin, _ = c1.StdoutPipe()
		c2.Stdout = os.Stdout
		_ = c2.Start()
		_ = c1.Run()
		_ = c2.Wait()

		out, err := c2.Output()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		fmt.Printf("combined out:\n%s\n", string(out))
	*/

	// I wasted 2 hours of my precious life looking for a wqy
	// to unmarshal JSON in golang
	// with an unknown JSON scheme
	// in combination with a top-level array

	return hash
}
