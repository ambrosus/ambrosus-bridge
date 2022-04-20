package contracts

func (b *Bridge) Raw() *BridgeRaw {
	return &BridgeRaw{Contract: b}
}
