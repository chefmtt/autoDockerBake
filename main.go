package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func generateDockerBakeHCL(username string, registryPrefix string, targetsDict map[string][]string) {
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
	groupBlock.Body().SetAttributeValue("targets", cty.StringVal("docker-bake"))

	// Create targets block

	for key, value := range targetsDict {
		for _, dockerfile := range value {

			d := strings.Split(dockerfile, ".")

			targetName := ""

			if len(d) > 1 {
				targetName = fmt.Sprintf("%s_%s", key, strings.Join(d[1:], "_"))
			} else {
				targetName = key
			}

			targetBlock := f.Body().AppendNewBlock("target", []string{targetName})
			targetBlock.Body().SetAttributeValue("dockerfile", cty.StringVal(dockerfile))
			targetBlock.Body().SetAttributeValue("context", cty.StringVal("./modules/message_monitoring"))

			platforms := []cty.Value{cty.StringVal("linux/amd64"), cty.StringVal("linux/arm64/v8")}
			targetBlock.Body().SetAttributeValue("platforms", cty.ListVal(platforms))

			image := "${DOCKER_USERNAME}/${DOCKER_REGISTRY_PREFIX}-message_monitoring:${TAG}"
			tags := []cty.Value{cty.StringVal(image)}
			targetBlock.Body().SetAttributeValue("tags", cty.ListVal(tags))
		}
	}

	// Set additional attributes or blocks as needed

	// Write the generated HCL to a file
	bakeFile.Write(f.Bytes())

}

func parseModules(modulesPath string) map[string][]string {

	entries, err := ioutil.ReadDir(modulesPath)
	if err != nil {
		log.Fatal(err)
	}

	targetsDict := make(map[string][]string)

	for _, e := range entries {

		if e.IsDir() {
			moduleName := e.Name()

			// Find Dockerfile(s) within the module directory
			modulePath := filepath.Join(modulesPath, moduleName)
			dockerfiles, err := filepath.Glob(filepath.Join(modulePath, "Dockerfile*"))
			if err != nil || dockerfiles == nil {
				log.Printf("Failed to find Dockerfile(s) in module '%s': %v", moduleName, err)
				continue
			}

			for i, d := range dockerfiles {
				dockerfiles[i] = filepath.Base(d)
			}

			// Add Dockerfile paths to the array
			targetsDict[moduleName] = dockerfiles
		}
	}
	return targetsDict
}

func main() {

	username := os.Args[1]
	registryPrefix := os.Args[2]
	modulesPath := "./modules"

	targetsDict := parseModules(modulesPath)

	generateDockerBakeHCL(username, registryPrefix, targetsDict)
}
