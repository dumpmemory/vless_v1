package vless

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"sync"
	"time"
	"unsafe"

	"github.com/hahahrfool/v2simple/proxy"
)

func init() {
	proxy.RegisterServer(Name, NewVlessServer)
}

type Server struct {
	addr  string
	users []*proxy.ID

	userHashes   map[[16]byte]*proxy.ID
	userCRUMFURS map[[16]byte]*CRUMFURS
	mux4Hashes   sync.RWMutex
}

func NewVlessServer(url *url.URL) (proxy.Server, error) {

	addr := url.Host
	uuidStr := url.User.Username()
	id, err := proxy.NewID(uuidStr)
	if err != nil {
		return nil, err
	}
	s := &Server{
		addr:         addr,
		userHashes:   make(map[[16]byte]*proxy.ID),
		userCRUMFURS: make(map[[16]byte]*CRUMFURS),
	}
	s.users = append(s.users, id)

	for _, user := range s.users {
		s.userHashes[user.UUID] = user

	}

	return s, nil
}

func (s *Server) Name() string { return Name }

func (s *Server) Addr() string { return s.addr }

//see https://github.com/v2fly/v2ray-core/blob/master/proxy/vless/inbound/inbound.go
func (s *Server) Handshake(underlay net.Conn) (io.ReadWriter, *proxy.TargetAddr, error) {

	if err := underlay.SetReadDeadline(time.Now().Add(time.Second * 4)); err != nil {
		return nil, nil, err
	}
	defer underlay.SetReadDeadline(time.Time{})

	var auth [17]byte
	num, err := underlay.Read(auth[:])
	if err != nil {
		return nil, nil, err
	}

	if num < 17 {
		

		return nil, nil, errors.New("fallback, reason 1")

	}
	

	version := auth[0]
	if version > 1 {
		return nil, nil, errors.New("invalid request version")
	}

	idBytes := auth[1:17]

	s.mux4Hashes.RLock()

	thisUUIDBytes := *(*[16]byte)(unsafe.Pointer(&idBytes[0]))

	if user := s.userHashes[thisUUIDBytes]; user != nil {
		s.mux4Hashes.RUnlock()
	} else {
		s.mux4Hashes.RUnlock()
		return nil, nil, errors.New("invalid user")
	}

	if version == 0 {
		var addonLenBytes [1]byte
		_, err := underlay.Read(addonLenBytes[:])
		if err != nil {
			return nil, nil, err
		}
		addonLenByte := addonLenBytes[0]
		if addonLenByte != 0 {
			//v2ray的vless中没有对应的任何处理。
			log.Println("potential illegal client")
		}
	}

	var commandBytes [1]byte
	num, err = underlay.Read(commandBytes[:])

	if err != nil {
		return nil, nil, errors.New("fallback, reason 2")
	}

	commandByte := commandBytes[0]

	addr := &proxy.TargetAddr{}

	switch commandByte {
		case proxy.CmdMux: //TODO: 实际目前暂时v2simple还未实现mux，先这么写

		addr.Port = 0
		addr.Name = "v1.mux.cool"

	case Cmd_CRUMFURS:
		if version != 1 {
			return nil, nil, errors.New("在vless的vesion不为1时使用了 CRUMFURS 命令")
		}

		_, err = underlay.Write([]byte{CRUMFURS_ESTABLISHED})
		if err != nil {
			return nil, nil, err
		}

		addr.Name = CRUMFURS_Established_Str // 使用这个特殊的办法来告诉调用者，预留了 CRUMFURS 信道，防止其关闭上层连接导致 CRUMFURS 信道 被关闭。

		theCRUMFURS := &CRUMFURS{
			Conn: underlay,
		}

		s.mux4Hashes.Lock()

		s.userCRUMFURS[thisUUIDBytes] = theCRUMFURS

		s.mux4Hashes.Unlock()

		return nil, addr, nil

	case proxy.CmdTCP, proxy.CmdUDP:

		var portbs [2]byte

		num, err = underlay.Read(portbs[:])

		if err != nil || num != 2 {
			return nil, nil, errors.New("fallback, reason 3")
		}

		addr.Port = int(binary.BigEndian.Uint16(portbs[:]))

		if commandByte == proxy.CmdUDP {
			addr.IsUDP = true
		}

		var ip_or_domain_bytesLength byte = 0

		var addrTypeBytes [1]byte
		_, err = underlay.Read(addrTypeBytes[:])

		if err != nil {
			return nil, nil, errors.New("fallback, reason 4")
		}

		addrTypeByte := addrTypeBytes[0]

		switch addrTypeByte {
		case proxy.AtypIP4:

			ip_or_domain_bytesLength = net.IPv4len
			addr.IP = make(net.IP, net.IPv4len)
		case proxy.AtypDomain:
			

			var domainNameLenBytes [1]byte
			_, err = underlay.Read(domainNameLenBytes[:])

			if err != nil {
				return nil, nil, errors.New("fallback, reason 5")
			}

			domainNameLenByte := domainNameLenBytes[0]

			ip_or_domain_bytesLength = domainNameLenByte
		case proxy.AtypIP6:

			ip_or_domain_bytesLength = net.IPv6len
			addr.IP = make(net.IP, net.IPv6len)
		default:
			return nil, nil, fmt.Errorf("unknown address type %v", addrTypeByte)
		}

		ip_or_domain := make([]byte, ip_or_domain_bytesLength)

		_, err = underlay.Read(ip_or_domain[:])

		if err != nil {
			return nil, nil, errors.New("fallback, reason 6")
		}

		if addr.IP != nil {
			copy(addr.IP, ip_or_domain)
		} else {
			addr.Name = string(ip_or_domain)
		}

	default:
		return nil, nil, errors.New("invalid vless command")
	}

	return &UserConn{
		Conn:    underlay,
		uuid:    thisUUIDBytes,
		version: int(version),
	}, addr, nil

}
func (s *Server) Stop() {

}

func (s *Server) Get_CRUMFURS(id string) *CRUMFURS {
	bs, err := proxy.StrToUUID(id)
	if err != nil {
		return nil
	}
	return s.userCRUMFURS[bs]
}

type UserConn struct {
	net.Conn
	uuid         [16]byte
	convertedStr string
	version      int
}

func (uc *UserConn) GetProtocolVersion() int {
	return uc.version
}
func (uc *UserConn) GetIdentityStr() string {
	if uc.convertedStr == "" {
		uc.convertedStr = proxy.UUIDToStr(uc.uuid)
	}

	return uc.convertedStr
}

type CRUMFURS struct {
	net.Conn
}

func (c *CRUMFURS) WriteUDPResponse(a *net.UDPAddr, b []byte) (err error) {
	atype := proxy.AtypIP4
	if len(a.IP) > 4 {
		atype = proxy.AtypIP6
	}
	buf := &bytes.Buffer{}

	buf.WriteByte(atype)
	buf.Write(a.IP)
	buf.WriteByte(byte(int16(a.Port) >> 8))
	buf.WriteByte(byte(int16(a.Port) << 8 >> 8))
	buf.Write(b)

	_, err = c.Write(buf.Bytes())
	return
}
