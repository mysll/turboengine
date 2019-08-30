package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Mailbox uint64

const (
	ID_MAX = 0xFFFFFFFFFFFF
)

func generate(appId uint16, flag int8, id uint64) Mailbox {
	return Mailbox(((uint64(appId) << 52) & 0xFFF0000000000000) | ((uint64(flag) & 0xF) << 48) | (id & ID_MAX))
}

func (m Mailbox) String() string {
	return fmt.Sprintf("mb://%x", uint64(m))
}

func (m Mailbox) Service() Mailbox {
	return m & 0xFFF0000000000000
}

// 是否为空
func (m Mailbox) IsNil() bool {
	return m == 0
}

// 获取服务编号
func (m Mailbox) ServiceId() uint16 {
	return uint16((m & 0xFFF0000000000000) >> 52)
}

// 获取标志位
func (m Mailbox) Flag() int8 {
	return int8((m >> 48) & 0xF)
}

// 获取id
func (m Mailbox) Id() uint64 {
	return uint64(m & ID_MAX)
}

// 获取uid
func (m Mailbox) Uid() uint64 {
	return uint64(m)
}

// 通过字符串生成mailbox
func NewMailboxFromStr(mb string) (Mailbox, error) {
	mbox := Mailbox(0)
	if !strings.HasPrefix(mb, "mb://") {
		return mbox, errors.New("mailbox string error")
	}
	vals := strings.Split(mb, "/")
	if len(vals) != 3 {
		return mbox, errors.New("mailbox string error")
	}

	var val uint64
	var err error

	val, err = strconv.ParseUint(vals[2], 16, 64)
	if err != nil {
		return mbox, err
	}
	mbox = Mailbox(val)
	return mbox, nil
}

// 通过uid生成mailbox
func NewMailboxFromUid(val uint64) Mailbox {
	return Mailbox(val)
}

// 通过服务编号获取mailbox
func GetServiceMailbox(appId uint16) Mailbox {
	m := generate(appId, 0, 0)
	return m
}

// 生成mailbox
func NewMailbox(appId uint16, flag int8, id uint64) Mailbox {
	if id > ID_MAX {
		panic("id is wrong")
	}
	m := generate(appId, flag, id)
	return m
}
