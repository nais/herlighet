package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var logger *zap.Logger

type pgStartupMessage []byte

// ref: https://www.postgresql.org/docs/9.3/protocol-error-fields.html
const (
	PG_ERR_SEVERITY byte = 'S'
	PG_ERR_SQLSTATE      = 'C'
	PG_ERR_MESSAGE       = 'M'
	PG_ERR_DETAIL        = 'D'
)

func main() {
	logger, _ = zap.NewDevelopment()
        smellTest()
        http.Handle("/metrics", promhttp.Handler())
        go http.ListenAndServe(":8080", nil)
	//logger, _ = zap.NewProduction()
	defer logger.Sync()
	logger.Info("Assuming the position...")
	ln, err := net.Listen("tcp", ":5432")
	if err != nil {
		logger.Panic("Could not start server", zap.Error(err))
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Error("Could not accept incoming request", zap.Error(err))
		}
		logger.Info("Incoming connection", zap.String("client_addr", conn.RemoteAddr().String()))
		go handleConnection(conn)
	}
}

func mkPgErrMsg(code byte, message string) []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(code)
	buf.WriteString(message)
	buf.WriteByte(0)
	return buf.Bytes()
}

func sendErrorPacket(error string, conn net.Conn) {
	errorBuf := &bytes.Buffer{}

	errorBuf.WriteByte('E')
	errorBuf.Write([]byte{0x00, 0x00, 0x00, 0x00}) // placeholder for nBytes
	errorBuf.Write(mkPgErrMsg(PG_ERR_SEVERITY, "FATAL"))
	errorBuf.Write(mkPgErrMsg(PG_ERR_MESSAGE, error))
	// The 'E' is not part of the packet length
	binary.BigEndian.PutUint32(errorBuf.Bytes()[1:], uint32(errorBuf.Len()))
	errorBuf.WriteByte(0)

	conn.Write(errorBuf.Bytes())
}

// Try to make a connection between the front and the rear end
func openWide(frontConn net.Conn, servAddr string) (net.Conn, error) {
	logger.Info("Opening wide...", zap.String("rear_end", servAddr))

	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	if err != nil {
		logger.Error("Failed to resolve server adress", zap.Error(err))
		return nil, errors.New("Failed to resolve server address")
	}

	rearConn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		logger.Error("Failed to connect to rear end", zap.Error(err))
		return nil, errors.New("Failed to connect with rear end")
	}

	return rearConn, nil
}

func handleConnection(frontConn net.Conn) {
	startupMessage, err := readStartupMessage(frontConn)
	if err != nil {
		logger.Error("Received invalid startup message", zap.Error(err))
		frontConn.Close()
		return
	}
	startupFields, _ := startupMessage.parse()
        if startupFields["encryptionRequest"] == "true" {
	    logger.Debug("Sending encryption downgrade request packet")
            frontConn.Write([]byte{0x4e})
            handleConnection(frontConn)
            return
        }
        servAddr, err := getRearEnd(startupFields["database"]) 
        if err != nil {
            sendErrorPacket(err.Error(), frontConn)
            return
        }
	rearConn, err := openWide(frontConn, servAddr)
	if err != nil {
		sendErrorPacket("Rear end not ready", frontConn)
		return
	}

	rearConn.Write(startupMessage)
	go io.Copy(frontConn, rearConn)
	go io.Copy(rearConn, frontConn)
}

func readStartupMessage(conn net.Conn) (pgStartupMessage, error) {
	var msgLen uint32
	clientReader := bufio.NewReader(conn)

	err := binary.Read(clientReader, binary.BigEndian, &msgLen)
	if err != nil {
		return nil, errors.New("Empty startup packet")
	}

	// ref: https://github.com/postgres/postgres/blob/master/src/include/libpq/pqcomm.h#L160
	if msgLen > 10000 {
		return nil, errors.New("Startup packet too big")
	}

	startupPacket := make([]byte, msgLen)
	binary.BigEndian.PutUint32(startupPacket[0:], msgLen)
	io.ReadFull(clientReader, startupPacket[4:])

	return startupPacket, nil
}

func (startupMessage *pgStartupMessage) parse() (map[string]string, error) {
	var msgLen uint32
	var protoVer uint32
	res := make(map[string]string)

	reader := bufio.NewReader(bytes.NewReader(*startupMessage))

	err := binary.Read(reader, binary.BigEndian, &msgLen)
	if err != nil {
		return res, errors.New("Failed to read startup message")
	}

	err = binary.Read(reader, binary.BigEndian, &protoVer)
	if err != nil {
		return res, errors.New("Could not read protocol version")
	}

	res["protoMajor"] = strconv.FormatUint(uint64(protoVer>>16), 10)
	res["protoMinor"] = strconv.FormatUint(uint64(protoVer&0xFFFF), 10)
        if (res["protoMajor"] == "1234") && (res["protoMinor"] == "5679") {
	    logger.Debug("Startup message is SSLRequest")
            res["encryptionRequest"] = "true"
            return res, nil
        } else if (res["protoMajor"] == "1234") && (res["protoMinor"] == "5680") {
	    logger.Debug("Startup message is GSSEncRequest")
            res["encryptionRequest"] = "true"
            return res, nil
        } else {
            res["encryptionRequest"] = "false"
        }

	for {
		key, _ := reader.ReadString(0)
		if len(key) == 1 {
			break
		} // "\x00" is the key-value-pair terminator
		value, _ := reader.ReadString(0)
		res[strings.TrimSuffix(key, "\x00")] = strings.TrimSuffix(value, "\x00")
	}

	return res, nil
}
