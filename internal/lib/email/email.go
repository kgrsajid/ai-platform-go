package email

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type Config struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

type Sender struct {
	cfg Config
}

func New(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) SendResetCode(to, code string) error {
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.SMTPHost)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
<body style="margin:0;padding:0;background:#f4f6f9;font-family:Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f9;padding:40px 0;">
    <tr>
      <td align="center">
        <table width="480" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);">

          <!-- Header -->
          <tr>
            <td style="background:#4F46E5;padding:32px 40px;text-align:center;">
              <h1 style="margin:0;color:#ffffff;font-size:22px;letter-spacing:0.5px;">Password Reset</h1>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:40px 40px 32px;">
              <p style="margin:0 0 16px;color:#374151;font-size:15px;line-height:1.6;">
                Hi, we received a request to reset your password. Use the code below:
              </p>

              <!-- Code block -->
              <div style="margin:24px 0;text-align:center;">
                <span style="display:inline-block;background:#F3F4F6;border:2px dashed #4F46E5;border-radius:10px;padding:18px 48px;font-size:36px;font-weight:700;letter-spacing:10px;color:#1F2937;font-family:'Courier New',monospace;">
                  %s
                </span>
              </div>

              <p style="margin:0 0 8px;color:#6B7280;font-size:13px;text-align:center;">
                ⏱ This code expires in <strong>15 minutes</strong>.
              </p>
              <p style="margin:16px 0 0;color:#9CA3AF;font-size:12px;text-align:center;">
                If you didn't request a password reset, you can safely ignore this email.
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background:#F9FAFB;padding:20px 40px;text-align:center;border-top:1px solid #E5E7EB;">
              <p style="margin:0;color:#9CA3AF;font-size:12px;">
                This is an automated message — please do not reply.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, code)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Password Reset Code\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.cfg.From, to, html,
	)

	addr := s.cfg.SMTPHost + ":" + strconv.Itoa(s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg))
}
