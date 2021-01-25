package types

type Row struct {
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
	Service string
	Count   int
}
