package command

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"go.mercari.io/hcledit"
)

type CreateOptions struct {
	Type  string
	After string
}

func NewCmdCreate() *cobra.Command {
	opts := &CreateOptions{}
	cmd := &cobra.Command{
		Use:   "create <query> <value> <file>",
		Short: "Create a new field",
		Long:  `Runs an address query on a hcl file and create new field with given value.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := runCreate(opts, args); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Type, "type", "t", "string", "Type of the value")
	cmd.Flags().StringVarP(&opts.After, "after", "a", "", "Field key which before the value will be created")
	return cmd
}

func runCreate(opts *CreateOptions, args []string) error {
	query := args[0]
	valueStr := args[1]
	filePath := args[2]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	value, err := convert(valueStr, opts.Type)
	if err != nil {
		return fmt.Errorf("failed to convert input to specific type: %s", err)
	}

	if err := editor.Create(query, value, hcledit.WithAfter(opts.After)); err != nil {
		return fmt.Errorf("failed to create: %s", err)
	}

	return editor.OverWriteFile()
}

func convert(inputStr, typeStr string) (interface{}, error) {
	switch typeStr {
	case "string":
		return inputStr, nil
	case "int":
		return strconv.Atoi(inputStr)
	case "bool":
		return strconv.ParseBool(inputStr)
	case "raw":
		return hcledit.RawVal(inputStr), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeStr)
	}
}
