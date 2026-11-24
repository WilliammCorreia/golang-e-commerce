package ids

import "crypto/rand"

const (
	Product = "PDT"
	Basket  = "BSK"
	Order   = "CMD"
)
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	out := make([]byte, 0, n)
	buf := make([]byte, 1)

	for len(out) < n {
		rand.Read(buf)

		// 252 = 7 × 36. On jette 252 à 255 : sinon, 0, 1, 2 et 3 sortiraient
		// un peu plus souvent que les autres caractères.
		if buf[0] >= 252 {
			continue
		}
		out = append(out, alphabet[buf[0]%36])
	}

	return string(out)
}

func New(prefix string) string {
	return prefix + "-" + randomString(6)
}

func Code() string {
	return randomString(6)
}
