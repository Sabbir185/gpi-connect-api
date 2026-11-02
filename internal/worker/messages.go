package worker

import (
	"fmt"
	"html"
	"strings"
)

type welcomeCopy struct {
	PushTitle    string
	PushBody     string
	EmailSubject string
	Greeting     string
	GreetingName string
	EmailBody    string
	Signature    string
	RTL          bool
}

var welcomeMessages = map[string]welcomeCopy{
	"en": {
		PushTitle:    "Welcome to GPI Connect",
		PushBody:     "Your account is ready. Finish your profile to get started.",
		EmailSubject: "Welcome to GPI Connect",
		Greeting:     "Welcome!",
		GreetingName: "Welcome, %s!",
		EmailBody:    "Your GPI Connect account is ready. Finish your profile in the app to join your team, mark attendance and stay connected with your people.",
		Signature:    "The GPI Connect team",
	},
	"bn": {
		PushTitle:    "GPI Connect-এ স্বাগতম",
		PushBody:     "আপনার অ্যাকাউন্ট তৈরি হয়েছে। শুরু করতে প্রোফাইলটি সম্পূর্ণ করুন।",
		EmailSubject: "GPI Connect-এ স্বাগতম",
		Greeting:     "স্বাগতম!",
		GreetingName: "স্বাগতম, %s!",
		EmailBody:    "আপনার GPI Connect অ্যাকাউন্ট তৈরি হয়েছে। অ্যাপে প্রোফাইলটি সম্পূর্ণ করে আপনার টিমে যোগ দিন, হাজিরা দিন এবং সবার সাথে যুক্ত থাকুন।",
		Signature:    "GPI Connect টিম",
	},
	"ar": {
		PushTitle:    "مرحبًا بك في GPI Connect",
		PushBody:     "حسابك جاهز. أكمل ملفك الشخصي للبدء.",
		EmailSubject: "مرحبًا بك في GPI Connect",
		Greeting:     "مرحبًا بك!",
		GreetingName: "مرحبًا بك يا %s!",
		EmailBody:    "حسابك في GPI Connect جاهز. أكمل ملفك الشخصي في التطبيق للانضمام إلى فريقك وتسجيل الحضور والبقاء على تواصل مع زملائك.",
		Signature:    "فريق GPI Connect",
		RTL:          true,
	},
}

// messagesFor picks by language subtag ("bn-BD" -> "bn") and falls back to English.
func messagesFor(locale string) welcomeCopy {
	language, _, _ := strings.Cut(strings.ToLower(locale), "-")
	if text, ok := welcomeMessages[language]; ok {
		return text
	}
	return welcomeMessages["en"]
}

func welcomeEmailHTML(text welcomeCopy, name string) string {
	greeting := text.Greeting
	if name != "" {
		greeting = fmt.Sprintf(text.GreetingName, html.EscapeString(name))
	}
	dir, align := "ltr", "left"
	if text.RTL {
		dir, align = "rtl", "right"
	}
	return fmt.Sprintf(`<!doctype html>
<html dir="%[1]s">
<body style="margin:0;padding:24px;background:#F6F1EA;font-family:Arial,Helvetica,sans-serif;color:#1F1F1F;">
  <div style="max-width:520px;margin:0 auto;background:#FFFFFF;border-radius:16px;padding:32px;text-align:%[2]s;">
    <div style="font-size:22px;font-weight:bold;color:#CF4A0C;margin-bottom:24px;">GPI Connect</div>
    <h1 style="font-size:20px;margin:0 0 12px;">%[3]s</h1>
    <p style="font-size:15px;line-height:24px;margin:0 0 24px;">%[4]s</p>
    <p style="font-size:14px;color:#6B6B6B;margin:0;">%[5]s</p>
  </div>
</body>
</html>`, dir, align, greeting, html.EscapeString(text.EmailBody), html.EscapeString(text.Signature))
}
