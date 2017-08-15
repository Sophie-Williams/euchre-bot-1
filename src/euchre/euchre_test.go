package euchre

import (
    "deck"
    "testing"
)

// Test euchre#Beat. A helper to determine which card wins in a head to head
// faceoff where one card is chosen to be leading and there is a defined
// trump suit.
// -----------------------------------------------------------------------------
func TestBeatTrumpNonTrump(t *testing.T) {
    trump := deck.H
    card1 := deck.Card{ deck.H, deck.Ten }
    card2 := deck.Card{ deck.C, deck.A }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}

func TestBeatNonTrumpTrump(t *testing.T) {
    trump := deck.C
    card1 := deck.Card{ deck.S, deck.Q }
    card2 := deck.Card{ deck.C, deck.A }

    if Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card2, card1, trump)
    }
}

func TestBeatDifferentNonTrumpFirstLess(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.S, deck.Nine }
    card2 := deck.Card{ deck.C, deck.A }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}

func TestBeatDifferentNonTrumpFirstGreater(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.S, deck.A }
    card2 := deck.Card{ deck.C, deck.Nine }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}

func TestBeatSameSuit(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.S, deck.A }
    card2 := deck.Card{ deck.S, deck.Nine }

    if !Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card1, card2, trump)
    }
}

func TestBeatBowers(t *testing.T) {
    trump := deck.D
    card1 := deck.Card{ deck.H, deck.J }
    card2 := deck.Card{ deck.D, deck.J }

    if Beat(card1, card2, trump) {
        t.Errorf("Expected %s to beat %s with trump %s", card2, card1, trump)
    }
}

// euchre#Winner. A helper to determine which player is the winner given the
// player who led, the list of cards, and the trump suit.
// -----------------------------------------------------------------------------
func TestWinnerOneTrump1(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.D, deck.A },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.D, deck.Ten },
        deck.Card{ deck.D, deck.Q },
    }

    trump := deck.H
    led := 2
    answer := 3

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}

func TestWinnerOneTrump2(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.Nine },
        deck.Card{ deck.S, deck.A },
        deck.Card{ deck.D, deck.Q },
    }

    trump := deck.D
    led := 2
    answer := 1

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}

func TestWinnerOneTrump3(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.D, deck.Q },
        deck.Card{ deck.D, deck.K },
        deck.Card{ deck.C, deck.Ten },
        deck.Card{ deck.D, deck.A },
    }

    led := 1
    trump := deck.C
    answer := 3

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}

func TestWinnerOneTrump4(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.J },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.C, deck.Ten },
    }

    led := 2
    trump := deck.C
    answer := 1

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}

func TestWinnerOneTrump5(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.J },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.C, deck.Ten },
    }

    led := 2
    trump := deck.C
    answer := 1

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}

func TestWinnerNonTrump(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.H, deck.Ten },
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.H, deck.Ten },
    }

    led := 3
    trump := deck.C
    answer := 0

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}

func TestWinnerManyTrump(t *testing.T) {
    played := []deck.Card {
        deck.Card{ deck.D, deck.J },
        deck.Card{ deck.C, deck.Ten },
        deck.Card{ deck.S, deck.J },
        deck.Card{ deck.C, deck.A },
    }

    led := 1
    trump := deck.C
    answer := 3

    res := Winner(played, trump, led)
    if answer != res {
        t.Errorf("Expected winner to be %d but got %d instead.", answer, res)
    }
}

// Test euchre#Possible. Returns the possible cards in a hand given the
// currently played cards.
// -----------------------------------------------------------------------------
func TestPossibleCantFollow(t *testing.T) {
    hand := []deck.Card {
        deck.Card{ deck.H, deck.Ten },
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.C, deck.J },
        deck.Card{ deck.C, deck.K },
        deck.Card{ deck.D, deck.J },
    }
    played := []deck.Card {
        deck.Card{ deck.S, deck.Ten },
    }
    trump := deck.H

    res := Possible(hand, played, trump)

    if len(res) != 5 {
        t.Errorf("Expected all cards to be possible.")
    }
}

func TestPossibleCanFollow(t *testing.T) {
    hand := []deck.Card {
        deck.Card{ deck.H, deck.Ten },
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.C, deck.J },
        deck.Card{ deck.C, deck.K },
        deck.Card{ deck.S, deck.J },
    }
    played := []deck.Card {
        deck.Card{ deck.S, deck.Ten },
    }
    trump := deck.H

    res := Possible(hand, played, trump)

    if len(res) != 1 || res[0] != 4 {
        t.Errorf("Expected only #4 not %d option(s).", len(res))
    }
}


/* Test euchre#noSuits.
 *
 * Checks whether we correctly learn what suits a player has based on their
 * inability to follow.
 */


/*
 * Test whether no players are identified as having no suits, since no cards
 * were played yet.
 */
func TestNoSuitsEmpty(t *testing.T) {
    prior := make([]Trick, 0)
    trump := deck.H

    res := noSuits(prior, trump)

    if len(res) > 0 {
        t.Errorf("Expected no players to lack a suit, but %d do.\n", len(res))
    }
}


/*
 * Test when only one player does not a given suit.
 */
func TestNoSuitsOne(t *testing.T) {
    prior := make([]Trick, 2)
    player := 1
    trump := deck.C

    cards1 := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.C, deck.K },
        deck.Card{ deck.H, deck.K },
        deck.Card{ deck.H, deck.Q },
    }
    cards2 := []deck.Card {
        deck.Card{ deck.S, deck.K },
        deck.Card{ deck.S, deck.Ten },
        deck.Card{ deck.S, deck.Q },
        deck.Card{ deck.S, deck.A },
    }

    trick1 := Trick {
        cards1,
        0,
        trump,
    }
    trick2 := Trick {
        cards2,
        1,
        trump,
    }

    prior[0] = trick1
    prior[1] = trick2

    res := noSuits(prior, trump)
    if len(res) != 1 {
        t.Errorf("Expected 1 player to not have suits, but %d do\n", len(res))
    }

    playerRes := res[player]
    if len(playerRes) != 1 || (len(playerRes) >= 1 && playerRes[0] != deck.H) {
        t.Errorf("Expected only H to be impossible, but these are: %d\n",
                 len(playerRes))

        for _, suit := range playerRes {
            t.Errorf(" %s ", suit)
        }
    }
}


/*
 * Test if unpresent suits can be detected for player 3, who is last in the
 * modulo order.
 */
func TestNoSuitsThree(t *testing.T) {
    prior := make([]Trick, 2)
    player := 3
    trump := deck.C

    cards1 := []deck.Card {
        deck.Card { deck.H, deck.A },
        deck.Card { deck.C, deck.K },
        deck.Card { deck.H, deck.K },
        deck.Card { deck.H, deck.Q },
    }
    cards2 := []deck.Card {
        deck.Card { deck.S, deck.K },
        deck.Card { deck.S, deck.Ten },
        deck.Card { deck.C, deck.Ten },
        deck.Card { deck.C, deck.A },
    }

    trick1 := Trick {
        cards1,
        0,
        trump,
    }
    trick2 := Trick {
        cards2,
        1,
        trump,
    }

    prior[0] = trick1
    prior[1] = trick2

    res := noSuits(prior, trump)
    if len(res) != 3 {
        t.Errorf("Expected 3 players to not have some suit, but %d do.\n",
                 len(res))
    }

    playerRes := res[player]
    if len(playerRes) != 1 || (len(playerRes) == 1 && playerRes[0] != deck.S) {
        t.Errorf("Expected only S to be impossible, but these are:")

        for _, suit := range playerRes {
            t.Errorf(" %s ", suit)
        }
    }
}


/*
 * Test that the module arithmetic works and that tricks whose player numbers
 * wrap accurately keep track of what player followed or did not follow suit.
 */
func TestNoSuitsPlayerWraps(t *testing.T) {
    prior := make([]Trick, 2)
    player := 0
    trump := deck.C

    cards1 := []deck.Card {
        deck.Card{ deck.H, deck.A },
        deck.Card{ deck.C, deck.K },
        deck.Card{ deck.H, deck.K },
        deck.Card{ deck.H, deck.Q },
    }
    cards2 := []deck.Card {
        deck.Card{ deck.S, deck.K },
        deck.Card{ deck.S, deck.Ten },
        deck.Card{ deck.C, deck.Ten },
        deck.Card{ deck.C, deck.A },
    }

    trick1 := Trick {
        cards1,
        3,
        trump,
    }
    trick2 := Trick {
        cards2,
        0,
        trump,
    }

    prior[0] = trick1
    prior[1] = trick2

    res := noSuits(prior, trump)
    if len(res) != 3 {
        t.Errorf("Expected 3 players to not have some suit, but %d do.\n",
                 len(res))
    }

    playerRes := res[player]
    if len(playerRes) != 1 || (len(playerRes) >= 1 && playerRes[0] != deck.H) {
        t.Errorf("Expected only H to be impossible, but these are:")

        for _, suit := range res {
            t.Errorf(" %s ", suit)
        }
    }
}


/*
 * Test that multiple no suits for one player can be detected.
 */
func TestNoSuitsMultiple(t *testing.T) {
    prior  := make([]Trick, 4)
    player := 1
    trump  := deck.S

    cards1 := []deck.Card {
        deck.Card{ deck.D, deck.A },
        deck.Card{ deck.S, deck.J },
        deck.Card{ deck.D, deck.Ten },
        deck.Card{ deck.D, deck.J },
    }
    cards2 := []deck.Card {
        deck.Card{ deck.C, deck.A },
        deck.Card{ deck.C, deck.Q },
        deck.Card{ deck.C, deck.K },
        deck.Card{ deck.S, deck.Ten },
    }
    cards3 := []deck.Card {
        deck.Card{ deck.H, deck.K },
        deck.Card{ deck.H, deck.Ten },
        deck.Card{ deck.H, deck.Nine },
        deck.Card{ deck.H, deck.A },
    }
    cards4 := []deck.Card {
        deck.Card{ deck.H, deck.Q },
        deck.Card{ deck.S, deck.A },
        deck.Card{ deck.C, deck.Ten },
        deck.Card{ deck.S, deck.K },
    }

    trick1 := Trick {
        cards1,
        0,
        trump,
    }
    trick2 := Trick {
        cards2,
        1,
        trump,
    }
    trick3 := Trick {
        cards3,
        0,
        trump,
    }
    trick4 := Trick {
        cards4,
        3,
        trump,
    }

    prior[0] = trick1
    prior[1] = trick2
    prior[2] = trick3
    prior[3] = trick4

    res := noSuits(prior, trump)
    if len(res) != 3 {
        t.Errorf("Expected 3 players to not have some suit, but %d do.\n",
                 len(res))
    }

    playerRes := res[player]
    if len(playerRes) != 2 {
        if (playerRes[0] == deck.D && playerRes[1] == deck.H) ||
           (playerRes[0] == deck.H && playerRes[1] == deck.D) {
            t.Errorf("Expected H and D to not be possible but got:")

            for _, suit := range playerRes {
                t.Errorf(" %s ", suit)
            }
        }
    }
}
