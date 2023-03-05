package main

import (
	"fmt"
	"net"

	vnc "github.com/kward/go-vnc"
)

var vncserver *vnc.ServerConn

func main() {
	// TCP 서버 생성
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Listening on 0.0.0.0:8000")

	// 연결 대기
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}

		// 암호 검증
		go func(conn net.Conn) {
			// 암호 수신
			passwordBytes := make([]byte, 1024)
			n, err := conn.Read(passwordBytes)
			if err != nil {
				fmt.Println("Error reading password:", err.Error())
				conn.Close()
				return
			}
			password := string(passwordBytes[:n])

			// 암호 검증
			if password != "mypassword\n" {
				conn.Write([]byte("Invalid password"))
				conn.Close()
				return
			}

			// 연결 허용
			fmt.Println("Connected by", conn.RemoteAddr().String())

			// VNC 서버 시작
			server := vnc.NewServer(true)
			vncserver, _ = server.ListenAndServe(":5900")

			// 마우스 이벤트 및 키 이벤트 처리 함수 등록
			vncserver.HandleKeyboardEvent = handleKeyEvent
			vncserver.HandlePointerEvent = handleMouseEvent

			fmt.Println("Starting VNC server...")
		}(conn)
	}
}

func handleMouseEvent(event vnc.PointerEvent) {
	// 마우스 이벤트 처리
	if event.Buttons&vnc.PointerButton1Mask != 0 {
		fmt.Println("Mouse button 1 clicked")
		// 마우스 왼쪽 버튼이 클릭되었을 때, 클릭 이벤트 발생
		vncserver.SendPointerEvent(vnc.PointerEvent{X: event.X, Y: event.Y, Buttons: 1})
	} else {
		// 마우스 왼쪽 버튼이 클릭되지 않았을 때, 이동 이벤트 발생
		vncserver.SendPointerEvent(vnc.PointerEvent{X: event.X, Y: event.Y})
	}
}

func handleKeyEvent(event vnc.KeyEvent) {
	// 키 이벤트 처리
	if event.Down {
		fmt.Println("Key pressed:", event.Key)
		vncserver.SendKeyEvent(event)
	} else {
		fmt.Println("Key released:", event.Key)
		vncserver.SendKeyEvent(event)
	}
}
