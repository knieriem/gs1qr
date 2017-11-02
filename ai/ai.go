package ai

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type AppID struct {
	ID       int
	Variable bool
}

func (ai AppID) String() string {
	digits := 2
	if ai.ID >= 1000 {
		digits = 4
	} else if ai.ID >= 100 {
		digits = 3
	}
	return fmt.Sprintf("%0*d", digits, ai.ID)
}

var (
	GTIN = &AppID{
		ID: 1,
	}

	BatchLot = &AppID{
		ID:       10,
		Variable: true,
	}
	ProdDate = &AppID{
		ID: 11,
	}
	BestBefore = &AppID{
		ID: 15,
	}

	Variant = &AppID{
		ID: 20,
	}
	Serial = &AppID{
		ID:       21,
		Variable: true,
	}

	VarCount = &AppID{
		ID:       30,
		Variable: true,
	}

	ShipToLoc = &AppID{
		ID: 410,
	}

	CPID = &AppID{
		ID:       8010,
		Variable: true,
	}
	ProductURL = &AppID{
		ID:       8200,
		Variable: true,
	}
)

type Elem struct {
	AI   *AppID
	Data string
}

var AppIDs = []*AppID{
	GTIN,
	BatchLot,
	ProdDate,
	BestBefore,
	Variant,
	Serial,
	VarCount,
	ShipToLoc,
	CPID,
	ProductURL,
}

var ErrSyntax = errors.New("syntax error")

func ParseSeq(s string) ([]Elem, error) {
	el := make([]Elem, 0, strings.Count(s, "("))
L:
	for s != "" {
		if s[0] != '(' {
			return nil, ErrSyntax
		}
		s = s[1:]
		i := strings.IndexByte(s, ')')
		if i <= 0 {
			return nil, ErrSyntax
		}
		id, err := strconv.Atoi(s[:i])
		if err != nil {
			return nil, err
		}
		s = s[i+1:]
		for _, ai := range AppIDs {
			if ai.ID == id {
				i := strings.IndexByte(s, '(')
				if i < 0 {
					el = append(el, Elem{AI: ai, Data: s})
					return el, nil
				}
				el = append(el, Elem{AI: ai, Data: s[:i]})
				s = s[i:]
				continue L
			}
		}
		return nil, errors.New("unknown AI")
	}
	return nil, nil
}
