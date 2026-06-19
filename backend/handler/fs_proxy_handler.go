package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func (h *FSProxyHandler) Mounts(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)

	resp, err := h.Client.Do("GET", "/api/mounts", claims.UID, claims.GID, groups, nil, "", nil)
	if err != nil {
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
	resp, err := h.Client.Do("GET", upstream, claims.UID, claims.GID, groups, nil, "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	proxyResponse(c, resp)
}

func (h *FSProxyHandler) Download(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/file", mount, path)
	resp, err := h.Client.Do("GET", upstream, claims.UID, claims.GID, groups, nil, "", nil)
	if err != nil {
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

	resp, err := h.Client.Do("POST", upstream, claims.UID, claims.GID, groups, c.Request.Body, contentType, extra)
	if err != nil {
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
	resp, err := h.Client.Do("DELETE", upstream, claims.UID, claims.GID, groups, nil, "", nil)
	if err != nil {
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
	resp, err := h.Client.Do("POST", upstream, claims.UID, claims.GID, groups, nil, "", nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer resp.Body.Close()

	proxyResponse(c, resp)
}

func proxyResponse(c *gin.Context, resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func (h *FSProxyHandler) Watch(c *gin.Context) {
	claims := middleware.GetClaims(c)
	groups := h.parseGroups(claims.Groups)
	mount := c.Query("mount")
	path := c.Query("path")

	upstream := buildUpstreamPath("/api/fs/watch", mount, path)
	resp, err := h.Client.DoStream("GET", upstream, claims.UID, claims.GID, groups)
	if err != nil {
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
