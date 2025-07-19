package controllers

import (
	"context"
	"errors"
	"firewall/services"
	"firewall/validation"
	"net/http"
	"strings"
	"time"

	"firewall/models"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FilterRequest defines the structure for the incoming JSON request
type FilterRequest struct {
	IP        string `json:"ip" binding:"omitempty,max=45"`
	Email     string `json:"email" binding:"omitempty,email,max=254"`
	UserAgent string `json:"user_agent" binding:"omitempty,max=500"`
	Country   string `json:"country" binding:"omitempty,len=2,alpha"`
	Content   string `json:"content" binding:"omitempty,max=10000"` // optional
	Username  string `json:"username" binding:"omitempty,max=100"`  // optional
}

// normalizeEmail removes dots from the local part of Gmail addresses
// This handles Gmail's behavior where dots are ignored in the local part
// Examples: test@gmail.com = t.e.s.t@gmail.com = t.e.s.t@gmail.de
func normalizeEmail(email string) string {
	if email == "" {
		return email
	}

	// Convert to lowercase
	email = strings.ToLower(email)

	// Split email into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // Invalid email format, return as-is
	}

	localPart := parts[0]
	domain := parts[1]

	// Check if it's a Gmail domain (gmail.com, gmail.de, gmail.co.uk, etc.)
	if strings.HasSuffix(domain, "gmail.com") ||
		strings.HasSuffix(domain, "gmail.de") ||
		strings.HasSuffix(domain, "gmail.co.uk") ||
		strings.HasSuffix(domain, "gmail.fr") ||
		strings.HasSuffix(domain, "gmail.it") ||
		strings.HasSuffix(domain, "gmail.es") ||
		strings.HasSuffix(domain, "gmail.nl") ||
		strings.HasSuffix(domain, "gmail.se") ||
		strings.HasSuffix(domain, "gmail.no") ||
		strings.HasSuffix(domain, "gmail.dk") ||
		strings.HasSuffix(domain, "gmail.fi") ||
		strings.HasSuffix(domain, "gmail.pl") ||
		strings.HasSuffix(domain, "gmail.cz") ||
		strings.HasSuffix(domain, "gmail.hu") ||
		strings.HasSuffix(domain, "gmail.ro") ||
		strings.HasSuffix(domain, "gmail.bg") ||
		strings.HasSuffix(domain, "gmail.hr") ||
		strings.HasSuffix(domain, "gmail.si") ||
		strings.HasSuffix(domain, "gmail.sk") ||
		strings.HasSuffix(domain, "gmail.lt") ||
		strings.HasSuffix(domain, "gmail.lv") ||
		strings.HasSuffix(domain, "gmail.ee") ||
		strings.HasSuffix(domain, "gmail.pt") ||
		strings.HasSuffix(domain, "gmail.gr") ||
		strings.HasSuffix(domain, "gmail.at") ||
		strings.HasSuffix(domain, "gmail.ch") ||
		strings.HasSuffix(domain, "gmail.be") ||
		strings.HasSuffix(domain, "gmail.lu") ||
		strings.HasSuffix(domain, "gmail.ie") ||
		strings.HasSuffix(domain, "gmail.mt") ||
		strings.HasSuffix(domain, "gmail.cy") ||
		strings.HasSuffix(domain, "gmail.is") ||
		strings.HasSuffix(domain, "gmail.li") ||
		strings.HasSuffix(domain, "gmail.mc") ||
		strings.HasSuffix(domain, "gmail.ad") ||
		strings.HasSuffix(domain, "gmail.va") ||
		strings.HasSuffix(domain, "gmail.sm") ||
		strings.HasSuffix(domain, "gmail.by") ||
		strings.HasSuffix(domain, "gmail.md") ||
		strings.HasSuffix(domain, "gmail.ua") ||
		strings.HasSuffix(domain, "gmail.ge") ||
		strings.HasSuffix(domain, "gmail.am") ||
		strings.HasSuffix(domain, "gmail.az") ||
		strings.HasSuffix(domain, "gmail.kz") ||
		strings.HasSuffix(domain, "gmail.kg") ||
		strings.HasSuffix(domain, "gmail.tj") ||
		strings.HasSuffix(domain, "gmail.tm") ||
		strings.HasSuffix(domain, "gmail.uz") ||
		strings.HasSuffix(domain, "gmail.mn") ||
		strings.HasSuffix(domain, "gmail.kr") ||
		strings.HasSuffix(domain, "gmail.jp") ||
		strings.HasSuffix(domain, "gmail.cn") ||
		strings.HasSuffix(domain, "gmail.hk") ||
		strings.HasSuffix(domain, "gmail.tw") ||
		strings.HasSuffix(domain, "gmail.sg") ||
		strings.HasSuffix(domain, "gmail.my") ||
		strings.HasSuffix(domain, "gmail.th") ||
		strings.HasSuffix(domain, "gmail.vn") ||
		strings.HasSuffix(domain, "gmail.ph") ||
		strings.HasSuffix(domain, "gmail.id") ||
		strings.HasSuffix(domain, "gmail.in") ||
		strings.HasSuffix(domain, "gmail.pk") ||
		strings.HasSuffix(domain, "gmail.bd") ||
		strings.HasSuffix(domain, "gmail.lk") ||
		strings.HasSuffix(domain, "gmail.np") ||
		strings.HasSuffix(domain, "gmail.mm") ||
		strings.HasSuffix(domain, "gmail.kh") ||
		strings.HasSuffix(domain, "gmail.la") ||
		strings.HasSuffix(domain, "gmail.br") ||
		strings.HasSuffix(domain, "gmail.ar") ||
		strings.HasSuffix(domain, "gmail.cl") ||
		strings.HasSuffix(domain, "gmail.co") ||
		strings.HasSuffix(domain, "gmail.pe") ||
		strings.HasSuffix(domain, "gmail.ve") ||
		strings.HasSuffix(domain, "gmail.ec") ||
		strings.HasSuffix(domain, "gmail.bo") ||
		strings.HasSuffix(domain, "gmail.py") ||
		strings.HasSuffix(domain, "gmail.uy") ||
		strings.HasSuffix(domain, "gmail.gy") ||
		strings.HasSuffix(domain, "gmail.sr") ||
		strings.HasSuffix(domain, "gmail.gf") ||
		strings.HasSuffix(domain, "gmail.mx") ||
		strings.HasSuffix(domain, "gmail.ca") ||
		strings.HasSuffix(domain, "gmail.us") ||
		strings.HasSuffix(domain, "gmail.au") ||
		strings.HasSuffix(domain, "gmail.nz") ||
		strings.HasSuffix(domain, "gmail.fj") ||
		strings.HasSuffix(domain, "gmail.pg") ||
		strings.HasSuffix(domain, "gmail.sb") ||
		strings.HasSuffix(domain, "gmail.vu") ||
		strings.HasSuffix(domain, "gmail.nc") ||
		strings.HasSuffix(domain, "gmail.pf") ||
		strings.HasSuffix(domain, "gmail.ws") ||
		strings.HasSuffix(domain, "gmail.to") ||
		strings.HasSuffix(domain, "gmail.ck") ||
		strings.HasSuffix(domain, "gmail.nu") ||
		strings.HasSuffix(domain, "gmail.tk") ||
		strings.HasSuffix(domain, "gmail.wf") ||
		strings.HasSuffix(domain, "gmail.as") ||
		strings.HasSuffix(domain, "gmail.gu") ||
		strings.HasSuffix(domain, "gmail.mp") ||
		strings.HasSuffix(domain, "gmail.pr") ||
		strings.HasSuffix(domain, "gmail.vi") ||
		strings.HasSuffix(domain, "gmail.um") ||
		strings.HasSuffix(domain, "gmail.af") ||
		strings.HasSuffix(domain, "gmail.ir") ||
		strings.HasSuffix(domain, "gmail.iq") ||
		strings.HasSuffix(domain, "gmail.sa") ||
		strings.HasSuffix(domain, "gmail.ae") ||
		strings.HasSuffix(domain, "gmail.om") ||
		strings.HasSuffix(domain, "gmail.qa") ||
		strings.HasSuffix(domain, "gmail.bh") ||
		strings.HasSuffix(domain, "gmail.kw") ||
		strings.HasSuffix(domain, "gmail.ye") ||
		strings.HasSuffix(domain, "gmail.jo") ||
		strings.HasSuffix(domain, "gmail.lb") ||
		strings.HasSuffix(domain, "gmail.sy") ||
		strings.HasSuffix(domain, "gmail.il") ||
		strings.HasSuffix(domain, "gmail.ps") ||
		strings.HasSuffix(domain, "gmail.eg") ||
		strings.HasSuffix(domain, "gmail.ly") ||
		strings.HasSuffix(domain, "gmail.tn") ||
		strings.HasSuffix(domain, "gmail.dz") ||
		strings.HasSuffix(domain, "gmail.ma") ||
		strings.HasSuffix(domain, "gmail.mr") ||
		strings.HasSuffix(domain, "gmail.sn") ||
		strings.HasSuffix(domain, "gmail.gm") ||
		strings.HasSuffix(domain, "gmail.gw") ||
		strings.HasSuffix(domain, "gmail.gn") ||
		strings.HasSuffix(domain, "gmail.sl") ||
		strings.HasSuffix(domain, "gmail.lr") ||
		strings.HasSuffix(domain, "gmail.ci") ||
		strings.HasSuffix(domain, "gmail.gh") ||
		strings.HasSuffix(domain, "gmail.tg") ||
		strings.HasSuffix(domain, "gmail.bj") ||
		strings.HasSuffix(domain, "gmail.ne") ||
		strings.HasSuffix(domain, "gmail.bf") ||
		strings.HasSuffix(domain, "gmail.ml") ||
		strings.HasSuffix(domain, "gmail.gn") ||
		strings.HasSuffix(domain, "gmail.cf") ||
		strings.HasSuffix(domain, "gmail.cm") ||
		strings.HasSuffix(domain, "gmail.td") ||
		strings.HasSuffix(domain, "gmail.cg") ||
		strings.HasSuffix(domain, "gmail.ga") ||
		strings.HasSuffix(domain, "gmail.gq") ||
		strings.HasSuffix(domain, "gmail.st") ||
		strings.HasSuffix(domain, "gmail.ao") ||
		strings.HasSuffix(domain, "gmail.cd") ||
		strings.HasSuffix(domain, "gmail.zr") ||
		strings.HasSuffix(domain, "gmail.rw") ||
		strings.HasSuffix(domain, "gmail.bi") ||
		strings.HasSuffix(domain, "gmail.mw") ||
		strings.HasSuffix(domain, "gmail.zm") ||
		strings.HasSuffix(domain, "gmail.zw") ||
		strings.HasSuffix(domain, "gmail.na") ||
		strings.HasSuffix(domain, "gmail.bw") ||
		strings.HasSuffix(domain, "gmail.ls") ||
		strings.HasSuffix(domain, "gmail.sz") ||
		strings.HasSuffix(domain, "gmail.ke") ||
		strings.HasSuffix(domain, "gmail.tz") ||
		strings.HasSuffix(domain, "gmail.ug") ||
		strings.HasSuffix(domain, "gmail.et") ||
		strings.HasSuffix(domain, "gmail.so") ||
		strings.HasSuffix(domain, "gmail.dj") ||
		strings.HasSuffix(domain, "gmail.km") ||
		strings.HasSuffix(domain, "gmail.mg") ||
		strings.HasSuffix(domain, "gmail.mu") ||
		strings.HasSuffix(domain, "gmail.sc") ||
		strings.HasSuffix(domain, "gmail.re") ||
		strings.HasSuffix(domain, "gmail.yt") ||
		strings.HasSuffix(domain, "gmail.com") {
		// Remove all dots from the local part for Gmail addresses
		localPart = strings.ReplaceAll(localPart, ".", "")
	}

	// Reconstruct the email
	return localPart + "@" + domain
}

// Helper: Charset-Erkennung (sehr einfach, z.B. ASCII, Latin, Cyrillic, etc.)
func detectCharset(s string) string {
	if s == "" {
		return "ASCII"
	}
	ascii := true
	latin := true
	cyrillic := true
	for _, r := range s {
		if r > 127 {
			ascii = false
		}
		if !(r >= 0x0020 && r <= 0x007E) && !(r >= 0x00A0 && r <= 0x00FF) {
			latin = false
		}
		if !(r >= 0x0400 && r <= 0x04FF) {
			cyrillic = false
		}
	}
	if ascii {
		return "ASCII"
	}
	if latin {
		return "Latin"
	}
	if cyrillic {
		return "Cyrillic"
	}
	if utf8.ValidString(s) {
		return "UTF-8"
	}
	return "Other"
}

// FilterRequestHandler prüft jetzt auch Charset-Regeln
func FilterRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var input FilterRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
			return
		}

		// Comprehensive validation using our validation package
		validationResult := validation.ValidateFilterRequest(input.IP, input.Email, input.UserAgent, input.Country, input.Username, input.Content)
		if !validationResult.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationResult.Errors,
			})
			return
		}

		// Validate that at least one filter field is provided
		if input.IP == "" && input.Email == "" && input.UserAgent == "" && input.Country == "" && input.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one filter field must be provided (ip, email, user_agent, country, or username)"})
			return
		}

		// Normalize email address (remove dots for Gmail addresses)
		normalizedEmail := normalizeEmail(input.Email)

		// Generate a cache key based on the normalized filter input
		cache := services.GetCacheFactory()
		cacheKey := "filter:" + input.IP + ":" + normalizedEmail + ":" + input.UserAgent + ":" + input.Country + ":" + input.Username

		// Track cache hit status BEFORE processing
		cacheHit := false
		if cached, exists, _ := cache.Get(cacheKey); exists {
			cacheHit = true

			// Log cache hit asynchronously before returning
			go func() {
				trafficLogging := services.NewTrafficLoggingService(db)

				// Convert to traffic logging format
				trafficReq := services.FilterRequest{
					IPAddress: input.IP,
					Email:     input.Email,
					UserAgent: input.UserAgent,
					Username:  input.Username,
					Country:   input.Country,
					Content:   input.Content,
				}

				// Create filter result for cache hit
				var trafficResult services.TrafficFilterResult

				// Handle different cache result types
				if filterResult, ok := cached.(services.FilterResult); ok {
					// Cache contains FilterResult
					trafficResult = services.TrafficFilterResult{
						FinalResult: filterResult.Result,
						FilterResults: map[string]interface{}{
							"result": filterResult.Result,
							"reason": filterResult.Reason,
							"field":  filterResult.Field,
							"value":  filterResult.Value,
						},
						ResponseTime: time.Since(startTime),
						CacheHit:     true,
					}
				} else if ginResult, ok := cached.(gin.H); ok {
					// Cache contains gin.H (from charset rules)
					trafficResult = services.TrafficFilterResult{
						FinalResult: ginResult["result"].(string),
						FilterResults: map[string]interface{}{
							"result": ginResult["result"],
							"reason": ginResult["reason"],
							"field":  ginResult["field"],
							"value":  ginResult["value"],
						},
						ResponseTime: time.Since(startTime),
						CacheHit:     true,
					}
				} else {
					// Fallback for unknown cache types
					trafficResult = services.TrafficFilterResult{
						FinalResult: "allowed",
						FilterResults: map[string]interface{}{
							"result": "allowed",
							"reason": "cached",
						},
						ResponseTime: time.Since(startTime),
						CacheHit:     true,
					}
				}

				// Create metadata
				metadata := map[string]string{
					"client_ip":      c.ClientIP(),
					"user_agent_raw": c.GetHeader("User-Agent"),
					"session_id":     c.GetHeader("X-Session-ID"),
				}

				// Log the request
				trafficLogging.LogFilterRequest(trafficReq, trafficResult, metadata)
			}()

			c.JSON(http.StatusOK, cached)
			return
		}

		// Lade alle Charset-Regeln
		var charsetRules []models.CharsetRule
		db.Find(&charsetRules)

		// Prüfe Email, UserAgent, Content, Username auf Charset-Regeln
		fields := map[string]string{
			"email":      input.Email, // Use original email for charset detection
			"user_agent": input.UserAgent,
			"content":    input.Content,
			"username":   input.Username,
		}
		for field, value := range fields {
			if value == "" {
				continue
			}
			cs := detectCharset(value)
			for _, rule := range charsetRules {
				if rule.Charset == cs {
					if rule.Status == "denied" {
						cache.Set(cacheKey, gin.H{"result": "denied", "reason": "charset denied", "field": field, field: value}, 5*time.Minute)
						c.JSON(200, gin.H{"result": "denied", "reason": "charset denied", "field": field, field: value})
						return
					}
					if rule.Status == "whitelisted" {
						cache.Set(cacheKey, gin.H{"result": "whitelisted", "reason": "charset whitelisted", "field": field, field: value}, 5*time.Minute)
						c.JSON(200, gin.H{"result": "whitelisted", "reason": "charset whitelisted", "field": field, field: value})
						return
					}
				}
			}
		}

		// Username-Filter: Prüfe gegen UsernameRule-Liste
		// Note: Username filtering is now handled by Elasticsearch regex filtering
		// This allows for both exact matches and regex patterns

		// Timeout for the entire operation (e.g., 5 seconds)
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Call the service to evaluate filters with normalized email
		finalResult, err := services.EvaluateFilters(ctx, input.IP, normalizedEmail, input.UserAgent, input.Country, input.Username)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		// Cache the result for 5 minutes
		cache.Set(cacheKey, finalResult, 5*time.Minute)

		// Log the traffic asynchronously
		go func() {
			trafficLogging := services.NewTrafficLoggingService(db)

			// Convert to traffic logging format
			trafficReq := services.FilterRequest{
				IPAddress: input.IP,
				Email:     input.Email,
				UserAgent: input.UserAgent,
				Username:  input.Username,
				Country:   input.Country,
				Content:   input.Content,
			}

			// Create filter result
			trafficResult := services.TrafficFilterResult{
				FinalResult: finalResult.Result,
				FilterResults: map[string]interface{}{
					"result": finalResult.Result,
					"reason": finalResult.Reason,
					"field":  finalResult.Field,
					"value":  finalResult.Value,
				},
				ResponseTime: time.Since(startTime),
				CacheHit:     cacheHit, // Use the cacheHit variable from above
			}

			// Create metadata
			metadata := map[string]string{
				"client_ip":      c.ClientIP(),
				"user_agent_raw": c.GetHeader("User-Agent"),
				"session_id":     c.GetHeader("X-Session-ID"),
			}

			// Log the request
			trafficLogging.LogFilterRequest(trafficReq, trafficResult, metadata)
		}()

		c.JSON(http.StatusOK, finalResult)
	}
}
