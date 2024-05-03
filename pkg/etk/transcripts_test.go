package etk_test

import (
	"embed"
	"testing"

	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/etk/etktest"
	"src.elv.sh/pkg/eval/evaltest"
)

//go:embed *.elvts
var transcripts embed.FS

func TestTranscripts(t *testing.T) {
	evaltest.TestTranscriptsInFS(t, transcripts,
		"code-area-fixture", etktest.MakeFixture(etk.CodeArea),
	)
}
