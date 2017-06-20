package update

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/weaveworks/flux"
)

const (
	ServiceSpecAll       = ServiceSpec("<all>")
	ServiceSpecAutomated = ServiceSpec("<automated>")
	ImageSpecLatest      = ImageSpec("<all latest>")
)

var (
	ErrInvalidReleaseKind = errors.New("invalid release kind")
)

// ReleaseKind says whether a release is to be planned only, or planned then executed
type ReleaseKind string

const (
	ReleaseKindPlan    ReleaseKind = "plan"
	ReleaseKindExecute             = "execute"
)

func ParseReleaseKind(s string) (ReleaseKind, error) {
	switch s {
	case string(ReleaseKindPlan):
		return ReleaseKindPlan, nil
	case string(ReleaseKindExecute):
		return ReleaseKindExecute, nil
	default:
		return "", ErrInvalidReleaseKind
	}
}

const UserAutomated = "<automated>"

// NB: these get sent from fluxctl, so we have to maintain the json format of
// this. Eugh.
type ReleaseSpec struct {
	ServiceSpecs []ServiceSpec
	ImageSpec    ImageSpec
	Kind         ReleaseKind
	Excludes     []flux.ServiceID
}

// ReleaseType gives a one-word description of the release, mainly
// useful for labelling metrics or log messages.
func (s ReleaseSpec) ReleaseType() string {
	switch {
	case s.ImageSpec == ImageSpecLatest:
		return "latest_images"
	default:
		return "specific_image"
	}
}

type ServiceSpec string // ServiceID or "<all>"

func ParseServiceSpec(s string) (ServiceSpec, error) {
	switch s {
	case string(ServiceSpecAll):
		return ServiceSpecAll, nil
	case string(ServiceSpecAutomated):
		return ServiceSpecAutomated, nil
	}
	id, err := flux.ParseServiceID(s)
	if err != nil {
		return "", errors.Wrap(err, "invalid service spec")
	}
	return ServiceSpec(id), nil
}

func (s ServiceSpec) AsID() (flux.ServiceID, error) {
	return flux.ParseServiceID(string(s))
}

func (s ServiceSpec) String() string {
	return string(s)
}

// ImageSpec is an ImageID, or "<all latest>" (update all containers
// to the latest available), or "<no updates>" (do not update any
// images)
type ImageSpec string

func ParseImageSpec(s string) (ImageSpec, error) {
	if s == string(ImageSpecLatest) {
		return ImageSpec(s), nil
	}

	parts := strings.Split(s, ":")
	if len(parts) != 2 || parts[1] == "" {
		return "", errors.Wrap(flux.ErrInvalidImageID, "blank tag (if you want latest, explicitly state the tag :latest)")
	}

	id, err := flux.ParseImageID(s)
	return ImageSpec(id.String()), err
}

func (s ImageSpec) String() string {
	return string(s)
}

func (s ImageSpec) AsID() (flux.ImageID, error) {
	return flux.ParseImageID(s.String())
}

func ImageSpecFromID(id flux.ImageID) ImageSpec {
	return ImageSpec(id.String())
}