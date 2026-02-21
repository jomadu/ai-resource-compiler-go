Configuration
Configuration file paths

You can configure Kiro CLI to match your development preferences and team standards. You can set configuration in one of three scopes:

Global - Configuration that is applied across all the projects where Kiro is used - <user-home>/.kiro/
Project - Configuration specific to a project - <project-root>/.kiro
Agent - Configuration defined in the agent configuration file - <user-home | project-root>/.kiro/agents
Configuration	Global Scope	Project Scope
MCP servers	~/.kiro/settings/mcp.json	.kiro/settings/mcp.json
Prompts	~/.kiro/prompts	.kiro/prompts
Custom agents	~/.kiro/agents	.kiro/agents
Steering	~/.kiro/steering	.kiro/steering
Settings	~/.kiro/settings/cli.json	
What can you configure at these scopes

Configuration	User Scope	Project Scope	Agent Scope
MCP servers	Yes	Yes	Yes
Prompts	Yes	Yes	No
Custom agents	Yes	Yes	N/A
Steering	Yes	Yes	Yes
Settings	Yes	N/A	N/A
Resolving configuration conflicts

Configuration conflicts are resolved by selecting the configuration that is closest to where you are interacting with Kiro CLI. For example, if you have a MCP configuration in both global and project mcp.json files, when you are chatting with Kiro in the project folder, the MCP configuration from the project folder will be applied.

Since you can also define a custom agents at a global and project scope, if there is a conflict between at the same level with the agent configuration, then Kiro CLI will choose the configuration from the agent.

Here's the priority order of how configuration is rationalized:

Configuration	Priority
MCP servers	Agent > Project > Global
Prompts	Project > Global
Custom agents	Project > Global
Steering	Project > Global
Since MCP servers can be configured in three scopes and there is a includeMcpJson setting in an agent configuration, MCP servers are handled slightly differently. Refer MCP server loading priority
