package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Header struct {
	Key string
	Value string
}

type Context struct {
	conn net.Conn
	Request *Request
	Response *Response
}

type Request struct {
	Method string
	Path string
	Headers []Header
}

type Response struct {
	headers []Header
	body []byte
	statusCode int
}

func (res *Response) addHeader(headers ...Header) {
	res.headers = append(res.headers, headers...)
}

func (res *Response) setStatusCode(code int) {
	res.statusCode = code
}

func (ctx *Context) response(body []byte)  {
	headers := ""
	for _, header := range ctx.Response.headers {
		headers += fmt.Sprintf("%s: %s\n", header.Key, header.Value)
	}
	r := []byte(fmt.Sprintf("HTTP/1.1 %d OK\n%s\n%s", ctx.Response.statusCode, headers, string(body)))
	ctx.conn.Write(r)
}

func main()  {
	l, err := net.Listen("tcp", ":80")
	if err != nil { log.Println(err) }

	for {
		conn, err := l.Accept()
		if err != nil { log.Println(err) }
		data := make([]byte, 1024)

		_, err = conn.Read(data)
		if err != nil { continue }
		request := string(data)

		if len(strings.ReplaceAll(request, "\x00", "")) == 0 {
			log.Println("ПУСТОЙ ЗАПРОС БЛИН БЛЯТЬ")
			continue
		}

		main := strings.Split(strings.Split(request, "\r\n")[0], " ")

		baseHeaders := make([]Header, 0)
		baseHeaders = append(baseHeaders, Header{
			Key:   "Content-Type",
			Value: "text/html; charset=utf-8",
		}, Header{
			Key:   "Server",
			Value: "TestHTTP",
		}, Header{
			Key:   "Connection",
			Value: "close",
		})
		ctx := &Context{
			conn: conn,
			Request: &Request{
				Method: main[0],
				Path: main[1],
			},
			Response: &Response{
				statusCode: 200,
				headers: baseHeaders,
			},
		}

		for _, header := range strings.Split(request, "\r\n")[1:] {
			if len(header) == 0 { continue }
			parsedHeader := strings.Split(header, ":")
			if len(parsedHeader) != 2 { continue }
			ctx.Request.Headers = append(ctx.Request.Headers, Header{
				Key: parsedHeader[0],
				Value: parsedHeader[1][1:],
			})
		}

		ctx.response([]byte("<html><head><title>Hello</title></head><body>Test</body></html>"))
		conn.Close()
	}
}