package i18n

// CLIENT-TODO: Should be updated to use url of the implementing project,
// so should not be left as astrolib. (this should be set by auto-check)
const PantsSourceID = "github.com/snivilised/pants"

type pantsTemplData struct{}

func (td pantsTemplData) SourceID() string {
	return PantsSourceID
}
