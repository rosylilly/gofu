package gofu

import (
	"github.com/gographics/imagick/imagick"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RequestContext struct {
	Path    string
	err     error
	Image   *Image
	Request *http.Request
	Query   *Query
	Wand    *imagick.MagickWand
}

func (r *RequestContext) Execute(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
		r.writeResponse(w)
	}()

	bench("init", func() {
		r.init(req)
	})
	bench("getImage", func() {
		r.getImage()
	})
	bench("parseQuery", func() {
		r.parseQuery()
	})

	bench("processImage", func() {
		r.processImage()
	})
}

func (r *RequestContext) init(req *http.Request) {
	r.Request = req
	r.Path = req.URL.Path[3:]
	r.Image = nil
	r.Query = nil
	r.err = nil
}

func (r *RequestContext) getImage() {
	r.Image, r.err = GetImage(r.Path)
	if r.err != nil {
		panic(r.err)
	}
}

func (r *RequestContext) parseQuery() {
	r.Query = ParseQuery(r.Request.URL.Query())
}

func (r *RequestContext) blob() []byte {
	if r.err != nil {
		return []byte(string(r.err.Error()))
	}

	var blob []byte

	if r.Image.Processed {
		blob = r.Image.Blob
	} else {
		blob, r.err = ioutil.ReadFile(r.Image.Path)
		if r.err != nil {
			return []byte(string(r.err.Error()))
		}
	}

	return blob
}

func (r *RequestContext) processImage() {
	r.Image.Process(r.Wand, r.Query)
}

func (r *RequestContext) writeResponse(w http.ResponseWriter) {
	blob := r.blob()

	if r.err == nil {
		w.WriteHeader(http.StatusOK)
	} else if r.err.Error() == "The specified key does not exist." {
		w.WriteHeader(http.StatusNotFound)
		blob = make([]byte, 0)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Add("Content-Type", http.DetectContentType(blob))
	w.Header().Add("Content-Length", strconv.Itoa(len(blob)))
	w.Write(blob)
}
