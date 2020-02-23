package setting

import "net/textproto"

type imageFile struct {
	Header   textproto.MIMEHeader
	Filename string
	Size     int64
	Body     []byte
}
