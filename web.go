package raccoon

import (
	"net"
	"net/http"
	"path/filepath"
)

func getWebSource(root string, addr string) (string, net.Listener, error) {
	webroot := filepath.Join(GetProgramPath(), root)
	listener, err := getListener(addr)
	return webroot, listener, err
}
func getListener(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func webServer(root string, listener net.Listener) error {
	return http.Serve(listener, http.FileServer(
		http.Dir(root),
	))
}
