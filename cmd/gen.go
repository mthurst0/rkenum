package cmd

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mthurst0/rkstrings"
	"github.com/spf13/cobra"
)

func parseAliases(aliases []string) (map[string][]string, error) {
	result := make(map[string][]string)
	uniq := make(map[string]bool)
	for _, alias := range aliases {
		parts := strings.Split(alias, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid alias: %s", alias)
		}
		a := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if result[v] == nil {
			result[v] = make([]string, 0)
		}
		result[v] = append(result[v], a)
		if uniq[a] {
			return nil, fmt.Errorf("duplicate alias: %s", a)
		}
		uniq[a] = true
	}
	return result, nil
}

func validAliases(aliases map[string][]string, values map[string]bool) error {
	for a := range aliases {
		if _, ok := values[a]; !ok {
			return fmt.Errorf("alias %s must reference a valid value", a)
		}
	}
	for _, as := range aliases {
		for _, a := range as {
			if _, ok := values[a]; ok {
				return fmt.Errorf("alias %s cannot be a value", a)
			}
		}
	}
	return nil
}

type brevBuilder struct {
	strings.Builder
	firstError error
}

func (b *brevBuilder) w(s string) {
	_, err := b.WriteString(s)
	if err != nil && b.firstError == nil {
		b.firstError = err
	}
}

func (b *brevBuilder) f(format string, args ...any) {
	_, err := fmt.Fprintf(b, format, args...)
	if err != nil && b.firstError == nil {
		b.firstError = err
	}
}

func GenerateEnum(
	pkg string,
	enumName string,
	values []string,
	noUnknown bool,
	aliases []string) (string, error) {

	aliasesMap, err := parseAliases(aliases)
	if err != nil {
		return "", err
	}
	valuesMap := make(map[string]bool)
	for _, v := range values {
		valuesMap[v] = true
	}
	if err := validAliases(aliasesMap, valuesMap); err != nil {
		return "", err
	}

	firstRune, _ := utf8.DecodeRuneInString(enumName[:1])
	firstRuneUpper := rkstrings.IsUpperRune(firstRune)
	nChar := "n"
	aChar := "a"
	sChar := "s"
	if firstRuneUpper {
		nChar = "N"
		aChar = "A"
		sChar = "S"
	}
	enumNameTitleCase := strings.ToTitle(enumName[0:1]) + enumName[1:]

	var b brevBuilder
	b.w("// Generated code - DO NOT EDIT\n\n")
	b.f("package %s\n\n", pkg)
	b.w("import (\n")
	b.w("\t\"fmt\"\n")
	b.w("\t\"strings\"\n")
	b.w(")\n\n")
	b.f("type %s int\n\n", enumName)
	b.w("const (\n")
	if !noUnknown {
		b.f("\t%sUnknown = %s(iota)\n", enumName, enumName)
		for _, v := range values {
			v = rkstrings.ToCamelCase(v)
			b.f("\t%s%s\n", enumName, v)
		}
	} else {
		firstValue := true
		for _, v := range values {
			v = rkstrings.ToCamelCase(v)
			if firstValue {
				b.f("\t%s%s = %s(iota)\n", enumName, v, enumName)
			} else {
				b.f("\t%s%s\n", enumName, v)
			}
			firstValue = false
		}
	}
	b.f("\t%sMax\n", enumName)
	b.w(")\n\n")

	b.f("func %sew%sFromString(s string) (%s, error) {\n",
		nChar, enumNameTitleCase, enumName)
	b.w("\tswitch strings.ToLower(s) {\n")
	for _, v := range values {
		b.f("\tcase \"%s\":\n", v)
		b.f("\t\treturn %s%s, nil\n", enumName, rkstrings.ToCamelCase(v))
	}
	b.w("\tdefault:\n")
	b.f("\t\treturn %s(0), fmt.Errorf(\"could not convert string to %s: %%s\", s)\n",
		enumName, enumName)
	b.w("\t}\n")
	b.w("}\n\n")

	b.f("func (v %s) %string() string {\n", enumName, sChar)
	b.f("\tswitch v {\n")
	for _, v := range values {
		b.f("\tcase %s%s:\n", enumName, rkstrings.ToCamelCase(v))
		b.f("\t\treturn \"%s\"\n", v)
	}
	b.w("\tdefault:\n")
	b.f("\t\treturn fmt.Sprintf(\"%s(%%d)\", v)\n", enumName)
	b.w("\t}\n")
	b.w("}\n")

	if len(aliasesMap) > 0 {
		b.f("\nfunc %sew%sFromAlias(s string) (%s, error) {\n",
			nChar, enumNameTitleCase, enumName)
		b.w("\tswitch strings.ToLower(s) {\n")
		for _, v := range values {
			if as, ok := aliasesMap[v]; ok {
				for _, a := range as {
					b.f("\tcase \"%s\":\n", a)
					b.f("\t\treturn %s%s, nil\n",
						enumName, rkstrings.ToCamelCase(v))
				}
			}
		}
		b.w("\tdefault:\n")
		b.f("\t\treturn %s(0), fmt.Errorf(\"could not convert string to %s: %%s\", s)\n",
			enumName, enumName)
		b.w("\t}\n")
		b.w("}\n\n")

		b.f("\nfunc (v %s) %sliases() []string {\n", enumName, aChar)
		b.w("\tswitch v {\n")
		for _, v := range values {
			b.f("\tcase %s%s:\n", enumName, rkstrings.ToCamelCase(v))
			b.w("\t\treturn []string{")
			for _, a := range aliasesMap[v] {
				b.f("\"%s\", ", a)
			}
			b.f("}\n")
		}
		b.w("\tdefault:\n")
		b.w("\t\treturn nil\n")
		b.w("\t}\n")
		b.w("}\n")
	}
	b.f("\nfunc %sew%s(s string) (%s, error) {\n", nChar, enumNameTitleCase, enumName)
	b.f("\tv, err := %sew%sFromString(s)\n", nChar, enumNameTitleCase)
	b.w("\tif err != nil {\n")
	b.f("\t\tv, err = %sew%sFromAlias(s)\n", nChar, enumNameTitleCase)
	b.w("\t}\n")
	b.w("\treturn v, err\n")
	b.w("}\n\n")

	return b.String(), nil
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

func createGenCmd() *cobra.Command {
	var packageName string
	var enumName string
	var enumValues []string
	var noUnknown bool
	var aliases []string
	var enumCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate an enum",
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
			s, err := GenerateEnum(
				packageName, enumName, parseStrings(enumValues), noUnknown, aliases)
			if err != nil {
				return err
			}
			fmt.Println(s)
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
	enumCmd.PersistentFlags().StringSliceVarP(
		&aliases, "alias", "a", []string{}, "Alias for a value")
	return enumCmd
}

func init() {
	rootCmd.AddCommand(createGenCmd())
}
