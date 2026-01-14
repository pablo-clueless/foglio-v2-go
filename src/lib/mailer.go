package lib

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"foglio/v2/src/config"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	templates    map[string]*template.Template
	mu           sync.RWMutex
	dialer       *gomail.Dialer
	templatesDir string
	pool         chan *gomail.SendCloser
	poolSize     int
	timeout      time.Duration
}

type EmailDto struct {
	To       []string
	Subject  string
	Template string
	Data     interface{}
}

type EmailConfig struct {
	SmtpHost        string
	SmtpPort        int
	SmtpUser        string
	SmtpPassword    string
	Timeout         time.Duration
	MaxRetries      int
	RetryDelay      time.Duration
	PoolSize        int
	UseTLS          bool
	InsecureSkipTLS bool
}

var (
	service *EmailService
	once    sync.Once
)

func getTemplatesDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "..", "templates")
}

func GetEmailService() *EmailService {
	once.Do(func() {
		emailConfig := EmailConfig{
			SmtpHost:        config.AppConfig.SmtpHost,
			SmtpPort:        config.AppConfig.SmtpPort,
			SmtpUser:        config.AppConfig.SmtpUser,
			SmtpPassword:    config.AppConfig.SmtpPassword,
			Timeout:         30 * time.Second,
			MaxRetries:      3,
			RetryDelay:      2 * time.Second,
			PoolSize:        5,
			UseTLS:          true,
			InsecureSkipTLS: false,
		}

		service = NewEmailService(emailConfig)
	})
	return service
}

func NewEmailService(cfg EmailConfig) *EmailService {
	dialer := gomail.NewDialer(
		cfg.SmtpHost,
		cfg.SmtpPort,
		cfg.SmtpUser,
		cfg.SmtpPassword,
	)

	if cfg.UseTLS {
		dialer.TLSConfig = &tls.Config{
			ServerName:         cfg.SmtpHost,
			InsecureSkipVerify: cfg.InsecureSkipTLS,
		}
	}

	dialer.LocalName = "localhost"

	es := &EmailService{
		templates:    make(map[string]*template.Template),
		templatesDir: getTemplatesDir(),
		dialer:       dialer,
		pool:         make(chan *gomail.SendCloser, cfg.PoolSize),
		poolSize:     cfg.PoolSize,
		timeout:      cfg.Timeout,
	}

	return es
}

func (es *EmailService) dialWithTimeout(ctx context.Context) (gomail.SendCloser, error) {
	type result struct {
		sender gomail.SendCloser
		err    error
	}

	resultCh := make(chan result, 1)

	go func() {
		sender, err := es.dialer.Dial()
		resultCh <- result{sender: sender, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("dial timeout: %w", ctx.Err())
	case res := <-resultCh:
		return res.sender, res.err
	}
}

func (es *EmailService) getOrCreateSender(ctx context.Context) (*gomail.SendCloser, error) {
	select {
	case sender := <-es.pool:
		return sender, nil
	default:
		sender, err := es.dialWithTimeout(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to dial SMTP: %w", err)
		}
		return &sender, nil
	}
}

func (es *EmailService) returnSender(sender *gomail.SendCloser) {
	select {
	case es.pool <- sender:
	default:
		err := (*sender).Close()
		if err != nil {
			fmt.Printf("failed to close sender: %v\n", err)
		}
	}
}

func (es *EmailService) getTemplate(templateName string) (*template.Template, error) {
	es.mu.RLock()
	tmpl, exists := es.templates[templateName]
	es.mu.RUnlock()

	if exists {
		return tmpl, nil
	}

	es.mu.Lock()
	defer es.mu.Unlock()

	if tmpl, exists = es.templates[templateName]; exists {
		return tmpl, nil
	}

	fullPath := filepath.Join(es.templatesDir, templateName+".html")
	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	es.templates[templateName] = tmpl
	return tmpl, nil
}

func (es *EmailService) renderTemplate(templateName string, data interface{}) (string, error) {
	tmpl, err := es.getTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (es *EmailService) testConnection(_ context.Context) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", es.dialer.Host, es.dialer.Port), es.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect to SMTP server: %w", err)
	}
	defer conn.Close()
	if es.dialer.TLSConfig != nil {
		if err := conn.(*net.TCPConn).SetDeadline(time.Now().Add(es.timeout)); err != nil {
			return fmt.Errorf("failed to set deadline: %w", err)
		}
	}
	return nil
}

func (es *EmailService) SendEmail(payload EmailDto) error {
	ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
	defer cancel()
	return es.SendEmailWithContext(ctx, payload)
}

func (es *EmailService) SendEmailWithContext(ctx context.Context, payload EmailDto) error {
	if err := es.testConnection(ctx); err != nil {
		return fmt.Errorf("SMTP connection test failed: %w", err)
	}

	html, err := es.renderTemplate(payload.Template, payload.Data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.AppConfig.AppEmail)
	msg.SetHeader("To", payload.To...)
	msg.SetHeader("Subject", payload.Subject)
	msg.SetBody("text/html", html)

	maxRetries := 3
	retryDelay := 2 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
			}
		}

		sender, err := es.getOrCreateSender(ctx)
		if err != nil {
			lastErr = err
			continue
		}

		sendDone := make(chan error, 1)
		go func() {
			sendDone <- gomail.Send(*sender, msg)
		}()

		select {
		case <-ctx.Done():
			err := (*sender).Close()
			if err != nil {
				fmt.Printf("failed to close sender: %v\n", err)
			}
			return fmt.Errorf("send timeout: %w", ctx.Err())
		case err := <-sendDone:
			if err != nil {
				err := (*sender).Close()
				if err != nil {
					fmt.Printf("failed to close sender: %v\n", err)
				}
				lastErr = fmt.Errorf("attempt %d failed: %w", attempt+1, err)
				continue
			}
			es.returnSender(sender)
			return nil
		}
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}

func (es *EmailService) SendEmailSimple(payload EmailDto) error {
	html, err := es.renderTemplate(payload.Template, payload.Data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.AppConfig.AppEmail)
	msg.SetHeader("To", payload.To...)
	msg.SetHeader("Subject", payload.Subject)
	msg.SetBody("text/html", html)

	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(2 * time.Second)
		}

		ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
		done := make(chan error, 1)

		go func() {
			done <- es.dialer.DialAndSend(msg)
		}()

		select {
		case <-ctx.Done():
			cancel()
			lastErr = fmt.Errorf("timeout on attempt %d", attempt+1)
			continue
		case err := <-done:
			cancel()
			if err != nil {
				lastErr = err
				continue
			}
			return nil
		}
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func (es *EmailService) SendBulkEmails(ctx context.Context, payloads []EmailDto) []error {
	errors := make([]error, len(payloads))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, es.poolSize)

	for i, payload := range payloads {
		wg.Add(1)
		go func(index int, p EmailDto) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			errors[index] = es.SendEmailWithContext(ctx, p)
		}(i, payload)
	}

	wg.Wait()
	return errors
}

func (es *EmailService) Close() {
	close(es.pool)
	for sender := range es.pool {
		err := (*sender).Close()
		if err != nil {
			fmt.Printf("failed to close sender: %v\n", err)
		}
	}
}

func SendEmail(payload EmailDto) error {
	return GetEmailService().SendEmail(payload)
}

func SendEmailWithTimeout(payload EmailDto, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return GetEmailService().SendEmailWithContext(ctx, payload)
}
