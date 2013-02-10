package webby

import (
	"regexp"
	"sort"
	"strings"
	"sync"
)

type vHost struct {
	name string
	boot BootRoute
}

// Use host name as string (e.g example.com)
type VHost struct {
	sync.RWMutex
	hosts []*vHost
}

// Use host name as string (e.g example.com)
type VHostMap map[string]BootRoute

func NewVHost(hosts VHostMap) *VHost {
	v := &VHost{}

	for host, bootroute := range hosts {
		v.hosts = append(v.hosts, &vHost{host, bootroute})
	}

	return v
}

func (v *VHost) getHosts() []*vHost {
	v.RLock()
	defer v.RUnlock()
	hosts := []*vHost{}
	hosts = append(hosts, v.hosts...)
	return hosts
}

func (v *VHost) View(w *Web) {
	for _, host := range v.getHosts() {
		if len(host.name) > len(w.Req.Host) {
			continue
		}
		if strings.ToLower(host.name) == strings.ToLower(w.Req.Host[:len(host.name)]) {
			host.boot.View(w)
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
	sync.RWMutex
	vhost vHostRegs
}

func NewVHostRegExp(hostmap VHostRegExpMap) *VHostRegExp {
	vh := &VHostRegExp{}
	vh.registerMap(hostmap)
	return vh
}

func (vh *VHostRegExp) getHosts() vHostRegs {
	vh.RLock()
	defer vh.RUnlock()
	hosts := vHostRegs{}
	hosts = append(hosts, vh.vhost...)
	return hosts
}

func (vh *VHostRegExp) register(RegExpRule string, bootroute BootRoute) {
	for _, host := range vh.vhost {
		if host.RegExp == RegExpRule {
			host.BootRoute = bootroute
			return
		}
	}

	vh.vhost = append(vh.vhost, &vHostRegExpItem{RegExpRule, regexp.MustCompile(RegExpRule), bootroute})
}

func (vh *VHostRegExp) registerMap(hostmap VHostRegExpMap) {
	if vh.vhost == nil {
		vh.vhost = vHostRegs{}
	}

	for rule, bootroute := range hostmap {
		vh.register(rule, bootroute)
	}

	sort.Sort(vh.vhost)
}

func (vh *VHostRegExp) AddMap(hostmap VHostRegExpMap) {
	vh.registerMap(hostmap)
}

func (vh *VHostRegExp) View(w *Web) {
	for _, host := range vh.getHosts() {
		if !host.RegExpComplied.MatchString(w.Req.Host) {
			continue
		}

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

	w.Error404()
}
