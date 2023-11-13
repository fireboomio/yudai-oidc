package object

type Response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type CodeResponse struct {
	Errors   []Res    `json:"errors"`
	Response Response `json:"data"`
}

type Res struct {
	Message  string   `json:"message"`
	Location []string `json:"locations"`
	Path     []string `json:"path"`
}
