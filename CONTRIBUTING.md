* [Git branching workflow](#branching)
* [Commit conventions](#commit)
* [Coding conventions](#code)

## <a name="branching"></a>Branching

### Naming

All branches shall have one of these names.

* `master`: the default branch. This is a clean branch where reviewed, approved and CI passed pull requests are merged into. Merging to this branch is restricted to project maintainers
* `fix/#-short-name` or `fix/short-name`: refers to a fix, preferably with issue number. The short name describes the bug or issue
* `feature/#-short-name` or `feature/short-name`: (main) feature branch, preferably with issue number. The short name describes the feature
* `feature/#-short-name-part`: a sub scope of the feature in a separate branch, that is intended to merge into the main feature branch before the main feature branch is merged into `master`

### Scope

A fix, feature or issue branch should be **small and focused** and should be scoped to a **single specific task**. Do not combine new features and refactoring of existing code.

### Pull requests and rebasing

* **Before** a reviewer is assigned, rebasing the branch to reduce the number of commits is highly advised. Self-review your own pull request: making the [commit](#commit) history clean, check for typos or incoherences, and make sure Continuous Integration passes.

Keep the commits to be merged clean: adhere to the commit message format defined below and instead of adding and deleting files within a pull request, drop or fix the concerning commit that added the file.

Interactive rebase (`git rebase -i`) can be used to rewrite commit messages that do not follow these contribution guidelines.

## <a name="commit"></a>Commit Messages

The first line of a commit message is the subject. The commit message may contain a body, separated from the subject by an empty line.

### Subject

The subject contains the concerning component or topic and a concise message in [the imperative mood](https://chris.beams.io/posts/git-commit/#imperative), starting with a capital. The subject may also contain references to issues or other resources.

The component or topic is typically a few characters long and should always be present. Component names are e.g.:

* `srr`: Soccer robot remote binary
* `util`: utilities
* `ci`: Continuous Integration instructions, e.g. Travis file
* `doc`: documentation
* `dev`: other non-functional development changes, e.g. Makefile, .gitignore, editor config
* `*`: changes affecting all code, e.g. primitive types

Changes that affect multiple components can be comma separated.

Good commit messages:

* `srr: Support -42 flag`
* `make: Set version from git tag, close #123`
* `srr,util: Fix WebSocket authentication`

Make sure that commits are scoped to something meaningful and could, potentially, be merged individually.

### Body

The body may contain a more detailed description of the commit, explaining what it changes and why. The "how" is less relevant, as this should be obvious from the diff.

## <a name="code"></a>Code

### Formatting

We want our code to be consistent, so we'll have to agree on a number of formatting rules. These rules should usually usually be applied by your editor. Make sure to install the [editorconfig](https://editorconfig.org) plugin for your editor.

Go code can be automatically formatted using the [`gofmt`](https://golang.org/cmd/gofmt/) tool. There are many editor plugins that call `gofmt` when you save your files.

#### General

Use **utf-8**, **LF** line endings, a **final newline** and **trim whitespace** from the end of the line (except in Markdown).

#### Tabs vs Spaces

Many developers have strong opinions about using tabs vs spaces. We apply the following rules:

* All `.go` files are indented using **tabs**
* The `Makefile` and all `.make` files are indented using **tabs**
* All other files are indented using **two spaces**

#### Line length

* If a line is longer than 80 columns, try to find a "natural" break
* If a line is longer than 120 columns, insert a line break
* In very special cases, longer lines are tolerated

### Linting

Use [`golint`](github.com/golang/lint/golint) to lint `.go` files and [`eslint`](https://eslint.org) to lint `.js` files. These tools should automatically be installed when initializing your development environment.

### API methods naming

All API method names should follow the naming convention of `VerbNoun` in upper camel case, where the verb uses the imperative mood and the noun is the resource type.

The following snippet defines the basic CRUD definitions for a resource named `Type`.
Note also that the order of the methods is defined by CRUD.

```
CreateType
GetType
ListTypes (returns slice)
UpdateType
DeleteType

AddTypeAttribute
SetTypeAttribute
GetTypeAttribute
ListTypeAttributes (returns slice)
RemoveTypeAttribute
```

### Variable naming

Variable names should be short and concise.

Follow the [official go guidelines](https://github.com/golang/go/wiki/CodeReviewComments#variable-names) and try to be consistent with Go standard library as much as possible, everything not defined in the tables below should follow Go standard library naming scheme. In general, variable names are English and descriptive, as well as putting adjectives and adverbs before the noun and verb respectively.

#### Single-word entities

|    entity     |  name  |  example type   |
| :-----------: | :----: | :-------------: |
|    context    |  ctx   | context.Context |
|     mutex     |   mu   |   sync.Mutex    |
| configuration |  conf  |                 |
|    logger     | logger |   log.Logger    |
|    message    |  msg   |                 |
|    status     |   st   |                 |
|    server     |  srv   |                 |
|      ID       |   id   |     string      |
|    counter    |  cnt   |       int       |

#### 2-word entities

In case both of the words have an implementation-specific meaning, the variable name is the combination of first letter of each word.

|   entity   | name |
| :--------: | :--: |
| wait group |  wg  |

In case one of the words specifies the meaning of the variable in a specific language construct context, the variable name is the combination of abbrevations of the words.

### Comments

Code should be as self-explanatory as possible. However, comments should be used to respect Go formatting guidelines and to explain what can not be expressed by pure code. Comments should be English sentences, and documentation-generating comments should be closed by a period. Comments can also be used to indicate steps to take in the future (_TODOs_).

* In **Go files**, comments should be added according to `golint` requirements and [Effective Go guidelines](https://golang.org/doc/effective_go.html#commentary), especially in regards to commenting exported packages, types and variables.
