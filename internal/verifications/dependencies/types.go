package dependencies

import (
	"github.com/fdaines/arch-go/pkg/archgo/configuration"
)

type RulesResult struct {
	Results []*RuleResult `json:"results"`
	Passes  bool          `json:"passes"`
}

type RuleResult struct {
	Rule          configuration.DependenciesRule `json:"rule"`
	Description   string                         `json:"description"`
	Verifications []Verification                 `json:"verifications"`
	Passes        bool                           `json:"passes"`
}

type Verification struct {
	Package string   `json:"package"`
	Details []string `json:"details"`
	Passes  bool     `json:"passes"`
}
