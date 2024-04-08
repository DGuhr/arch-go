package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"

	monkey "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"

	"github.com/fdaines/arch-go/old/config"
)

func TestDescribeCommand(t *testing.T) {
	viper.AddConfigPath("../test/")

	t.Run("when arch-go.yaml has no rules", func(t *testing.T) {
		cmd := NewDescribeCommand()
		patch := monkey.ApplyFuncReturn(config.LoadConfig, &config.Config{}, nil)
		defer patch.Reset()

		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, _ := io.ReadAll(b)

		expected := `Dependency Rules
	* No rules defined
Function Rules
	* No rules defined
Content Rules
	* No rules defined
Naming Rules
	* No rules defined
`
		assert.Equal(t, expected, string(out))
	})

	t.Run("when arch-go.yaml has rules", func(t *testing.T) {
		cmd := NewDescribeCommand()
		patch := monkey.ApplyFuncReturn(config.LoadConfig, configLoaderMockWithRules, nil)
		defer patch.Reset()

		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, _ := io.ReadAll(b)

		expected := `Dependency Rules
	* Packages that match pattern 'foobar',
		* Should only depends on:
			* Internal Packages that matches:
				- 'foo'
			* External Packages that matches:
				- 'bar'
			* StandardLib Packages that matches:
				- 'foobar'
		* Should not depends on:
			* Internal Packages that matches:
				- 'oof'
			* External Packages that matches:
				- 'rab'
			* StandardLib Packages that matches:
				- 'raboof'
Function Rules
	* Packages that match pattern 'function-package' should comply with the following rules:
		* Functions should not have more than 3 lines
		* Functions should not have more than 1 parameters
		* Functions should not have more than 2 return values
		* Files should not have more than 4 public functions
Content Rules
	* Packages that match pattern 'package1' should only contain interfaces
	* Packages that match pattern 'package2' should only contain structs
	* Packages that match pattern 'package3' should only contain functions
	* Packages that match pattern 'package4' should only contain methods
	* Packages that match pattern 'package5' should not contain interfaces
	* Packages that match pattern 'package6' should not contain structs
	* Packages that match pattern 'package7' should not contain functions
	* Packages that match pattern 'package8' should not contain methods

Naming Rules
	* Packages that match pattern 'foobar' should comply with:
		* Structs that implement interfaces matching name 'foo' should have simple name starting with 'bla'
	* Packages that match pattern 'barfoo' should comply with:
		* Structs that implement interfaces matching name 'foo' should have simple name ending with 'anything'

`
		assert.Equal(t, expected, string(out))
	})

	t.Run("when arch-go.yaml does not exist", func(t *testing.T) {
		exitCalls := 0
		fakeExit := func(int) {
			exitCalls++
		}
		patch := monkey.ApplyFunc(os.Exit, fakeExit)
		defer patch.Reset()
		configLoaderPatch := monkey.ApplyFuncReturn(config.LoadConfig, nil, fmt.Errorf("Error details"))
		defer configLoaderPatch.Reset()

		cmd := NewDescribeCommand()

		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, _ := io.ReadAll(b)

		expected := `Error: Error details
`
		assert.Equal(t, expected, string(out))
		assert.Equal(t, 1, exitCalls)
	})
}

var configLoaderMockWithRules = &config.Config{
	DependenciesRules: []*config.DependenciesRule{
		{
			Package: "foobar",
			ShouldOnlyDependsOn: &config.Dependencies{
				Internal: []string{"foo"},
				External: []string{"bar"},
				Standard: []string{"foobar"},
			},
			ShouldNotDependsOn: &config.Dependencies{
				Internal: []string{"oof"},
				External: []string{"rab"},
				Standard: []string{"raboof"},
			},
		},
	},
	ContentRules: []*config.ContentsRule{
		{
			Package:                     "package1",
			ShouldOnlyContainInterfaces: true,
		},
		{
			Package:                  "package2",
			ShouldOnlyContainStructs: true,
		},
		{
			Package:                    "package3",
			ShouldOnlyContainFunctions: true,
		},
		{
			Package:                  "package4",
			ShouldOnlyContainMethods: true,
		},
		{
			Package:                    "package5",
			ShouldNotContainInterfaces: true,
		},
		{
			Package:                 "package6",
			ShouldNotContainStructs: true,
		},
		{
			Package:                   "package7",
			ShouldNotContainFunctions: true,
		},
		{
			Package:                 "package8",
			ShouldNotContainMethods: true,
		},
	},
	FunctionsRules: []*config.FunctionsRule{
		{
			Package:                  "function-package",
			MaxParameters:            1,
			MaxReturnValues:          2,
			MaxLines:                 3,
			MaxPublicFunctionPerFile: 4,
		},
	},
	CyclesRules: []*config.CyclesRule{
		{
			Package:                "foobar",
			ShouldNotContainCycles: true,
		},
	},
	NamingRules: []*config.NamingRule{
		{
			Package: "foobar",
			InterfaceImplementationNamingRule: &config.InterfaceImplementationRule{
				StructsThatImplement:             "foo",
				ShouldHaveSimpleNameStartingWith: "bla",
			},
		},
		{
			Package: "barfoo",
			InterfaceImplementationNamingRule: &config.InterfaceImplementationRule{
				StructsThatImplement:           "foo",
				ShouldHaveSimpleNameEndingWith: "anything",
			},
		},
	},
}
