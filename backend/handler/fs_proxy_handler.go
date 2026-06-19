package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"posixfy-cloud/backend/middleware"
	"posixfy-cloud/backend/service"

	"github.com/gin-gonic/gin"
)

type FSProxyHandler struct {
	Client *service.FSClient
}

func NewFSProxyHandler(client *service.FSClient) *FSProxyHandler {
	return &FSProxyHandler{Client: client}
}

// buildUpstreamPath creates a URL path with properly encoded query parameters.
func buildUpstreamPath(endpoint string, mount string, path string) string {
	q := url.Values{}
	q.Set("mount", mount)
	q.Set("path", path)
	return endpoint + "?" + q.Encode()
}

func (h *FSProxyHandler) parseGroups(groupsJSON string) string {
	var groups []int
	if err := json.Unmarshal([]byte(groupsJSON), &groups); err != nil {
		return ""
	}
	parts := make([]string, len(groups))
	for i, g := range groups {
		parts[i] = strconv.Itoa(g)
	}
	return strings.Join(parts, ",")
}

// upstreamHeaders returns the extra headers forwarded to the bridge, always
// including the request correlation id so a single user action can be traced
// across cloud and bridge logs.
func upstreamHeaders(c *gin.Context, base map[string]string) map[string]string {
	if base == nil {
		base = map[string]string{}
	}
	if rid := middleware.GetRequestID(c); rid != "" {
		base[middleware.RequestIDHeader] = rid
	}
	return base
}

// logUpstreamFailure records a failed call to the bridge. These errors were
// previously swallowed, which made file-operation failures undiagnosable.
func logUpstreamFailure(c *gin.Context, op, upstream string, err error) {
	slog.Error("upstream request failed",
		"op", op,
		"upstream", upstream,
		"request_id", middleware.GetRequestID(c),
		"err", err,
	)
}

func (h *FSProxyHandler) Mounts(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)

	resp, err := h.Client.Do("GET", "/api/mounts", claims.UID, claims.GID, groups, nil, "", upstreamHeaders(c, nil))
	if err != nil {
		logUpstreamFailure(c, "mounts", "/api/mounts", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	proxyResponse(c, resp)
}

func (h *FSProxyHandler) List(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/list", mount, path)
	resp, err := h.Client.Do("GET", upstream, claims.UID, claims.GID, groups, nil, "", upstreamHeaders(c, nil))
	if err != nil {
		logUpstreamFailure(c, "list", upstream, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	// Never cache directory listings: a stale listing can make a just-deleted
	// file appear to "come back" after a refresh.
	c.Header("Cache-Control", "no-store")
	proxyResponse(c, resp)
}

func (h *FSProxyHandler) Download(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/file", mount, path)
	resp, err := h.Client.Do("GET", upstream, claims.UID, claims.GID, groups, nil, "", upstreamHeaders(c, nil))
	if err != nil {
		logUpstreamFailure(c, "download", upstream, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	// stream file download
	for _, key := range []string{"Content-Type", "Content-Disposition", "Content-Length"} {
		if v := resp.Header.Get(key); v != "" {
			c.Header(key, v)
		}
	}
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func (h *FSProxyHandler) Upload(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/upload", mount, path)
	contentType := c.GetHeader("Content-Type")

	// Forward OCC headers if present
	extra := map[string]string{}
	if v := c.GetHeader("X-Expected-MTime"); v != "" {
		extra["X-Expected-MTime"] = v
	}
	if v := c.GetHeader("X-Expected-Size"); v != "" {
		extra["X-Expected-Size"] = v
	}

	resp, err := h.Client.Do("POST", upstream, claims.UID, claims.GID, groups, c.Request.Body, contentType, upstreamHeaders(c, extra))
	if err != nil {
		logUpstreamFailure(c, "upload", upstream, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	proxyResponse(c, resp)
}

func (h *FSProxyHandler) Delete(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/delete", mount, path)
	resp, err := h.Client.Do("DELETE", upstream, claims.UID, claims.GID, groups, nil, "", upstreamHeaders(c, nil))
	if err != nil {
		logUpstreamFailure(c, "delete", upstream, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	proxyResponse(c, resp)
}

func (h *FSProxyHandler) Mkdir(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/mkdir", mount, path)
	resp, err := h.Client.Do("POST", upstream, claims.UID, claims.GID, groups, nil, "", upstreamHeaders(c, nil))
	if err != nil {
		logUpstreamFailure(c, "mkdir", upstream, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	proxyResponse(c, resp)
}

func proxyResponse(c *gin.Context, resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read upstream response body",
			"request_id", middleware.GetRequestID(c),
			"err", err,
		)
	}
	if resp.StatusCode >= 400 {
		slog.Warn("upstream returned error status",
			"status", resp.StatusCode,
			"request_id", middleware.GetRequestID(c),
			"body", string(body),
		)
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func (h *FSProxyHandler) Watch(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/watch", mount, path)
	resp, err := h.Client.DoStream("GET", upstream, claims.UID, claims.GID, groups, upstreamHeaders(c, nil))
	if err != nil {
		logUpstreamFailure(c, "watch", upstream, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(c.Writer, "%s\n", line)
		// Flush on empty line (SSE event boundary) or comment lines
		if line == "" || strings.HasPrefix(line, ":") {
			flusher.Flush()
		}
	}
}
