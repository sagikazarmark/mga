package dispatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sagikazarmark.dev/mga/internal/generate/gentypes"
)

func TestParse(t *testing.T) {
	expected := Events{
		TypeRef: gentypes.TypeRef{
			Name: "Events",
			Package: gentypes.PackageRef{
				Name: "parser",
				Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
			},
		},
		Methods: []EventMethod{
			{
				Name: "Event",
				Event: gentypes.TypeRef{
					Name: "Event",
					Package: gentypes.PackageRef{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: false,
				ReturnsError:    false,
			},
			{
				Name: "EventEmbedded",
				Event: gentypes.TypeRef{
					Name: "Event",
					Package: gentypes.PackageRef{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "EventEmbeddedFromUnexportedInterface",
				Event: gentypes.TypeRef{
					Name: "Event",
					Package: gentypes.PackageRef{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "EventWithContext",
				Event: gentypes.TypeRef{
					Name: "Event",
					Package: gentypes.PackageRef{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    false,
			},
			{
				Name: "EventWithContextAndError",
				Event: gentypes.TypeRef{
					Name: "Event",
					Package: gentypes.PackageRef{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "EventWithError",
				Event: gentypes.TypeRef{
					Name: "Event",
					Package: gentypes.PackageRef{
						Name: "parser",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
					},
				},
				ReceivesContext: false,
				ReturnsError:    true,
			},
			{
				Name: "ImportedAliasedEvent",
				Event: gentypes.TypeRef{
					Name: "ImportedEvent",
					Package: gentypes.PackageRef{
						Name: "imports",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "ImportedEvent",
				Event: gentypes.TypeRef{
					Name: "ImportedEvent",
					Package: gentypes.PackageRef{
						Name: "imports",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "ImportedEventDispatch",
				Event: gentypes.TypeRef{
					Name: "ImportedEvent",
					Package: gentypes.PackageRef{
						Name: "imports",
						Path: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
		},
	}

	evdisp, err := Parse("./testdata/parser", "Events")
	require.NoError(t, err)

	assert.Equal(t, expected, evdisp, "the parsed interface does not match the expected one")
}
