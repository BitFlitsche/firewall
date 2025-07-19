package utils

import (
	"fmt"
	"net"
	"strings"
)

// CIDRInfo contains information about a CIDR block
type CIDRInfo struct {
	Network  string
	Mask     int
	StartIP  net.IP
	EndIP    net.IP
	TotalIPs uint64
	IsValid  bool
}

// ParseCIDR parses a CIDR string and returns CIDRInfo
func ParseCIDR(cidr string) (*CIDRInfo, error) {
	// Handle single IP addresses (not CIDR)
	if !strings.Contains(cidr, "/") {
		ip := net.ParseIP(cidr)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", cidr)
		}
		return &CIDRInfo{
			Network:  cidr + "/32", // Single IP as /32
			Mask:     32,
			StartIP:  ip,
			EndIP:    ip,
			TotalIPs: 1,
			IsValid:  true,
		}, nil
	}

	// Parse CIDR notation
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %s", err)
	}

	// Calculate network information
	mask, _ := ipNet.Mask.Size()
	startIP := ipNet.IP
	endIP := make(net.IP, len(startIP))
	copy(endIP, startIP)

	// Calculate end IP
	for i := range endIP {
		endIP[i] = startIP[i] | ^ipNet.Mask[i]
	}

	// Calculate total IPs in range
	totalIPs := uint64(1) << (32 - mask)

	return &CIDRInfo{
		Network:  cidr,
		Mask:     mask,
		StartIP:  startIP,
		EndIP:    endIP,
		TotalIPs: totalIPs,
		IsValid:  true,
	}, nil
}

// IsIPInCIDR checks if an IP address falls within a CIDR block
func IsIPInCIDR(ipStr, cidrStr string) (bool, error) {
	// Parse the IP to check
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	// Parse the CIDR block
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR notation: %s", err)
	}

	// Check if IP is in the network
	return ipNet.Contains(ip), nil
}

// ValidateCIDR checks if a string is valid CIDR notation
func ValidateCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// IsSingleIP checks if a string represents a single IP (not CIDR)
func IsSingleIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && !strings.Contains(ipStr, "/")
}

// IsCIDRNotation checks if a string is CIDR notation
func IsCIDRNotation(ipStr string) bool {
	return strings.Contains(ipStr, "/") && ValidateCIDR(ipStr)
}

// GetCIDRRange returns the start and end IPs of a CIDR block
func GetCIDRRange(cidr string) (net.IP, net.IP, error) {
	info, err := ParseCIDR(cidr)
	if err != nil {
		return nil, nil, err
	}
	return info.StartIP, info.EndIP, nil
}

// FormatCIDRInfo returns a human-readable description of a CIDR block
func FormatCIDRInfo(cidr string) (string, error) {
	info, err := ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	if info.TotalIPs == 1 {
		return fmt.Sprintf("Single IP: %s", info.StartIP.String()), nil
	}

	return fmt.Sprintf("Range: %s to %s (%d IPs)",
		info.StartIP.String(),
		info.EndIP.String(),
		info.TotalIPs), nil
}

// IPToUint32 converts an IP address to uint32 for comparison
func IPToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

// Uint32ToIP converts uint32 back to IP address
func Uint32ToIP(ipInt uint32) net.IP {
	return net.IPv4(byte(ipInt>>24), byte(ipInt>>16), byte(ipInt>>8), byte(ipInt))
}
