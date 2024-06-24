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

The Go community often uses highly abbreviated names for variables, types, and
functions, following the personal preference of some of Go's creators. While
brevity can be useful, our primary focus for code quality in this codebase is
clarity and ease of understanding, i.e. self-describing code. Our open-source
tool has a wide user base and a small group of maintainers, with many
contributors adding just a single feature. To ensure our code is accessible to
everyone, we consistently use descriptive identifiers. For more context, please
refer to [this article](https://michaelwhatcott.com/familiarity-admits-brevity).

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

Pointers in Go serve various orthogonal purposes. One is expressing optionality.
The simplest way to create a variable that can either have a value or not is
with a pointer. Here, `nil` signifies the absence of a value, while a non-nil
pointer indicates the presence of a value. However, Go does not enforce checks
for absent values, which leads to runtime panics when attempting to access an
uninitialized variable. The Git Town codebase wraps optional values in a generic
Option type. This approach makes it explicit to both human and machine readers
whether a type is optional. It also enforces optionality checks or at least
makes their absence obvious.

Another use of pointers in Go is for performance optimization: if a variable is
too large to pass by value, it can be passed by reference. The Git Town codebase
does not employ this optimization as it is not necessary in our use case.

Pointers can also indicate mutability. To mutate variables passed as function
arguments they must be pointers. The challenge here is that it can be unclear
why a function argument is a pointer: it is optional, mutable, or simply too
large to pass by value? To clarify this, the Git Town codebase wraps mutable
function arguments and struct fields in the generic `Mutable` type to denote
mutability. Any variable not wrapped in a `Mutable` should be considered
immutable.

While this practice introduces a small amount of additional complexity, it is
justified by the increased robustness of the codebase. It has eliminated entire
categories of bugs that occurred relatively frequently before. We have adopted
the naming conventions from the Rust programming language as they have proven
effective in that community.

#### One concept per file

In the Git Town codebase each concept (such as type definitions, functions, or
constants) is located in its own file. This organization simplifies the process
of locating specific concepts by opening the file with the matching name.
