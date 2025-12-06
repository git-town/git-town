import { exec } from "node:child_process"
import { promisify } from "node:util"
import { HelpOutput } from "./help-output"
const execAsync = promisify(exec)

/** GitTownCommand represents a specific Git Town command, like "append" or "sync" */
export class GitTownCommand {
  name: string

  constructor(name: string) {
    this.name = name
  }

  /** provides the actual arguments that this Git Town command accepts, determined by calling it with --help and parsing the output */
  async actualArgs(): Promise<string[][]> {
    const result = await execAsync(`git town ${this.name} --help`)
    const output = new HelpOutput(result.stdout)
    return output.flags()
  }
}
