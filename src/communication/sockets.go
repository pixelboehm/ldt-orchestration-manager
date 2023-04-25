package communication

import (
	"fmt"
	"net"
)

func SendResultToSocket(out net.Conn, res string) {
	_, err := out.Write([]byte(res))
	if err != nil {
		panic(fmt.Sprintf("error writing to socket: %v\n", err))
	}
}

func GetCommandFromSocket(in net.Conn) string {
	for {
		buf := make([]byte, 512)
		nr, err := in.Read(buf)
		if err != nil {
			return "help"
		}

		data := buf[0:nr]
		return string(data)
	}
}
