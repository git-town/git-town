# Git Town architecture

### Design goals

The major design goals of the Git Town codebase are:

1. Execute a number of Git operations depending on conditions in the Git repo.
   Some of these conditions might change at runtime.
2. Allow the end user to resolve problems in the same terminal window and shell
   environment that Git Town executes in.
3. Reliably undo anything that Git Town has done upon request.

### General structure

To keep the amount of code manageable, the Git Town codebase separates
functionality into subsystems for parsing configuration data, syncing branches,
undoing changes, interacting with the CLI, interacting with external hosting
services, etc.

Each subsystem defines its own domain concepts, helpers, and business logic. To
prevent cyclic package dependencies, subsystems define concepts and data types
in dedicated `*domain` packages.

Higher-level subsystems like syncing branches use lower-level subsystems like
Git and configuration access but never the other way around.

### Execution framework

Git Town addresses requirements 1 and 2 via an
[interpreter](https://en.wikipedia.org/wiki/Interpreter_(computing)) that
executes programs consisting of using Git-related opcodes. Each Git Town
command:

- inspects the state of the Git repo
- assembles a program that implements the Git operations that Git Town needs to
  perform
  - this program consists of opcodes that the Git Town interpreter can execute
- starts the Git Town interpreter engine to execute this program

If there are issues that require the user to resolve in a terminal window, the
interpreter:

- persists the current interpreter state (runstate) to disk
- exits the running Git Town process to lets the user use the terminal window
  and shell environment that they used to call Git Town to resolve the problems
- prints an explanation of the problem and what the user needs to do

After resolving the problems and restarting Git Town, the interpreter recognizes
and loads the persisted state from disk and resumes executing it.

### Undo framework

To undo a previously run Git Town command (requirement 3), Git Town:

- compares snapshots of the affected Git repository before and after the command
  ran
- determines the changes that the Git Town command made to the repo
- creates a program that reverses these changes
- starts the interpreter to execute this program

### Custom code style

The Git Town codebase deviates in some areas from the recommended Go coding
style. These decisions weren't easy. Here is some background why we did them.

#### Favor descriptive naming over brevity

Many Go codebases, including Go's standard library, use heavily abbreviated
identifier names. Git Town's code base favors self-describing identifier names
over short ones because that's often quicker, less ambiguous, and leads to
better readable and understandable code with fewer bugs. This is especially true
for an open-source codebase that many readers aren't familiar with.

Code with descriptive identifier names is easier to work with because it doesn't
require keeping the mapping of several concepts to their abbreviations in one's
head while thinking about the code. See
https://michaelwhatcott.com/familiarity-admits-brevity for more background.

#### Use `self` for method receivers

The
[Go code review comments wiki page](https://go.dev/wiki/CodeReviewComments#receiver-names)
recommends avoiding generic names like `this` or `self` for method receivers and
instead use short one or two letter names. After doing this for many years we
find this approach unhelpful in practice. Git Town uses `self` for method
receivers and enforces this using a linter. This decision, while costly in terms
of going against a pretty widespread convention, has been worthwile because it
made an entire array of inconveniences and headaches disappear.

The Go review comments wiki page is incorrect that the method receiver is just
another function argument. The method receiver gets defined separate from the
other arguments. It is the central data element in the method, otherwise that
method shouldn't be a method but a function. The method receiver is the only
argument of which one can safely access private properties without violating
abstraction and encapsulation boundaries.

Go doesn't provide a clear convention for exactly how to name the method
receiver. A number of alternatives exist, each with their distinct pros and
cons, and none working for all cases. This leads to time and energy wasted
figuring out the right method receiver name and justifying it in code reviews.
The only option that works in all cases without any bikeshedding is `self`.

The Go recommendation leads to excessive churn. Renaming a type now also
requires renaming the receiver in all its methods. This leads to changes on
dozens more lines for simple rename refactors. This isn't tool supported because
of the lack of a convention, so has to be done manually and reviewed with some
level of care.

This makes refactoring unnecessarily costly, noisy, and thereby sometimes not
worth the effort. It is critical that code smells and drift get cleaned up
regularly without one having to justify it because refactoring is essential for
maintaining code health. A healthy codebase is the most important asset in the
21st century because it enables product and business agility and thereby
success.

`self` is pretty short, it's only four characters.

The names of the other function arguments also shouldn't be abbreviated. All
identifiers need to be descriptive.

https://michaelwhatcott.com/receiver-names-in-go and
https://dev.to/codypotter/the-case-for-self-receivers-in-go-3h7f provide more
background.
