package power

import (
	"get.porter.sh/porter/pkg/exec/builder"
)

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []Step // using UnmarshalYAML so that we don't need a custom type per action
}

// MarshalYAML converts the action back to a YAML representation
// install:
//   power:
//     ...
func (a Action) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{a.Name: a.Steps}, nil
}

// MakeSteps builds a slice of Step for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]Step{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - power: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			step := result.(*[]Step)
			a.Steps = append(a.Steps, *step...)
		}
		break // There is only 1 action
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	// Go doesn't have generics, nothing to see here...
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

type Step struct {
	Instruction `yaml:"power"`
}

var _ builder.ExecutableStep = Step{}
var _ builder.StepWithOutputs = Step{}
var _ builder.SuppressesOutput = Step{}

// Actions is a set of actions, and the steps, passed from Porter.
type Actions []Action

// UnmarshalYAML takes chunks of a porter.yaml file associated with this mixin
// and populates it on the current action set.
// install:
//   power:
//     ...
//   power:
//     ...
// upgrade:
//   power:
//     ...
func (a *Actions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, Action{})
	if err != nil {
		return err
	}

	for actionName, action := range results {
		for _, result := range action {
			s := result.(*[]Step)
			*a = append(*a, Action{
				Name:  actionName,
				Steps: *s,
			})
		}
	}
	return nil
}

type Instruction struct {
	Description    string        `yaml:"description"`
	Service        string        `yaml:"group"`
	Operation      string        `yaml:"operation"`
	Arguments      []string      `yaml:"arguments,omitempty"`
	Flags          builder.Flags `yaml:"flags,omitempty"`
	Outputs        []Output      `yaml:"outputs,omitempty"`
	SuppressOutput bool          `yaml:"suppress-output,omitempty"`
}

func (s Step) GetCommand() string {
	return "Microsoft.CompositeMixin.Power.Deployment.Client"
}

func (s Step) GetWorkingDir() string {
	return ""
}

func (s Step) GetArguments() []string {
	args := make([]string, 0, len(s.Arguments)+2)

	// Specify the Service and Operation
	args = append(args, s.Service)
	args = append(args, s.Operation)

	// Append the positional arguments
	args = append(args, s.Arguments...)

	return args
}

func (s Step) GetFlags() builder.Flags {
	// Always request json formatted output
	return append(s.Flags, builder.NewFlag("output", "json"))
}

func (s Step) GetOutputs() []builder.Output {
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}

func (s Step) SuppressesOutput() bool {
	return s.SuppressOutput
}

var _ builder.OutputJsonPath = Output{}

type Output struct {
	Name     string `yaml:"name"`
	JsonPath string `yaml:"jsonPath"`
}

func (o Output) GetName() string {
	return o.Name
}

func (o Output) GetJsonPath() string {
	return o.JsonPath
}
