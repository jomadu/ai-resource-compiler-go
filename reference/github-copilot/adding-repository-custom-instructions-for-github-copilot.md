Adding repository custom instructions for GitHub Copilot
Create repository custom instructions files that give Copilot additional context on how to understand your project and how to build, test and validate its changes.
Tool navigation
Visual Studio Code
JetBrains IDEs
Visual Studio
Web browser
Eclipse
Xcode
In this article
This version of this article is for using repository custom instructions and prompt files in VS Code. Click the tabs above for instructions on using custom instructions in other environments.

Introduction

Repository custom instructions let you provide Copilot with repository-specific guidance and preferences. For more information, see About customizing GitHub Copilot responses.

Prerequisites for repository custom instructions

You must have a custom instructions file (see the instructions below).

Custom instructions must be enabled. This feature is enabled by default. See Enabling or disabling repository custom instructions later in this article.

Creating custom instructions

VS Code supports three types of repository custom instructions. For details of which GitHub Copilot features support these types of instructions, see About customizing GitHub Copilot responses.

Repository-wide custom instructions, which apply to all requests made in the context of a repository.

These are specified in a copilot-instructions.md file in the .github directory of the repository. See Creating repository-wide custom instructions.

Path-specific custom instructions, which apply to requests made in the context of files that match a specified path.

These are specified in one or more NAME.instructions.md files within or below the .github/instructions directory in the repository. See Creating path-specific custom instructions.

If the path you specify matches a file that Copilot is working on, and a repository-wide custom instructions file also exists, then the instructions from both files are used.

Agent instructions are used by AI agents.

You can create one or more AGENTS.md files, stored anywhere within the repository. When Copilot is working, the nearest AGENTS.md file in the directory tree will take precedence. For more information, see the openai/agents.md repository.

Note

Support of AGENTS.md files outside of the workspace root is currently turned off by default. For details of how to enable this feature, see Use custom instructions in VS Code in the VS Code documentation.
Creating repository-wide custom instructions

In the root of your repository, create a file named .github/copilot-instructions.md.

Create the .github directory if it does not already exist.

Add natural language instructions to the file, in Markdown format.

Whitespace between instructions is ignored, so the instructions can be written as a single paragraph, each on a new line, or separated by blank lines for legibility.

Creating path-specific custom instructions

Create the .github/instructions directory if it does not already exist.

Optionally, create subdirectories of .github/instructions to organize your instruction files.

Create one or more NAME.instructions.md files, where NAME indicates the purpose of the instructions. The file name must end with .instructions.md.

At the start of the file, create a frontmatter block containing the applyTo keyword. Use glob syntax to specify what files or directories the instructions apply to.

For example:

---
applyTo: "app/models/**/*.rb"
---
You can specify multiple patterns by separating them with commas. For example, to apply the instructions to all TypeScript files in the repository, you could use the following frontmatter block:

---
applyTo: "**/*.ts,**/*.tsx"
---
Glob examples:

* - will all match all files in the current directory.
** or **/* - will all match all files in all directories.
*.py - will match all .py files in the current directory.
**/*.py - will recursively match all .py files in all directories.
src/*.py - will match all .py files in the src directory. For example, src/foo.py and src/bar.py but not src/foo/bar.py.
src/**/*.py - will recursively match all .py files in the src directory. For example, src/foo.py, src/foo/bar.py, and src/foo/bar/baz.py.
**/subdir/**/*.py - will recursively match all .py files in any subdir directory at any depth. For example, subdir/foo.py, subdir/nested/bar.py, parent/subdir/baz.py, and deep/parent/subdir/nested/qux.py, but not foo.py at a path that does not contain a subdir directory.
Optionally, to prevent the file from being used by either Copilot coding agent or Copilot code review, add the excludeAgent keyword to the frontmatter block. Use either "code-review" or "coding-agent".

For example, the following file will only be read by Copilot coding agent.

---
applyTo: "**"
excludeAgent: "code-review"
---
If the excludeAgent keyword is not included in the front matterblock, both Copilot code review and Copilot coding agent will use your instructions.

Add your custom instructions in natural language, using Markdown format. Whitespace between instructions is ignored, so the instructions can be written as a single paragraph, each on a new line, or separated by blank lines for legibility.

Did you successfully add a custom instructions file to your repository?

 
Custom instructions in use

The instructions in the file(s) are available for use by Copilot as soon as you save the file(s). Instructions are automatically added to requests that you submit to Copilot.

Custom instructions are not visible in the Chat view or inline chat, but you can verify that they are being used by Copilot by looking at the References list of a response in the Chat view. If custom instructions were added to the prompt that was sent to the model, the .github/copilot-instructions.md file is listed as a reference. You can click the reference to open the file.

Screenshot of an expanded References list, showing the 'copilot-instructions.md' file highlighted with a dark orange outline.

Enabling or disabling repository custom instructions

You can choose whether or not you want Copilot to use repository-based custom instructions.

Enabling or disabling custom instructions for Copilot Chat

Custom instructions are enabled for Copilot Chat by default but you can disable, or re-enable, them at any time. This applies to your own use of Copilot Chat and does not affect other users.

Open the Setting editor by using the keyboard shortcut Command+, (Mac) / Ctrl+, (Linux/Windows).
Type instruction file in the search box.
Select or clear the checkbox under Code Generation: Use Instruction Files.
Enabling or disabling custom instructions for Copilot code review

Custom instructions are enabled for Copilot code review by default but you can disable, or re-enable, them in the repository settings on GitHub.com. This applies to Copilot's use of custom instructions for all code reviews it performs in this repository.

On GitHub, navigate to the main page of the repository.

Under your repository name, click  Settings. If you cannot see the "Settings" tab, select the  dropdown menu, then click Settings.

Screenshot of a repository header showing the tabs. The "Settings" tab is highlighted by a dark orange outline.
In the "Code & automation" section of the sidebar, click  Copilot, then Code review.

Toggle the “Use custom instructions when reviewing pull requests” option on or off.

Enabling and using prompt files

Note

Copilot prompt files are in public preview and subject to change. Prompt files are only available in VS Code, Visual Studio, and JetBrains IDEs. See About customizing GitHub Copilot responses.
For community-contributed examples of prompt files for specific languages and scenarios, see the Awesome GitHub Copilot Customizations repository.
Prompt files let you build and share reusable prompt instructions with additional context. A prompt file is a Markdown file, stored in your workspace, that mimics the existing format of writing prompts in Copilot Chat (for example, Rewrite #file:x.ts). You can have multiple prompt files in your workspace, each of which defines a prompt for a different purpose.

Enabling prompt files

To enable prompt files, configure the workspace settings.

Open the command palette by pressing Ctrl+Shift+P (Windows/Linux) / Command+Shift+P (Mac).
Type "Open Workspace Settings (JSON)" and select the option that's displayed.
In the settings.json file, add "chat.promptFiles": true to enable the .github/prompts folder as the location for prompt files. This folder will be created if it does not already exist.
Creating prompt files

Open the command palette by pressing Ctrl+Shift+P (Windows/Linux) / Command+Shift+P (Mac).

Type "prompt" and select Chat: Create Prompt.

Enter a name for the prompt file, excluding the .prompt.md file name extension. The name can contain alphanumeric characters and spaces and should describe the purpose of the prompt information the file will contain.

Write the prompt instructions, using Markdown formatting.

You can reference other files in the workspace by using Markdown links—for example, [index](../../web/index.ts)—or by using the #file:../../web/index.ts syntax. Paths are relative to the prompt file. Referencing other files allows you to provide additional context, such as API specifications or product documentation.

Using prompt files

At the bottom of the Copilot Chat view, click the Attach context icon ().

In the dropdown menu, click Prompt... and choose the prompt file you want to use.

Optionally, attach additional files, including prompt files, to provide more context.

Optionally, type additional information in the chat prompt box.

Whether you need to do this or not depends on the contents of the prompt you are using.

Submit the chat prompt.

For more information about prompt files, see Use prompt files in Visual Studio Code in the Visual Studio Code documentation.

Further reading

Support for different types of custom instructions
Customization library—a curated collection of examples
Using custom instructions to unlock the power of Copilot code review