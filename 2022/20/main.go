package main

import (
	"fmt"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Grove Positioning System", part1, part2).Run()
}

func part1() error {
	payload, err := parse(InputFile)
	if err != nil {
		return err
	}

	crypto := NewMixingCrypto(payload)
	decrypted := crypto.Decrypt(1)

	var zeroIndex int
	for index, num := range decrypted {
		if num == 0 {
			zeroIndex = index
			break
		}
	}

	var groveCoordinates []int
	for _, n := range []int{1000, 2000, 3000} {
		groveCoordinates = append(
			groveCoordinates,
			decrypted[maths.Mod(zeroIndex+n, len(decrypted))],
		)
	}
	fmt.Println(maths.Sum(groveCoordinates...))
	return nil
}

func part2() error {
	payload, err := parse(InputFile)
	if err != nil {
		return err
	}

	// apply decryption key
	for i := range payload {
		payload[i] = payload[i] * 811589153
	}

	crypto := NewMixingCrypto(payload)
	decrypted := crypto.Decrypt(10)

	var zeroIndex int
	for index, num := range decrypted {
		if num == 0 {
			zeroIndex = index
			break
		}
	}

	var groveCoordinates []int
	for _, n := range []int{1000, 2000, 3000} {
		groveCoordinates = append(
			groveCoordinates,
			decrypted[maths.Mod(zeroIndex+n, len(decrypted))],
		)
	}
	fmt.Println(maths.Sum(groveCoordinates...))
	return nil
}

func parse(path string) ([]int, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var nums []int
	for _, line := range lines {
		if line == "" {
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}

type MixingCrypto struct {
	payload   []int
	decrypted []int

	payloadToDecrypted []int
	decryptedToPayload []int
}

func NewMixingCrypto(payload []int) *MixingCrypto {
	crypto := &MixingCrypto{
		payload:            append([]int{}, payload...),
		decrypted:          append([]int{}, payload...),
		payloadToDecrypted: make([]int, len(payload)),
		decryptedToPayload: make([]int, len(payload)),
	}

	for i := range crypto.payload {
		crypto.payloadToDecrypted[i] = i
		crypto.decryptedToPayload[i] = i
	}

	return crypto
}

func (crypto *MixingCrypto) Decrypt(rounds int) []int {
	for i := 0; i < rounds; i++ {
		for index, num := range crypto.payload {
			crypto.applyMix(index, num)
		}
	}
	return crypto.decrypted
}

func (crypto *MixingCrypto) applyMix(index, num int) {
	// A mix equal to the length of payload means we end up back where we started,
	// so we can discount anything more than or equal to the length of the payload
	limit := maths.Mod(maths.Abs(num), len(crypto.payload)-1)
	
	alreadyIncreasedLimit := false
	for i := 0; i < limit; i++ {
		direction := 1
		if num < 0 {
			direction = -1
		}

		from := crypto.payloadToDecrypted[index]
		to := maths.Mod(from+direction, len(crypto.payload))

		// This is a horrible hack to make sure we do 1 extra if we've moved it past either end
		// this simulates the 'circular' bahaviour
		// i.e. If we have 1,2,3 and we've move 2 one space then we get 1,3,2 but we want it to look like 2,1,3
		//      so we need an extra iteration to do the additional swap
		//
		// The flag is so we only do it once per "move"
		if to == 0 || to == len(crypto.payload)-1 {
			if alreadyIncreasedLimit {
				alreadyIncreasedLimit = false
				} else {
					alreadyIncreasedLimit = true
					limit += 1
				}
			}
		// special case if we go from start to end
		if from == 0 && to == len(crypto.payload)-1 {
			crypto.decrypted = append(crypto.decrypted[1:], crypto.decrypted[0])

			// recompute indices
			for idx := range crypto.payload {
				crypto.payloadToDecrypted[idx] = maths.Mod(crypto.payloadToDecrypted[idx]-1, len(crypto.payload))
			}
			crypto.payloadToDecrypted[index] = len(crypto.payload) - 1
			for pidx, dinx := range crypto.payloadToDecrypted {
				crypto.decryptedToPayload[dinx] = pidx
			}

			continue
		}
		// special case if we go from end to start
		if from == len(crypto.payload)-1 && to == 0 {
			crypto.decrypted = append([]int{crypto.decrypted[len(crypto.payload)-1]}, crypto.decrypted[:len(crypto.payload)-1]...)

			// recompute indices
			for idx := range crypto.payload {
				crypto.payloadToDecrypted[idx] = maths.Mod(crypto.payloadToDecrypted[idx]+1, len(crypto.payload))
			}
			crypto.payloadToDecrypted[index] = 0
			for pidx, dinx := range crypto.payloadToDecrypted {
				crypto.decryptedToPayload[dinx] = pidx
			}

			continue
		}

		fromPayloadIndex := index
		toPayloadIndex := crypto.decryptedToPayload[to]

		// update payload to decrypted
		crypto.payloadToDecrypted[fromPayloadIndex] = to
		crypto.payloadToDecrypted[toPayloadIndex] = from

		// update decrypted to payload
		crypto.decryptedToPayload[to] = fromPayloadIndex
		crypto.decryptedToPayload[from] = toPayloadIndex

		// swap
		crypto.decrypted[from], crypto.decrypted[to] = crypto.decrypted[to], crypto.decrypted[from]
	}
}
