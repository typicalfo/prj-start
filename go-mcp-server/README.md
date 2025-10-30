# MCP Go SDK Documentation

This comprehensive guide covers creating MCP (Model Context Protocol) servers using the official Go SDK. The documentation is based on the official GitHub repository at https://github.com/modelcontextprotocol/go-sdk.

## Table of Contents

1. [Installation and Setup](#installation-and-setup)
2. [Basic Server Creation](#basic-server-creation)
3. [Tool Implementation](#tool-implementation)
4. [Resource Handling](#resource-handling)
5. [Prompts](#prompts)
6. [Sampling](#sampling)
7. [Error Handling](#error-handling)
8. [Testing and Deployment](#testing-and-deployment)
9. [Advanced Features](#advanced-features)
10. [Examples](#examples)

## Installation and Setup

### Go Version Requirements

The SDK requires Go 1.23.0 or later:

```go
module github.com/your-project

go 1.23.0

require (
    github.com/modelcontextprotocol/go-sdk v0.5.0
)
```

### Dependencies

The SDK has the following key dependencies:

- `github.com/golang-jwt/jwt/v5 v5.2.2` - JWT token handling
- `github.com/google/go-cmp v0.7.0` - Comparison utilities
- `github.com/google/jsonschema-go v0.3.0` - JSON Schema validation
- `github.com/yosida95/uritemplate/v3 v3.0.2` - URI template handling
- `golang.org/x/oauth2 v0.30.0` - OAuth2 support
- `golang.org/x/tools v0.34.0` - Go tooling

### Basic Setup

```bash
# Initialize a new Go module
go mod init your-mcp-server

# Add the MCP SDK dependency
go get github.com/modelcontextprotocol/go-sdk

# Download dependencies
go mod tidy
```

## Basic Server Creation

### Simple Server

Here's the most basic MCP server that provides a greeting tool:

```go
package main

import (
    "context"
    "log"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
    Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type Output struct {
    Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (
    *mcp.CallToolResult,
    Output,
    error,
) {
    return nil, Output{Greeting: "Hi " + input.Name}, nil
}

func main() {
    // Create a server with implementation info
    server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

    // Add a tool
    mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)

    // Run the server over stdin/stdout
    if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
        log.Fatal(err)
    }
}
```

### Server with Options

```go
server := mcp.NewServer(&mcp.Implementation{
    Name:    "my-server",
    Version: "1.0.0",
}, &mcp.ServerOptions{
    InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
        log.Println("Client initialized connection")
    },
    KeepAlive: 30 * time.Second, // Send ping every 30 seconds
})
```

### Client Connection

To connect to an MCP server from a client:

```go
package main

import (
    "context"
    "log"
    "os/exec"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
    ctx := context.Background()

    // Create a client
    client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

    // Connect to a server over stdin/stdout
    transport := &mcp.CommandTransport{Command: exec.Command("myserver")}
    session, err := client.Connect(ctx, transport, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    // Call a tool
    params := &mcp.CallToolParams{
        Name:      "greet",
        Arguments: map[string]any{"name": "you"},
    }
    res, err := session.CallTool(ctx, params)
    if err != nil {
        log.Fatalf("CallTool failed: %v", err)
    }
    if res.IsError {
        log.Fatal("tool failed")
    }
    for _, c := range res.Content {
        log.Print(c.(*mcp.TextContent).Text)
    }
}
```

## Tool Implementation

### Basic Tool with Schema Inference

The SDK can automatically infer JSON schemas from Go types:

```go
type CalculatorInput struct {
    X int `json:"x" jsonschema:"first number"`
    Y int `json:"y" jsonschema:"second number"`
    Operation string `json:"operation" jsonschema:"operation to perform (add, subtract, multiply, divide)"`
}

type CalculatorOutput struct {
    Result int `json:"result" jsonschema:"calculation result"`
}

func Calculate(ctx context.Context, req *mcp.CallToolRequest, input CalculatorInput) (
    *mcp.CallToolResult,
    CalculatorOutput,
    error,
) {
    var result int
    switch input.Operation {
    case "add":
        result = input.X + input.Y
    case "subtract":
        result = input.X - input.Y
    case "multiply":
        result = input.X * input.Y
    case "divide":
        if input.Y == 0 {
            return nil, CalculatorOutput{}, fmt.Errorf("division by zero")
        }
        result = input.X / input.Y
    default:
        return nil, CalculatorOutput{}, fmt.Errorf("unknown operation: %s", input.Operation)
    }

    return nil, CalculatorOutput{Result: result}, nil
}

// Add to server
mcp.AddTool(server, &mcp.Tool{
    Name:        "calculate",
    Description: "Perform basic arithmetic operations",
}, Calculate)
```

### Custom Schema Definition

For more complex schemas, define them manually:

```go
import "github.com/google/jsonschema-go/jsonschema"

customSchema := &jsonschema.Schema{
    Type: "object",
    Properties: map[string]*jsonschema.Schema{
        "name": {
            Type: "string",
            MinLength: jsonschema.Ptr(1),
            MaxLength: jsonschema.Ptr(100),
        },
        "age": {
            Type: "integer",
            Minimum: jsonschema.Ptr(0.0),
            Maximum: jsonschema.Ptr(150.0),
        },
    },
    Required: []string{"name"},
}

mcp.AddTool(server, &mcp.Tool{
    Name:        "custom_tool",
    Description: "Tool with custom schema",
    InputSchema: customSchema,
}, handler)
```

### Tool Handler Patterns

#### Simple Handler (Inferred Types)

```go
func SimpleHandler(ctx context.Context, req *mcp.CallToolRequest, input MyInput) (
    *mcp.CallToolResult,
    MyOutput,
    error,
) {
    // Process input
    result := MyOutput{...}

    // Return nil for CallToolResult to use default formatting
    return nil, result, nil
}
```

#### Advanced Handler (Manual Control)

```go
func AdvancedHandler(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Manual parsing and validation
    var input MyInput
    if err := json.Unmarshal(req.Params.Arguments, &input); err != nil {
        return &mcp.CallToolResult{
            Content: []mcp.Content{&mcp.TextContent{Text: "Invalid input: " + err.Error()}},
            IsError: true,
        }, nil
    }

    // Process
    result := process(input)

    // Return structured result
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: fmt.Sprintf("Result: %v", result)},
        },
        StructuredContent: result,
    }, nil
}
```

### Tool Management

```go
// Add multiple tools
server.AddTool(tool1, handler1)
server.AddTool(tool2, handler2)

// Remove tools
server.RemoveTools("tool1", "tool2")

// List tools (client-side)
for tool, err := range session.Tools(ctx, nil) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Tool:", tool.Name, "-", tool.Description)
}
```

## Resource Handling

### Adding Resources

```go
// Static resource
server.AddResource(&mcp.Resource{
    URI:         "file:///static.txt",
    Name:        "Static File",
    Description: "A static text file",
    MIMEType:    "text/plain",
}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
    return &mcp.ReadResourceResult{
        Contents: []*mcp.ResourceContents{{
            URI:      req.Params.URI,
            Text:     "This is static content",
            MIMEType: "text/plain",
        }},
    }, nil
})

// Dynamic resource
server.AddResource(&mcp.Resource{
    URI:         "file:///timestamp.txt",
    Name:        "Current Time",
    Description: "Current timestamp",
}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
    timestamp := time.Now().Format(time.RFC3339)
    return &mcp.ReadResourceResult{
        Contents: []*mcp.ResourceContents{{
            URI:  req.Params.URI,
            Text: timestamp,
        }},
    }, nil
})
```

### Resource Templates

```go
// Template for numbered files
server.AddResourceTemplate(&mcp.ResourceTemplate{
    URITemplate: "file:///logs/{date}/{hour}.log",
    Name:        "Hourly Logs",
    Description: "Log files by date and hour",
}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
    // Extract parameters from URI
    // Implementation would parse date/hour from req.Params.URI

    logContent := fmt.Sprintf("Logs for %s", req.Params.URI)
    return &mcp.ReadResourceResult{
        Contents: []*mcp.ResourceContents{{
            URI:  req.Params.URI,
            Text: logContent,
        }},
    }, nil
})
```

### Client-Side Resource Access

```go
// List resources
for resource, err := range session.Resources(ctx, nil) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Resource: %s - %s\n", resource.URI, resource.Name)
}

// Read a resource
result, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
    URI: "file:///example.txt",
})
if err != nil {
    log.Fatal(err)
}
for _, content := range result.Contents {
    fmt.Println("Content:", content.Text)
}
```

### Resource Subscriptions

```go
// Server-side subscription handling
server := mcp.NewServer(&mcp.Implementation{Name: "server", Version: "1.0.0"}, &mcp.ServerOptions{
    SubscribeHandler: func(ctx context.Context, req *mcp.ServerRequest[*mcp.SubscribeParams]) error {
        log.Printf("Client subscribed to: %s", req.Params.URI)
        return nil
    },
    UnsubscribeHandler: func(ctx context.Context, req *mcp.ServerRequest[*mcp.UnsubscribeParams]) error {
        log.Printf("Client unsubscribed from: %s", req.Params.URI)
        return nil
    },
})

// Client-side subscription
err := session.Subscribe(ctx, &mcp.SubscribeParams{URI: "file:///dynamic.txt"})
if err != nil {
    log.Fatal(err)
}

// Handle updates
client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, &mcp.ClientOptions{
    ResourceUpdatedHandler: func(ctx context.Context, req *mcp.ResourceUpdatedRequest) {
        log.Printf("Resource updated: %s", req.Params.URI)
    },
})
```

## Prompts

### Adding Prompts

```go
type CodeReviewInput struct {
    Code     string `json:"code" jsonschema:"the code to review"`
    Language string `json:"language" jsonschema:"programming language"`
}

func CodeReviewPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    var input CodeReviewInput
    if err := json.Unmarshal(req.Params.Arguments, &input); err != nil {
        return nil, err
    }

    return &mcp.GetPromptResult{
        Description: "Code review prompt",
        Messages: []*mcp.PromptMessage{
            {
                Role: "system",
                Content: &mcp.TextContent{
                    Text: "You are an expert code reviewer. Provide constructive feedback.",
                },
            },
            {
                Role: "user",
                Content: &mcp.TextContent{
                    Text: fmt.Sprintf("Please review this %s code:\n\n%s", input.Language, input.Code),
                },
            },
        },
    }, nil
}

// Add prompt with arguments
server.AddPrompt(&mcp.Prompt{
    Name:        "code_review",
    Description: "Generate a code review prompt",
    Arguments: []*mcp.PromptArgument{
        {
            Name:        "code",
            Description: "The code to review",
            Required:    true,
        },
        {
            Name:        "language",
            Description: "Programming language",
            Required:    false,
        },
    },
}, CodeReviewPrompt)
```

### Client-Side Prompt Usage

```go
// List prompts
for prompt, err := range session.Prompts(ctx, nil) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Prompt: %s - %s\n", prompt.Name, prompt.Description)
}

// Get a prompt
result, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
    Name: "code_review",
    Arguments: map[string]string{
        "code":     "func hello() { fmt.Println(\"Hello\") }",
        "language": "Go",
    },
})
if err != nil {
    log.Fatal(err)
}

for _, msg := range result.Messages {
    fmt.Printf("%s: %s\n", msg.Role, msg.Content.(*mcp.TextContent).Text)
}
```

## Sampling

### Server-Side Sampling Request

```go
// Server requesting sampling from client
result, err := session.CreateMessage(ctx, &mcp.CreateMessageParams{
    Messages: []*mcp.SamplingMessage{
        {
            Role: "user",
            Content: &mcp.TextContent{Text: "Hello, how are you?"},
        },
    },
    ModelPreferences: &mcp.ModelPreferences{
        Hints: []*mcp.ModelHint{
            {Name: "gpt-4"},
        },
        IntelligencePriority: 0.8,
    },
    SystemPrompt: "You are a helpful assistant.",
    IncludeContext: "thisRequest",
    Temperature:    0.7,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Sampled message:", result.Content.(*mcp.TextContent).Text)
```

### Client-Side Sampling Handler

```go
client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, &mcp.ClientOptions{
    CreateMessageHandler: func(ctx context.Context, req *mcp.CreateMessageRequest) (*mcp.CreateMessageResult, error) {
        // This would integrate with your AI service
        // For example, call OpenAI, Anthropic, etc.

        // Mock response
        return &mcp.CreateMessageResult{
            Content: &mcp.TextContent{
                Text: "I'm doing well, thank you for asking!",
            },
            Model: "gpt-4",
            StopReason: "endTurn",
        }, nil
    },
})
```

## Error Handling

### Tool Errors

```go
func RiskyTool(ctx context.Context, req *mcp.CallToolRequest, input Input) (
    *mcp.CallToolResult,
    Output,
    error,
) {
    if input.Value < 0 {
        // Return error - this becomes a tool error
        return nil, Output{}, fmt.Errorf("value must be non-negative")
    }

    return nil, Output{Result: input.Value * 2}, nil
}
```

### Protocol Errors

Protocol errors are handled automatically by the SDK. When a server handler returns an error, it's converted to a JSON-RPC error response.

```go
// Server handler that might fail
func UnreliableTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    if rand.Float32() < 0.1 { // 10% chance
        return nil, mcp.NewError(mcp.InvalidParams, "random failure")
    }

    return &mcp.CallToolResult{
        Content: []mcp.Content{&mcp.TextContent{Text: "Success!"}},
    }, nil
}
```

### Client Error Handling

```go
result, err := session.CallTool(ctx, &mcp.CallToolParams{
    Name:      "unreliable_tool",
    Arguments: map[string]any{},
})
if err != nil {
    // Check if it's a protocol error
    var jsonErr *mcp.JSONRPCError
    if errors.As(err, &jsonErr) {
        log.Printf("Protocol error %d: %s", jsonErr.Code, jsonErr.Message)
    } else {
        log.Printf("Other error: %v", err)
    }
    return
}

if result.IsError {
    // Tool returned an error
    for _, content := range result.Content {
        if text, ok := content.(*mcp.TextContent); ok {
            log.Printf("Tool error: %s", text.Text)
        }
    }
}
```

## Testing and Deployment

### Unit Testing Tools

```go
func TestCalculatorTool(t *testing.T) {
    // Create a test server
    server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
    mcp.AddTool(server, &mcp.Tool{Name: "calculate"}, Calculate)

    // Create in-memory transport
    t1, t2 := mcp.NewInMemoryTransports()

    // Connect server
    session, err := server.Connect(context.Background(), t1, nil)
    if err != nil {
        t.Fatal(err)
    }
    defer session.Wait()

    // Create client
    client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
    clientSession, err := client.Connect(context.Background(), t2, nil)
    if err != nil {
        t.Fatal(err)
    }
    defer clientSession.Close()

    // Test the tool
    result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
        Name: "calculate",
        Arguments: map[string]any{
            "x": 5,
            "y": 3,
            "operation": "add",
        },
    })
    if err != nil {
        t.Fatal(err)
    }

    if result.IsError {
        t.Fatal("Tool returned error")
    }

    // Parse result
    var output CalculatorOutput
    if err := json.Unmarshal(result.StructuredContent, &output); err != nil {
        t.Fatal(err)
    }

    if output.Result != 8 {
        t.Errorf("Expected 8, got %d", output.Result)
    }
}
```

### HTTP Server Deployment

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
    // Create server
    server := mcp.NewServer(&mcp.Implementation{
        Name:    "my-mcp-server",
        Version: "1.0.0",
    }, nil)

    // Add tools, resources, prompts...

    // Create HTTP handler
    handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
        return server
    }, &mcp.StreamableHTTPOptions{
        Stateless: false, // Use sessions for stateful operations
    })

    // Create HTTP server
    httpServer := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    log.Println("Starting MCP server on :8080")
    log.Fatal(httpServer.ListenAndServe())
}
```

### Docker Deployment

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Stdio Deployment

For command-line tools that communicate via stdin/stdout:

```go
func main() {
    server := mcp.NewServer(&mcp.Implementation{Name: "cli-tool", Version: "1.0.0"}, nil)
    // Add tools...

    if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
        log.Fatal(err)
    }
}
```

## Advanced Features

### Middleware

```go
// Logging middleware
func loggingMiddleware(handler mcp.MethodHandler) mcp.MethodHandler {
    return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
        log.Printf("Request: %s", method)
        result, err := handler(ctx, method, req)
        log.Printf("Response: %v, Error: %v", result, err)
        return result, err
    }
}

// Rate limiting middleware
func rateLimitMiddleware(handler mcp.MethodHandler) mcp.MethodHandler {
    // Implementation with token bucket or similar
    return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
        // Check rate limit
        return handler(ctx, method, req)
    }
}

// Apply middleware
server.AddReceivingMiddleware(loggingMiddleware, rateLimitMiddleware)
```

### Authentication

```go
// JWT Authentication
jwtAuth := auth.RequireBearerToken(verifyJWT, &auth.RequireBearerTokenOptions{
    Scopes: []string{"read", "write"},
})

// API Key Authentication
apiKeyAuth := auth.RequireBearerToken(verifyAPIKey, &auth.RequireBearerTokenOptions{
    Scopes: []string{"read"},
})

// Apply to HTTP handler
authenticatedHandler := jwtAuth(mcpHandler)
```

### Logging

```go
// Server-side logging to client
logger := slog.New(mcp.NewLoggingHandler(session, &mcp.LoggingHandlerOptions{
    LoggerName:   "my-server",
    MinInterval:  time.Second,
}))
logger.Info("Processing request", "user", "alice")
```

### Progress Reporting

```go
func LongRunningTool(ctx context.Context, req *mcp.CallToolRequest, input Input) (
    *mcp.CallToolResult,
    Output,
    error,
) {
    // Check if client requested progress
    if token := req.Params.GetProgressToken(); token != nil {
        for i := 0; i < 10; i++ {
            req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
                ProgressToken: token,
                Progress:      float64(i),
                Total:         10,
                Message:       fmt.Sprintf("Step %d/10", i+1),
            })
            time.Sleep(100 * time.Millisecond)
        }
    }

    return nil, Output{Result: "Done"}, nil
}
```

### Cancellation

```go
func CancellableTool(ctx context.Context, req *mcp.CallToolRequest, input Input) (
    *mcp.CallToolResult,
    Output,
    error,
) {
    // Long operation that can be cancelled
    select {
    case <-time.After(5 * time.Second):
        return nil, Output{Result: "Completed"}, nil
    case <-ctx.Done():
        return nil, Output{}, ctx.Err()
    }
}
```

## Examples

### Complete Server Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type WeatherInput struct {
    Location string `json:"location" jsonschema:"city or location name"`
    Days     int    `json:"days" jsonschema:"number of days to forecast"`
}

type WeatherOutput struct {
    Location string    `json:"location"`
    Forecast []string  `json:"forecast"`
    Updated  time.Time `json:"updated"`
}

func GetWeather(ctx context.Context, req *mcp.CallToolRequest, input WeatherInput) (
    *mcp.CallToolResult,
    WeatherOutput,
    error,
) {
    // Mock weather data
    forecast := make([]string, input.Days)
    for i := range forecast {
        forecast[i] = fmt.Sprintf("Day %d: Sunny, 72°F", i+1)
    }

    return nil, WeatherOutput{
        Location: input.Location,
        Forecast: forecast,
        Updated:  time.Now(),
    }, nil
}

func main() {
    server := mcp.NewServer(&mcp.Implementation{
        Name:    "weather-server",
        Version: "1.0.0",
    }, nil)

    // Add weather tool
    mcp.AddTool(server, &mcp.Tool{
        Name:        "get_weather",
        Description: "Get weather forecast for a location",
    }, GetWeather)

    // Add a static resource
    server.AddResource(&mcp.Resource{
        URI:         "weather://help",
        Name:        "Weather Help",
        Description: "Help information for weather tools",
        MIMEType:    "text/plain",
    }, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
        return &mcp.ReadResourceResult{
            Contents: []*mcp.ResourceContents{{
                URI:      req.Params.URI,
                Text:     "Use get_weather tool with location and days parameters",
                MIMEType: "text/plain",
            }},
        }, nil
    })

    log.Println("Weather MCP server starting...")
    if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
        log.Fatal(err)
    }
}
```

### Distributed Server Example

```go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "sync"
    "sync/atomic"
    "time"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
    httpAddr   = flag.String("http", "", "HTTP address")
    childPorts = flag.String("child_ports", "", "comma-separated child ports")
)

func main() {
    flag.Parse()

    if v := os.Getenv("MCP_CHILD_PORT"); v != "" {
        child(v)
    } else {
        parent()
    }
}

func parent() {
    exe, err := os.Executable()
    if err != nil {
        log.Fatal(err)
    }

    ports := strings.Split(*childPorts, ",")
    var wg sync.WaitGroup

    // Start child processes
    for i, port := range ports {
        cmd := exec.CommandContext(context.Background(), exe, os.Args[1:]...)
        cmd.Env = append(os.Environ(), fmt.Sprintf("MCP_CHILD_PORT=%s", port))
        cmd.Stderr = os.Stderr

        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := cmd.Run(); err != nil {
                log.Printf("Child failed: %v", err)
            }
        }()
    }

    // Reverse proxy
    var nextBackend atomic.Int64
    proxy := &httputil.ReverseProxy{
        Rewrite: func(r *httputil.ProxyRequest) {
            child := int(nextBackend.Add(1)) % len(ports)
            r.SetURL(&url.URL{
                Scheme: "http",
                Host:   fmt.Sprintf("localhost:%s", ports[child]),
            })
        },
    }

    log.Printf("Proxy listening on %s", *httpAddr)
    log.Fatal(http.ListenAndServe(*httpAddr, proxy))
}

func child(port string) {
    server := mcp.NewServer(&mcp.Implementation{Name: "counter"}, nil)

    var count atomic.Int64
    mcp.AddTool(server, &mcp.Tool{Name: "inc"}, func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, struct{ Count int64 }, error) {
        n := count.Add(1)
        req.Session.Log(ctx, &mcp.LoggingMessageParams{
            Data:  fmt.Sprintf("request %d", n),
            Level: "info",
        })
        return nil, struct{ Count int64 }{n}, nil
    })

    handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
        return server
    }, &mcp.StreamableHTTPOptions{Stateless: true})

    log.Printf("Child listening on localhost:%s", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), handler))
}
```

This documentation provides a comprehensive guide to building MCP servers with the Go SDK. The examples demonstrate various patterns and best practices for creating robust, feature-rich MCP implementations.