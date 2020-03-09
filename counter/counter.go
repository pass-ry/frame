package counter

type Counter interface {
	Inc(status bool, msg ...string)
}

type (
	Kind uint8
)

const (
	Kind_E_XiaoBao     Kind = 1
	Kind_JD_URL_Import Kind = 2
	Kind_Bind          Kind = 3
	Kind_ICDC          Kind = 4

	statusFail    int = 0
	statusSuccess int = 2
)
