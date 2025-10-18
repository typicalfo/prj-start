# Quick Start: Upstash MCP for Development Agents

## Your Mission
You have access to a Upstash Vector database containing Go development examples, patterns, and best practices. Use this knowledge to accelerate development and ensure code quality.

## Immediate Actions

### 1. Explore the Database (5 minutes)
```
"What Go Fiber recipes are available?"
"Show me clean architecture examples"
"What database patterns are documented?"
"How is authentication implemented?"
```

### 2. Understand the Data Structure
- **Web Development**: Go Fiber v2 examples
- **Architecture**: Clean architecture patterns
- **Data**: PostgreSQL, ORM patterns
- **APIs**: REST, OpenAPI implementations
- **Testing**: Unit tests, mocking strategies
- **Deployment**: Vercel configurations

### 3. Create Your Instruction File
After exploration, create `[PROJECT_NAME]_INSTRUCTIONS.md` with:

```markdown
# [Project] Development Instructions

## Project Context
[Project type, tech stack, requirements]

## Relevant Data
[Which database sections apply to this project]

## Query Templates
- "Show me [feature] implementation in Go Fiber"
- "What are the patterns for [concept]?"
- "How do examples handle [problem]?"

## Development Workflow
1. Query for existing implementations
2. Analyze and adapt patterns
3. Apply to project context
4. Validate against best practices
```

## Effective Query Patterns

**For Features:**
- "Show me how to implement [feature] with [requirements]"

**For Architecture:**
- "What are the clean architecture patterns for [domain]?"

**For Problems:**
- "How do examples solve [specific issue] with [constraints]?"

## Best Practices
1. **Query First**: Always search before implementing
2. **Adapt, Don't Copy**: Modify patterns for your context
3. **Validate**: Ensure solutions fit project requirements
4. **Document**: Note successful approaches for future use

## Example Workflow
```
User: "Need user authentication"
Agent: "Searching authentication patterns..."
Query: "Show me Go Fiber authentication with JWT"
Results: [Relevant examples]
Agent: "Based on examples, here are three approaches..."
```

## Next Steps
1. Run exploration queries
2. Analyze relevant patterns
3. Create project-specific instructions
4. Begin development with MCP assistance

Your goal: Leverage existing knowledge to build faster, better, and more consistent Go applications.