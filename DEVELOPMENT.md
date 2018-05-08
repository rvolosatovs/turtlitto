# Soccer Robot Remote Development

## Development Environment

The development environment heavily relies on [`make`](https://www.gnu.org/software/make/). Under the hood, `make` calls other tools such as `go`, `yarn` etc. Let's first make sure you have `go`, `node` and `yarn`:

### MacOS

Using [Homebrew](https://brew.sh):

```sh
brew install go node yarn
```

### Linux

On Ubuntu (or Ubuntu [using the Windows 10 Subsystem for Linux](https://www.microsoft.com/nl-NL/store/p/ubuntu/9nblggh4msv6?rtc=1)):

```sh
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list

curl -sS https://deb.nodesource.com/gpgkey/nodesource.gpg.key | sudo apt-key add -
echo "deb https://deb.nodesource.com/node_8.x xenial main" | sudo tee /etc/apt/sources.list.d/nodesource.list
echo "deb-src https://deb.nodesource.com/node_8.x xenial main" | sudo tee -a /etc/apt/sources.list.d/nodesource.list

sudo apt-get update
sudo apt-get install build-essential nodejs yarn

curl -sS https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz -o go1.10.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.1.linux-amd64.tar.gz
sudo ln -s /usr/local/go/bin/* /usr/local/bin
```

### Windows

* Install Go from [official website](https://golang.org/dl/).
* Install Make from e.g. [here](https://sourceforge.net/projects/gnuwin32/files/make/3.81/make-3.81.exe/download?use_mirror=datapacket&download=).
* Install Node.js from [official website](https://nodejs.org/en/download/current/).
* Install Yarn from [official website](https://yarnpkg.com/lang/en/docs/install/#windows-stable).

### Getting started with Go Development

_Note, that the commands should be executed in a **bash** shell(it is installed by default with git on Windows)_

We will first need a Go workspace. The Go workspace is a folder that contains the following sub-folders:

* `src` which contains all source files
* `pkg` which contains compiled package objects
* `bin` which contains binary executables

From now on this folder is referred to as `$GOPATH`. By default, Go assumes that it's in `$HOME/go`.
Execute this to explicitly setup `$GOPATH` and add `$GOPATH/bin` to your `$PATH`.

```sh
printf 'export GOPATH="$(go env GOPATH)"\nexport PATH="$PATH:$GOPATH/bin"' >> ~/.profile
source ~/.profile
```

Now that your Go development environment is ready, it strongly recommended to get familiar with Go by following the [Tour of Go](https://tour.golang.org/).

### Getting started with development

_Note the `--recursive` flag!_

```sh
git clone --recursive git@github.com:rvolosatovs/turtlitto.git $GOPATH/src/github.com/rvolosatovs/turtlitto
```

All development is done in this directory.

```sh
cd $GOPATH/src/github.com/rvolosatovs/turtlitto
```

If you run Windows, execute the following command as well to account for the difference in the newline policy for Windows compared to Unix based systems:

```sh
git config --global core.autocrlf true
```

Ensure the dependencies of the project are installed:

```sh
make deps
```

#### Folder Structure

```
.
├── STYLE.md     guidelines for contributing: branching, commits, code style, etc.
├── DEVELOPMENT.md      guide for setting up your development environment
├── Gopkg.lock          dependency lock file managed by golang/dep
├── Gopkg.toml          dependency file managed by golang/dep
├── Makefile            dev/test/build tooling
├── README.md           general information about this project
│   ...
├── cmd                 contains the different binaries
│   └── soccer-robot-remote          contains the Soccer Robot Remote
├── docs                contains the documentation
├── front               contains the frontend of the project
├── pkg                 contains all libraries used in the backend
├── release             binaries will be compiled to this folder - not added to git
└── vendor              dependencies managed by golang/dep - not added to git
```

#### Testing

For backend:

```sh
make go.test
```

For frontend:

```sh
make js.test
```

For testing everything:

```sh
make test
```

#### Building

There's one binary to be built: the `soccer-robot-remote-linux-amd64`, which holds the remote control for the soccer robots and the frontend of the application.

To build those run:

```sh
make
```

This will result in `release/soccer-robot-remote` and `release/front`generated.

To build a Docker container run:

```sh
make docker
```

You can later run the project using `docker-compose up` from the root of the project. It binds the web interface on `:4242` and assumes an active TRC socket at `.trc/trc.sock`.

#### Local development

The application consists of two modules, namely the go backend server and react application for the client side. There are several ways to run the application on your machine, but in order to make debugging easier, we will deploy them separately.

1.  Start the Go server using `go run cmd/soccer-robot-remote/main.go -socket <unix-socket>` or `docker-compose up`.
2.  Open another terminal and move to the project folder.
3.  Start webpack development server and host the React app: `yarn start`.
4.  Open `http://localhost:3000` in your browser(should happen automatically).

This approach makes development and debugging much easier. Webpack development server has the hot module replacement feature, that will automatically show all changes in the browser without the need to rebuild the project. Also, the source code wont be obfuscated and minified.

## Frontend guidelines

This section is intended to give an overview of chosen tooling, provide examples of using those and define coding standards.

### styled-components

Check [the official website](https://www.styled-components.com/).
If you have any experience in building websites, you might wonder why not styling the application just by writing css declarations in `*.css` files and import them into the React app. There are many reasons for that, but the main one - styling React applications with just `*.css` files is messy. There is no notion of a component in `css`, as it is in React (everything is a component), just a set of selectors.

`styled-components`, `sc` for short, allows creating styled React components right in `*.js` files, preferably next to your "real" React components. Consider this case:

```
 // style.css

 .button {
   ...
 }

 .button-primary {
   ...
 }

 .button-error {
   ...
 }

 // App.js

 import "./styles/style.css"

 export default () => {
   return (
     <div>
       <button className="button button-primary">Ok</button>
       <button className="button button-error">Cancel</button>
     </div>

   )
 }
```

vs

```
 // App.js

 import styled from "styled-components";

 const Button = styled.button`...`;

 const ButtonPrimary = Button.extend`...`;

 const ButtonError = Button.extend`...`;

 export default () => {
   return (
     <div>
       <ButtonPrimary />
       <ButtonError />
     </div>
   )
 }
```

As one can notice, in the second example all styles for a specific component are stored in the same file with the component itself. This approach makes it easy to change styles, avoid typo's in css class names and so on. Of course, we can store styles as javascript objects and then pass them to the components, but this is not much different from using an external stylesheet. We can benefit more from `sc`. `sc` returns real React components that can receive properties, can be reused and tested. It is worth noting, that you can use a regular css syntax in `sc`. Here is an example component from the demo:

```
const Message = styled.li`
 font-size: 1.5rem;
 padding: 0.5rem 0;
`;
```

#### Usage

Do not use element selectors in `sc`, as well as ids.

```
  // both are super bad
  const ProductList = styled.ul`
    ...
    > li {
      ...
    }
  `;

  const ProductCard = styled.div`
    ...
    > #product-price {
      ...
    }
  `;
```

Create a separate component for list items.

```
  // good
  const ProductList = styled.ul`...`;
  const Product = styled.li`...`;
```

Note, that you can also nest components in `sc`.

```
  const ProductList = styled.ul`
    ...
    ${Product} {
      ...
    }
  `
```

Use pseudo-classes the same way.

```
  const ProductList = styled.ul`
    ...
    &:hover ${Product} {
      ...
    }
  `
```

If there are too many components used for styling, feel free to move it to a separate `js` file within the same folder.

```
  // styles.js

  export const ButtonPrimary = styled.button`...`;
  ...
  export const ButtonError = styled.button`...`;

  // App.js
  import { ButtonPrimary, ButtonError } from "./styled";

  // use buttons
```

You are not limited to using only React components, however you can use regular html elements for small layout changes.

```
const Description = styled.p`...`;

export default () => {
  return (
    <Description>description here and <strong>some text in bold</strong>.</Description>
  )
}
```

Do not enclose your `sc` components within the components class declaration!

### React.js

Please take a look on the following topics in the official documentation:

* [Controlled components](https://reactjs.org/docs/forms.html#controlled-components)
* [Conditional Rendering](https://reactjs.org/docs/conditional-rendering.html)
* [Handling events](https://reactjs.org/docs/handling-events.html)

Note, that some of the code snippets use older versions of React and specify what `this` should point to via the `bind()` function - you do not need to do that. The version we use will always point to the enclosing class - the react component.

There are only 2 reasons to use React class component:

1.  The component will keep its internal `state`
2.  The component will need to hook up into components lifecycle methods, e.g. `componentDidMount()`

If none of those apply for the component you create - use a stateless functional component.

```
 // bad
 class FlashMessage extends Component {
   render(props) {
     return (
       <span>props.message</span>
     )
   }
 }

 // good
 export default props => <span>props.message</span>
```

Use `this.state = {}` in the body of the component to set initial state and use it **ONLY ONCE**. Update the state via `this.setState()`. If you set the state directly, the UI will not be updated by React.

Do not pass location specific props to components.

```
// dont do that
export default () => <TurtleCard isMobileVersion={true}>
export default () => <TurtleCard isMainSection={true}>
```

### Naming

* Use PascalCase for components, e.g. `NavigationBar`
* Use camelCase for functions and variables, e.g. `getAvailableTurtles()` or `availableTurtles`
* Use SCREAMING_SNAKE_CASE for constants, e.g. `API_ENDPOINT`
* Use kebab-case for the rest, e.g. images `flag-en.png`, but import with camelCase - `flagEn`

Do not specify the `.js` extension when importing javascript, webpack does it for you.

```
 // bad
 import Button from "./Button.js";

 // good
 import Button from "./Button";
```

Do not prefix/postfix the names of your components according with their underlying html element (there are some exceptions though, `Button`, `Form`, `List` - are ok). Always find a business entity that is represented by every component you use. If you cannot find such, most probably you do not need the component.

```
 // bad
 <FlexBox />
 <Div />
 <TurtleDiv />
 <Circle/>
 <Item />
 <ImgDiv />
 <Parent />
 <Child />

 // good
 <FlightSearchForm />
 <TurtleCard />
 <StartButton />
 <DownloadIcon />
 <PortConnectionInput />
 <FormWrapper />
 <RegistrationForm />
```

Do not worry about extra bytes, webpack will minify the resulting css and js files, so use full variable/function/class names.

```
 // bad! items? x? i?
 items.forEach((x, i) => {
   ...
 });

 // good
 turtles.forEach((turtle, index) => {
   ...
 });
```

### Javascript

Mutability is bad, especially when building UI. It hard to test and make any assumptions. Deal only with immutable data:

* Use pure functions, e.g. `map`, `filter`
* Forget about `for (let i = 0; i < arr.length; i++) {...}` use `forEach` or other array functions depending on what you want to achieve
* Use arrow functions
* NEVER use `var`
* think twice before using `let`, usually you can decompose your code to have `const` everywhere
* ALWAYS use `const` !!!

Read [var vs let](http://www.jstips.co/en/javascript/keyword-var-vs-let/).

Keep an eye where does `this` point to. If you are getting errors, use the arrow function and do not use `bind()`. See [js function vs arrow function](https://stackoverflow.com/questions/34361379/arrow-function-vs-function-declaration-expressions-are-they-equivalent-exch)

```
 import TextInput from "./TextInput";

 class MessageList extends Component {
   handleInputChange(value) {
     this.setState({value: value});
   }

   render() {
     return (
       <TextInput onChange={this.handleInputChange} /> // will produce an error
       <TextInput onChange={(value) => this.handleInputChange(value)} /> // is ok, this points to the MessageList class
     )
   }
 }
```

Despite the fact that webpack has the dead code elimination procedure, be explicit on the parts of code you import from external libraries.

```
// bad
import lodash from "lodash";

lodash.minBy(...);

// good
import minBy from "lodash/minBy";

minBy(...);
```

### General

Do not base the layout only on `div`'s (despite the fact it is possible), use elements that semantically feel right. If you are working on a new section - use `<section>`, if you need a list - use `<ul>`.

NEVER set `outline: none`, see [outlinenone](http://www.outlinenone.com/).

NEVER set explicit width/height to any elements except images. There are some exceptions for dirty hacks, e.g. to make something scrollable.

Build the layout as [mobile-first](https://getflywheel.com/layout/start-practicing-mobile-first-development/).
