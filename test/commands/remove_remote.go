package commands

// RemoveRemote deletes the Git remote with the given name.
func RemoveRemote(cmds *TestCommands, name string) error {
	cmds.Config.RemotesCache.Invalidate()
	_, err := cmds.Run("git", "remote", "rm", name)
	return err
}
