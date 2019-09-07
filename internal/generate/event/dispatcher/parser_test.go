package dispatcher

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		pkg           string
		interfaceName string
		spec          Spec
	}{
		{
			pkg:           "events",
			interfaceName: "Events",
			spec: Spec{
				Package: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/events",
				Name:    "Events",
				Events: []EventSpec{
					{
						DispatchName:    "Event",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/events",
						EventName:       "Event",
						ReceivesContext: false,
						ReturnsError:    false,
					},
					{
						DispatchName:    "EventWithContext",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/events",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    false,
					},
					{
						DispatchName:    "EventWithContextAndError",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/events",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    true,
					},
					{
						DispatchName:    "EventWithError",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/events",
						EventName:       "Event",
						ReceivesContext: false,
						ReturnsError:    true,
					},
				},
			},
		},
		{
			pkg:           "embeds",
			interfaceName: "Events",
			spec: Spec{
				Package: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds",
				Name:    "Events",
				Events: []EventSpec{
					{
						DispatchName:    "Event",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    true,
					},
					{
						DispatchName:    "MoreEvent",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    true,
					},
					{
						DispatchName:    "OtherEvent",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    true,
					},
				},
			},
		},
		{
			pkg:           "imports",
			interfaceName: "Events",
			spec: Spec{
				Package: "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports",
				Name:    "Events",
				Events: []EventSpec{
					{
						DispatchName:    "Event",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    true,
					},
					{
						DispatchName:    "MoreEvent",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    true,
					},
					{
						DispatchName:    "OtherEvent",
						EventPackage:    "sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds",
						EventName:       "Event",
						ReceivesContext: true,
						ReturnsError:    true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.pkg, func(t *testing.T) {
			spec, err := Parse("./"+filepath.Join("testdata/parser", test.pkg), test.interfaceName)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(spec, test.spec) {
				t.Error("the parsed spec does not match the expected one")
			}
		})
	}
}
