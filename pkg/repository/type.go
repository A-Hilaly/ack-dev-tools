package repository

type RepositoryType int

const (
	RepositoryTypeCore RepositoryType = iota
	RepositoryTypeTooling
	RepositoryTypeController
	RepositoryTypeUnknown
)

func (rt RepositoryType) String() string {
	switch rt {
	case RepositoryTypeCore:
		return "core"
	case RepositoryTypeController:
		return "controller"
	case RepositoryTypeUnknown:
		return "UNKNOWN"
	default:
		panic("unsupported repository type")
	}
}

func repositoryTypeFromString(s string) RepositoryType {
	switch s {
	case "core":
		return RepositoryTypeCore
	case "controller":
		return RepositoryTypeController
	default:
		panic("unsupported repository type")
	}
}
