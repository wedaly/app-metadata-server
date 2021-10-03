package server

// App stores metadata about an application.
type App struct {
	Title       string
	Version     string
	Maintainers []Maintainer
	Company     string
	Website     string
	Source      string
	License     string
	Description string
}

// Validate checks that the application metadata is valid.
func (app App) Validate(errs *ValidationErrors) {
	ValidateStringNonEmpty(errs, "app.title", app.Title)
	ValidateStringNonEmpty(errs, "app.version", app.Version)
	ValidateStringNonEmpty(errs, "app.company", app.Company)
	ValidateURL(errs, "app.website", app.Website)
	ValidateURL(errs, "app.source", app.Source)
	ValidateStringNonEmpty(errs, "app.license", app.License)
	ValidateStringNonEmpty(errs, "app.description", app.Description)
	app.validateMaintainers(errs)
}

func (app App) validateMaintainers(errs *ValidationErrors) {
	if len(app.Maintainers) == 0 {
		errs.Append("app.maintainers", "At least one maintainer must be specified")
	}

	for _, m := range app.Maintainers {
		m.Validate(errs)
	}
}

// Maintainer stores metadata about an application maintainer.
type Maintainer struct {
	Name  string
	Email string
}

// Validate checks that the maintainer is valid.
func (m Maintainer) Validate(errs *ValidationErrors) {
	ValidateStringNonEmpty(errs, "maintainer.name", m.Name)
	ValidateEmailAddress(errs, "maintainer.email", m.Email)
}
