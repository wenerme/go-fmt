package metfile

import (
	"github.com/wenerme/goform"
	"io"
	"net"
)

type MetWriter struct {
	w    io.Writer
	buf  []byte
	form *MetForm
}

func NewMetWriter(form *MetForm, w io.Writer) *MetWriter {
	return &MetWriter{
		w:    w,
		form: form,
		buf:  make([]byte, 4),
	}
}

func (mw *MetWriter) WriteForm() (err error) {
	if err = mw.writeHeader(); err != nil {
		return
	}
	if err = mw.writeServers(); err != nil {
		return
	}
	return
}

func (mw *MetWriter) writeHeader() (err error) {
	w := mw.w

	if _, err = w.Write([]byte{0x0E}); err != nil {
		return
	}

	return
}

func (mw *MetWriter) writeServers() (err error) {
	w := mw.w
	buf := mw.buf
	form := mw.form

	Endian.PutUint32(buf, uint32(len(form.Servers)))
	if _, err = w.Write(buf); err != nil {
		return
	}

	for _, server := range form.Servers {
		if _, err = w.Write(server.IP[net.IPv6len-4:]); err != nil {
			return
		}

		Endian.PutUint16(buf, uint16(server.Port))
		if _, err = w.Write(buf[:2]); err != nil {
			return
		}

		Endian.PutUint32(buf, uint32(len(server.Tags)))
		if _, err = w.Write(buf); err != nil {
			return
		}

		for _, tag := range server.Tags {
			if _, err = w.Write([]byte{byte(tag.ValueType)}); err != nil {
				return
			}

			if len(tag.Name) < 2 {
				buf[0] = 0x01
				buf[1] = 0x00
				buf[2] = byte(tag.Code)
				if _, err = w.Write(buf[:3]); err != nil {
					return
				}
			} else {
				Endian.PutUint16(buf, uint16(len(tag.Name)))
				if _, err = w.Write(buf[:2]); err != nil {
					return
				}
				if _, err = w.Write([]byte(tag.Name)); err != nil {
					return
				}
			}

			switch tag.ValueType {
			case TagValueTypeString:
				Endian.PutUint16(buf, uint16(len(tag.StringValue)))
				if _, err = w.Write(buf[:2]); err != nil {
					return
				}
				if _, err = w.Write([]byte(tag.StringValue)); err != nil {
					return
				}
			case TagValueTypeUint32:
				Endian.PutUint32(buf, tag.Uint32Value)
				if _, err = w.Write(buf); err != nil {
					return
				}
			default:
				err = goform.EInvalidType
				return
			}
		}
	}

	return
}
