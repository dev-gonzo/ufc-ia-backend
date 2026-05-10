package ufcstats

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	scoreRe  = regexp.MustCompile(`(\d+)\s*-\s*(\d+)`)
	landedRe = regexp.MustCompile(`(\d+)\s+of\s+(\d+)`)
	roundRe  = regexp.MustCompile(`(?i)round\s+(\d+)`)
)

func parseFightDetails(html string, url string) (*FightDetailsScrape, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}

	fight := &Fight{URL: url}
	result := &FightDetailsScrape{
		Fight: fight,
		RoundStats: FightRoundStats{
			Red:  []RoundFighterStat{},
			Blue: []RoundFighterStat{},
		},
	}

	titleSel := document.Find(".b-fight-details__fight-title").First()
	titleText := normalizeWhitespace(titleSel.Text())
	bonusTypes := parseBonusTypes(titleSel)
	if strings.Contains(strings.ToLower(titleText), "title bout") || titleSel.Find("img[src*='belt']").Length() > 0 {
		result.IsTitleBout = true
	}

	if titleText != "" && strings.TrimSpace(fight.WeightClass) == "" {
		t := strings.TrimSpace(strings.TrimSuffix(titleText, "Bout"))
		fight.WeightClass = t
	}

	document.Find(".b-fight-details__persons .b-fight-details__person").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".b-fight-details__person-name a").First().Text())
		href, _ := s.Find(".b-fight-details__person-name a").First().Attr("href")
		status := strings.TrimSpace(s.Find(".b-fight-details__person-status").First().Text())

		f := &Fighter{Name: name, URL: strings.TrimSpace(href), Record: ""}
		if i == 0 {
			result.Red = f
			if status == "W" {
				fight.Winner = "red"
			}
		}
		if i == 1 {
			result.Blue = f
			if status == "W" {
				fight.Winner = "blue"
			}
		}
	})
	result.Bonuses = buildBonuses(bonusTypes, fight.Winner)

	document.Find(".b-fight-details__text-item, .b-fight-details__text-item_first").Each(func(_ int, s *goquery.Selection) {
		t := normalizeWhitespace(s.Text())
		switch {
		case strings.HasPrefix(t, "Weight:"):
			if strings.TrimSpace(fight.WeightClass) == "" {
				fight.WeightClass = strings.TrimSpace(strings.TrimPrefix(t, "Weight:"))
			}
		case strings.HasPrefix(t, "Method:"):
			fight.Method = strings.TrimSpace(strings.TrimPrefix(t, "Method:"))
		case strings.HasPrefix(t, "Round:"):
			roundStr := strings.TrimSpace(strings.TrimPrefix(t, "Round:"))
			if n, err := strconv.Atoi(roundStr); err == nil {
				fight.Round = n
			}
		case strings.HasPrefix(t, "Time:"):
			fight.Time = strings.TrimSpace(strings.TrimPrefix(t, "Time:"))
		case strings.HasPrefix(t, "Referee:"):
			result.RefereeName = strings.TrimSpace(strings.TrimPrefix(t, "Referee:"))
		}
	})

	document.Find(".b-fight-details__text").Each(func(_ int, p *goquery.Selection) {
		if !strings.Contains(normalizeWhitespace(p.Text()), "Details:") {
			return
		}

		p.Find(".b-fight-details__text-item").Each(func(_ int, item *goquery.Selection) {
			name := normalizeWhitespace(item.Find("span").First().Text())
			text := normalizeWhitespace(item.Text())
			if name == "" {
				name = strings.TrimSpace(strings.ReplaceAll(text, scoreRe.FindString(text), ""))
				name = strings.TrimSpace(strings.TrimSuffix(name, "-"))
				name = strings.TrimSpace(name)
			}

			m := scoreRe.FindStringSubmatch(text)
			if name == "" || len(m) != 3 {
				return
			}
			red, _ := strconv.Atoi(m[1])
			blue, _ := strconv.Atoi(m[2])
			result.Judges = append(result.Judges, JudgeScore{Name: name, RedScore: red, BlueScore: blue})
		})
	})

	stats, rounds := parseFightRoundStats(document, fight.Round)
	result.RoundStats = stats
	result.Rounds = rounds
	if result.Rounds == 0 && fight.Round > 0 {
		result.Rounds = fight.Round
	}

	return result, nil
}

func parseFightRoundStats(document *goquery.Document, fallbackRounds int) (FightRoundStats, int) {
	redByRound := map[int]*RoundFighterStat{}
	blueByRound := map[int]*RoundFighterStat{}
	maxRound := 0

	document.Find("table").Each(func(_ int, table *goquery.Selection) {
		header := strings.ToLower(normalizeWhitespace(table.Find("th").Text()))
		isMainStats :=
			strings.Contains(header, "kd") &&
				strings.Contains(header, "sig. str") &&
				strings.Contains(header, "total str") &&
				strings.Contains(header, "sub. att") &&
				strings.Contains(header, "ctrl") &&
				!strings.Contains(header, "head") &&
				!strings.Contains(header, "body") &&
				!strings.Contains(header, "leg") &&
				!strings.Contains(header, "distance")
		isSigStats :=
			strings.Contains(header, "sig. str") &&
				strings.Contains(header, "head") &&
				strings.Contains(header, "body") &&
				strings.Contains(header, "leg") &&
				strings.Contains(header, "distance") &&
				strings.Contains(header, "clinch") &&
				strings.Contains(header, "ground")

		if !isMainStats && !isSigStats {
			return
		}

		if table.HasClass("js-fight-table") {
			table.Find("thead.b-fight-details__table-row_type_head").Each(func(_ int, thead *goquery.Selection) {
				headText := normalizeWhitespace(thead.Text())
				m := roundRe.FindStringSubmatch(headText)
				if len(m) != 2 {
					return
				}
				round, err := strconv.Atoi(m[1])
				if err != nil || round <= 0 {
					return
				}
				if round > maxRound {
					maxRound = round
				}

				row := findNextRoundRow(thead)
				if row == nil || row.Length() == 0 {
					return
				}

				if isMainStats {
					addMainStatsRow(redByRound, blueByRound, round, row)
				}
				if isSigStats {
					addSigStatsRow(redByRound, blueByRound, round, row)
				}
			})
			return
		}

		tbody := table.Find("tbody").First()
		if tbody.Length() == 0 {
			return
		}
		row := tbody.Find("tr.b-fight-details__table-row").First()
		if row.Length() == 0 {
			return
		}

		if isMainStats {
			addMainStatsRow(redByRound, blueByRound, 0, row)
		}
		if isSigStats {
			addSigStatsRow(redByRound, blueByRound, 0, row)
		}
	})

	if maxRound == 0 && fallbackRounds > 0 {
		maxRound = fallbackRounds
	}

	for r := 0; r <= maxRound; r++ {
		ensureRoundStat(redByRound, r)
		ensureRoundStat(blueByRound, r)
	}

	roundKeys := map[int]bool{}
	for r := range redByRound {
		roundKeys[r] = true
	}
	for r := range blueByRound {
		roundKeys[r] = true
	}

	var rounds []int
	for r := range roundKeys {
		rounds = append(rounds, r)
	}
	sort.Ints(rounds)

	red := make([]RoundFighterStat, 0, len(rounds))
	blue := make([]RoundFighterStat, 0, len(rounds))
	for _, r := range rounds {
		if v := redByRound[r]; v != nil {
			red = append(red, *v)
		} else {
			red = append(red, RoundFighterStat{Round: r})
		}
		if v := blueByRound[r]; v != nil {
			blue = append(blue, *v)
		} else {
			blue = append(blue, RoundFighterStat{Round: r})
		}
	}

	return FightRoundStats{Red: red, Blue: blue}, maxRound
}

func findNextRoundRow(thead *goquery.Selection) *goquery.Selection {
	next := thead.Next()
	for next != nil && next.Length() > 0 {
		name := goquery.NodeName(next)
		if name == "tr" && next.HasClass("b-fight-details__table-row") {
			return next
		}
		if name == "tbody" {
			row := next.Find("tr.b-fight-details__table-row").First()
			if row.Length() > 0 {
				return row
			}
		}
		next = next.Next()
	}
	return nil
}

func ensureRoundStat(m map[int]*RoundFighterStat, round int) *RoundFighterStat {
	if v, ok := m[round]; ok && v != nil {
		return v
	}
	v := &RoundFighterStat{Round: round}
	m[round] = v
	return v
}

func addMainStatsRow(redByRound map[int]*RoundFighterStat, blueByRound map[int]*RoundFighterStat, round int, row *goquery.Selection) {
	tds := row.Find("td")
	if tds.Length() < 10 {
		return
	}

	red := ensureRoundStat(redByRound, round)
	blue := ensureRoundStat(blueByRound, round)

	red.KD, blue.KD = parseIntPair(tds.Eq(1))
	red.SigLanded, red.SigAttempted, blue.SigLanded, blue.SigAttempted = parseLandedPair(tds.Eq(2))
	red.TotalLanded, red.TotalAttempted, blue.TotalLanded, blue.TotalAttempted = parseLandedPair(tds.Eq(4))
	red.TDLanded, red.TDAttempted, blue.TDLanded, blue.TDAttempted = parseLandedPair(tds.Eq(5))
	red.SubAtt, blue.SubAtt = parseIntPair(tds.Eq(7))
	red.Rev, blue.Rev = parseIntPair(tds.Eq(8))
	red.CTRL, blue.CTRL = parseStringPair(tds.Eq(9))
}

func addSigStatsRow(redByRound map[int]*RoundFighterStat, blueByRound map[int]*RoundFighterStat, round int, row *goquery.Selection) {
	tds := row.Find("td")
	if tds.Length() < 9 {
		return
	}

	red := ensureRoundStat(redByRound, round)
	blue := ensureRoundStat(blueByRound, round)

	if red.SigLanded == 0 && red.SigAttempted == 0 && blue.SigLanded == 0 && blue.SigAttempted == 0 {
		red.SigLanded, red.SigAttempted, blue.SigLanded, blue.SigAttempted = parseLandedPair(tds.Eq(1))
	}

	red.HeadLanded, red.HeadAttempted, blue.HeadLanded, blue.HeadAttempted = parseLandedPair(tds.Eq(3))
	red.BodyLanded, red.BodyAttempted, blue.BodyLanded, blue.BodyAttempted = parseLandedPair(tds.Eq(4))
	red.LegLanded, red.LegAttempted, blue.LegLanded, blue.LegAttempted = parseLandedPair(tds.Eq(5))
	red.DistanceLanded, red.DistanceAttempted, blue.DistanceLanded, blue.DistanceAttempted = parseLandedPair(tds.Eq(6))
	red.ClinchLanded, red.ClinchAttempted, blue.ClinchLanded, blue.ClinchAttempted = parseLandedPair(tds.Eq(7))
	red.GroundLanded, red.GroundAttempted, blue.GroundLanded, blue.GroundAttempted = parseLandedPair(tds.Eq(8))
}

func parseStringPair(td *goquery.Selection) (string, string) {
	a := normalizeWhitespace(td.Find("p.b-fight-details__table-text").Eq(0).Text())
	b := normalizeWhitespace(td.Find("p.b-fight-details__table-text").Eq(1).Text())
	return a, b
}

func parseIntPair(td *goquery.Selection) (int, int) {
	aStr, bStr := parseStringPair(td)
	a, _ := strconv.Atoi(strings.TrimSpace(aStr))
	b, _ := strconv.Atoi(strings.TrimSpace(bStr))
	return a, b
}

func parseLandedPair(td *goquery.Selection) (int, int, int, int) {
	aStr, bStr := parseStringPair(td)
	aL, aA := parseLanded(aStr)
	bL, bA := parseLanded(bStr)
	return aL, aA, bL, bA
}

func parseLanded(s string) (int, int) {
	if m := landedRe.FindStringSubmatch(s); len(m) == 3 {
		l, _ := strconv.Atoi(m[1])
		a, _ := strconv.Atoi(m[2])
		return l, a
	}
	return 0, 0
}

func parseBonusTypes(titleSel *goquery.Selection) []string {
	seen := map[string]bool{}
	var out []string
	has := func(substr string) bool {
		return titleSel.Find("img").FilterFunction(func(_ int, s *goquery.Selection) bool {
			return strings.Contains(strings.ToLower(s.AttrOr("src", "")), substr)
		}).Length() > 0
	}

	if has("perf.png") {
		seen["PERF"] = true
	}
	if has("fight.png") {
		seen["FIGHT"] = true
	}
	if has("sub.png") {
		seen["SUB"] = true
	}
	if has("ko.png") {
		seen["KO"] = true
	}

	for _, t := range []string{"FIGHT", "PERF", "SUB", "KO"} {
		if seen[t] {
			out = append(out, t)
		}
	}
	return out
}

func buildBonuses(types []string, winner string) []FightBonus {
	var out []FightBonus
	winner = strings.TrimSpace(winner)

	add := func(t string, recipients []string) {
		for _, r := range recipients {
			out = append(out, FightBonus{Type: t, Recipient: r})
		}
		if len(recipients) == 0 {
			out = append(out, FightBonus{Type: t})
		}
	}

	for _, t := range types {
		switch t {
		case "FIGHT":
			add(t, []string{"red", "blue"})
		case "PERF", "SUB", "KO":
			if winner == "red" || winner == "blue" {
				add(t, []string{winner})
			} else {
				add(t, nil)
			}
		default:
			add(t, nil)
		}
	}

	return out
}
