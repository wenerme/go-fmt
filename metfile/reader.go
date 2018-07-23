package metfile

import (
	"github.com/wenerme/goform"
	"io"
	"net"
)

func NewMetReader(r io.Reader) *MetReader {
	return &MetReader{
		r:   r,
		buf: make([]byte, 4),
	}
}

func (mr *MetReader) ReadForm() (met *MetForm, err error) {
	if err = mr.readHeader(); err != nil {
		return
	}
	if err = mr.readServers(); err != nil {
		return
	}

	met = mr.form
	return
}
func (mr *MetReader) readTags(server *MetServer, count int) (err error) {
	r := mr.r
	buf := mr.buf

	for i := 0; i < count; i++ {
		tag := &MetServerTag{}

		if _, err = r.Read(buf[:1]); err != nil {
			return
		}
		tag.ValueType = TagValueType(uint8(buf[0]))

		if _, err = r.Read(buf[:2]); err != nil {
			return
		}
		if nameLen := int(Endian.Uint16(buf)); nameLen == 1 {
			if _, err = r.Read(buf[:1]); err != nil {
				return
			}
			tag.Code = TagCode(uint8(buf[0]))
			neoType := tag.Code.ValueType(tag.ValueType)
			if neoType != tag.ValueType {
				mr.log("value type mismatch %v => %v", tag.ValueType, neoType)
			}
			tag.ValueType = neoType
		} else {
			nameBuf := make([]byte, nameLen)
			if _, err = io.ReadFull(r, nameBuf); err != nil {
				return
			}
			tag.Name = string(nameBuf)
		}

		switch tag.ValueType {
		case TagValueTypeString:
			{
				if _, err = r.Read(buf[:2]); err != nil {
					return
				}
				l := int(Endian.Uint16(buf))
				b := make([]byte, l)
				if _, err = io.ReadFull(r, b); err != nil {
					return
				}
				tag.StringValue = string(b)
			}
		case TagValueTypeUint32:
			if _, err = r.Read(buf); err != nil {
				return
			}
			tag.Uint32Value = Endian.Uint32(buf)
		default:
			err = goform.EInvalidType
			return
		}

		server.Tags = append(server.Tags, tag)
		mr.log("Tag#%d %v", i, tag)
	}

	return
}

func (mr *MetReader) readServers() (err error) {
	r := mr.r
	buf := mr.buf

	_, err = r.Read(buf)
	if err != nil {
		return
	}

	serverCount := int(Endian.Uint32(buf))
	mr.log("Server count %d", serverCount)

	mr.form = &MetForm{
		Servers: make([]*MetServer, 0),
	}

	for i := 0; i < serverCount; i++ {
		if _, err = r.Read(buf); err != nil {
			return
		}
		ip := net.IPv4(buf[0], buf[1], buf[2], buf[3])

		if _, err = r.Read(buf[:2]); err != nil {
			return
		}
		port := Endian.Uint16(buf)

		if _, err = r.Read(buf); err != nil {
			return
		}
		tagCount := int(Endian.Uint32(buf))

		server := &MetServer{
			IP:   ip,
			Port: int(port),
			Tags: make([]*MetServerTag, 0),
		}
		mr.form.Servers = append(mr.form.Servers, server)

		mr.log("Server#%v ip=%v port=%v tagCount=%v", i, ip, port, tagCount)

		if err = mr.readTags(server, tagCount); err != nil {
			return
		}
	}

	return err
}

func (mr *MetReader) readHeader() (err error) {
	r := mr.r
	buf := mr.buf

	_, err = r.Read(buf[:1])
	if err != nil {
		return
	}

	// Ox0E or OxE0
	if !(buf[0] == 0x0E || buf[0] == 0xE0) {
		err = goform.EInvalidFileFormat
		return
	}

	return
}

func (mr *MetReader) log(fmt string, args ...interface{}) {
	if mr.Logger != nil {
		mr.Logger(fmt, args...)
	}
}

type MetReader struct {
	r      io.Reader
	buf    []byte
	form   *MetForm
	Logger func(format string, args ...interface{})
}
