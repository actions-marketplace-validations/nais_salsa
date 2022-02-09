package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan/common"
	"github.com/nais/salsa/pkg/scan/deps"
	"github.com/nais/salsa/pkg/utils"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var project string
var inputContext string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan files and dependencies for a given project",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		if len(args) == 1 {
			artifact = args[0]
		}

		if PathFlags.Repo == "" {
			return errors.New("repo name must be specified")
		}

		log.Infof("prepare to scan path %s for project %s...", PathFlags.WorkDir(), project)
		workDir := PathFlags.WorkDir()
		deps, err := deps.Dependencies(workDir)
		err = GenerateProvenance(workDir, PathFlags.Repo, deps, &inputContext)
		if err != nil {
			return err
		}
		return nil
	},
}

func GenerateProvenance(workDir, project string, dependencies *common.ArtifactDependencies, inputContext *string) error {
	context, err := vcs.CreateCIContext(inputContext)
	if err != nil {
		return err
	}
	provenanceArtifact := intoto.CreateProvenanceArtifact(project, dependencies).WithRunnerContext(context)
	predicate := intoto.GenerateSlsaPredicate(provenanceArtifact)
	statement, err := json.Marshal(predicate)
	if err != nil {
		return fmt.Errorf("marshal: %v\n", err)
	}

	provenanceFileName := utils.ProvenanceFile(project)
	err = os.WriteFile(workDir+"/"+provenanceFileName, statement, 0644)
	if err != nil {
		return fmt.Errorf("write to file: %v\n", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&inputContext, "context", "", "context of build environment")
}
