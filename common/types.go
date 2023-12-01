package common

import (
	"net"
	"time"
)

type Creation struct {
	Ctime *time.Time `json:"ctime"`
}

type Exclude struct {
	Exclude []Glob `json:"exclude"`
}

type Glob string

func (glob Glob) String() string {
	return string(glob)
}

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

type User struct {
	UID  *uint32 `json:"uid"`
	User string  `json:"user"`
}
