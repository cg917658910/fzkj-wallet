package handlers

import (
	"context"

	"net/http"

	"github.com/cg917658910/fzkj-wallet/notify-service/lib/codes"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
	pb "github.com/cg917658910/fzkj-wallet/notify-service/proto"
	"github.com/cg917658910/fzkj-wallet/notify-service/svc"
	"github.com/go-kit/kit/endpoint"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

var logger = log.DLogger()

// WrapEndpoints accepts the service's entire collection of endpoints, so that a
// set of middlewares can be wrapped around every middleware (e.g., access
// logging and instrumentation), and others wrapped selectively around some
// endpoints and not others (e.g., endpoints requiring authenticated access).
// Note that the final middleware wrapped will be the outermost middleware
// (i.e. applied first)
func WrapEndpoints(in svc.Endpoints) svc.Endpoints {

	// Pass a middleware you want applied to every endpoint.
	// optionally pass in endpoints by name that you want to be excluded
	// e.g.
	// in.WrapAllExcept(authMiddleware, "Status", "Ping")

	// Pass in a svc.LabeledMiddleware you want applied to every endpoint.
	// These middlewares get passed the endpoints name as their first argument when applied.
	// This can be used to write generic metric gathering middlewares that can
	// report the endpoint name for free.
	// github.com/metaverse/truss/_example/middlewares/labeledmiddlewares.go for examples.
	// in.WrapAllLabeledExcept(errorCounter(statsdCounter), "Status", "Ping")

	// How to apply a middleware to a single endpoint.
	// in.ExampleEndpoint = authMiddleware(in.ExampleEndpoint)
	in.WrapAllExcept(logAllRequest())
	in.WrapAllExcept(convertError())

	return in
}

func WrapService(in pb.NotifyServiceServer) pb.NotifyServiceServer {
	return in
}

func logAllRequest() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Infoln("request", request)
			response, err := next(ctx, request)
			if err != nil {
				logger.Infoln("response", "error", err)
			} else {
				logger.Infoln("response", response)
			}
			return response, err
		}
	}
}

func convertError() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {

			ret, err := next(ctx, request)
			if err != nil {

				code, ok := err.(codes.Code)
				if ok {
					grpccode := grpccodes.Unknown
					switch code.HTTPStatus {
					case http.StatusUnauthorized:
						grpccode = grpccodes.Unauthenticated
					case http.StatusNotFound:
						grpccode = grpccodes.NotFound
					case http.StatusUnprocessableEntity:
						grpccode = grpccodes.AlreadyExists
					case http.StatusForbidden:
						grpccode = grpccodes.PermissionDenied
					case http.StatusBadRequest:
						grpccode = grpccodes.InvalidArgument
					}
					err = grpcstatus.Errorf(grpccode, code.Msg)
				}
			}
			return ret, err
		}
	}
}
