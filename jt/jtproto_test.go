package jt

import (
	"log"
	"testing"

	"github.com/xaces/xproto"
)

func TestP(t *testing.T) {
	s, err := xproto.NewServer(&xproto.Options{
		Port: 22000,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Release()
	s.Handler.Status = xproto.LogStatus
	s.ListenTCPAndServe()
}
