package common

import (
	"net"
	"time"
)

type Creation struct {
	Ctime *time.Time `json:"ctime"`
}

type Env map[string]string

type Exclude struct {
	Exclude []Glob `json:"exclude"`
}

type Glob string

type Group struct {
	GID   *uint32 `json:"gid"`
	Group string  `json:"group"`
}

type Host struct {
	IP   net.IP `json:"ip"`
	Host string `json:"host"`
}

type Include struct {
	Include []Glob `json:"include"`
}

type Labels map[string]string

type Mode struct {
	Mode uint32 `json:"value"`
}

type Platform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Variant      string `json:"variant"`
}

type User struct {
	UID  *uint32 `json:"uid"`
	User string  `json:"user"`
}
