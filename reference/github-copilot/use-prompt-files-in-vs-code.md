Use prompt files in VS Code
Prompt files, also known as slash commands, let you simplify prompting for common tasks by encoding them as standalone Markdown files that you can invoke directly in chat. Each prompt file includes task-specific context and guidelines about how the task should be performed.

Unlike custom instructions that are applied automatically, you invoke prompt files manually in chat.

Use prompt files to:

Simplify prompting for common tasks, such as scaffolding a new component, running and fixing tests, or preparing a pull request
Override default behavior of a custom agent, such as creating a minimal implementation plan or generating mockups for API calls
Prompt file locations
You can define prompt files for a specific workspace or at the user level, where they are available across all your workspaces.

Expand table
Scope	Default file location
Workspace	.github/prompts folder
User profile	prompts folder of the current VS Code profile
You can configure additional file locations for workspace prompt files with the 
chat.promptFilesLocations

 setting.
Prompt file format
Prompt files are Markdown files with the .prompt.md extension. The optional YAML frontmatter header configures the prompt's behavior:

Expand table
Field	Required	Description
description	No	A short description of the prompt.
name	No	The name of the prompt, used after typing / in chat. If not specified, the file name is used.
argument-hint	No	Hint text shown in the chat input field to guide users on how to interact with the prompt.
agent	No	The agent used for running the prompt: ask, agent, plan, or the name of a custom agent. By default, the current agent is used. If tools are specified, the default agent is agent.
model	No	The language model used when running the prompt. If not specified, the currently selected model in model picker is used.
tools	No	A list of tool or tool set names that are available for this prompt. Can include built-in tools, tool sets, MCP tools, or tools contributed by extensions. To include all tools of an MCP server, use the <server name>/* format.
Learn more about tools in chat.
Note
If a given tool is not available when running the prompt, it is ignored.

The body contains the prompt text in Markdown format. Provide specific instructions, guidelines, or any other relevant information that you want the AI to follow.

You can reference other workspace files by using Markdown links. Use relative paths to reference these files, and ensure that the paths are correct based on the location of the prompt file.

To reference agent tools in the body text, use the #tool:<tool-name> syntax. For example, to reference the githubRepo tool, use #tool:githubRepo.

Within a prompt file, you can reference variables by using the ${variableName} syntax. You can reference the following variables:

Workspace variables - ${workspaceFolder}, ${workspaceFolderBasename}
Selection variables - ${selection}, ${selectedText}
File context variables - ${file}, ${fileBasename}, ${fileDirname}, ${fileBasenameNoExtension}
Input variables - ${input:variableName}, ${input:variableName:placeholder} (pass values to the prompt from the chat input field)
The following examples demonstrate how to use prompt files. For more community-contributed examples, see the Awesome Copilot repository.

Example: generate a React form component
Example: using variables
Example: perform a security review of a REST API
Create a prompt file
When you create a prompt file, choose whether to store it in your workspace or user profile. Workspace prompt files apply only to that workspace, while user prompt files are available across multiple workspaces.

To create a prompt file:

Tip
Type /prompts in the chat input to quickly open the Configure Prompt Files menu.

In the Chat view, select Configure Chat (gear icon) > Prompt Files, and then select New prompt file.

Screenshot showing the Chat view, and Configure Chat menu, highlighting the Configure Chat button.

Alternatively, use the Chat: New Prompt File or Chat: New Untitled Prompt File command from the Command Palette (⇧⌘P).

Choose the scope of the prompt file:

Workspace: creates the prompt file in the .github/prompts folder of your workspace to only use it within that workspace. Add more prompt folders for your workspace with the 
chat.promptFilesLocations

 setting.
User profile: creates the prompt file in the current profile folder to use it across all your workspaces.
Enter a file name for your prompt file. This is the default name that appears when you type / in chat.

Author the chat prompt by using Markdown formatting.

Fill in the YAML frontmatter at the top of the file to configure the prompt's description, agent, tools, and other settings.
Add instructions for the prompt in the body of the file.
To modify an existing prompt file, in the Chat view, select Configure Chat > Prompt Files, and then select a prompt file from the list. Alternatively, use the Chat: Configure Prompt Files command from the Command Palette (⇧⌘P) and select the prompt file from the Quick Pick.

Use a prompt file in chat
You have multiple options to run a prompt file:

In the Chat view, type / followed by the prompt name in the chat input field. Agent skills also appear as slash commands alongside prompt files.
You can add extra information in the chat input field. For example, /create-react-form formName=MyForm or /create-api for listing customers.

Run the Chat: Run Prompt command from the Command Palette (⇧⌘P) and select a prompt file from the Quick Pick.
Open the prompt file in the editor, and press the play button in the editor title area. You can choose to run the prompt in the current chat session or open a new chat session.
This option is useful for quickly testing and iterating on your prompt files.

Tip
Use the 
chat.promptFilesRecommendations

 setting to show prompts as recommended actions when starting a new chat session.
Screenshot showing a "explain" prompt file recommendation in the Chat view.

Tool list priority
You can specify the list of available tools for both a custom agent and prompt file by using the tools metadata field. Prompt files can also reference a custom agent by using the agent metadata field.

The list of available tools in chat is determined by the following priority order:

Tools specified in the prompt file (if any)
Tools from the referenced custom agent in the prompt file (if any)
Default tools for the selected agent (if any)
Sync user prompt files across devices
VS Code can sync your user prompt files across multiple devices by using Settings Sync.

To sync your user prompt files, enable Settings Sync for prompt and instruction files:

Make sure you have Settings Sync enabled.

Run Settings Sync: Configure from the Command Palette (⇧⌘P).

Select Prompts and Instructions from the list of settings to sync.

Tips for writing effective prompts
Clearly describe what the prompt should accomplish and what output format is expected.
Provide examples of the expected input and output to guide the AI's responses.
Use Markdown links to reference custom instructions rather than duplicating guidelines in each prompt.
Take advantage of built-in variables like ${selection} and input variables to make prompts more flexible.
Use the editor play button to test your prompts and refine them based on the results.
Frequently asked questions
How do I know where a prompt file comes from?
Prompt files can come from different sources: built-in, user-defined in your profile, workspace-defined prompts in your current workspace, or extension-contributed prompts.

To identify the source of a prompt file:

Select Chat: Configure Prompt Files from the Command Palette (⇧⌘P).
Hover over the prompt file in the list. The source location is displayed in a tooltip.
Tip
Use the chat customization diagnostics view to see all loaded prompt files and any errors. Right-click in the Chat view and select Diagnostics. Learn more about troubleshooting AI in VS Code.

Related resources
Create custom instructions
Configure tools in chat
Community contributed instructions, prompts, and custom agents