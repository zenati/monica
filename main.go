package main

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"bytes"
	"os/exec"
	"strings"
)

/*
  Config
  Structure mirroring the format of a valid .monica.yml file.
  Consists of an engine name, associated version and an array of reactions.
*/
type Config struct {
	Engine string
	Reactions []Reaction
}

/*
  Reaction
  Structure mirroring the format of a valid reaction if a config file.
  Consists of a name and an array of reaction commands.
*/
type Reaction struct {
	Name string
	Desc string
	Content []ReactionCommand
	Arguments []ReactionArgument
}

/*
  ReactionArgument
*/
type ReactionArgument struct {
	Name string
	Flag *string
}

/*
  ConfigArguments
*/
type ConfigArguments struct {
	Name string
	Arguments []string
}

/*
  ReactionCommand
  Structure mirroring the format of a valid command if a config file.
  Consists of a type, path, command, source, destination, variable,
  path and a value.
*/
type ReactionCommand struct {
	Command string
	ReactionName string `yaml:"reaction"`
}

func main() {
	config := unmarshalConfig()
	kingpin.Flag("debug", "Enable debug mode.").Bool()
	kingpin.CommandLine.HelpFlag.Short('h')

	processConfig(&config)
	kingpin.Version("0.0.1")

	chosenReaction := kingpin.Parse()
	processReactions(&config, &chosenReaction)
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
		text("File .monica.yml not detected.", color.FgRed)
		os.Exit(0)
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		text(err.Error(), color.FgRed)
		os.Exit(0)
	}

	return config
}

/*
  processConfig
  Takes a Config pointer in argument and loops through the list
  of reactions and commands, executing one after another in a
  thread safe executeCommand function.
*/
func processConfig(config *Config) {
	for i := 0; i < len(config.Reactions); i++ {
		reaction := &config.Reactions[i]
		argsList := extractArguments(&reaction.Content)
		cmdFlags := kingpin.Command(reaction.Name, reaction.Desc)

		for j := 0; j < len(argsList); j++ {
			cmdFlag := cmdFlags.Flag(argsList[j], "").Short(argsList[j][0]).Required().String()

			argument := ReactionArgument{}
			argument.Name = argsList[j]
			argument.Flag = cmdFlag

			reaction.Arguments = append(reaction.Arguments, argument)
		}
	}
}

/*
  extractArguments
*/
func extractArguments(reactionCommands *[]ReactionCommand) []string {
	var arguments []string

	for index := 0; index < len(*reactionCommands); index++ {
		command := (*reactionCommands)[index].Command

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
  processReactions
*/
func processReactions(config *Config, reaction *string) {
	for index := 0; index < len(config.Reactions); index++ {
		if *reaction == config.Reactions[index].Name {
			processReaction(&config.Reactions[index])
		}
	}
}

/*
  processReaction
  Takes a Reaction as a parameter
*/
func processReaction(reaction *Reaction) {
	text(fmt.Sprintf("executing: %s", reaction.Name), color.FgGreen)

	for j := 0; j < len(reaction.Content); j++ {
		processCommand(reaction, j)
	}
}

/*
  processCommand
  Takes a Reaction as a parameter
*/
func processCommand(reaction *Reaction, index int) {
	command := reaction.Content[index].Command

	for j := 0; j < len(reaction.Arguments); j++ {
		varName := reaction.Arguments[j].Name
		varValue := reaction.Arguments[j].Flag
		varToChange := fmt.Sprintf("${%s}", varName)

		command = strings.Replace(command, varToChange, *varValue, -1)
	}

	coloredContent := fmt.Sprintf("\t-> %s", command)
	text(coloredContent, color.FgGreen)

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

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		text(stderr.String(), color.FgRed)
		os.Exit(0)
	}
}

/*
  Map
  Returns a new slice containing the results of applying the function f
  to each string in the original slice.
*/
func Map(vs []string, f func(string, int) string) []string {
	vsm := make([]string, len(vs))

	for i, v := range vs {
		vsm[i] = f(v, i)
	}

	return vsm
}

/*
  prefix
  Displays a prefix to all engine related messages
*/
func prefix() string {
	return fmt.Sprintf("monica")
}

/*
  text
  Displays a message on the screen using a particular color
*/
func text(content string, attribute color.Attribute, returnOperator ...bool) {
	returnLine := true
	var printfContent string

	if len(returnOperator) > 0 {
		returnLine = returnOperator[0]
	}

	if returnLine {
		printfContent = "%s %s\n"
	} else {
		printfContent = "\r%s %s"
	}

	fmt.Printf(printfContent, colored(prefix(), attribute), content)
}

/*
  colored
  Displays a message on the screen using a particular color
*/
func colored(text string, attribute color.Attribute) string {
	return color.New(attribute).SprintFunc()(text)
}
