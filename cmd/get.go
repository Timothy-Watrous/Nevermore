package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strconv"
)

func init() {
	addHandler(get{},
		"Usage:  get [container_name] itemName # \n \n Get the specified item.",
		permissions.Player,
		"GET", "TAKE", "G")
}

type get cmd

func (get) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("You go to get.. uh??")
		return
	}

	targetStr := s.words[0]
	targetNum := 1
	whereStr := ""
	whereNum := 1

	if len(s.words) > 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		} else {
			whereStr = s.words[1]
		}
	}
	if len(s.words) > 2 {
		if whereStr != "" {
			if val2, err2 := strconv.Atoi(s.words[2]); err2 == nil {
				whereNum = val2
			} else {
				whereStr = s.words[1]
			}
		}
	}

	if len(s.words) > 3 {
		if val3, err3 := strconv.Atoi(s.words[3]); err3 == nil {
			whereNum = val3
		}
	}

	if whereStr == "" {
		roomInventory := s.where.Items.Search(targetStr, targetNum)
		if roomInventory != nil {
			if roomInventory.ItemType == 10 {
				s.actor.RunHook("act")
				s.where.Items.Remove(roomInventory)
				s.actor.Gold.Add(roomInventory.Value)
				s.msg.Actor.SendGood("You picked up ", roomInventory.Name, " and put it in your gold pouch.")
				s.msg.Observers.SendInfo("You see ", s.actor.Name, " get ", roomInventory.Name, ".")
				return
			} else if (s.actor.Inventory.TotalWeight + roomInventory.GetWeight()) <= s.actor.MaxWeight() {
				s.actor.RunHook("act")
				s.actor.Inventory.Lock()
				s.where.Items.Remove(roomInventory)
				s.actor.Inventory.Add(roomInventory)
				s.actor.Inventory.Unlock()
				s.msg.Actor.SendGood("You get ", roomInventory.Name, ".")
				s.msg.Observers.SendInfo(s.actor.Name, " takes ", roomInventory.Name, ".")
				return
			} else {
				s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
				return
			}
		}
	}

	// Try to find the where if it's not the room
	if whereStr != "" {
		log.Println("Looking elsewhere for ", targetStr)
		where := s.where.Items.Search(whereStr, whereNum)
		if where != nil && where.ItemType == 9 {
			whereInventory := where.Storage.Search(targetStr, targetNum)
			if whereInventory != nil {
				if whereInventory.ItemType == 10 {
					s.actor.RunHook("act")
					where.Storage.Lock()
					where.Storage.Remove(whereInventory)
					s.actor.Gold.Add(whereInventory.Value)
					s.msg.Actor.SendGood("You take ", whereInventory.Name, " from ", where.Name, " and put it in your gold pouch.")
					s.msg.Observers.SendInfo("You see ", s.actor.Name, " take ", whereInventory.Name, " from ", where.Name, ".")
					where.Storage.Unlock()
					return
				} else if (s.actor.Inventory.TotalWeight + whereInventory.GetWeight()) <= s.actor.MaxWeight() {
					s.actor.RunHook("act")
					where.Storage.Lock()
					s.actor.Inventory.Lock()
					where.Storage.Remove(whereInventory)
					s.actor.Inventory.Add(whereInventory)
					where.Storage.Unlock()
					s.actor.Inventory.Unlock()
					s.msg.Actor.SendGood("You take ", whereInventory.Name, " from ", where.Name, ".")
					s.msg.Observers.SendInfo("You see ", s.actor.Name, " take ", whereInventory.Name, " take ", where.Name, ".")
					return
				} else {
					s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
					return
				}
			}
		}
	}

	s.msg.Actor.SendBad("Get what?")
	s.ok = true
}
