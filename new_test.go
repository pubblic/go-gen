package gen

import "golang.org/x/exp/rand"

func Example_New() {
	source := rand.NewSource(987654321)

	// 10 digits random number string
	_ = New(source,
		Number(10),
	)

	// 10 ~ 15 digits random number string
	_ = New(source,
		Number(10),
		Option{Number(1)},
		Option{Number(1)},
		Option{Number(1)},
		Option{Number(1)},
		Option{Number(1)},
	)

	// 10 ~ 15 digits random number string
	_ = New(source,
		Number(10),
		Repeat{Option{Number(1)}, 5},
	)

	// 5 random alphabet
	_ = New(source,
		Alphabet(5),
	)

	// 5 ~ 7 random lowercased alphabet
	_ = New(source,
		LowerAlphabet(5),
		Option{LowerAlphabet(1)},
		Option{LowerAlphabet(1)},
	)

	// Literal string or bytes
	_ = New(source,
		String("hello"),
		Bytes([]byte("world")),
	)

	// User-defined ASCII table
	// Choose randomly 10 character in "hello world"
	_ = New(source,
		Repeat{ByteTable("hello world"), 10},
	)

	// Unicode table
	_ = New(source,
		Repeat{RuneTable([]rune("가나다")), 10},
	)

	// Choose randomly
	// Result is one of 111, 222 and 333.
	_ = New(source,
		Pick{
			String("111"),
			String("222"),
			String("333"),
		},
	)

	// Random shuffle
	// Result is one of the permutations of 111, 222 and 333.
	_ = New(source,
		Shuffle{
			String("111"),
			String("222"),
			String("333"),
		},
	)
}
