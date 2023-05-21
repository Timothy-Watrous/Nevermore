package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(equip{},
		"Usage:  equip item # \n\n Try to equip an item from your inventory",
		permissions.Player,
		"equip", "wield", "wear")
}

type equip cmd

func (equip) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendBad("What did you want to equip?")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	what := s.actor.Inventory.Search(name, nameNum)
	if what != nil {
		s.actor.RunHook("combat")
		if ok, msg := s.actor.CanEquip(what); !ok {
			s.msg.Actor.SendBad(msg)
			s.ok = true
			return
		}

		if s.actor.Equipment.Equip(what) {
			s.msg.Actor.SendGood("You equip " + what.DisplayName())
			s.msg.Observers.SendInfo(s.actor.Name + " equips " + what.DisplayName())
			s.actor.Inventory.Remove(what)
		} else {
			s.msg.Actor.SendBad("You cannot equip that.")
		}

		s.ok = true
		return
	}
	s.msg.Actor.SendInfo("What did you want to equip?")
	s.ok = true
}
