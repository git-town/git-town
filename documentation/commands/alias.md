<h1 textrun="command-heading">Alias command</h1>

<blockquote textrun="command-summary">
Adds or removes default global aliases
</blockquote>

<a textrun="command-description">
Global aliases allow Git Town commands to be used like native Git commands.
When aliases are set, you can run "git hack" instead of having to run "git town hack".
Example: "git append" becomes equivalent to "git town append".

When adding aliases, no existing aliases will be overwritten.

Note that this can conflict with other tools that also define additional Git commands.
</a>

#### Usage

<pre textrun="command-usage">
git town alias (true | false)
</pre>
