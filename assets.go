package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _esc_localFS struct{}

var _esc_local _esc_localFS

type _esc_staticFS struct{}

var _esc_static _esc_staticFS

type _esc_file struct {
	compressed string
	size       int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_esc_localFS) Open(name string) (http.File, error) {
	f, present := _esc_data[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_esc_staticFS) Open(name string) (http.File, error) {
	f, present := _esc_data[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		gr, err = gzip.NewReader(bytes.NewBufferString(f.compressed))
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (f *_esc_file) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_esc_file
	}
	return &httpFile{
		Reader:    bytes.NewReader(f.data),
		_esc_file: f,
	}, nil
}

func (f *_esc_file) Close() error {
	return nil
}

func (f *_esc_file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_esc_file) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_esc_file) Name() string {
	return f.name
}

func (f *_esc_file) Size() int64 {
	return f.size
}

func (f *_esc_file) Mode() os.FileMode {
	return 0
}

func (f *_esc_file) ModTime() time.Time {
	return time.Time{}
}

func (f *_esc_file) IsDir() bool {
	return f.isDir
}

func (f *_esc_file) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _esc_local
	}
	return _esc_static
}

var _esc_data = map[string]*_esc_file{

	"/index.html": {
		local: "public/index.html",
		size:  2270,
		compressed: "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xecV\xcbn\xdb:\x10\xdd\xe7+\x06\n\xee.\xa2$۲o\x1c\xd9@\x81\x02\xed\xa2\xe8/\x1449\xb2\x88P\xa2@ұݠ\xff^Z\x0f\xeb\xe5\x18]\x17E\x16\x11y\xe6\xc59s\x06N2\x9b\xcb\xed\x03@\x92!\xe5\x97\x0f\xf7i\x85\x95\xb8\xfd\xac\xf6%e\xafIP\x1fk\xc8سD\xb0" +
			"\xe7\x127\x9eœ\r\x981^\x8d\x01\xec\x14?\xc3{spG\xe7\xbd\xd7\xeaPp\x9f)\xa9\xf4\x1a\x1e\xd34}\xb9\xe2\xa9*\xac\x9f\xd2\\\xc8\xf3\x1a\xbc\xaf(\xdf\xd0\nF\xe1;\x1e\xd0{\xba\x9e\x9f>iA員\x85\xf1\rj1\x8e`\xc4O\\C\xb4(O-\xf0\xab\xf9\xff\xc8\x1cNE\x81\xbaW\xd5Qp\x9b\xad\xe19" +
			"\xfc\xaf\x8b\x93ӓ\xdfܯ\xe2\xb0\v\x04PR\xceE\xb1\xf7\xb5\xd8g\xd6e\x89o\x81\x12\xd3\t\x96S\xbd\x17E\xebG\x0fVM\xb0ڭ\x0f\xb5\x85\x1fd\xaf`)\x8c{\xe4\xa5\xefk(T\x81\x93\xfck\bǱ\xd7\xe0\x8a\xe9\xae۸R\xc0;\xb4\\\xcc\xe7\xf3\x97^F\xe2z\xfb\x86\xda\xc0\xb62\x9b&\x89\xc2\xe9\xfb\xaa<" +
			"\xdd\xe5Ni\x8e.v\xe4\x92\x1b%\x05w\x1406}\xdd \x17\xd9k\xc4\xc2\x15V{7m\x89\xbb\x10\x11ݱg\xf6a\xad\xe4\x8cR\xaa\xe3\x9d\x00i\xc4\x16a\xfaq\x00\x8d\xfc\x8e7\xae\x16l~'\xfd^\xd3\xf3\x1d\xf7\x1dwޫΝ\xf6z{Q\x90ϑ)M\xadPŘ\xdfN@5\x02\x1aK\xa4\x16\f\xd3JJ" +
			"\bݟ\xd5N\x16%\xd5X\xd8)٤Rn/__/$\xc6|\xa4\xa4#\xd6\xd3\x1a\x87\xbd\x81\x92N?~\xd6 \x11\x89\xa6\xe3\xbfS֪|0\n\xbd\x1ajv\xbbb\xae\xf37\xa6\xb5c\xf2\x86\xf1\x98\u0086\xb5\x1b\x96c\xbaZ\x86n\x98\x0e\xa9I\x82Jd\xd52\f\xdam\x98\\vZ\xb3\xfb\xb8x\x03\xc17\xdeu\xab" +
			"\\\x17_\x92Eݺt\xdf\xed\xb5\xd31\x93Ԙ\x8d\xd7L\xcc\xd5ÁNe\rX\xb5\xa8\a9\x90B\xa61\xddx\x83\xdbADQ\xa4j\x84\x0e\x82V\x8f\xf5\xb6G\xc1^\xd1\x12\xe6\xba%L\xa6J\xc2TNv:\t\xa4\xb8Ἅ\xc8\xf3\n\xbe\xec\xc0\x87h5#\xb3\x90\xacH\x1c\u007fd\xfcM\x14\x87\x13\xccȒ\xccg" +
			"~\xecӜ/\x17p\xfa\u007f\xf9c\xb9\x98\xba$\xc1A\x0e\x9e\x18\xd0^3\x06\xe6\xbdW\xd4\x13\xf1\xa7\xbd\xf9\xeb\xdb\xe1\xca\xfe7'\x0fc\xc3$pҬu[\xcbՉ\xb0\xfaY\xf3;\x00\x00\xff\xff\xba\x1a2\x1c\xde\b\x00\x00",
	},

	"/": {
		isDir: true,
		local: "public",
	},
}
