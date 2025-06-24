<div align="center">
  <h1>️☄ meteor ☄</h1>
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/stefanlogue/meteor">
  <h5>Meteor is a simple, highly customisable CLI tool that helps you to write <a href="https://www.conventionalcommits.org/">conventional commits</a> with git.</h5>
</div>

You can call `meteor` where you'd normally type `git commit`. All flags
supported in `git commit` will still work.

![Demo](demos/demo-without-boards.gif)

## Installation

### Homebrew

```console
brew tap stefanlogue/tools
brew install --cask meteor
```

> [!IMPORTANT]
> If you previously installed `meteor` from the formula, you'll need to
> uninstall it before installing it from the cask

### Go

Install with Go (1.21+):

```console
go install github.com/stefanlogue/meteor@latest
```

Or grab a binary from [the latest release](https://github.com/stefanlogue/meteor/releases/latest).

## Customisation

You can customise the options available by creating a `.meteor.json` file
anywhere in the directory tree (at or above the current working directory). The
config file closest to the current working directory will be preferred. This
enables you to have different configs for different parent directories, such as
one for your personal work, one for your actual work, one for open source work
etc. For global configurations you can create a `config.json` file in the
`~/.config/meteor/` directory.

The content should be in the following format:

```json
{
  "boards": [
    { "name": "COMP" },
    { "name": "PERS" }
  ],
  "coauthors": [
    { "name": "John Doe", "email": "john.doe@email.com" }
  ],
  "prefixes": [
    { "type":  "feat", "description":  "a new feature"},
    { "type":  "fix", "description":  "a bug fix"},
    { "type":  "bug", "description":  "introducing a bug"}
  ],
  "commitTitleCharLimit": 60,
  "commitBodyCharLimit": 60,
  "commitBodyLineLength": 40
}
```

### Boards

![Demo with boards](demos/demo-with-boards.gif)

If you use boards (Jira etc) but need a way to have commits without one, add the
following to the `boards` array:

```json
{
  "boards": [
    { "name": "COMP" },
    { "name": "NONE" }
  ]
}
```

If you want to define a set of predefined scopes to select from rather than
typing the scope, a `scopes` array can be added to your config:

> [!WARNING]
> Setting predefined scopes removes the ability to type the scope

```json
{
  "scopes": [
    { "name": "scope1" },
    { "name": "scope2"}
  ]
}
```

### Line wrapping

To enforce line wrapping on the commit body, set the `commitBodyLineLength`
config option to any integer greater than or equal to 20.

### Message Templates

If the default commit message templates aren't exactly what you're looking for,
you can provide your own! The syntax can be seen in the defaults below:

```json
{
  "messageTemplate": "@type(@scope): @message",
  "messageWithTicketTemplate": "@ticket(@scope): <@type> @message"
}
```

`messageTemplate` needs to have:

- `@type`: the conventional commit type i.e. `feat`, `chore` etc.
- `@message`: the commit message
- `(@scope)`: (optional but recommended) the scope of the commit, must be within
parentheses

`messageWithTicketTemplate` also additionally takes `@ticket`

### Intro

If you want to skip the intro screen to save a keypress, add the following to
your config:

```json
{
  "showIntro": false
}
```
