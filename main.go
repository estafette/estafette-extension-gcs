package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alecthomas/kingpin"
)

var (
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	paramsJSON      = kingpin.Flag("params", "Extension parameters, created from custom properties.").Envar("ESTAFETTE_EXTENSION_CUSTOM_PROPERTIES").Required().String()
	credentialsJSON = kingpin.Flag("credentials", "GKE credentials configured at service level, passed in to this trusted extension.").Envar("ESTAFETTE_CREDENTIALS_KUBERNETES_ENGINE").Required().String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log to stdout and hide timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log startup message
	logInfo("Starting %v version %v...", app, version)

	logInfo("Unmarshalling credentials parameter...")
	var credentialsParam CredentialsParam
	err := json.Unmarshal([]byte(*paramsJSON), &credentialsParam)
	if err != nil {
		log.Fatal("Failed unmarshalling credential parameter: ", err)
	}

	logInfo("Validating required credential parameter...")
	valid, errors := credentialsParam.ValidateRequiredProperties()
	if !valid {
		log.Fatal("Not all valid fields are set: ", errors)
	}

	logInfo("Unmarshalling injected credentials...")
	var credentials []GKECredentials
	err = json.Unmarshal([]byte(*credentialsJSON), &credentials)
	if err != nil {
		log.Fatal("Failed unmarshalling injected credentials: ", err)
	}

	logInfo("Checking if credential %v exists...", credentialsParam.Credentials)
	credential := GetCredentialsByName(credentials, credentialsParam.Credentials)
	if credential == nil {
		log.Fatalf("Credential with name %v does not exist.", credentialsParam.Credentials)
	}

	logInfo("Unmarshalling parameters / custom properties...")
	var params Params
	err = json.Unmarshal([]byte(*paramsJSON), &params)
	if err != nil {
		log.Fatal("Failed unmarshalling parameters: ", err)
	}

	logInfo("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults()

	logInfo("Retrieving service account email from credentials...")
	var keyFileMap map[string]interface{}
	err = json.Unmarshal([]byte(credential.AdditionalProperties.ServiceAccountKeyfile), &keyFileMap)
	if err != nil {
		log.Fatal("Failed unmarshalling service account keyfile: ", err)
	}
	var saClientEmail string
	if saClientEmailIntfc, ok := keyFileMap["client_email"]; !ok {
		log.Fatal("Field client_email missing from service account keyfile")
	} else {
		if t, aok := saClientEmailIntfc.(string); !aok {
			log.Fatal("Field client_email not of type string")
		} else {
			saClientEmail = t
		}
	}

	logInfo("Storing gcs credential %v on disk...", credentialsParam.Credentials)
	err = ioutil.WriteFile("/key-file.json", []byte(credential.AdditionalProperties.ServiceAccountKeyfile), 0600)
	if err != nil {
		log.Fatal("Failed writing service account keyfile: ", err)
	}

	logInfo("Authenticating to google cloud")
	runCommand("gcloud", []string{"auth", "activate-service-account", saClientEmail, "--key-file", "/key-file.json"})

	logInfo("Setting gcloud account")
	runCommand("gcloud", []string{"config", "set", "account", saClientEmail})

	logInfo("Setting gcloud project")
	runCommand("gcloud", []string{"config", "set", "project", credential.AdditionalProperties.Project})

	// clean up old stuff
	switch params.Action {
	case "copy":

		cpArgs := []string{"cp", "-r"}
		if params.ACL != "" {
			cpArgs = append(cpArgs, "-a", params.ACL)
		}
		if params.Compress != nil && *params.Compress {
			cpArgs = append(cpArgs, "-Z")
		}
		if params.Parallel != nil && *params.Parallel {
			cpArgs = append(cpArgs, "-m")
		}

		cpArgs = append(cpArgs, params.Source, fmt.Sprintf("gs://%v/%v", params.Bucket, params.Destination))

		runCommand("gsutil", cpArgs)

		break
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand(command string, args []string) {
	err := runCommandExtended(command, args)
	handleError(err)
}

func runCommandExtended(command string, args []string) error {
	logInfo("Running command '%v %v'...", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Dir = "/estafette-work"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	log.Println("")
	return err
}

func getCommandOutput(command string, args []string) (string, error) {
	logInfo("Getting output for command '%v %v'...", command, strings.Join(args, " "))
	output, err := exec.Command(command, args...).Output()

	return string(output), err
}

func logInfo(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	log.Printf("%v\n\n", formattedMessage)
}
