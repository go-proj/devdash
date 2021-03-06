package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Phantas0s/devdash/internal/platform"
	"github.com/pkg/errors"
)

const (
	githubBoxStars          = "github.box_stars"
	githubBoxWatchers       = "github.box_watchers"
	githubBoxOpenIssues     = "github.box_open_issues"
	githubTableRepositories = "github.table_repositories"
	githubTableBranches     = "github.table_branches"
	githubTableIssues       = "github.table_issues"
	githubTablePullRequests = "github.table_pull_requests"
	githubBarViews          = "github.bar_views"
	githubBarCommits        = "github.bar_commits"
	githubBarStars          = "github.bar_stars"
)

type githubWidget struct {
	tui    *Tui
	client *platform.Github
}

// NewGithubWidget with all information necessary to connect to the Github API.
func NewGithubWidget(token string, owner string, repo string) (*githubWidget, error) {
	g, err := platform.NewGithubClient(token, owner, repo)
	if err != nil {
		return nil, err
	}
	return &githubWidget{
		client: g,
	}, nil
}

// CreateWidgets for the Github service.
func (g *githubWidget) CreateWidgets(widget Widget, tui *Tui) (err error) {
	g.tui = tui

	switch widget.Name {
	case githubBoxStars:
		err = g.boxStars(widget)
	case githubBoxWatchers:
		err = g.boxWatchers(widget)
	case githubBoxOpenIssues:
		err = g.boxOpenIssues(widget)
	case githubTableRepositories:
		err = g.tableRepo(widget)
	case githubTableBranches:
		err = g.tableBranches(widget)
	case githubTableIssues:
		err = g.tableIssues(widget)
	case githubTablePullRequests:
		err = g.tablePullRequests(widget)
	case githubBarViews:
		err = g.barViews(widget)
	case githubBarCommits:
		err = g.barCommits(widget)
	case githubBarStars:
		err = g.barStars(widget)
	default:
		return errors.Errorf("can't find the widget %s for service github", widget.Name)
	}

	return
}

func (g *githubWidget) boxStars(widget Widget) error {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := fmt.Sprintf(" Github Stars for %s", repo)
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	stars, err := g.client.TotalStars(repo)
	if err != nil {
		return err
	}

	s := strconv.FormatInt(int64(stars), 10)

	err = g.tui.AddTextBox(
		s,
		title,
		widget.Options,
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *githubWidget) boxWatchers(widget Widget) error {
	title := " Github Watchers "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	w, err := g.client.TotalWatchers(repo)
	if err != nil {
		return err
	}

	s := strconv.FormatInt(int64(w), 10)

	err = g.tui.AddTextBox(
		s,
		title,
		widget.Options,
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *githubWidget) boxOpenIssues(widget Widget) error {
	title := " Github Open Issues "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	w, err := g.client.TotalOpenIssues(repo)
	if err != nil {
		return err
	}

	s := strconv.FormatInt(int64(w), 10)

	err = g.tui.AddTextBox(
		s,
		title,
		widget.Options,
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *githubWidget) tableRepo(widget Widget) (err error) {
	title := " Github Repositories "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	metrics := []string{"name", "stars", "watchers", "forks", "open_issues"}
	if _, ok := widget.Options[optionMetrics]; ok {
		if len(widget.Options[optionMetrics]) > 0 {
			metrics = strings.Split(strings.TrimSpace(widget.Options[optionMetrics]), ",")
		}
	}

	order := "pushed"
	if _, ok := widget.Options[optionOrder]; ok {
		order = widget.Options[optionRowLimit]
	}

	rs, err := g.client.ListRepo(int(limit), order, metrics)
	if err != nil {
		return err
	}

	g.tui.AddTable(rs, title, widget.Options)

	return nil
}

func (g *githubWidget) tableBranches(widget Widget) (err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Branches "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	bs, err := g.client.ListBranches(repo, int(limit))
	if err != nil {
		return err
	}

	g.tui.AddTable(bs, title, widget.Options)

	return nil
}

func (g *githubWidget) tableIssues(widget Widget) (err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Issues "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	is, err := g.client.ListIssues(repo, int(limit))
	if err != nil {
		return err
	}

	g.tui.AddTable(is, title, widget.Options)

	return nil
}

func (g *githubWidget) tablePullRequests(widget Widget) (err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Pull Requests "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	var limit int64 = 5
	if _, ok := widget.Options[optionRowLimit]; ok {
		limit, err = strconv.ParseInt(widget.Options[optionRowLimit], 10, 0)
		if err != nil {
			return errors.Wrapf(err, "%s must be a number", widget.Options[optionRowLimit])
		}
	}

	is, err := g.client.ListPullRequests(repo, int(limit))
	if err != nil {
		return err
	}

	g.tui.AddTable(is, title, widget.Options)

	return nil
}

func (g *githubWidget) barViews(widget Widget) (err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Views "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	dim, counts, err := g.client.Views(repo, 0)
	if err != nil {
		return err
	}

	g.tui.AddBarChart(counts, dim, title, widget.Options)

	return nil
}

// TODO to refactor - transforming any date statement (weeks_ago, month_ago) into days weeks_ago in platform.date, and plugt it in.
func (g *githubWidget) barCommits(widget Widget) (err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Commit Per Week "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	sd := "7_weeks_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		sd = widget.Options[optionStartDate]
	}

	ed := "0_weeks_ago"
	if _, ok := widget.Options[optionEndDate]; ok {
		ed = widget.Options[optionEndDate]
	}

	scope := ownerScope
	if _, ok := widget.Options[optionScope]; ok {
		scope = widget.Options[optionScope]
	}

	if !strings.Contains(sd, "weeks_ago") || !strings.Contains(ed, "weeks_ago") {
		return errors.New("The widget github.bar_commits require you to indicate a week range, ie startDate: 5_weeks_ago, endDate: 1_weeks_ago ")
	}

	sw, err := platform.ExtractCountPeriod(sd)
	if err != nil {
		return err
	}

	ew, err := platform.ExtractCountPeriod(ed)
	if err != nil {
		return err
	}

	dim, counts, err := g.client.CountCommits(repo, scope, sw, ew, time.Now())
	if err != nil {
		return err
	}

	g.tui.AddBarChart(counts, dim, title, widget.Options)

	return nil
}

func (g *githubWidget) barStars(widget Widget) (err error) {
	var repo string
	if _, ok := widget.Options[optionRepository]; ok {
		repo = widget.Options[optionRepository]
	}

	title := " Github Stars "
	if _, ok := widget.Options[optionTitle]; ok {
		title = widget.Options[optionTitle]
	}

	startDate := "7_days_ago"
	if _, ok := widget.Options[optionStartDate]; ok {
		startDate = widget.Options[optionStartDate]
	}

	endDate := "today"
	if _, ok := widget.Options[optionEndDate]; ok {
		endDate = widget.Options[optionEndDate]
	}

	sd, ed, err := platform.ConvertDates(time.Now(), startDate, endDate)
	if err != nil {
		return err
	}

	dim, counts, err := g.client.CountStars(repo, sd, ed)
	if err != nil {
		return err
	}

	g.tui.AddBarChart(counts, dim, title, widget.Options)

	return nil
}
