package functions

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fdaines/arch-go/config"
	"github.com/fdaines/arch-go/impl/model"
	"strings"
)

type FunctionsRuleVerification struct {
	Module  	   string
	Description    string
	Rule    	   *config.FunctionsRule
	PackageDetails []model.PackageVerification
	Passes		   bool
}

func NewFunctionsRuleVerification(module string, rule *config.FunctionsRule) *FunctionsRuleVerification {
	var ruleDescriptions []string
	if rule.MaxParameters > 0 {
		ruleDescriptions = append(ruleDescriptions, fmt.Sprintf("'at most %d parameters'", rule.MaxParameters))
	}
	if rule.MaxReturnValues > 0 {
		ruleDescriptions = append(ruleDescriptions, fmt.Sprintf("'at most %d return values'", rule.MaxReturnValues))
	}
	if rule.MaxLines > 0 {
		ruleDescriptions = append(ruleDescriptions, fmt.Sprintf("'at most %d lines'", rule.MaxLines))
	}
	if rule.MaxPublicFunctionPerFile > 0 {
		ruleDescriptions = append(ruleDescriptions, fmt.Sprintf("'no more than %d public functions per file'", rule.MaxPublicFunctionPerFile))
	}
	description := fmt.Sprintf("Functions in packages matching pattern '%s' should have [%s]", rule.Package, strings.Join(ruleDescriptions, ","))

	return &FunctionsRuleVerification{
		Module: module,
		Rule: rule,
		Description: description,
		Passes: true,
	}
}

func (d *FunctionsRuleVerification) Verify() {
	result := true
	for index, pd := range d.PackageDetails {
		packagePasses := true

		ruleResult, ruleDetails := check_max_parameters(pd.Package, d.Module, d.Rule.MaxParameters)
		packagePasses = packagePasses && ruleResult
		d.PackageDetails[index].Details = append(d.PackageDetails[index].Details, ruleDetails...)

		ruleResult, ruleDetails = check_max_return_values(pd.Package, d.Module, d.Rule.MaxReturnValues)
		packagePasses = packagePasses && ruleResult
		d.PackageDetails[index].Details = append(d.PackageDetails[index].Details, ruleDetails...)

		ruleResult, ruleDetails = check_max_lines(pd.Package, d.Module, d.Rule.MaxLines)
		packagePasses = packagePasses && ruleResult
		d.PackageDetails[index].Details = append(d.PackageDetails[index].Details, ruleDetails...)

		ruleResult, ruleDetails = check_max_public_functions(pd.Package, d.Module, d.Rule.MaxPublicFunctionPerFile)
		packagePasses = packagePasses && ruleResult
		d.PackageDetails[index].Details = append(d.PackageDetails[index].Details, ruleDetails...)

		d.PackageDetails[index].Passes = packagePasses
		result = result && packagePasses
	}
	d.Passes = result
}

func (d *FunctionsRuleVerification) PrintResults() {
	if d.Passes {
		color.Green("[PASS] - %s\n", d.Description)
	} else {
		color.Red("[FAIL] - %s\n", d.Description)
	}
	for _, p := range d.PackageDetails {
		if p.Passes {
			color.Green("\tPackage '%s' passes\n", p.Package.Path)
		} else {
			color.Red("\tPackage '%s' fails\n", p.Package.Path)
			for _, str := range p.Details {
				color.Red("\t\t%s\n", str)
			}
		}
	}
}
