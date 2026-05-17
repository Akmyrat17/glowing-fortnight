package enums

type LogType string

const (
	LogTypeRequest  LogType = "request"
	LogTypeResponse LogType = "response"
	LogTypeSystem   LogType = "system"
	LogTypeAudit    LogType = "audit"
)

func (t LogType) String() string { return string(t) }
