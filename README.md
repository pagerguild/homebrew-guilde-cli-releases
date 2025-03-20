# guilde-cli-releases
A Private Repository For Releases of Guilde's CLI/MCP tool

## Early Access Installation
### Supported Platforms
- macOS Apple Silicon: Supported
- macOS Intel: Alpha
- Linux x86: Alpha
- Linux ARM: Alpha
- Windows: Coming soon

### Early Access
1. For each developer who will be using the Guilde CLI/MCP, send Guilde:
    - GitHub username
    - Slack handle
2. The Guilde team will provide you with a *one-use login code*.
3. Run `brew tap pagerguild/guilde-cli-releases`
4. Run `brew install guilde-cli`
5. Run `guilde` to verify that it is in your path - if you have not installed Guilde previously, you will see a message saying so.  If you have Guilde installed already, it may just wait for input.  In that case, hit `ctrl-c` to cancel.
6. Run `guilde login [one-use login code]`
    - Guilde will automatically create a config file located at `~/.config/guilde/guilde-cli.json` that will contain an API key and a pointer to your default Github repo (the repo that you currently ask questions about in Slack with Guilde already).
    - If this is your first time installing Guilde, Guilde will also automatically create or update the Cursor MCP config file at `~/.cursor/mcp.json` pointing to your `guilde` installation (typically `/opt/homebrew/bin/guilde`).
7. In Cursor, click the top right gear in the UI. This will open Cursor Settings.  In that window, choose MCP.  You should see Guilde in your list of MCP servers.  Click the reload icon next to it to be sure.  In the Tools list for Guilde, you should see `schema_object_detail` among others.

### Updating

brew update 
brew upgrade guilde-cli


### Troubleshooting

If you run into trouble please do reach out to the Guilde team.

In the meantime, you can check to make sure your config files are in the proper format.

~/.config/guilde/guilde-cli.json
```
{
  "apiKey": "[Guilde API Key, created when guilde login is run]",
  "repository": "[Your Github org name]/[Your Github repo name]",
  "logPath": "-",
  "apiHost": "https://api.guilde.ai"
}
```
~/.cursor/mcp.json
```
{
  "mcpServers": {
    "Guilde MCP": {
      "command": "/opt/homebrew/bin/guilde",
      "args": [
        "mcp"
      ]
    }
  }
}
```

