package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/nais/salsa/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AttestOptions struct {
	Key           string `mapstructure:"key"`
	NoUpload      bool   `mapstructure:"no-upload"`
	RekorURL      string `mapstructure:"rekor-url"`
	PredicateFile string `mapstructure:"predicate"`
	PredicateType string `mapstructure:"type"`
}

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "sign and upload in-toto attestation",
	RunE: func(cmd *cobra.Command, args []string) error {
		var attest AttestOptions
		err := viper.Unmarshal(&attest)
		if err != nil {
			return err
		}

		if PathFlags.Repo == "" {
			return errors.New("repo name must be specified")
		}

		s := utils.StartSpinner(fmt.Sprintf("attestation for %s finished", PathFlags.Repo))
		out, err := attest.Exec(args)
		if err != nil {
			return err
		}
		log.Infof("finished attestation %s\n", out)
		s.Stop()
		path := PathFlags.RepoDir

		file := path + "/" + attest.PredicateFile + ".att"
		err = os.WriteFile(file, []byte(out), os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("could not write file %s %w", file, err)
		}
		return nil
	},
}

func (o AttestOptions) Exec(a []string) (string, error) {
	err := utils.RequireCommand("cosign")
	if err != nil {
		return "", err
	}
	args := []string{
		"attest",
		"--type", o.PredicateType,
		"--predicate", o.PredicateFile,
		"--key", o.Key,
		"--rekor-url", o.RekorURL,
	}
	if o.NoUpload {
		args = append(args, "--no-upload")
	}
	args = append(args, a...)

	cmd := exec.Command(
		"cosign",
		args...,
	)

	cmd.Dir = PathFlags.WorkDir()
	return utils.Exec(cmd)
}

func init() {
	rootCmd.AddCommand(attestCmd)
	attestCmd.Flags().String("key", "",
		"path to the private key file, KMS URI or Kubernetes Secret")
	//cmd.Flags().
	attestCmd.Flags().Bool("no-upload", false,
		"do not upload the generated attestation")
	attestCmd.Flags().String("rekor-url", "https://rekor.sigstore.dev",
		"address of transparency log")
	attestCmd.Flags().String("predicate", "",
		"the predicate file used for attestation")
	attestCmd.Flags().String("type", "slsaprovenance",
		"specify a predicate type (slsaprovenance|link|spdx|custom) or an URI (default \"slsaprovenance\")\n")
	viper.BindPFlags(attestCmd.Flags())
}
