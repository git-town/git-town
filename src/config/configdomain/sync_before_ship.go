package configdomain

type SyncBeforeShip bool

func (self SyncBeforeShip) Bool() bool {
	return bool(self)
}
