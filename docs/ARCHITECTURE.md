# Git Town architecture

### Design goals

Complexity in the Git Town codebase arises from multiple conflicting design
goals:

1. Execute a highly variable set of Git operations depending on the current
   status of the repository. Git Town's business logic covers so many edge cases
   that most Git Town commands aren't just a simple scripts, they are complex
   programs.
2. When a step in these programs fails, terminate to allow the end user to
   resolve problems in the same terminal window and shell environment that they
   ran Git Town in and then resume execution.
3. Be able to reliably undo everything that Git Town has done.

### General structure

To keep the amount of code manageable, the Git Town codebase separates
functionality into orthogonal, composable subsystems. Subsystems exist for
parsing configuration data, syncing branches, calculating undo operations,
interacting with the CLI, interacting with external hosting services, etc.

Each subsystem defines its own domain concepts, types, business logic, and
helper functions. To prevent cyclic package dependencies, subsystems define
concepts and data types in dedicated `*domain` packages.

Higher-level subsystems like the business logic to sync branches use lower-level
subsystems for executing Git and access configuration. Low-level subsystems
don't have access to high-level subsystems.

### Execution framework

Git Town addresses requirements 1 and 2 via an
[interpreter](https://en.wikipedia.org/wiki/Interpreter_(computing)) that
executes Git-Town-specific programs consisting of using Git-related opcodes.
Each Git Town command:

- Inspects the state of the Git repo.
- Determines the Git operations that Git Town needs to perform and stores them
  as a Git Town program. This program consists of opcodes that the Git Town
  interpreter can execute. Most of these opcodes execute Git commands.
- Starts the Git Town interpreter engine to execute this program.

If there are issues that require the user to resolve in a terminal window, the
interpreter:

- Persists the current interpreter state (called "runstate") to disk.
- Exits the running Git Town process to give the user access to the shell to
  resolve the encountered problems.
- Prints a human-friendly explanation of the problem and what the user needs to
  do.

After resolving the problems and restarting Git Town, the interpreter recognizes
and loads the persisted runstate from disk and resumes executing it.

### Undo framework

To undo a previously run Git Town command (requirement 3), Git Town:

- compares snapshots of the affected Git repository before and after the command
  ran
- determines the changes that the Git Town command made to the repo
- creates a program that reverses these changes
- starts the interpreter to execute this program

### Code style

The Git Town codebase deviates in some areas from the "official" Go coding
style. The decisions to make these deviations wasn't easy but necessary after
trying the regular Go style for years. Here is some background what is different
and why.

#### Favor descriptive naming over brevity

There is a somewhat widespread convention in the Go community to use extremely
abbreviated names for variables, types, and functions. While we are not against
this, having short names has the lowest priority of all code quality best
practices in our book. Making code self-describing is way more important than
brevity. The Git Town codebase therefore uses self-describing varible names
whenever it helps. As an open-source codebase, we don't assume any familiarity
of any reader with the codebase. See
https://michaelwhatcott.com/familiarity-admits-brevity for more background.

#### Use `self` for method receivers

The
[Go code review comments wiki page](https://go.dev/wiki/CodeReviewComments#receiver-names)
recommends avoiding generic names like `this` or `self` for method receivers and
instead use short one or two letter names. After doing this for many years on
the Git Town codebase we find this approach unhelpful and unwieldy in practice.

The Go review comments wiki page is incorrect that the method receiver is just
another function argument. The method receiver gets defined separate from the
other arguments. It is the central data element in the method, otherwise that
method shouldn't be a method but a function. The method receiver is the only
argument of which one can safely access private properties without violating
abstraction and encapsulation boundaries.

Go doesn't provide a clear convention for exactly how to name the method
receiver. A number of alternatives exist, each with their distinct pros and
cons, and none working well in all situations. This leads to time and energy
wasted figuring out the right method receiver name and justifying it in code
reviews. The only option that works in all cases without any bikeshedding is
`self`.

The Go recommendation leads to excessive churn. Renaming a type now also
requires renaming the receiver in all its methods. This leads to shotgun changes
on dozens more lines for simple rename refactors. This isn't tool supported
because of the lack of a convention, so has to be done manually and reviewed
with some level of care.

This makes refactoring unnecessarily costly, noisy, and thereby sometimes not
worth the effort. That's a bad outcome in which everybody loses. The ability to
refactor trumps adherence to debatable community standards.

`self` is pretty short, it's only four characters.

The names of the other function arguments also shouldn't be abbreviated. All
identifiers need to be descriptive.

https://michaelwhatcott.com/receiver-names-in-go and
https://dev.to/codypotter/the-case-for-self-receivers-in-go-3h7f provide more
background.

#### Dedicated generic types for Optionality and Mutablitity

Go gives pointers several orthogonal meanings. In Go, pointers can be used to
express optionality. The easiest way to have a variable that can sometimes have
a value and sometimes not is making it a pointer. In this case, `nil` means
there is no value and not-nil means there is a value. The problem with this is
that Go doesn't help with checking for absent values in any way. This leads to
runtime panics when trying to use a variable that contains nothing.

The Git Town codebase therefore wraps optional values inside a generic `Option`
type. This makes clear that a type is optional, and enforces an optionality
check.

Another function of pointers in Go is a performance optimizations: If a variable
is too large to pass by value, one can pass it by reference. The Git Town
codebase doesn't use this performance optimization because it isn't needed.

The final function of pointers in Go is to express mutability. If you want to
mutate variables provided as function arguments, you must provide them as a
pointer. The problem with this approach is that it's not obvious why a function
argument was provided as a pointer. Is it optional? Is it mutable? Is it merely
too heavy to pass by value? The Git Town codebase uses the generic `Mutable`
type to express whether a variable is mutable or not. Any struct field or
function argument that isn't wrapped in a `Mutable` should be considered
immutable.

#### One concept per file

Go recommends a programming style where each Go file contains many different
concepts (type definitions, functions, constants). In contrast, in the Git Town
codebase each concept is located in its own file. This allows finding concepts
by filename.
