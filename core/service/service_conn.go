package service

import "turboengine/common/log"

type NetHandle struct {
	conn Conn
	svr  *service
}

func (h *NetHandle) Handle(conn Conn) {
	h.conn = conn
	log.Info("new conn ", conn.Addr())
}
