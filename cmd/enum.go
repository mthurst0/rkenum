package cmd

import (
	"fmt"
	"strings"

	"github.com/mthurst0/rkutil/rkstrings"
	"github.com/spf13/cobra"
)

func GenerateEnum(pkg string, enumName string, values []string, noUnknown bool) string {
	var sb strings.Builder
	sb.WriteString("// Generated code - DO NOT EDIT\n\n")
	sb.WriteString(fmt.Sprintf("package %s\n\n", pkg))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"strings\"\n")
	sb.WriteString(")\n\n")
	sb.WriteString(fmt.Sprintf("type %s int\n\n", enumName))
	sb.WriteString("const (\n")
	if !noUnknown {
		sb.WriteString(fmt.Sprintf("\t%sUnknown = %s(iota)\n", enumName, enumName))
		for _, v := range values {
			v = rkstrings.ToCamelCase(v)
			sb.WriteString(fmt.Sprintf("\t%s%s\n", enumName, v))
		}
	} else {
		firstValue := true
		for _, v := range values {
			v = rkstrings.ToCamelCase(v)
			if firstValue {
				sb.WriteString(fmt.Sprintf("\t%s%s = %s(iota)\n", enumName, v, enumName))
			} else {
				sb.WriteString(fmt.Sprintf("\t%s%s\n", enumName, v))
			}
			firstValue = false
		}
	}
	sb.WriteString(fmt.Sprintf("\t%sMax\n", enumName))
	sb.WriteString(")\n\n")

	sb.WriteString(fmt.Sprintf("func Parse%s(s string) (%s, error) {\n", enumName, enumName))
	sb.WriteString(fmt.Sprintf("\tswitch strings.ToLower(s) {\n"))
	for _, v := range values {
		sb.WriteString(fmt.Sprintf("\tcase \"%s\":\n", v))
		sb.WriteString(fmt.Sprintf("\t\treturn %s%s, nil\n", enumName, rkstrings.ToCamelCase(v)))
	}
	sb.WriteString(fmt.Sprint("\tdefault:\n"))
	sb.WriteString(fmt.Sprintf("\t\treturn %s(0), fmt.Errorf(\"could not convert string to %s: %%s\", s)\n", enumName, enumName))
	sb.WriteString(fmt.Sprint("\t}\n"))
	sb.WriteString(fmt.Sprint("}\n\n"))

	sb.WriteString(fmt.Sprintf("func (v %s) String() string {\n", enumName))
	sb.WriteString(fmt.Sprintf("\tswitch v {\n"))
	for _, v := range values {
		sb.WriteString(fmt.Sprintf("\tcase %s%s:\n", enumName, rkstrings.ToCamelCase(v)))
		sb.WriteString(fmt.Sprintf("\t\treturn \"%s\"\n", v))
	}
	sb.WriteString(fmt.Sprint("\tdefault:\n"))
	sb.WriteString(fmt.Sprintf("\t\treturn fmt.Sprintf(\"%s(%%d)\", v)\n", enumName))
	sb.WriteString(fmt.Sprint("\t}\n"))
	sb.WriteString(fmt.Sprint("}\n"))

	return sb.String()
}

func parseStrings(values []string) []string {
	result := make([]string, 0)
	for _, value := range values {
		for _, v1 := range strings.Split(value, ",") {
			v1 = strings.TrimSpace(v1)
			for _, v2 := range strings.Fields(v1) {
				result = append(result, v2)
			}
		}
	}
	return result
}

func createEnumCmd() *cobra.Command {
	var packageName string
	var enumName string
	var enumValues []string
	var noUnknown bool
	var enumCmd = &cobra.Command{
		Use:   "enum",
		Short: "Create an enum",
		RunE: func(cmd *cobra.Command, args []string) error {
			if packageName == "" {
				return fmt.Errorf("package name must be set")
			}
			if enumName == "" {
				return fmt.Errorf("enum name must be set")
			}
			if len(enumValues) == 0 {
				return fmt.Errorf("enum values must be set")
			}
			fmt.Println(GenerateEnum(packageName, enumName, parseStrings(enumValues), noUnknown))
			return nil
		},
	}
	enumCmd.PersistentFlags().StringVarP(
		&packageName, "package", "p", "", "Package name")
	enumCmd.PersistentFlags().StringVarP(
		&enumName, "name", "n", "", "Name of the enum")
	enumCmd.PersistentFlags().StringSliceVarP(
		&enumValues, "values", "v", []string{}, "Values of the enum")
	enumCmd.PersistentFlags().BoolVarP(
		&noUnknown, "no-unknown", "", false, "Don't generate Unknown value")
	return enumCmd
}

func init() {
	rootCmd.AddCommand(createEnumCmd())
}
