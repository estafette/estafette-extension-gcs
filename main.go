package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
)

var (
	appgroup  string
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

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(appgroup, app, version, branch, revision, buildDate)

	// create context to cancel commands on sigterm
	ctx := foundation.InitCancellationContext(context.Background())

	log.Info().Msg("Unmarshalling credentials parameter...")
	var credentialsParam CredentialsParam
	err := json.Unmarshal([]byte(*paramsJSON), &credentialsParam)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling credential parameter")
	}

	log.Info().Msg("Validating required credential parameter...")
	valid, errors := credentialsParam.ValidateRequiredProperties()
	if !valid {
		log.Fatal().Msgf("Not all valid fields are set: %v", errors)
	}

	log.Info().Msg("Unmarshalling injected credentials...")
	var credentials []GKECredentials
	err = json.Unmarshal([]byte(*credentialsJSON), &credentials)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling injected credentials")
	}

	log.Info().Msgf("Checking if credential %v exists...", credentialsParam.Credentials)
	credential := GetCredentialsByName(credentials, credentialsParam.Credentials)
	if credential == nil {
		log.Fatal().Msgf("Credential with name %v does not exist.", credentialsParam.Credentials)
	}

	log.Info().Msg("Unmarshalling parameters / custom properties...")
	var params Params
	err = json.Unmarshal([]byte(*paramsJSON), &params)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling parameters")
	}

	log.Info().Msg("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults()

	log.Info().Msg("Retrieving service account email from credentials...")
	var keyFileMap map[string]interface{}
	err = json.Unmarshal([]byte(credential.AdditionalProperties.ServiceAccountKeyfile), &keyFileMap)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling service account keyfile")
	}
	var saClientEmail string
	if saClientEmailIntfc, ok := keyFileMap["client_email"]; !ok {
		log.Fatal().Msg("Field client_email missing from service account keyfile")
	} else {
		if t, aok := saClientEmailIntfc.(string); !aok {
			log.Fatal().Msg("Field client_email not of type string")
		} else {
			saClientEmail = t
		}
	}

	log.Info().Msgf("Storing gcs credential %v on disk...", credentialsParam.Credentials)
	err = ioutil.WriteFile("/key-file.json", []byte(credential.AdditionalProperties.ServiceAccountKeyfile), 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed writing service account keyfile")
	}

	log.Info().Msg("Authenticating to google cloud")
	foundation.RunCommandWithArgs(ctx, "gcloud", []string{"auth", "activate-service-account", saClientEmail, "--key-file", "/key-file.json"})

	log.Info().Msg("Setting gcloud account")
	foundation.RunCommandWithArgs(ctx, "gcloud", []string{"config", "set", "account", saClientEmail})

	log.Info().Msg("Setting gcloud project")
	foundation.RunCommandWithArgs(ctx, "gcloud", []string{"config", "set", "project", credential.AdditionalProperties.Project})

	// clean up old stuff
	switch params.Action {
	case "copy":

		cpArgs := []string{}
		if params.Parallel != nil && *params.Parallel {
			cpArgs = append(cpArgs, "-m")
		}

		cpArgs = append(cpArgs, "cp")
		if params.ACL != "" {
			cpArgs = append(cpArgs, "-a", params.ACL)
		}
		if params.Compress != nil && *params.Compress {
			cpArgs = append(cpArgs, "-Z")
		}
		if params.Recursive != nil && *params.Recursive {
			cpArgs = append(cpArgs, "-r")
		}

		cpArgs = append(cpArgs, params.Source, fmt.Sprintf("gs://%v/%v", params.Bucket, params.Destination))

		foundation.RunCommandWithArgs(ctx, "gsutil", cpArgs)

		break
	}
}
