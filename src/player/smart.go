package player

import (
    "deck"
    "euchre"
    "fmt"
    "github.com/klaidliadon/next"
    "math/rand"
    "time"
)

type situation struct {
    player1 []deck.Card
    player2 []deck.Card
    player3 []deck.Card
    kitty   []deck.Card
}

type Decision struct {
    Move int
    Value int
}

type SmartPlayer struct {
}

func NewSmart() (*SmartPlayer) {
    return &SmartPlayer{ }
}

func (p *SmartPlayer) Pickup(hand [5]deck.Card, top deck.Card, who int) bool {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return r.Intn(2) == 1
}

func (p *SmartPlayer) Discard(hand [5]deck.Card,
                              top deck.Card) ([5]deck.Card, deck.Card) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    total := hand[:]
    total = append(total, top)

    i := r.Intn(len(total))
    chosen := total[i]
    total[i] = total[len(total) - 1]
    total = total[:len(total) - 1]

    copy(hand[:], total[:5])

    return hand, chosen
}

func (p *SmartPlayer) Call(hand [5]deck.Card, top deck.Card) (deck.Suit, bool) {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    s := deck.SUITS[r.Intn(len(deck.SUITS))]
    for s == top.Suit {
        s = deck.SUITS[r.Intn(len(deck.SUITS))]
    }

    return s, r.Intn(2) == 1
}

func (p *SmartPlayer) Play(setup euchre.Setup, hand, played []deck.Card,
                           prior []euchre.Trick) ([]deck.Card, deck.Card) {
    // There are two levels of reasoning in this method. When there are 5/4
    // cards in the player's hand, there are simply too many possibilities to
    // compute. So for this amount of cards, rules will be used. Otherwise, if
    // there are 3 or fewer cards, computational power can be used to figure out
    // the best move.

    if len(hand) <= 3 {
        winners := make(map[int]int)

        for situation := range situations(setup, hand, played, prior) {
            hands := [4][]deck.Card{hand, situation.player1, situation.player2, situation.player3}
            dec := minimax(hands, played, setup.Trump, 0)
            if dec.Value == 1 {
                winners[dec.Move]++
            }
        }

        chosen := 0
        winValue := -1
        for index, count := range winners {
            fmt.Printf("%d\t%d\n", index, count)
            if winValue > count {
                winValue = count
                chosen = index
            }
        }

        for _, card := range hand {
            fmt.Print(card)
            fmt.Print(" ")
        }
        fmt.Println()
        fmt.Println(chosen)

        // TODO: Remove repetition.
        final := hand[chosen]
        hand = append(hand[:chosen], hand[chosen + 1:]...)

        return hand, final
    } else {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))

        playable := euchre.Possible(hand, played, setup.Trump)
        chosen := playable[r.Intn(len(playable))]
        final := hand[chosen]
        hand = append(hand[:chosen], hand[chosen + 1:]...)

        return hand, final
    }
}

// A minimax implementation that will return what card to play and the resulting
// hand after playing that card.
// TODO: Move to AI package.
func minimax(hands [4][]deck.Card, played []deck.Card, trump deck.Suit,
             player int) Decision {
    hand := hands[player]

    // If there is only one card in the hand then simply return that.
    if len(hand) == 1 {
        hand = append(hand[:0], hand[1:]...)

        // TODO: what is played here?
        w := euchre.Winner(played, trump, (player + 1) % 4)
        v := 0
        if w == 0 || w == 2 {
            v = 1
        } else if w == 1 || w == 3 {
            v = 0
        }

        return Decision {
            0,
            v,
        }
    } else {
        var chosen int
        var bestValue int

        if player == 0 {
            bestValue = -1
        } else if player == 1 || player == 3 || player == 2 {
            bestValue = 2
        }

        poss := euchre.Possible(hand, played, trump)
        if player == 0 {
            for _, p := range poss {
                fmt.Print(p)
                fmt.Print(" ")
            }
        }
        fmt.Println()

        for _, i := range poss {
            card := hand[i]

            // Assume the current card is played on this players turn so add it
            // to the played list for the next player. The slice must be
            // duplicated since it is a pointer, and do not want duplicate
            // references in recursive call.
            newPlayed := make([]deck.Card, len(played))
            copy(newPlayed, played)
            newPlayed = append(newPlayed, card)

            // Then remove this card from the appropriate hand.
            newHand := make([]deck.Card, len(hand))
            copy(newHand, hand)
            newHand = append(newHand[:i], newHand[i + 1:]...)
            hands[player] = newHand

            var dec Decision
            if len(played) < 4 {
                dec = minimax(hands, newPlayed, trump, (player + 1) % 4)
            } else {
                w := euchre.Winner(played, trump, (player + 1) % 4)
                v := 0
                if w == 0 {
                    v = 1
                } else if w == 1 || w == 3 || w == 2 {
                    v = 0
                }

                dec = Decision {
                    3,
                    v,
                }
            }

            // If we are the current player, then try to maximize the results.
            if player == 0 {
                if dec.Value > bestValue {
                    bestValue = dec.Value
                    chosen = i
                }
            } else if player == 1 || player == 3 || player == 2 {
            // Otherwise we will try to minimize the results.
                if dec.Value < bestValue {
                    bestValue = dec.Value
                    chosen = i
                }
            }
        }

        return Decision {
            chosen,
            bestValue,
        }
    }
}

// A generator that iterates through all possible hands or situations given a
// player's current hand, the cards currently played, and the cards played in
// previous tricks. This is a generator for now for memory purposes. For example
// if everybody has 3 cards, there are about 369000 possilibilities since the
// kitty must be taken into account.
// hand     - The current hand of the player.
// played   - The cards that have already been played on this trick.
// prior    - The cards that have been played in previous tricks.
// discard  - The card that was discarded by the player if applicable.
// dealer   - The number designation for ther person who dealt the cards.
// pickedUp - Flag to designate if the top card was picked up by the dealer.
// top      - The card that was on top of the kitty.
// TODO: How permutation works?
func situations(setup euchre.Setup, hand []deck.Card, played []deck.Card,
                tricks []euchre.Trick) chan situation {
    var prior []deck.Card
    for i := 0; i < len(tricks); i++ {
        prior = append(prior, tricks[i].Cards[:]...)
    }

    // Figure out if the top card has been played, and is thus, a known card or
    // not.
    topPlayed := false
    for _, card := range append(played, prior...) {
        if setup.Top == card {
            topPlayed = true
            break
        }
    }

    // Figure out how many unknown cards are in each player's hand and in the
    // kitty.
    nums := [4]int{len(hand), len(hand), len(hand), 4}
    for i := 2; i > 2 - len(played); i-- {
        nums[i]--
    }
    if setup.PickedUp && setup.Dealer != 0 && !topPlayed {
    // If the top was picked up by somebody else and it has not been played yet
    // we know where it is.
        nums[setup.Dealer - 1]--
    } else if (setup.PickedUp && setup.Dealer == 0) || (!setup.PickedUp) {
    // If we picked up the top card or nobody picked it up, then we know about
    // one of the 4 cards not in play.
        nums[3]--
    }

    // Create a set-like structure that has all the cards currently in play. Do
    // this by adding all cards and then removing those that are in your hand
    // have been played, and the card that was picked up if any.
    unknowns := make(map[deck.Card]bool)
    for _, card := range deck.CARDS {
        unknowns[card] = true
    }

    // We know where the top is if wasn't picked up or if somebody else picked
    // it up and it hasn't been played yet.
    if !setup.PickedUp || (setup.PickedUp && setup.Dealer != 0 && !topPlayed) {
        delete(unknowns, setup.Top)
    } else if setup.PickedUp && setup.Dealer == 0 {
    // Similarly, we know one card is (the discarded card) if we picked it up.
        delete(unknowns, setup.Discard)
    }

    // Remove all other cards that we have already seen in some way.
    for _, card := range append(hand, append(played, prior...)...) {
        delete(unknowns, card)
    }

    // Set available to the keys of the unknowns map. These are the cards whose
    // distribution is unknown.
    available := make([]deck.Card, len(unknowns), len(unknowns))
    i := 0
    for card := range unknowns {
        available[i] = card
        i++
    }

    if nums[0] + nums[1] + nums[2] + nums[3] != len(available) {
        panic("Number of freely fluxing cards is not correct.")
    }

    c := make(chan situation)
    go func() {
        // TODO
        for multi := range multinomial(nums[0], nums[1], nums[2], nums[3]) {
            cards := make([][]deck.Card, 0)
            for i := 0; i < len(multi); i++ {
                cards = append(cards, make([]deck.Card, 0))
                for j := 0; j < len(multi[i]); j++ {
                    cards[i] = append(cards[i], available[multi[i][j].(int)])
                }
            }

            next := situation {
                cards[0],
                cards[1],
                cards[2],
                cards[3],
            }

            // If the top card has not been played yet and it was picked up by
            // somebody, else add it to their cards. Basically we know of a card
            // in somebody's hand, so it wasn't freely above but should be added
            // now. Similarly, the last two ifs add a card to the kitty if we
            // know about it, i.e. it isn't picked up or we discarded a card.
            if setup.PickedUp && setup.Dealer != 0 && !topPlayed {
                switch setup.Dealer {
                case 1:
                    next.player1 = append(next.player1, setup.Top)
                case 2:
                    next.player2 = append(next.player2, setup.Top)
                case 3:
                    next.player3 = append(next.player3, setup.Top)
                }
            } else if setup.PickedUp && setup.Dealer == 0 {
                next.kitty = append(next.kitty, setup.Discard)
            } else if !setup.PickedUp {
                next.kitty = append(next.kitty, setup.Top)
            }

            c <- next
        }
        close(c)
    }()

    return c
}

// Given a list of integer sizes for multinomial choosing, return a channel that
// gives a slice of integer slices. The sizes of the integer slices correspond
// to the varidic arguments. This is the same as providing the
// (sum(ks); ks[0], ks[1], ..., ks[n]) ways to choose ks[0], ks[1], ..., ks[n]
// integers from sum(ks) integers.
// ks - The arguments to the multinomial function.
// Returns a channel you can range over that provides a slice of slices, where
// each entry has n slices for each selection.
// TODO: Channels or generator-consumer pattern?
// TODO: Improve this runtime, probably by not relying on underlying combination
//       logic.
func multinomial(ks ...int) chan [][]interface{} {
    c := make(chan [][]interface{})

    sum := 0
    for _, k := range ks {
        sum += k
    }

    idxs := make([]interface{}, sum)
    for i := 0; i < len(idxs); i++ {
        idxs[i] = i
    }

    if len(ks) > 1 {
        go func() {
            defer close(c)

            for comb := range next.Combination(idxs, ks[0], false) {
                // TODO: Check if this assumption is right. And it can change
                // any moment probably so be careful.
                // TODO: Type assertions.
                purged := make([]interface{}, sum)
                copy(purged, idxs)
                for i := len(comb) - 1; i >= 0; i-- {
                    chosen := comb[i].(int)
                    purged = append(purged[:chosen], purged[chosen + 1:]...)
                }

                for multi := range multinomial(ks[1:]...) {
                    for i := 0; i < len(multi); i++ {
                        // Must copy this choice when reassigning indexes since
                        // modifying the slice inside of multi will modify all of
                        // that element's appeareances in later combinations
                        // since it modifies the underlying array.
                        old := multi[i]
                        multi[i] = make([]interface{}, len(multi[i]))
                        for j := 0; j < len(multi[i]); j++ {
                            multi[i][j] = purged[old[j].(int)]
                        }
                    }

                    next := make([][]interface{}, 0)
                    next = append(next, comb)
                    next = append(next, multi...)

                    c <- next
                }
            }
        }()
    } else {
        go func() {
            defer close(c)

            c <- [][]interface{}{ idxs }
        }()
    }

    return c
}