package webby

import (
	"regexp"
	"sort"
	"strings"
)

// Use host name as string (e.g example.com)
type VHost map[string]BootRoute

func (v VHost) View(w *Web) {
	for host, bootroute := range v {
		if len(host) > len(w.Req.Host) {
			continue
		}
		if strings.ToLower(host) == strings.ToLower(w.Req.Host[:len(host)]) {
			bootroute.View(w)
			return
		}
	}

	w.Error404()
}

// Use host name regexp as string (e.g. (?P<subdomain>[a-z0-9-_]+)\.example\.com)
type VHostRegExpMap map[string]BootRoute

type vHostRegExpItem struct {
	RegExp         string
	RegExpComplied *regexp.Regexp
	BootRoute      BootRoute
}

type vHostRegs []*vHostRegExpItem

func (vh vHostRegs) Len() int {
	return len(vh)
}

func (vh vHostRegs) Less(i, j int) bool {
	return vh[i].RegExp < vh[j].RegExp
}

func (vh vHostRegs) Swap(i, j int) {
	vh[i], vh[j] = vh[j], vh[i]
}

type VHostRegExp struct {
	vhost vHostRegs
}

func NewVHostMap(hostmap VHostRegExpMap) *VHostRegExp {
	vh := &VHostRegExp{}
	vh.registerMap(hostmap)
	return vh
}

func (vh *VHostRegExp) registerMap(hostmap VHostRegExpMap) {
	if vh.vhost == nil {
		vh.vhost = vHostRegs{}
	}

	for rule, bootroute := range hostmap {
		vh.vhost = append(vh.vhost, &vHostRegExpItem{rule, regexp.MustCompile(rule), bootroute})
	}

	sort.Sort(vh.vhost)
}

func (vh *VHostRegExp) AddMap(hostmap VHostRegExpMap) {
	vh.registerMap(hostmap)
}

func (vh *VHostRegExp) View(w *Web) {
	for _, host := range vh.vhost {
		if host.RegExpComplied.MatchString(w.Req.Host) {
			names := host.RegExpComplied.SubexpNames()
			matches := host.RegExpComplied.FindStringSubmatch(w.Req.Host)

			for key, name := range names {
				if name != "" {
					w.Param.Add(name, matches[key])
				}
			}

			host.BootRoute.View(w)

			return
		}
	}

	w.Error404()
}
