# Git Town architecture

### Design goals

The complexity in the Git Town codebase stems from balancing several challenging
design objectives:

1. **Extreme configurability:** Execute a highly variable set of Git operations
   depending on the current status of the repository. Git Town's business logic
   covers so many edge cases that most Git Town commands aren't just a simple
   hard-coded scripts, they are executable programs custom-built for your
   specific use case.
2. **Terminate and resume:** When any operation in these programs fails,
   terminate the entire application to allow the end user to resolve problems in
   the same terminal window and shell environment that they ran Git Town in, and
   then resume execution.
3. **Reliable undo:**Be able to reliably undo everything that Git Town has done.

### General structure

The Git Town codebase separates independent parts of the complex domain model
into orthogonal, composable _subsystems_. Subsystems define their own domain
model types, data structures, and business logic. Examples for subsystems are
configuration data, the interpreter that executes programs, executing shell
commands, interacting with code hosting APIs, determining undo operations, etc.

Subsystems define concepts and data types in dedicated `*domain` packages so
that they can all use each other's data types without introducing cyclic
dependencies.

### Execution framework

Git Town addresses requirements 1 and 2 via an
[interpreter](https://en.wikipedia.org/wiki/Interpreter_(computing)) that
executes programs consisting of Git-related _opcodes_. Most of these opcodes
execute Git commands. Some open browser windows or talk to codehosting APIs.
Each Git Town command:

- Inspects the state of the Git repo.
- Determines the Git operations that Git Town needs to perform and stores them
  as a Git Town program.
- Starts the Git Town interpreter to execute this program.

If there are issues that require the user to edit files or run Git commands, the
interpreter:

- Persists the current interpreter state (called "runstate") to disk.
- Exits Git Town to give the user access to the shell to resolve the encountered
  problems.
- Prints a human-friendly explanation of the problem and what the user needs to
  do.

After resolving the problems and restarting Git Town, the interpreter loads the
persisted runstate from disk and resumes executing it.

### Undo framework

To undo a previously run Git Town command (requirement 3), Git Town:

- compares snapshots of the affected Git repository before and after the command
  ran
- determines the changes that the Git Town command made to the Git repo
- creates a program that reverses these changes
- starts the interpreter to execute this program

### Code style

The Git Town codebase includes some intentional deviations from the "official"
Go coding style. These decisions were made after extensive experience with the
standard Go style revealed issues that could be avoided through these changes.
Below is an explanation of these differences and the rationale behind them.

#### Favor descriptive naming over brevity

The Go community often uses highly abbreviated names for variables, types, and
functions, following the personal preference of some of Go's creators. While
brevity can be useful, our primary focus for code quality in this codebase is
ease of understanding and achieving self-describing code. Our open-source tool
has a wide user base and a small group of maintainers, with many contributors
adding just a single feature. To ensure our code is accessible to everyone, we
consistently use descriptive identifiers even if they are longer than a few
characters. For more context, please refer to
[this article](https://michaelwhatcott.com/familiarity-admits-brevity).

#### Use `self` for method receivers

The
[Go code review comments wiki page](https://go.dev/wiki/CodeReviewComments#receiver-names)
suggests using short, one or two-letter names for method receivers, rather than
generic names like `this` or `self`. After many years of following this
guideline in the Git Town codebase, we found it to be impractical and
cumbersome.

Contrary to the wiki page, the method receiver is more than just another
function argument. It is defined separately and serves as the central data
element in the method. Otherwise, the method should simply be a function. The
method receiver is unique in that it allows safe access to private properties
without violating abstraction and encapsulation boundaries.

Go does not provide a clear convention for naming method receivers. Various
alternatives exist, each with its own pros and cons, but none are universally
effective. This ambiguity leads to wasted time and effort in determining the
appropriate receiver name and defending it during code reviews. Using `self`
consistently avoids these issues and prevents unnecessary bikeshedding.

Adhering to the Go recommendation can cause excessive churn. Renaming a type
necessitates renaming the receiver in all its methods, resulting in unnecessary
changes across many lines of code. This process is manual and time-consuming due
to the lack of a standard convention, making refactoring costly and noisy, often
deterring the effort altogether. This outcome is detrimental to the codebase, as
the ability to refactor efficiently is more important than strictly following
debatable community standards.

The name `self` is sufficiently concise, being only four characters long. For
more background please refer to https://michaelwhatcott.com/receiver-names-in-go
and https://dev.to/codypotter/the-case-for-self-receivers-in-go-3h7f.

#### Dedicated generic types for Optionality and Mutablitity

Pointers in Go serve various orthogonal purposes. One is expressing optionality.
The simplest way to create a variable that can either have a value or not is
with a pointer. Here, `nil` signifies the absence of a value, while a non-nil
pointer indicates the presence of a value. However, Go does not enforce checks
for absent values, which leads to runtime panics when attempting to access an
uninitialized variable. The Git Town codebase wraps optional values in a generic
`Option` type. This approach makes it explicit to both human and machine readers
whether a type is optional. It also enforces optionality checks.

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

While this practice introduces a thin layer of additional complexity, and more
code to handle edge cases correctly, this complexity is justified by the
drastically increased robustness of the codebase. It has eliminated entire
categories of bugs that had occurred relatively frequently before. We have
adopted the naming conventions from the Rust programming language as they have
proven effective in that community.

#### One concept per file

In the Git Town codebase each concept (such as type definitions, functions, or
constants) is located in its own file. This organization simplifies the process
of locating specific concepts by opening the file with the matching name.

#### Newtypes

Git Town's domain model has several dozen meanings for `string` and other basic
data types. Distinguishing them from each other eliminates an entire category of
bugs. Git Town's codebase therefore makes copious use of the newtype pattern,
i.e. it defines dedicated data types for each domain concepts.
