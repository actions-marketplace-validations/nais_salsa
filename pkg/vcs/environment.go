package vcs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Environment struct {
	GitHubContext      GitHubContext       `json:"github"`
	Event              *Event              `json:"event,omitempty"`
	RunnerContext      RunnerContext       `json:"context"`
	CurrentEnvironment *CurrentEnvironment `json:"env,omitempty"`
}

func CreateCIEnvironment(githubContext, runnerContext, envsContext *string) (*Environment, error) {
	// Required when creating CI environment
	if len(*githubContext) == 0 || len(*runnerContext) == 0 {
		return nil, nil
	}
	env := Environment{}

	if err := ParseGithub(githubContext, &env); err != nil {
		return nil, fmt.Errorf("parsing github: %w", err)
	}

	if err := ParseRunner(runnerContext, &env); err != nil {
		return nil, fmt.Errorf("parsing runner: %w", err)
	}

	if err := ParseEnv(envsContext, &env); err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return &env, nil
}

func (in *Environment) RepoUri() string {
	return fmt.Sprintf("%s/%s", in.GitHubContext.ServerUrl, in.GitHubContext.Repository)
}

func (in *Environment) BuildInvocationId() string {
	return fmt.Sprintf("%s/actions/runs/%s", in.RepoUri(), in.GitHubContext.RunId)
}

func (in *Environment) GithubSha() string {
	return in.GitHubContext.SHA
}

func (in *Environment) BuilderId() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return in.RepoUri() + GitHubHostedIdSuffix
	}
	return in.RepoUri() + GitHubHostedIdSuffix
}

func (in *Environment) EventInputs() json.RawMessage {
	if in.GitHubContext.EventName != "workflow_dispatch" {
		return nil
	}
	return in.Event.Inputs
}

func (in *Environment) FilteredEnvironment() *Environment {
	// Should also contain environment variables.
	// These are always set because it is not possible
	// to know whether they were referenced or not.
	return &Environment{
		GitHubContext: GitHubContext{
			RunId: in.GitHubContext.RunId,
		},
		Event:              nil,
		RunnerContext:      in.RunnerContext,
		CurrentEnvironment: in.CurrentEnvironment,
	}
}
