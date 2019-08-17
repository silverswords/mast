package mastgrpc

//Interceptors
//Please send a PR to add new interceptors or middleware to this list
//
//Auth
//grpc_auth - a customizable (via AuthFunc) piece of auth middleware
//Logging
//grpc_ctxtags - a library that adds a Tag map to context, with data populated from request body
//grpc_zap - integration of zap logging library into gRPC handlers.
//grpc_logrus - integration of logrus logging library into gRPC handlers.
//Monitoring
//grpc_prometheus⚡ - Prometheus client-side and server-side monitoring middleware
//otgrpc⚡ - OpenTracing client-side and server-side interceptors
//grpc_opentracing - OpenTracing client-side and server-side interceptors with support for streaming and handler-returned tags
//Client
//grpc_retry - a generic gRPC response code retry mechanism, client-side middleware
//Server
//grpc_validator - codegen inbound message validation from .proto options
//grpc_recovery - turn panics into gRPC errors
//ratelimit - grpc rate limiting by your own limiter
//Status
//This code has been running in production since May 2016 as the basis of the gRPC micro services stack at Improbable.
//
//Additional tooling will be added, and contributions are welcome.
//
