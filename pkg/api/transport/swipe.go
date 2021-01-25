package transport

import (
	"net/http"

	. "github.com/swipe-io/swipe/v2"

	"github.com/fesyunoff/availability/pkg/api/service"
)

func Swipe() {
	Build(
		Service(
			Interface((*service.ScraperRequest)(nil), ""),

			HTTPServer(),

			ClientsEnable([]string{"go", "js"}),

			OpenapiEnable(),
			OpenapiOutput("./"),

			ReadmeEnable(),

			MethodOptions(service.ScraperRequest.GetAvailability,
				RESTPath("/getAvailability"),
				RESTQueryVars([]string{"site", "site"}),
				RESTMethod(http.MethodGet),
				Logging(true),
			),

			MethodOptions(service.ScraperRequest.GetResponceTime,
				RESTPath("/getResponceTime"),
				RESTQueryVars([]string{"limit", "limit"}),
				RESTMethod(http.MethodGet),
				Logging(true),
			),

			MethodOptions(service.ScraperRequest.GetStatistics,
				RESTPath("/getStatistics"),
				RESTQueryVars([]string{"hours", "hours", "limit", "limit"}),
				RESTMethod(http.MethodGet),
				Logging(true),
			),

			MethodDefaultOptions(Logging(false), Instrumenting(true)),
		),
	)
}
