## Customisation

You can customise the options available by creating a `.meteor.json` in the root of your repository, or in your home directory. The repository-level config will be preferred if it exists.

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
    { "title":  "feat", "description":  "a new feature"},
    { "title":  "fix", "description":  "a bug fix"},
    { "title":  "bug", "description":  "introducing a bug"}
  ],
}
```

If you use boards (Jira etc) but need a way to have commits without one, add the following to the `boards` array:
```json
{
  "boards": [
    { "name": "COMP" },
    { "name": "NONE" }
  ]
}
```
