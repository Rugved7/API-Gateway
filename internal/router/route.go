// Package router is responsible for the routing engine config
package router

type Route struct {
	Prefix   string
	Upstream string
}
