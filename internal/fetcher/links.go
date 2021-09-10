package fetcher

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
	"github.com/hi20160616/ms-tibetpost/configs"
	"github.com/pkg/errors"
)

func fetchLinks() ([]string, error) {
	rt := []string{}

	for _, rawurl := range configs.Data.MS.URL {
		links, err := getLinks(rawurl)
		if err != nil {
			return nil, err
		}
		rt = append(rt, links...)
	}
	home := linksFilter(rt, `.*?/\d*-.*`)
	news := linksFilter(rt, `.*?/news/.*?/\d*-.*`)
	features := linksFilter(rt, `.*?/features/.*?/\d*-.*`)
	op_ed := linksFilter(rt, `.*?/op-ed/(.*?/)?\d*-.*`)
	ecosystem := linksFilter(rt, `.*?/ecosystem/(.*?/)?\d*-.*`)
	influence_impact := linksFilter(rt, `.*?/influence-impact/(.*?/)?\d*-.*`)
	tibet_facts := linksFilter(rt, `.*?/tibet-facts/(.*?/)?\d*-.*`)
	rt = []string{}
	rt = append(rt, home...)
	rt = append(rt, news...)
	rt = append(rt, features...)
	rt = append(rt, op_ed...)
	rt = append(rt, ecosystem...)
	rt = append(rt, influence_impact...)
	rt = append(rt, tibet_facts...)
	rt = gears.StrSliceDeDupl(rt)
	blackList := []string{
		`.*?/92-health-safety`,
		`.*?/84-letters-to-editor`,
		`.*?/93-glaciers-ice-caps`,
		`.*?/94-ocean-climate`, // TODO: this is a category
		`.*?/95-green-clean`,
		`.*?/97-buddhist-science`,
		`.*?/98-religious`,
		`.*?/99-spiritual-life`,
		`.*?/100-mind-brain`,
		`.*?/101-medicines-and-healthcare`,
	}
	for _, e := range blackList {
		rt = kickOut(rt, e)
	}
	return rt, nil
}

// getLinksJson get links from a url that return json data.
func getLinksJson(rawurl string) ([]string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	raw, _, err := exhtml.GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`"url":\s"(.*?)",`)
	rs := re.FindAllStringSubmatch(string(raw), -1)
	rt := []string{}
	for _, item := range rs {
		rt = append(rt, "https://"+u.Hostname()+item[1])
	}
	return gears.StrSliceDeDupl(rt), nil
}

func getLinks(rawurl string) ([]string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if links, err := exhtml.ExtractLinks(u.String()); err != nil {
		return nil, errors.WithMessagef(err, "[%s] cannot extract links from %s",
			configs.Data.MS.Title, rawurl)
	} else {
		return gears.StrSliceDeDupl(links), nil
	}
}

// kickOutLinksMatchPath will kick out the links match the path,
func kickOutLinksMatchPath(links []string, path string) []string {
	tmp := []string{}
	// path = "/" + url.QueryEscape(path) + "/"
	// path = url.QueryEscape(path)
	for _, link := range links {
		if !strings.Contains(link, path) {
			tmp = append(tmp, link)
		}
	}
	return tmp
}

// TODO: use point to impletement linksFilter
// linksFilter return all matched string list
func linksFilter(links []string, regex string) []string {
	flinks := []string{}
	re := regexp.MustCompile(regex)
	s := strings.Join(links, "\n")
	flinks = re.FindAllString(s, -1)
	return flinks
}

func kickOut(links []string, regex string) []string {
	blackList := linksFilter(links, regex)
	return gears.StrSliceDiff(links, blackList)
}
