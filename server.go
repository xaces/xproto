package xproto

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/xaces/xproto/internal"
)

// Serve 服务
type Server struct {
	listener *net.TCPListener //
	wg       sync.WaitGroup   //
	adapters []AdapterHandler
	Handler  Handler
	Result   *Result
	*Options
}

var defaultAdapters []AdapterHandler

func Use(codec AdapterHandler) {
	defaultAdapters = append(defaultAdapters, codec)
}

func (o *Server) isVaildOptions() error {
	if o.Options == nil || o.Port == 0 {
		return errors.New("invalid options")
	}
	if o.Host == "" {
		o.Host = internal.LocalIPAddr()
	}
	if o.RecvTimeout == 0 {
		o.RecvTimeout = 30
	}
	if o.RequestTimeout == 0 {
		o.RequestTimeout = 10
	}
	if o.adapters == nil {
		o.adapters = defaultAdapters
	}
	return nil
}

// NewServe 默认服务
func NewServer(opts *Options, adapters ...AdapterHandler) (*Server, error) {
	s := &Server{Options: opts, adapters: adapters}
	if err := s.isVaildOptions(); err != nil {
		return s, err
	}
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", opts.Port)) //获取一个tcpAddr
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return s, err
	}
	s.listener = listener
	s.Result = NewResult(s.RequestTimeout)
	return s, nil
}

// ListenTCPAndServe start server
func (s *Server) ListenTCPAndServe() {
	log.Println("xproto server listening at:", s.listener.Addr().String())
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		ctx := newConn(conn, s)
		if nil == ctx {
			continue
		}
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			ctx.start()
		}()
	}
}

// Release 关闭服务
func (s *Server) Release() {
	if s.listener != nil {
		s.listener.Close()
		SyncStopAll()
	}
	s.wg.Wait()
}
