package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/linuxunsw/vote/tui/internal/sdk"
	"github.com/linuxunsw/vote/tui/internal/tui/styles"
)

func Voting(data sdk.PublicBallot, vote map[string]string) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key(string(sdk.NominationExecutiveRolesPresident)).
				Title("president").
				Options(optionsForRole(data, sdk.NominationExecutiveRolesPresident)...),
			huh.NewSelect[string]().
				Key(string(sdk.NominationExecutiveRolesSecretary)).
				Title("secretary").
				Options(optionsForRole(data, sdk.NominationExecutiveRolesSecretary)...),
			huh.NewSelect[string]().
				Key(string(sdk.NominationExecutiveRolesTreasurer)).
				Title("treasurer").
				Options(optionsForRole(data, sdk.NominationExecutiveRolesTreasurer)...),
			huh.NewSelect[string]().
				Key(string(sdk.NominationExecutiveRolesArcDelegate)).
				Title("arc delegate").
				Options(optionsForRole(data, sdk.NominationExecutiveRolesArcDelegate)...),
			huh.NewSelect[string]().
				Key(string(sdk.NominationExecutiveRolesEdiOfficer)).
				Title("edi officer").
				Options(optionsForRole(data, sdk.NominationExecutiveRolesEdiOfficer)...),
			huh.NewSelect[string]().
				Key(string(sdk.NominationExecutiveRolesGrievanceOfficer)).
				Title("grievance officer").
				Options(optionsForRole(data, sdk.NominationExecutiveRolesGrievanceOfficer)...),
		),
	).WithTheme(styles.FormTheme())

	return form
}

func optionsForRole(data sdk.PublicBallot, role sdk.NominationExecutiveRoles) []huh.Option[string] {
	var opts []huh.Option[string]
	candidates, ok := data.Candidates[string(role)]

	if !ok {
		return []huh.Option[string]{huh.NewOption("no option", "")}
	}

	for _, candidate := range *candidates {
		opts = append(opts, huh.NewOption(candidate.CandidateName, candidate.NominationId))
	}

	return opts
}
