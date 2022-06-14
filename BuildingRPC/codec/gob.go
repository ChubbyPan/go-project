// one type of codec: gob's implement

package codec

import(
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn 		io.ReadWriteCloser	//instance (when TCP or socket connnected) 
	buf 		*bufio.Writer		//prevent block 
	dec 		*gob.Decoder
	enc			*gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec{
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: 	conn
		buf: 	buf
		dec:	gob.NewDecoder(conn)		// create conn's decoder,
		enc:	gob.NewEncoder(buf)			// create encoder
	}
}

/*
type Codec interface {
	io.Closer
	ReadHeader(*Header) 	error
	ReadBody(interface{})	error
	Write(*Header, interface{})	error
}
*/ 
// method ReadHeader, ReadBody, Write and Close
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{})(err error){
	defer func() {
		//  将任何缓冲的数据写入底层的 io.Writer; 无err时,关闭channel(当前对象c的conn fd)
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	// encode header write to channal
	if err := c.enc.Encode(h); err != nil{
		log.Println("rpc codec: gob error encoding header:",err)
		return err
	}
	// encode body, write to channal
	if err := c.enc.Encode(body); err != nil{
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}