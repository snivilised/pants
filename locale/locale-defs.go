package locale

const (
	// PantsSourceID source ID for the pants module
	PantsSourceID = "github.com/snivilised/pants"
)

type pantsTemplData struct{}

func (td pantsTemplData) SourceID() string {
	return PantsSourceID
}
