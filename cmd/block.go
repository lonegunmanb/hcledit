package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newBlockCmd())
}

func newBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block",
		Short: "Edit block",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newBlockGetCmd(),
		newBlockMvCmd(),
		newBlockListCmd(),
	)

	return cmd
}

func newBlockGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ADDRESS>",
		Short: "Get block",
		Long: `Get matched blocks at a given address

Arguments:
  ADDRESS          An address of block to get.
`,
		RunE: runBlockGetCmd,
	}

	return cmd
}

func runBlockGetCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 argument, but got %d arguments", len(args))
	}

	address := args[0]

	return editor.GetBlock(cmd.InOrStdin(), cmd.OutOrStdout(), "-", address)
}

func newBlockMvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mv <FROM_ADDRESS> <TO_ADDRESS>",
		Short: "Move block (Rename block type and labels)",
		Long: `Move block (Rename block type and labels)

Arguments:
  FROM_ADDRESS     An old address of block.
  TO_ADDRESS       A new address of block.
`,
		RunE: runBlockMvCmd,
	}

	flags := cmd.Flags()
	flags.StringP("file", "f", "", "A path of input file")
	flags.BoolP("write", "w", false, "Overwrite input file (default: false)")
	return cmd
}

func runBlockMvCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	from := args[0]
	to := args[1]

	filename, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	fs := afero.NewOsFs()
	var inStream io.Reader
	if len(filename) != 0 {
		file, err := fs.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open file: %s", err)
		}
		defer file.Close()
		inStream = file
	} else {
		inStream = cmd.InOrStdin()
		filename = "-"
	}

	write, err := cmd.Flags().GetBool("write")
	if err != nil {
		return err
	}
	var outStream io.Writer
	if write {
		if len(filename) == 0 {
			return errors.New("when using write option, a file name is requreid")
		}
		outStream = new(bytes.Buffer)
	} else {
		outStream = cmd.OutOrStdout()
	}
	if err := editor.RenameBlock(inStream, outStream, filename, from, to); err != nil {
		return err
	}
	if !write {
		return nil
	}

	if err = afero.WriteFile(fs, filename, outStream.(*bytes.Buffer).Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}

func newBlockListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List block",
		RunE:  runBlockListCmd,
	}

	flags := cmd.Flags()
	flags.StringP("file", "f", "", "A path of input file")
	return cmd
}

func runBlockListCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("expected 0 argument, but got %d arguments", len(args))
	}

	filename, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	fs := afero.NewOsFs()
	var inStream io.Reader
	if len(filename) != 0 {
		file, err := fs.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open file: %s", err)
		}
		defer file.Close()
		inStream = file
	} else {
		inStream = cmd.InOrStdin()
		filename = "-"
	}

	return editor.ListBlock(inStream, cmd.OutOrStdout(), filename)
}
