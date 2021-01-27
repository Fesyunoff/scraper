package types

type Config struct {
	FileName        string
	Host            string
	Port            int
	Timeout         int
	Proxy           string
	Time            int
	HostDB          string
	PortDB          int
	UserDB          string
	PasswordDB      string
	NameDB          string
	SchemaName      string
	RespTableName   string
	ReqTableName    string
	UsrTableName    string
	StatHoursBefore int64
	StatLimit       int
}

type Row struct {
	Id         int
	Service    string
	Date       int64
	Responce   bool
	StatusCode int
	Duration   int64
}

type Request struct {
	Id      int
	Service string
	Date    int64
}

type Stat struct {
	Service string `json:"service"`
	Count   int    `json:"count"`
}
