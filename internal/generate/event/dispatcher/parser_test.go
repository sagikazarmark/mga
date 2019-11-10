package dispatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	spec, err := Parse("./testdata/parser", "Events")
	if err != nil {
		t.Fatal(err)
	}

	expected := InterfaceSpec{
		Name: "Events",
		Package: PackageSpec{
			Name: "parser",
			Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
		},
		Methods: []MethodSpec{
			{
				Name: "Event",
				Event: TypeSpec{
					Name: "Event",
					Package: PackageSpec{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: false,
				ReturnsError:    false,
			},
			{
				Name: "EventEmbedded",
				Event: TypeSpec{
					Name: "Event",
					Package: PackageSpec{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "EventEmbeddedFromUnexportedInterface",
				Event: TypeSpec{
					Name: "Event",
					Package: PackageSpec{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "EventWithContext",
				Event: TypeSpec{
					Name: "Event",
					Package: PackageSpec{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    false,
			},
			{
				Name: "EventWithContextAndError",
				Event: TypeSpec{
					Name: "Event",
					Package: PackageSpec{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "EventWithError",
				Event: TypeSpec{
					Name: "Event",
					Package: PackageSpec{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: false,
				ReturnsError:    true,
			},
			{
				Name: "ImportedAliasedEvent",
				Event: TypeSpec{
					Name: "ImportedEvent",
					Package: PackageSpec{
						Name: "imports",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "ImportedEvent",
				Event: TypeSpec{
					Name: "ImportedEvent",
					Package: PackageSpec{
						Name: "imports",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "ImportedEventDispatch",
				Event: TypeSpec{
					Name: "ImportedEvent",
					Package: PackageSpec{
						Name: "imports",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
		},
	}

	assert.Equal(t, expected, spec, "the parsed spec does not match the expected one")
}
