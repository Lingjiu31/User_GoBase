package qqemail

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
)

type Service struct {
	host      string
	port      int
	email     string
	password  string
	emailSend *emailSend
}

type emailSend struct {
	subject string
	body    string
}

func NewService(host string, port int, email, password string) *Service {
	return &Service{
		host:     host,
		port:     port,
		email:    email,
		password: password,
	}
}

func (s *Service) Ready(subject string, body string) *Service {
	// 考虑并发问题, 发送前复制一份
	newService := *s
	newService.emailSend = &emailSend{
		subject: subject,
		body:    body,
	}
	return &newService
}

func (s *Service) Send(ctx context.Context, target ...string) error {
	if s.emailSend == nil {
		return fmt.Errorf("请先调用 Ready() 设置邮件内容")
	}
	header := make(map[string]string)
	header["From"] = "test " + "<" + s.email + ">"
	header["To"] = strings.Join(target, ", ")
	header["Subject"] = s.emailSend.subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	body := s.emailSend.body
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	// 注意：如果是 base64 传输，body 必须编码
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	auth := smtp.PlainAuth("", s.email, s.password, s.host)

	// 调用你写的 SendMailWithTLS
	err := s.sendMailWithTLS(
		fmt.Sprintf("%s:%d", s.host, s.port),
		auth,
		s.email,
		target,
		[]byte(message),
	)

	if err != nil {
		fmt.Println("发送失败:", err)
	} else {
		fmt.Println("发送成功!")
	}
	return err
}

func (s *Service) dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("tls.Dial Error:", err)
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func (s *Service) sendMailWithTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := s.dial(addr)
	if err != nil {
		log.Println("Create smtp client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
