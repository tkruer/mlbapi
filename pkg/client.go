package mlbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// JSON is the generic object shape returned by raw endpoint calls.
type JSON = map[string]any

// Params contains endpoint parameters passed to raw endpoint calls.
//
// Keys are matched against endpoint path parameters and query parameters.
// Values are stringified before being sent to the MLB Stats API.
type Params map[string]any

// PathParamDefinition describes how a path parameter is rendered for an endpoint.
type PathParamDefinition struct {
	Type          string
	Default       string
	HasDefault    bool
	TrueValue     string
	FalseValue    string
	LeadingSlash  bool
	TrailingSlash bool
	Required      bool
}

// EndpointDefinition describes a raw MLB Stats API endpoint.
type EndpointDefinition struct {
	Name           string
	PathTemplate   string
	PathParamOrder []string
	PathParams     map[string]PathParamDefinition
	QueryParams    []string
	RequiredParams [][]string
	Note           string
}

// Client is an MLB Stats API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Option configures a Client created by NewClient.
type Option func(*Client)

// WithBaseURL overrides the default Stats API base URL for a client.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if strings.TrimSpace(baseURL) != "" {
			c.BaseURL = baseURL
		}
	}
}

// WithHTTPClient sets the HTTP client used for requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.HTTPClient = httpClient
		}
	}
}

// NewClient creates a new MLB Stats API client.
func NewClient(opts ...Option) *Client {
	client := &Client{
		BaseURL:    DefaultBaseURL,
		HTTPClient: http.DefaultClient,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	return client
}

// DefaultClient is the package-level client used by package helper functions.
var DefaultClient = NewClient()

// Get performs a raw endpoint call using DefaultClient.
func Get(ctx context.Context, endpoint string, params Params) (JSON, error) {
	return DefaultClient.Get(ctx, endpoint, params)
}

// GetForce performs a raw endpoint call and forces unknown parameters into the query string.
func GetForce(ctx context.Context, endpoint string, params Params) (JSON, error) {
	return DefaultClient.GetForce(ctx, endpoint, params)
}

// Get performs a raw endpoint call using the receiver client.
func (c *Client) Get(ctx context.Context, endpoint string, params Params) (JSON, error) {
	return c.get(ctx, endpoint, params, false)
}

// GetForce performs a raw endpoint call using the receiver client and forces unknown parameters into the query string.
func (c *Client) GetForce(ctx context.Context, endpoint string, params Params) (JSON, error) {
	return c.get(ctx, endpoint, params, true)
}

func (c *Client) get(ctx context.Context, endpoint string, params Params, force bool) (JSON, error) {
	definition, ok := defaultEndpoints[endpoint]
	if !ok {
		return nil, fmt.Errorf("mlbapi: invalid endpoint %q", endpoint)
	}

	requestURL, err := c.buildURL(definition, params, force)
	if err != nil {
		return nil, err
	}

	var payload JSON
	if err := c.fetchJSON(ctx, endpoint, requestURL, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *Client) fetchJSON(ctx context.Context, endpoint, requestURL string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("mlbapi: create request: %w", err)
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("mlbapi: request %q: %w", endpoint, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("mlbapi: %s returned %d: %s", endpoint, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("mlbapi: decode %q response: %w", endpoint, err)
	}

	return nil
}

// Notes returns local documentation for a raw endpoint, including supported parameters and any endpoint notes.
func (c *Client) Notes(endpoint string) (string, error) {
	definition, ok := defaultEndpoints[endpoint]
	if !ok {
		return "", fmt.Errorf("mlbapi: invalid endpoint %q", endpoint)
	}

	requiredPathParams := make([]string, 0, len(definition.PathParamOrder))
	for _, name := range definition.PathParamOrder {
		if definition.PathParams[name].Required {
			requiredPathParams = append(requiredPathParams, name)
		}
	}
	if len(requiredPathParams) == 0 {
		requiredPathParams = []string{"None"}
	}

	requiredQueryParams := "None"
	if !satisfiesRequiredParams(definition.RequiredParams, map[string]struct{}{}) {
		requiredQueryParams = fmt.Sprintf("%v", definition.RequiredParams)
	}

	var builder strings.Builder
	builder.WriteString("Endpoint: ")
	builder.WriteString(endpoint)
	builder.WriteString(" \n")
	builder.WriteString("All path parameters: ")
	builder.WriteString(fmt.Sprintf("%v", definition.PathParamOrder))
	builder.WriteString(". \n")
	builder.WriteString("Required path parameters (note: ver will be included by default): ")
	builder.WriteString(fmt.Sprintf("%v", requiredPathParams))
	builder.WriteString(". \n")
	builder.WriteString("All query parameters: ")
	builder.WriteString(fmt.Sprintf("%v", definition.QueryParams))
	builder.WriteString(". \n")
	builder.WriteString("Required query parameters: ")
	builder.WriteString(requiredQueryParams)
	builder.WriteString(". \n")
	if contains(definition.QueryParams, "hydrate") {
		builder.WriteString("The hydrate function is supported by this endpoint. Call the endpoint with {'hydrate':'hydrations'} in the parameters to return a list of available hydrations. For example, statsapi.get('schedule',{'sportId':1,'hydrate':'hydrations','fields':'hydrations'})\n")
	}
	if definition.Note != "" {
		builder.WriteString("Developer notes: ")
		builder.WriteString(definition.Note)
	}

	return builder.String(), nil
}

// Notes returns local documentation for a raw endpoint using DefaultClient metadata.
func Notes(endpoint string) (string, error) {
	return DefaultClient.Notes(endpoint)
}

func (c *Client) buildURL(definition EndpointDefinition, params Params, force bool) (string, error) {
	pathParams := make(map[string]string)
	queryValues := url.Values{}
	providedQueryParams := make(map[string]struct{})

	for key, value := range params {
		if pathParam, ok := definition.PathParams[key]; ok {
			if pathParam.Type == "bool" {
				if boolValue, ok := parseBool(value); ok {
					if boolValue {
						pathParams[key] = pathParam.TrueValue
					} else {
						pathParams[key] = pathParam.FalseValue
					}
				}
			} else {
				pathParams[key] = stringify(value)
			}
			continue
		}

		if contains(definition.QueryParams, key) || force {
			queryValues.Set(key, stringify(value))
			providedQueryParams[key] = struct{}{}
		}
	}

	requestURL := joinBaseAndPath(c.BaseURL, definition.PathTemplate)
	for name, value := range pathParams {
		requestURL = strings.ReplaceAll(requestURL, "{"+name+"}", applySlashes(definition.PathParams[name], value))
	}

	for {
		start := strings.IndexByte(requestURL, '{')
		end := strings.IndexByte(requestURL, '}')
		if start == -1 || end == -1 || end < start {
			break
		}

		paramName := requestURL[start+1 : end]
		pathParam, ok := definition.PathParams[paramName]
		if !ok {
			requestURL = strings.Replace(requestURL, "{"+paramName+"}", "", 1)
			continue
		}

		replacement := ""
		switch {
		case !pathParam.Required:
			replacement = ""
		case pathParam.HasDefault && pathParam.Default != "":
			replacement = applySlashes(pathParam, pathParam.Default)
		case force:
			replacement = ""
		default:
			return "", fmt.Errorf("mlbapi: missing required path parameter {%s} for endpoint %q", paramName, definition.Name)
		}

		requestURL = strings.Replace(requestURL, "{"+paramName+"}", replacement, 1)
	}

	if !force && !satisfiesRequiredParams(definition.RequiredParams, providedQueryParams) {
		return "", fmt.Errorf("mlbapi: missing required query parameters for endpoint %q; required parameter sets: %v", definition.Name, definition.RequiredParams)
	}

	if encoded := queryValues.Encode(); encoded != "" {
		requestURL += "?" + encoded
	}

	return requestURL, nil
}

func applySlashes(param PathParamDefinition, value string) string {
	var builder strings.Builder
	if param.LeadingSlash {
		builder.WriteByte('/')
	}
	builder.WriteString(value)
	if param.TrailingSlash {
		builder.WriteByte('/')
	}
	return builder.String()
}

func joinBaseAndPath(baseURL, path string) string {
	trimmedBase := strings.TrimRight(baseURL, "/")
	trimmedPath := strings.TrimLeft(path, "/")
	return trimmedBase + "/" + trimmedPath
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func satisfiesRequiredParams(required [][]string, provided map[string]struct{}) bool {
	for _, group := range required {
		if len(group) == 0 {
			return true
		}

		satisfied := true
		for _, param := range group {
			if _, ok := provided[param]; !ok {
				satisfied = false
				break
			}
		}

		if satisfied {
			return true
		}
	}

	return len(required) == 0
}

func stringify(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case TeamID:
		return strconv.Itoa(int(typed))
	case time.Time:
		return typed.Format(time.RFC3339)
	case fmt.Stringer:
		return typed.String()
	case bool:
		return strconv.FormatBool(typed)
	case int:
		return strconv.Itoa(typed)
	case int8:
		return strconv.FormatInt(int64(typed), 10)
	case int16:
		return strconv.FormatInt(int64(typed), 10)
	case int32:
		return strconv.FormatInt(int64(typed), 10)
	case int64:
		return strconv.FormatInt(typed, 10)
	case uint:
		return strconv.FormatUint(uint64(typed), 10)
	case uint8:
		return strconv.FormatUint(uint64(typed), 10)
	case uint16:
		return strconv.FormatUint(uint64(typed), 10)
	case uint32:
		return strconv.FormatUint(uint64(typed), 10)
	case uint64:
		return strconv.FormatUint(typed, 10)
	case float32:
		return strconv.FormatFloat(float64(typed), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	default:
		return fmt.Sprint(value)
	}
}

func parseBool(value any) (bool, bool) {
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		parsed, err := strconv.ParseBool(strings.ToLower(strings.TrimSpace(typed)))
		if err != nil {
			return false, false
		}
		return parsed, true
	default:
		return false, false
	}
}

func (c *Client) httpClient() *http.Client {
	if c != nil && c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}
