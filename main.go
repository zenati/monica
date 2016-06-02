package main

import (
	"bytes"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

/*
  Config
  Structure mirroring the format of a valid .monica.yml file.
  Consists of an array of actions.
*/
type Config struct {
	Actions []Action
}

/*
  Action
  Structure mirroring the format of a valid action if a config file.
  Consists of a name and an array of action commands.
*/
type Action struct {
	Name      string
	Desc      string
	Content   []ActionContent
	Default   []map[string]string
	Arguments []ActionArgument
}

/*
  ActionArgument
*/
type ActionArgument struct {
	Name string
	Flag *string
}

/*
  ConfigArguments
*/
type ConfigArguments struct {
	Name      string
	Arguments []string
}

/*
  ActionContent
  Structure mirroring the format of a valid command if a config file.
  Consists of a type, path, command, source, destination, variable,
  path and a value.
*/
type ActionContent struct {
	Action  string
	Command string
}

func main() {
	config := unmarshalConfig()
	kingpin.Flag("debug", "Enable debug mode.").Bool()
	kingpin.CommandLine.HelpFlag.Short('h')

	processConfig(&config)
	kingpin.Version("0.1")

	chosenAction := kingpin.Parse()
	processActions(&config, &chosenAction)
}

/*
  unmarshalConfig
  Reads the .monica.yml config file and extracts content
  to the Config struct defined above. One extracted, content is parsed
  and falls through the execution process.
*/
func unmarshalConfig() Config {
	config := Config{}
	content, err := ioutil.ReadFile(".monica.yml")

	if err != nil {
		fmt.Println("File .monica.yml not detected.")
		os.Exit(0)
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	return config
}

/*
  processConfig
  Takes a Config pointer in argument and loops through the list
  of actions and commands, executing one after another in a
  thread safe executeCommand function.
*/
func processConfig(config *Config) {
	for i := 0; i < len(config.Actions); i++ {
		var cmdFlag *string

		action := &config.Actions[i]
		cmdFlags := kingpin.Command(action.Name, action.Desc)

		argsList := extractArguments(&action.Content)
		defsList := extractDefaults(&action.Default)

		for j := 0; j < len(argsList); j++ {
			if defs, exists := defsList[argsList[j]]; exists {
				cmdFlag = cmdFlags.Flag(argsList[j], "").Default(defs).String()
			} else {
				cmdFlag = cmdFlags.Flag(argsList[j], "").Required().String()
			}

			argument := ActionArgument{}
			argument.Name = argsList[j]
			argument.Flag = cmdFlag

			action.Arguments = append(action.Arguments, argument)
		}
	}
}

/*
	extractDefaults
*/
func extractDefaults(actionDefault *[]map[string]string) map[string]string {
	defaults := map[string]string{}

	for _, mapData := range *actionDefault {
		for key, value := range mapData {
			defaults[key] = value
		}
	}

	return defaults
}

/*
  extractArguments
*/
func extractArguments(actionContent *[]ActionContent) []string {
	var arguments []string

	for index := 0; index < len(*actionContent); index++ {
		command := (*actionContent)[index].Command

		re := regexp.MustCompile(`\$\{([^}]+)\}`)
		match := re.FindAllStringSubmatch(command, -1)

		for j := 0; j < len(match); j++ {
			arguments = appendIfMissing(arguments, match[j][1])
		}
	}

	return arguments
}

/*
  appendIfMissing
*/
func appendIfMissing(data []string, i string) []string {
	for _, element := range data {
		if element == i {
			return data
		}
	}

	return append(data, i)
}

/*
  processActions
*/
func processActions(config *Config, action *string) {
	for index := 0; index < len(config.Actions); index++ {
		if *action == config.Actions[index].Name {
			processAction(&config.Actions[index])
		}
	}
}

/*
  processAction
  Takes a Action as a parameter
*/
func processAction(action *Action) {
	fmt.Printf("executing: %s\n", action.Name)

	for j := 0; j < len(action.Content); j++ {
		processCommand(action, j)
	}
}

/*
  processCommand
  Takes a Action as a parameter
*/
func processCommand(action *Action, index int) {
	command := action.Content[index].Command

	for j := 0; j < len(action.Arguments); j++ {
		varName := action.Arguments[j].Name
		varValue := action.Arguments[j].Flag
		varToChange := fmt.Sprintf("${%s}", varName)

		command = strings.Replace(command, varToChange, *varValue, -1)
	}

	fmt.Printf("-> %s\n", command)

	executableCommand := strings.Split(command, " ")
	executeCommand(executableCommand[0], executableCommand[1:]...)
}

/*
  executeCommand
  Executes a kernel thread safe command with associated arguments
  defined as a vector of infinite sub-components. This displays the
  stdout in case the debug mode is enabled, and omit otherwize.
*/
func executeCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		fmt.Println(stderr.String())
		os.Exit(0)
	} else if (stdout.String() != "") {
		fmt.Print(stdout.String())
	}
}
