package utils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIP extracts the real client IP address from the request
// It checks various headers in order of preference to handle proxies and load balancers
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (most common for proxies)
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if isValidIP(ip) {
				return ip
			}
		}
	}

	// Check X-Real-IP header (used by some proxies)
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" && isValidIP(xRealIP) {
		return xRealIP
	}

	// Check X-Client-IP header
	xClientIP := c.GetHeader("X-Client-IP")
	if xClientIP != "" && isValidIP(xClientIP) {
		return xClientIP
	}

	// Check CF-Connecting-IP header (Cloudflare)
	cfConnectingIP := c.GetHeader("CF-Connecting-IP")
	if cfConnectingIP != "" && isValidIP(cfConnectingIP) {
		return cfConnectingIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}

// isValidIP checks if the given string is a valid IP address
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsPrivateIP checks if the IP address is in a private range
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check for private IP ranges
	privateRanges := []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link-local
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(parsedIP) {
			return true
		}
	}

	return false
}