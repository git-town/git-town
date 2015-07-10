Git is an intentionally generic and low-level tool.
It is designed to support many different ways of using it equally well.
It is a really wonderful *foundation* for flexible and robust source code management.


### Problem

Most teams develop code in [feature branches](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow)
that are cut from and merged into a main branch (typically "development" or "master").
They use a central repository like [GitHub](http://github.com/), [Bitbucket](https://bitbucket.org/), or their own Git server.

Using Git directly for such high-level development workflows is possible, but it can be cumbersome and repetitive.

* Creating a new feature branch in the middle of active development requires up to 6 individual Git commands.
* Updating a feature branch with the latest changes from the main branch takes up to 7.
* Even merging a finished feature into the main branch can require 5 commands.

As busy developers, we don't have the time and energy to type and run these commands all the time.
Omitting certain steps from these workflows on larger development teams leads to problems like

* frequently outdated feature branches that conflict with the main branch
* lengthy, and therefore hard to review, pull requests containing several independent changes that should be in their own branches
* broken main branches after merging in conflicting features
* a convoluted Git history that is hard to read
* many other headaches that diminish productivity and happiness as the development team grows


### Solution

The underlying problem with all these issues is that an intentionally low-level and generic source
code management tool (Git) is used directly for specific high-level development workflows.
This will always be inefficient.

_Git Town_ solves these issue by adding [high-level workflow commands](/README.md#commands) to Git.
These commands make it super easy to create, merge, synchronize, and clean up feature branches.
They are robust, safe, and [well-tested](https://github.com/Originate/git-town/tree/master/features).



### Advantages

Git Town increases team productivity and fun by helping achieve

* smaller and more focused feature branches that are easier and faster to review
* less frequent and severe merge conflicts, thanks to regular feature branch synchronization
* a more stable main branch, thanks to resolving issues before merging feature branches
* a more manageable repository, thanks to automatic removal of stale branches
* less time spent using Git and more time spent coding, thanks to
  * minimizing network requests: each command performs just a single fetch and skips unnecessary pushes
  * executing several Git commands as a batch, rather than manually typing and executing commands sequentially

Git Town empowers the user to reach these goals in a natural and intuitive way.
Git is not modified or interfered with in any way, as Git Town simply an extension of Git.

Git Town reduces the Git learning curve for beginners. By doing the complicated stuff under the hood, it allows
them to use Git like experts. However, because Git Town is automating basic Git commands, it's also completely transparent.
Users can see exactly what commands are being issued as part of Git Town's higher-level commands.

It also allows experts use Git more efficiently by automating what they would have manually typed and mentally
processed tens or even hundreds of times a day!
