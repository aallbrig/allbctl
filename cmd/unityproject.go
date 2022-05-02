package cmd

import (
	"fmt"
	"github.com/aallbrig/allbctl/pkg/osagnostic"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	FullscreenTemplateHTML = "https://gist.githubusercontent." +
		"com/aallbrig/2d07e3bbf03da818705db3215216e5cf/raw/752a534f7193cbd2c2b3a8929d5c0115d06adbb8/index.html"
	FullscreenTemplateJS = "https://gist.githubusercontent." +
		"com/aallbrig/2c243ce8b3d39bff2a0674744585d2e2/raw/a684ad3f108ede8a7e963300785967f3ed2c5a11/main.js"
	FullscreenTemplateCSS = "https://gist.githubusercontent." +
		"com/aallbrig/f51e371876df31830ef03c10bc192b50/raw/de8129c867c9e8007bf3227f6a02b1e6515fb5ba/style.css"
)

func newInitializeGitRepoCommand(path string) *exec.Cmd {
	return exec.Command("git", "-C", path, "init", "-b", "main")
}

func macFindUnityLTSExecutable() string {
	// Find latest version of unity in the 2020 class
	directoryInfo, _ := ioutil.ReadDir("/Applications/Unity/Hub/Editor")
	latest := ""
	for _, file := range directoryInfo {
		if file.IsDir() && strings.HasPrefix(file.Name(), "2020") && file.Name() > latest {
			latest = file.Name()
		}
	}
	return fmt.Sprintf("/Applications/Unity/Hub/Editor/%s/Unity.app/Contents/MacOS/Unity", latest)
}

func findUnityLTSExecutable() string {
	// TODO: its lame to have this hardcoded to MAC
	// it doesn't work on windows hosts nor *NIX or containers ☹️
	return macFindUnityLTSExecutable()
}

func newInitializeUnityProjectCommand(path string) *exec.Cmd {
	return exec.Command(findUnityLTSExecutable(), "-createProject", path, "-quit")
}

func newCurlCopy(url string, destinationFile string) *exec.Cmd {
	if _, err := os.Stat(destinationFile); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(destinationFile), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(destinationFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	return exec.Command("curl", url, "--output", destinationFile)
}

func newCopyUnityGitignore(path string) *exec.Cmd {
	webResource := "https://raw.githubusercontent.com/github/gitignore/main/Unity.gitignore"
	localDestination := filepath.Join(path, ".gitignore")
	return newCurlCopy(webResource, localDestination)
}

func newCopyFullscreenWebGLTemplateHTML(path string) *exec.Cmd {
	localDestination := filepath.Join(path, "Assets", "WebGLTemplates", "Fullscreen", "index.html")
	return newCurlCopy(FullscreenTemplateHTML, localDestination)
}
func newCopyFullscreenWebGLTemplateJS(path string) *exec.Cmd {
	localDestination := filepath.Join(path, "Assets", "WebGLTemplates", "Fullscreen", "TemplateData", "main.js")
	return newCurlCopy(FullscreenTemplateJS, localDestination)
}
func newCopyFullscreenWebGLTemplateCSS(path string) *exec.Cmd {
	localDestination := filepath.Join(path, "Assets", "WebGLTemplates", "Fullscreen", "TemplateData", "style.css")
	return newCurlCopy(FullscreenTemplateCSS, localDestination)
}

func copyFullscreenWebGLTemplate(path string) error {
	err := newCopyFullscreenWebGLTemplateHTML(path).Run()
	if err != nil {
		return err
	}
	err = newCopyFullscreenWebGLTemplateJS(path).Run()
	if err != nil {
		return err
	}
	err = newCopyFullscreenWebGLTemplateCSS(path).Run()
	return err
}

var projectNamePrompt = promptui.Prompt{
	Label: "Project Name",
	Validate: func(input string) error {
		return nil
	},
	Default: "new-unity-project",
}
var operatingSystem = osagnostic.NewOperatingSystem()
var projectName string
var ignoreUnityCommands bool
var installWebGLFullscreenTemplate bool

func NewUnityProjectCommand() *cobra.Command {
	var unityProjectCommand = &cobra.Command{
		Use:     "new-unity-project",
		Aliases: []string{},
		Short:   "Run new unity project setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current working directory
			// Check (or obtain) the name of the source code repository directory
			if projectName == "" {
				result, err := projectNamePrompt.Run()
				if err != nil {
					// TODO: What happens when an error is returned?
				}
				projectName = result
			}
			sourceCodePath := filepath.Join(operatingSystem.CurrentWorkingDirectory, projectName)
			// Make source code directory
			operatingSystem.CreateDirectory(sourceCodePath)

			if err := newInitializeGitRepoCommand(sourceCodePath).Run(); err != nil {
				return err
			}
			// Create a unity subdirectory
			unityProjectPath := filepath.Join(sourceCodePath, "unity", projectName)
			operatingSystem.CreateDirectory(unityProjectPath)

			if ignoreUnityCommands == false {
				// Run equivalent to
				// /Applications/Unity/Hub/Editor/2021.2.10f1/Unity.app/Contents/MacOS/Unity -createProject $(pwd)/unity/ + projectName
				if err := newInitializeUnityProjectCommand(unityProjectPath).Run(); err != nil {
					return err
				}
			}

			if err := newCopyUnityGitignore(unityProjectPath).Run(); err != nil {
				return err
			}

			if installWebGLFullscreenTemplate {
				if err := copyFullscreenWebGLTemplate(unityProjectPath); err != nil {
					return err
				}
			}
			// if UNITY_LICENSE envvar is set
			// gh secret set UNITY_LICENSE --body "${UNITY_LICENSE}"
			// if UNITY_EMAIL envvar is set
			// gh secret set UNITY_EMAIL --body "${UNITY_EMAIL}"
			// if UNITY_PASSWORD envvar is set
			// gh secret set UNITY_PASSWORD --body "${UNITY_PASSWORD}"
			// curl activate-unity-license.yml from gist to .github/activate-unity-license.yml
			// https://gist.githubusercontent.
			//com/aallbrig/915341c99b9f73f03c922a7f94e47041/raw/1cd37fc5930c8d61622d8c96f78245983b6caf8b/activate-unity-license.yml
			// curl unity.yml from gist to .github/unity.yml
			// https://gist.githubusercontent.com/aallbrig/c54066dfcb6e2cd527c9313f396c7f48/raw/7f5d3397db772aa05ad84901b9aaadfd5150bcb4/unity.yml

			// if this is a webGL project
			// make unity project subdirectory to hold webGL template
			// (optional) ask for google analytics ID
			// curl index.html down from webGL fullscreen gist
			// https://gist.githubusercontent.com/aallbrig/2d07e3bbf03da818705db3215216e5cf/raw/752a534f7193cbd2c2b3a8929d5c0115d06adbb8/index.html
			// curl style.css down from webGL fullscreen gist
			// https://gist.githubusercontent.com/aallbrig/f51e371876df31830ef03c10bc192b50/raw/de8129c867c9e8007bf3227f6a02b1e6515fb5ba/style.css
			// curl main.js down from webGL fullscreen gist
			// https://gist.githubusercontent.com/aallbrig/2c243ce8b3d39bff2a0674744585d2e2/raw/a684ad3f108ede8a7e963300785967f3ed2c5a11/main.js
			return nil
		},
	}
	unityProjectCommand.Flags().StringVar(
		&projectName,
		"project",
		"",
		"Name of both the source code directory and unity project",
	)
	unityProjectCommand.Flags().BoolVar(
		&ignoreUnityCommands,
		"ignore-unity-commands",
		false,
		"Optional flag to disable running unity commands",
	)
	unityProjectCommand.Flags().BoolVar(
		&installWebGLFullscreenTemplate,
		"install-webgl-fullscreen-template",
		false,
		"Optional flag to install a sensible fullscreen WebGL template",
	)
	return unityProjectCommand
}
