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

// ConflictInfo contains information about IP/CIDR conflicts
type ConflictInfo struct {
	Type        string   `json:"type"`        // "ip_in_cidr", "cidr_overlaps", "exact_match"
	Message     string   `json:"message"`     // Human-readable message
	Conflicting []string `json:"conflicting"` // List of conflicting entries
	Severity    string   `json:"severity"`    // "warning", "error", "info"
	Status      string   `json:"status"`      // Status of conflicting entry
}

// CheckIPConflicts checks if an IP address conflicts with existing CIDR ranges
func CheckIPConflicts(ip string, existingCIDRs []string, existingStatuses map[string]string, newStatus string) ([]ConflictInfo, error) {
	var conflicts []ConflictInfo

	// Parse the IP to check
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	for _, cidr := range existingCIDRs {
		// Check if IP is within this CIDR range
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue // Skip invalid CIDR
		}

		if ipNet.Contains(ipAddr) {
			cidrStatus := existingStatuses[cidr]
			severity := "warning"

			// If same status, prevent creation (redundant rule)
			if cidrStatus == newStatus {
				severity = "error"
			}

			conflicts = append(conflicts, ConflictInfo{
				Type:        "ip_in_cidr",
				Message:     fmt.Sprintf("IP %s is already covered by CIDR range %s (status: %s)", ip, cidr, cidrStatus),
				Conflicting: []string{cidr},
				Severity:    severity,
				Status:      cidrStatus,
			})
		}
	}

	return conflicts, nil
}

// CheckCIDRConflicts checks if a CIDR range conflicts with existing IPs or CIDR ranges
func CheckCIDRConflicts(newCIDR string, existingIPs []string, existingCIDRs []string, existingStatuses map[string]string, newStatus string) ([]ConflictInfo, error) {
	var conflicts []ConflictInfo

	// Parse the new CIDR
	_, err := ParseCIDR(newCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %s", err)
	}

	// Check conflicts with existing individual IPs
	for _, ip := range existingIPs {
		ipAddr := net.ParseIP(ip)
		if ipAddr == nil {
			continue // Skip invalid IPs
		}

		// Check if this IP is within the new CIDR range
		_, ipNet, err := net.ParseCIDR(newCIDR)
		if err != nil {
			continue
		}

		if ipNet.Contains(ipAddr) {
			ipStatus := existingStatuses[ip]
			severity := "warning"

			// If same status, prevent creation (redundant rule)
			if ipStatus == newStatus {
				severity = "error"
			}

			conflicts = append(conflicts, ConflictInfo{
				Type:        "cidr_covers_ip",
				Message:     fmt.Sprintf("CIDR range %s would cover existing IP %s (status: %s)", newCIDR, ip, ipStatus),
				Conflicting: []string{ip},
				Severity:    severity,
				Status:      ipStatus,
			})
		}
	}

	// Check conflicts with existing CIDR ranges
	for _, existingCIDR := range existingCIDRs {
		if existingCIDR == newCIDR {
			// Exact match
			cidrStatus := existingStatuses[existingCIDR]
			conflicts = append(conflicts, ConflictInfo{
				Type:        "exact_match",
				Message:     fmt.Sprintf("CIDR range %s already exists (status: %s)", newCIDR, cidrStatus),
				Conflicting: []string{existingCIDR},
				Severity:    "error",
				Status:      cidrStatus,
			})
			continue
		}

		// Check for overlap
		overlap, err := CheckCIDROverlap(newCIDR, existingCIDR)
		if err != nil {
			continue // Skip invalid CIDR
		}

		if overlap {
			cidrStatus := existingStatuses[existingCIDR]
			severity := "error"

			// If same status, prevent creation (redundant rule)
			if cidrStatus == newStatus {
				severity = "error"
			} else {
				severity = "warning" // Different status, allow with warning
			}

			conflicts = append(conflicts, ConflictInfo{
				Type:        "cidr_overlaps",
				Message:     fmt.Sprintf("CIDR range %s overlaps with existing range %s (status: %s)", newCIDR, existingCIDR, cidrStatus),
				Conflicting: []string{existingCIDR},
				Severity:    severity,
				Status:      cidrStatus,
			})
		}
	}

	return conflicts, nil
}

// CheckCIDROverlap checks if two CIDR ranges overlap
func CheckCIDROverlap(cidr1, cidr2 string) (bool, error) {
	// Parse both CIDR ranges
	_, ipNet1, err := net.ParseCIDR(cidr1)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR 1: %s", err)
	}

	_, ipNet2, err := net.ParseCIDR(cidr2)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR 2: %s", err)
	}

	// Check if either network contains the other's start IP
	return ipNet1.Contains(ipNet2.IP) || ipNet2.Contains(ipNet1.IP), nil
}

// GetConflictingEntries returns all entries that would conflict with a new IP or CIDR
func GetConflictingEntries(newEntry string, existingEntries []string) ([]string, error) {
	var conflicts []string

	// Determine if new entry is IP or CIDR
	if IsSingleIP(newEntry) {
		// Check if IP conflicts with existing CIDR ranges
		for _, existing := range existingEntries {
			if IsCIDRNotation(existing) {
				_, ipNet, err := net.ParseCIDR(existing)
				if err != nil {
					continue
				}

				ip := net.ParseIP(newEntry)
				if ip != nil && ipNet.Contains(ip) {
					conflicts = append(conflicts, existing)
				}
			}
		}
	} else if IsCIDRNotation(newEntry) {
		// Check if CIDR conflicts with existing entries
		for _, existing := range existingEntries {
			if existing == newEntry {
				// Exact match
				conflicts = append(conflicts, existing)
				continue
			}

			if IsCIDRNotation(existing) {
				// Check CIDR overlap
				overlap, err := CheckCIDROverlap(newEntry, existing)
				if err == nil && overlap {
					conflicts = append(conflicts, existing)
				}
			} else {
				// Check if existing IP is within new CIDR
				_, ipNet, err := net.ParseCIDR(newEntry)
				if err != nil {
					continue
				}

				ip := net.ParseIP(existing)
				if ip != nil && ipNet.Contains(ip) {
					conflicts = append(conflicts, existing)
				}
			}
		}
	}

	return conflicts, nil
}
