/*
Package rpc is a JSONRPC 2.0 server implementation.

See the JSONRPC 2.0 specification at https://www.jsonrpc.org/specification

This is not a 100% compliant implementation. What it lacks, when compared to the spec, makes no
practical difference when used in a normal web app, imo. It differs from the spec in these ways:

  - Request parameter 'by-name' structures are not implemented
  - Request parameters can be any JSON or custom data type, except that for slices, only strings
    will work. The workaround is define a custom type to contain the slice.
  - Notifications are not implemented

To use it, define a struct that implements the API interface, then register the handler function
with your server.

The handler will forward RPC commands to matching functions on the API object. A matching API
implementation has the following form:

	type MyAPI struct {
	}

	// GetAuth implements API.GetAuth
	func (api MyAPI) GetAuth(w http.ResponseWriter, r *http.Request, rpcreq rpc.JSONRPCRequest)
		(MySession, error) {
		...
	}

	func (api MyAPI) ExampleRPCFunction(
		w http.ResponseWriter, // Optional
		r *http.Request, // Optional
		session WebSession, // Optional
		rpcParam1: type,
		...
		rpcParamN: type,
	) (response[, error]) {
		...
	}

GetAuth is called before a function invocation. If no authentication is required to execute
a function, then GetAuth should return nil, nil.
*/
package rpc

import "github.com/petchells/nrtm4tools/internal/nrtm4/util"

var logger = util.Logger
