package urlshort

import "net/http"

type ServerWriter struct {
	w           http.ResponseWriter
	location    string
	wroteHeader bool
}

func (s *ServerWriter) Header() http.Header {
	return s.w.Header()
}

func (s *ServerWriter) WriteHeader(code int) http.Header {
	if s.wroteHeader == false {
		s.w.Header().Set("Location", s.location)
		s.wroteHeader = true
	}
	s.w.WriteHeader(code)
	return s.w.Header()
}

func (s *ServerWriter) Write(b []byte) (int, error) {
	if s.wroteHeader == false {
		// We hit this case if user never calls WriteHeader (default 200)
		s.w.Header().Set("Location", s.location)
		s.wroteHeader = true
	}
	return s.w.Write(b)
}
