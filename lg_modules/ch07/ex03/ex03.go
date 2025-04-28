package main

import (
	"io"
	"os"
	"sort"
	"strings"
)

type Team struct {
	name        string
	playerNames []string
}

type League struct {
	Teams map[string]Team
	Wins  map[string]int
}

func (l *League) MatchResult(team1 string, score1 int, team2 string, score2 int) {
	if _, ok := l.Teams[team1]; !ok {
		return
	}
	if _, ok := l.Teams[team2]; !ok {
		return
	}
	switch {
	case score1 > score2:
		l.Wins[team1]++

	case score1 < score2:
		l.Wins[team2]++

	default:
		return
	}
}

func (l *League) Ranking() []string {
	teamNames := make([]string, 0, len(l.Teams))
	for teamName := range l.Teams {
		teamNames = append(teamNames, teamName)
	}
	sort.SliceStable(teamNames, func(i, j int) bool {
		return l.Wins[teamNames[i]] > l.Wins[teamNames[j]]
	})
	return teamNames
}

type Ranker interface {
	Ranking() []string
}

func RankPrinter(r Ranker, w io.Writer) error {
	result := strings.Join(r.Ranking(), "\n") + "\n"
	_, err := io.WriteString(w, result)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	l := League{
		Teams: map[string]Team{
			"GER": {
				name:        "GER",
				playerNames: []string{"Frank", "Bernhard", "Raimund"},
			},
			"FRA": {
				name:        "FRA",
				playerNames: []string{"François", "Bernard", "Raymond"},
			},
			"IT": {
				name:        "IT",
				playerNames: []string{"Franco", "Bernardo", "Raimondo"},
			},
		},
		Wins: map[string]int{},
	}
	l.MatchResult("GER", 5, "FRA", 3)
	l.MatchResult("GER", 1, "IT", 1)
	l.MatchResult("FR", 3, "IT", 4)
	l.MatchResult("IT", 3, "GER", 5)
	RankPrinter(&l, os.Stdout)
}
