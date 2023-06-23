// TODO: Support matrix syntax ? https://github.com/docker/buildx/pull/1690
// Limitations (or features ?): Only searches in subfolders of the modules folder, not recursively
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zclconf/go-cty/cty"
)

type DockerBakeTarget struct {
	name              string
	module            string
	context           string
	dockerfile        string
	dockerfilePurpose string
}

func generateDockerBakeHCL(username string, registryPrefix string, targetsStruct []DockerBakeTarget) {

	f := hclwrite.NewEmptyFile()

	bakeFile, err := os.Create("docker-bake.hcl")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create variables block
	variablesBlock := f.Body().AppendNewBlock("variable", []string{"DOCKER_USERNAME"})
	variablesBlock.Body().SetAttributeValue("default", cty.StringVal(username))

	variablesBlock = f.Body().AppendNewBlock("variable", []string{"DOCKER_REGISTRY_PREFIX"})
	variablesBlock.Body().SetAttributeValue("default", cty.StringVal(registryPrefix))

	variablesBlock = f.Body().AppendNewBlock("variable", []string{"TAG"})
	variablesBlock.Body().SetAttributeValue("default", cty.StringVal("latest"))

	groupName := fmt.Sprintf("%s-modules", registryPrefix)
	groupBlock := f.Body().AppendNewBlock("group", []string{groupName})

	targets := make([]cty.Value, 0, len(targetsStruct)) // Computed below

	// Create targets block

	for _, target := range targetsStruct {

		targets = append(targets, cty.StringVal(target.name))

		targetBlock := f.Body().AppendNewBlock("target", []string{target.name})
		targetBlock.Body().SetAttributeValue("dockerfile", cty.StringVal(target.dockerfile))
		targetBlock.Body().SetAttributeValue("context", cty.StringVal(target.context))

		platforms := []cty.Value{cty.StringVal("linux/amd64"), cty.StringVal("linux/arm64/v8")}
		targetBlock.Body().SetAttributeValue("platforms", cty.ListVal(platforms))

		tag := ":${TAG}"
		if target.dockerfilePurpose != "" {
			tag = tag + "-" + target.dockerfilePurpose
		}
		image := "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-" + target.module + tag
		tokens := hclwrite.Tokens{
			{Type: hclsyntax.TokenOQuote, Bytes: []byte(`"`)},
			{Type: hclsyntax.TokenQuotedLit, Bytes: []byte(image)},
			{Type: hclsyntax.TokenCQuote, Bytes: []byte(`"`)},
		}
		targetBlock.Body().SetAttributeRaw("tags", tokens)

	}

	groupBlock.Body().SetAttributeValue("targets", cty.ListVal(targets))

	// Write the generated HCL to a file
	bakeFile.Write(f.Bytes())

}

func parseModules(modulesPath string) map[string][]string {

	entries, err := os.ReadDir(modulesPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to read modules directory")
	}

	targetsDict := make(map[string][]string)

	for _, e := range entries {

		if e.IsDir() {
			moduleName := e.Name()

			log.Info().
				Str("folder", moduleName).
				Msg("Exploring folder")

			// Find Dockerfile(s) within the module directory
			modulePath := filepath.Join(modulesPath, moduleName)
			pattern := "(^Dockerfile(\\.[[:alnum:]_-]+)*$)|(^(([[:alnum:]_-]+)\\.)*Dockerfile$)"
			r, _ := regexp.Compile(pattern)
			dockerfiles, err := filepath.Glob(filepath.Join(modulePath, "*Dockerfile*"))
			if err != nil || dockerfiles == nil {
				log.Debug().
					Str("folder", moduleName).
					Err(err).
					Msg("Failed to find Dockerfile(s) in folder")
				continue
			}

			validDockerfiles := make([]string, 0, len(dockerfiles))
			for _, d := range dockerfiles {
				dName := filepath.Base(d)

				if r.MatchString(dName) {
					validDockerfiles = append(validDockerfiles, dName)

					log.Debug().
						Str("dockerfile", dName).
						Str("folder", moduleName).
						Msgf("Found Dockerfile in folder")
				}
			}

			// Add Dockerfile paths to the array
			targetsDict[modulePath] = validDockerfiles
		}
	}
	return targetsDict
}

func computeTargets(targetsDict map[string][]string) []DockerBakeTarget {

	DockerBakeTargets := make([]DockerBakeTarget, 0, len(targetsDict))

	for key, value := range targetsDict {

		modulePath := key
		moduleName := filepath.Base(modulePath)

		for _, dockerfile := range value {

			d := strings.Split(dockerfile, ".")

			targetName := moduleName
			purpose := ""

			if len(d) > 1 && d[0] == "Dockerfile" {
				purpose = strings.Join(d[1:], "-")
				targetName = targetName + "-" + purpose
			} else if len(d) > 1 && d[len(d)-1] == "Dockerfile" {
				purpose = strings.Join(d[:len(d)-1], "-")
				targetName = targetName + "-" + purpose
			}

			DockerBakeTarget := DockerBakeTarget{
				name:              targetName,
				module:            moduleName,
				context:           modulePath,
				dockerfile:        dockerfile,
				dockerfilePurpose: purpose,
			}

			DockerBakeTargets = append(DockerBakeTargets, DockerBakeTarget)
		}
	}
	return DockerBakeTargets
}

func InitLogger(logLevel string) {

	switch logLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

func main() {

	usernamePtr := flag.String("username", "", "Docker username")
	registryPrefixPtr := flag.String("registry_prefix", "", "Docker registry prefix")
	modulesPathPtr := flag.String("modules_path", "", "Path to modules directory")
	logLevelPtr := flag.String("log", "info", "Log level")
	flag.Parse()
	username := *usernamePtr
	registryPrefix := *registryPrefixPtr
	modulesPath := *modulesPathPtr
	if username == "" || registryPrefix == "" || modulesPath == "" {
		log.Fatal().
			Msg("Missing required arguments")
	}

	InitLogger(*logLevelPtr)

	log.Info().
		Str("folder", modulesPath).
		Msg("Parsing possible targets in folder")

	targetsDockerfiles := parseModules(modulesPath)
	DockerBakeTargets := computeTargets(targetsDockerfiles)

	log.Printf("Found %d targets across %d modules", len(DockerBakeTargets), len(targetsDockerfiles))

	log.Info().
		Int("targets", len(DockerBakeTargets)).
		Int("modules", len(targetsDockerfiles)).
		Msg("Found targets in explored folders")

	generateDockerBakeHCL(username, registryPrefix, DockerBakeTargets)

	log.Info().
		Msg("Generated DockerBake HCL file")
}
