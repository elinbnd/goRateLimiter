package rateLimiter

import (
	"fmt"
	"net"
	"sync"
	"strings"
	"time"
)
type client struct {
	rps int
	rpm int
	rph int
	rpd int
	lastSec  int64
	lastMin  int64
	lastHour int64
	lastDay  int64
	updLimit int
	downloadLimit int
	trafLimit int
	
}
type RateLimiter struct {
	mu sync.Mutex
	clients map[string]*client
	subClients map[string]*client
	rps int
	rpm int
	rph int
	rpd int
	uploadClient int 
	downloadClient int
	totalClient int
	maxConnections int 
}
func subNetworking(ip string) string {
	parsed := net.ParseIP(ip).To4()
	if parsed == nil {
		return ip
	}
	return fmt.Sprintf("%d.%d.%d.0", parsed[0], parsed[1], parsed[2])
}
func NewRateLimiter(rps, rpm, rph, rpd, updLimit, downloadLimit, trafLimit, maxConnections int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*client),
		subClients: make(map[string]*client),
		rps:     rps,
		rpm:     rpm,
		rph:     rph,
		rpd:     rpd,
		uploadClient: updLimit,
		downloadClient: downloadLimit,
		totalClient: trafLimit,
		maxConnections: maxConnections,
		}
	}
func (rateLimiter *RateLimiter) CheckUpload(ip string, size int64) bool {
	rateLimiter.mu.Lock()
	defer rateLimiter.mu.Unlock()
	c, ok := rateLimiter.clients[ip]
	if !ok {
		c = &client{
		}
		rateLimiter.clients[ip] = c
	}
	c.updLimit = c.updLimit + int(size) 
	c.trafLimit = c.trafLimit + int(size)
	if c.updLimit > rateLimiter.uploadClient {
		return false
	}
	if c.trafLimit > rateLimiter.totalClient {
		return false
	}
	return true
}
func (rateLimiter *RateLimiter) CheckDownload(ip string, size int64) bool {
	rateLimiter.mu.Lock()
	defer rateLimiter.mu.Unlock()
	c, ok := rateLimiter.clients[ip]
	if !ok {
		c = &client{}
		rateLimiter.clients[ip] = c
	}
	c.downloadLimit = c.downloadLimit + int(size)
	c.trafLimit = c.trafLimit + int(size)
	if c.downloadLimit > rateLimiter.downloadClient {
		return false
	}
	if c.trafLimit > rateLimiter.totalClient {
		return false
	}
	return true

}
func (r *RateLimiter) Adobe(ip string) bool{
	return r.Trust(ip)
}

func (r *RateLimiter) Trust(ip string) bool {

	r.mu.Lock()
	defer r.mu.Unlock()
	subNetwork := subNetworking(ip)
	now := time.Now()
	c, ok := r.clients[ip]
	if !ok {
		c = &client{}
		r.clients[ip] = c
	}
	subClient, ok := r.subClients[subNetwork]
	if !ok {
		subClient = &client{}
		r.subClients[subNetwork] = subClient
	}
	sec := now.Unix()
	min := sec / 60
	hour := sec / 3600
	day := sec / 86400
	if c.lastSec != sec {
		c.rps = 0
		c.lastSec = sec
	}
	if c.lastMin != min {
		c.rpm = 0
		c.lastMin = min
	}
	if c.lastHour != hour {
		c.rph = 0
		c.lastHour = hour
	}

	if c.lastDay != day {
		c.rpd = 0
		c.lastDay = day
		c.updLimit = 0
		c.downloadLimit = 0
		c.trafLimit = 0
	}
	if subClient.lastSec != sec {
		subClient.rps = 0 
		subClient.lastSec = sec
	}
	if subClient.lastMin != min {
		subClient.rpm = 0 
		subClient.lastMin = min
	}
	if subClient.lastHour != hour {
		subClient.rph = 0
		subClient.lastHour = hour 
	}
	if subClient.lastDay != day {
		subClient.rpd = 0 
		subClient.lastDay = day
	}
	if c.rps >= r.rps ||
		c.rpm >= r.rpm ||
		c.rph >= r.rph ||
		c.rpd >= r.rpd {
		return false
	}
	if subClient.rps >= r.rps * 5 || subClient.rpm >= r.rpm * 5 || subClient.rph >= r.rph * 5 || subClient.rpd >= r.rpd * 5 {
		return false
	}
	subClient.rps = subClient.rps + 1
	subClient.rpm = subClient.rpm + 1
	subClient.rph = subClient.rph + 1
	subClient.rpd = subClient.rpd + 1
	c.rps = c.rps + 1
	c.rpm = c.rpm + 1
	c.rph = c.rph + 1
	c.rpd = c.rpd + 1

	return true
}
func (rateLimiter *RateLimiter) CheckMaxConnections(currConns int) bool {
	return currConns < rateLimiter.maxConnections
}
func (r *RateLimiter) SetLimits(rps, rpm, rph, rpd int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rps = rps
	r.rpm = rpm
	r.rph = rph
	r.rpd = rpd
}