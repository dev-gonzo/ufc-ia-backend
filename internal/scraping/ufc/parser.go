package ufc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseAthleteDetails(
	html string,
	athleteURL string,
	athleteSlug string,
) (*AthleteDetails, string, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}

	d := &AthleteDetails{
		AthleteURL:  athleteURL,
		AthleteSlug: athleteSlug,
	}

	ogImage := strings.TrimSpace(document.Find(`meta[property="og:image"]`).AttrOr("content", ""))
	if ogImage != "" {
		ogImage = strings.TrimSpace(strings.Split(ogImage, "?")[0])
	}

	document.Find(".c-bio__field").Each(func(_ int, s *goquery.Selection) {
		label := normalizeWhitespace(s.Find(".c-bio__label").First().Text())
		value := normalizeWhitespace(s.Find(".c-bio__text").First().Text())
		if label == "" || value == "" {
			return
		}

		label = strings.ToLower(label)

		switch label {
		case "status":
			v := strings.ToLower(value)
			if strings.Contains(v, "ativo") {
				b := true
				d.IsActive = &b
			} else if strings.Contains(v, "inativo") {
				b := false
				d.IsActive = &b
			}
		case "cidade natal", "hometown":
			v := value
			d.Hometown = &v
		case "estilo de luta", "fighting style":
			v := value
			d.FightingStyle = &v
		case "altura", "height":
			v := value
			d.Height = &v
		case "peso", "weight":
			v := value
			d.Weight = &v
		case "estreia no ufc", "ufc debut":
			v := value
			d.UFCDebut = &v
		case "envergadura", "reach":
			v := value
			d.Reach = &v
		case "alcance das pernas", "leg reach":
			v := value
			d.LegReach = &v
		}
	})

	if v := extractBlockText(document, ".field--name-qna-facts"); v != "" {
		d.FighterFacts = &v
	}
	if v := extractBlockText(document, ".field--name-qna-ufc"); v != "" {
		d.UFCHistory = &v
	}
	if v := extractBlockText(document, ".field--name-qna"); v != "" {
		d.QA = &v
	}

	return d, ogImage, nil
}

func extractBlockText(document *goquery.Document, selector string) string {
	sel := document.Find(selector).First()
	if sel.Length() == 0 {
		return ""
	}

	html, err := sel.Html()
	if err != nil {
		return normalizeWhitespace(sel.Text())
	}

	replacer := strings.NewReplacer(
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"</p>", "\n",
		"</li>", "\n",
		"&nbsp;", " ",
	)

	html = replacer.Replace(html)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return normalizeWhitespace(sel.Text())
	}

	text := doc.Text()
	text = normalizeNewlines(text)
	return strings.TrimSpace(text)
}

var newlineRe = regexp.MustCompile(`\n{3,}`)

func normalizeNewlines(s string) string {
	lines := strings.Split(s, "\n")
	var out []string
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			out = append(out, "")
			continue
		}
		out = append(out, normalizeWhitespace(ln))
	}
	res := strings.Join(out, "\n")
	res = newlineRe.ReplaceAllString(res, "\n\n")
	return strings.TrimSpace(res)
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
