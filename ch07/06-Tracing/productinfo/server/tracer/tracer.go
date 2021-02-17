//---------------------------------------------------------
// 코드 7-13 부분
//---------------------------------------------------------
package tracer

import (
	"log"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

func InitTracing() {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: collectorEndpointURI,
		AgentEndpoint:     agentEndpointURI,
		ServiceName:       "product_info",
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
}

//---------------------------------------------------------
