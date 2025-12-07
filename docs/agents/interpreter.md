# VM-Based Execution Framework

Git Town uses an interpreter that executes self-modifying code consisting of
Git-related opcodes:

- Commands inspect Git repo state and generate a program of opcodes
- The interpreter (`internal/vm/`) executes these programs
- Programs can modify themselves at runtime based on repo state
- Runstate is persisted to disk for resume capability
