package golzjd

const c1 int64 = 0xcc9e2d51
const c2 int64 = 0x1b873593
const c3 int64 = 0x85ebca6b // fmix
const c4 int64 = 0xc2b2ae35 // fmix

type Datablock struct {
	I int32
	C [4]int8
}

type MurmurHash3 struct  {
	Len uint32
	H1 int32
	seed int32
	Data Datablock
}

func NewDataBlock() Datablock {
	return Datablock{0, [4]int8{0,0,0,0}}
}

func (b *Datablock) Reset() {
	b.I = 0
	b.C = [4]int8{0,0,0,0}
}

func NewMurmurHash3() MurmurHash3 {
	return NewMurmurHash3FromSeed(0)
}

func NewMurmurHash3FromSeed(s int32) MurmurHash3 {
	var m MurmurHash3
	m.seed = s
	m.Reset()
	return m
}

func (mh* MurmurHash3) Reset() {
	mh.Len = 0
	mh.H1 = mh.seed
	mh.Data = NewDataBlock()
}

func (mh* MurmurHash3) PushByte(b int8) int32 {
	//store the current byte of input
	mh.Data.C[mh.Len % 4] = b
	mh.Len++

	var h1_as_if_done int32 = 0
	if (mh.Len > 0 && mh.Len % 4 == 0) { //we have a valid history of 4 items!
		// little endian load order
		k1 := mh.Data.I
		k1 = int32(int64(k1) * c1)
		k1 = int32(rotl32(uint32(k1),15))
		k1 = int32(int64(k1) * c2)

		mh.H1 ^= k1
		mh.H1 = int32(rotl32(uint32(mh.H1),13))
		mh.H1 = int32(int64(mh.H1) * int64(5) + 0xe6546b64)

		h1_as_if_done = mh.H1
		mh.Data.I = 0 //data is out the window now
	} else {
		k1 := mh.Data.I
		h1_as_if_done = mh.H1

		k1 = int32(int64(k1) * c1)
		k1 = int32(rotl32(uint32(k1), 15))
		k1 = int32(int64(k1) * c2)
		h1_as_if_done ^= k1
	}

	h1_as_if_done ^= int32(mh.Len)
	h1_as_if_done = fmix32(h1_as_if_done)
	return h1_as_if_done
}

func rotl32(x uint32, r int8) uint32 {
	return uint32(int32(x << r) | int32(x >> (32 - r)))
}

func fmix32(h int32) int32 {
	h ^= int32(uint32(h) >> 16) // similar to >>> in Java
	h = int32(int64(h) * c3)
	h = int32(int64(h) ^ int64(uint32(h) >> 13))
	h = int32(int64(h) * c4)
	h = int32(int64(h) ^ int64(uint32(h) >> 16))
	return h
}