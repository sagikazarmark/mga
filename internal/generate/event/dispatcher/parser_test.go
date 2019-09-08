package dispatcher

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	spec, err := Parse("./testdata/parser", "Events")
	if err != nil {
		t.Fatal(err)
	}

	expected := Spec{
		Package: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
		Name:    "Events",
		Events: []EventSpec{
			{
				DispatchName:    "Event",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
				EventName:       "Event",
				ReceivesContext: false,
				ReturnsError:    false,
			},
			{
				DispatchName:    "EventEmbedded",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
				EventName:       "Event",
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				DispatchName:    "EventEmbeddedFromUnexportedInterface",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
				EventName:       "Event",
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				DispatchName:    "EventWithContext",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
				EventName:       "Event",
				ReceivesContext: true,
				ReturnsError:    false,
			},
			{
				DispatchName:    "EventWithContextAndError",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
				EventName:       "Event",
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				DispatchName:    "EventWithError",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser",
				EventName:       "Event",
				ReceivesContext: false,
				ReturnsError:    true,
			},
			{
				DispatchName:    "ImportedAliasedEvent",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
				EventName:       "ImportedEvent",
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				DispatchName:    "ImportedEvent",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
				EventName:       "ImportedEvent",
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				DispatchName:    "ImportedEventDispatch",
				EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
				EventName:       "ImportedEvent",
				ReceivesContext: true,
				ReturnsError:    true,
			},
		},
	}

	if !reflect.DeepEqual(expected, spec) {
		t.Error("the parsed spec does not match the expected one")
	}
}
