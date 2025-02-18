package telegramhtml

import "strings"

func ToTelegramHTML(html string) string {
	r1 := strings.NewReplacer("\t", "", "\n", "", "  ", "", "   ", "", "    ", "")
	text_content := r1.Replace(html)
	r2 := strings.NewReplacer("<br>", "\n", "<br/>", "\n")
	text_content = r2.Replace(text_content)
	return text_content
}
