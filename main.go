package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func fetch(url *neturl.URL) (string, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func writeTempScript(filename, content string) (string, error) {
	tmpfile, err := ioutil.TempFile("", filename)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write([]byte(content))
	if err != nil {
		return "", err
	}

	err = os.Chmod(tmpfile.Name(), 0755)
	if err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func runCommand(name string, arg ...string) error {
	cmd := exec.CommandContext(context.Background(), name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ask(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message + " [y/n]: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return (text == "yes" || text == "y")
}

func curlSh() error {
	var (
		integrity *SRIFlag
		pager     = getEnv("CURLSH_PAGER", getEnv("PAGER", "less"))
		trusted   bool
		url       = new(neturl.URL)
		useSudo   bool
	)

	flag.Var(integrity, "hash", "SRI hash")
	flag.StringVar(&pager, "pager", pager, "select pager (CURLSH_PAGER, PAGER)")
	flag.BoolVar(&useSudo, "sudo", false, "run the script with sudo")
	flag.BoolVar(&trusted, "trusted", false, "whenver the script is trusted")
	flag.Var(URLFlag{url}, "url", "URL to fetch")

	flag.Parse()

	filename := path.Base(url.Path)

	// Fetch script
	scriptContent, err := fetch(url)
	if err != nil {
		return err
	}

	// Write to temporary place
	scriptPath, err := writeTempScript(filename, scriptContent)
	if err != nil {
		return err
	}
	defer os.Remove(scriptPath)
	log.Println("DEBUG: script path:", scriptPath)

	// Check integrity
	if len(integrity.List) > 0 {
		// FIXME: pick the best algo
		if !integrity.List[0].Check([]byte(scriptContent)) {
			return fmt.Errorf("SRI check failed")
		}
		log.Println("SRI check passed")
	}
	log.Println("XXXXX", integrity.List)

	// Verify the script in a pager for extra eyeballing
	if trusted {
		if len(integrity.List) == 0 {
			if url.Scheme == "http" {
				log.Println("U FUCKING WOT MATE")
				log.Println("You disabled every security possible")
				log.Println("Any middleman can inject a new vrsion of the script")
				log.Println("Don't use HTTP without integrity checks")
				time.Sleep(5 * time.Second)
				log.Println("Deleting all files in the system")
				time.Sleep(2 * time.Second)
				log.Println("Ok I'm joking")
				os.Exit(99)
			} else {
				log.Println("WARNING: script integrity not verified. The script could change at any time and you would not know about it")
			}
		} else {
			log.Println("DEBUG: assuming the integrity to be intact")
		}
	} else {
		// FIXME: security hole, escape the scriptPath
		err = runCommand("sh", "-c", fmt.Sprintf("exec %s %s", pager, scriptPath))
		if err != nil {
			return err
		}
		if !ask("Did you read this script carefully and understand what it does?") {
			return fmt.Errorf("Aborting, script not trusted")
		}
	}

	// And finally, run the script
	if useSudo {
		err = runCommand("sudo", scriptPath)
	} else {
		err = runCommand(scriptPath)
	}

	return err
}

func main() {
	err := curlSh()
	if err != nil {
		log.Fatalln(err)
	}
}
