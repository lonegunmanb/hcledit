package cmd

import (
	"fmt"
	"io"

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

	return cmd
}

func runBlockMvCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	from := args[0]
	to := args[1]

	return editor.RenameBlock(cmd.InOrStdin(), cmd.OutOrStdout(), "-", from, to)
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
