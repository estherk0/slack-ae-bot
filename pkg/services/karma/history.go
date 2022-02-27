package karma

import (
	"bytes"
	"context"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
)

const (
	defaultDaysOfHistory = 7
	historyTemplateStr   = `:eyes: You received total {{ .Count }} points for {{ .Days }} days.
	{{- range $userID, $point := .Summary }}
			{{ $point }} points from <@{{ $userID }}>
	{{- end }}
	`
)

var historyTmpl = template.Must(template.New("history template").
	Funcs(templateFuncs).Parse(historyTemplateStr))

func (s *service) GetHistories(event *slackevents.AppMentionEvent) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	season, err := s.karmaRepository.GetCurrentSeason(ctx)
	if err != nil {
		logrus.Errorln("[GetHistories] failed to get season error: ", err.Error())
		return nil
	}

	logs, err := s.karmaRepository.SearchLogs(ctx, season.SeasonID, event.User, defaultDaysOfHistory)
	if err != nil {
		logrus.Errorln("[GetHistories] error: ", err.Error())
		return nil
	}

	/* calculate karma points by giverID */
	summary := map[string]int{}
	for _, log := range logs {
		summary[log.GiverID]++
	}

	var buffer bytes.Buffer
	params := struct {
		Count   int
		Days    int
		Summary map[string]int
	}{
		len(logs),
		defaultDaysOfHistory,
		summary,
	}
	historyTmpl.Execute(&buffer, params)
	s.slackapiService.PostMessage(event.Channel, buffer.String())
	return nil
}
