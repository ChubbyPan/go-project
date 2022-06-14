// code and encode message

package codec

import(
	"io"
)

// RPC call: func(service.Method, args, &reply)
// client request: service, method, args
// server response: error , reply 
type Header struct {
	ServiceMethod 	string 		// name of service and method
	Seq				uint64  	// request ID
	Error			string		// error msg, set by server
}

// define codec interface, each instance can call this interface(both gob type and json type)
type Codec interface {
	io.Closer
	ReadHeader(*Header) 	error
	ReadBody(interface{})	error
	Write(*Header, interface{})	error
}

type NewCodecFunc func(io.ReadWriterCloser) Codec

type Type string

// two type of codec
const(
	GobType 		Type = "application/gob"
	JsonType		Type = "application/json" 	// not implement
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewCodecFunc
}