package tools

// Helper function to read a uint32 from a byte slice
func readUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

// Helper function to read an int32 from a byte slice
func readInt32(b []byte) int32 {
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

// Helper function to read a uint16 from a byte slice
func readUint16(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}
