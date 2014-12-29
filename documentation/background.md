# The Big Picture

Git is an intentionally generic and basic tool.
It is designed to support many different ways of using it equally well.
It is a really wonderful *foundation* for flexible and robust source code management.


### Problem

Most teams develop code in [feature branches](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow) that are cut from and merged into a main branch (typically "development" or "master").
They use a central repository like [GitHub](http://github.com/), [Bitbucket](https://bitbucket.org/), or their own Git server.

Using Git as an intentionally generic and low-level tool directly for such high-level development workflows is possible but cumbersome and repetitive.

* Creating a new feature branch in the middle of active development requires up to 6 individual Git commands.
* Updating a feature branch with the latest changes from the main branch takes up to 7.
* Even merging a finished feature into the main branch requires 5 commands when done correctly.

Almost nobody has the time and energy to manually type and run all these commands at all times.
Omitting them on larger development teams leads to problems like

* frequently outdated feature branches that conflict with the main branch
* sizable and therefore hard to review pull requests that contain several independent changes that shoud be in their own branch
* breakage on the main branches after merging in conflicting features
* a convoluted Git history that is hard to use
* many other headaches that diminish productivity and happiness as the development teams grow.


### Solution
The underlying problem of all these issues is that an intentionally low-level and generic source code management tool (Git) is used directly for specific high-level development workflows.
This will always be inefficient.

_Git Town_ solves these issue by adding [high-level workflow commands](../readme.md#commands) to Git.
These commands make it super easy to create, merge, synchronize, and clean up feature branches.
They are robust, safe, and [well tested](https://github.com/Originate/git-town/tree/master/features).



### Advantages

Git Town increases team productivity and fun by helping achieve

* smaller and more focussed feature branches that are easier and faster to review
* less frequent and severe merge conflicts thanks to regular feature branch synchronization
* less brokenness on the main branch thanks to resolving issues before merging feature branches
* a more manageable repository thanks to automatic removal of old branches
* less time spent using Git and more time spent coding, thanks to
  * minimizing network requests: each command performs just a single fetch and skips unnecessary pushes
  * executing several Git commands as a batch rather than them being typed and run interactively

Git Town achieves all this in a completely natural and non-intrusive way that uses no special or magic tricks.
It allows to use the rest of Git normally.

Git Town reduces the Git learning curve for beginners. By doing the complicated stuff under the hood it allows them to use Git like experts.
It also allows experts use Git more efficiently, by automating what they would have typed and processed manually over and over again.
