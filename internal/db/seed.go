package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"math/rand"

	"github.com/guisithos/go-social/internal/store"
)

var usernames = []string{
	"aaron", "blake", "cindy", "dylan", "ella", "freddie", "gabby", "holly",
	"isaac", "jackie", "kyle", "lily", "mason", "nora", "owen", "paula",
	"quincy", "ruby", "sam", "tina", "ulysses", "vincent", "willa", "xavi",
	"yara", "zane", "andrea", "brad", "chloe", "dean", "elsa", "felix",
	"gwen", "harvey", "iris", "james", "kim", "logan", "mia", "naomi",
	"oswald", "phoebe", "quinnley", "ryan", "sara", "tony", "ursel", "veronica",
	"will", "ximena", "yosef", "zara",
}

var titles = []string{
	"The Sorcerer's Awakening", "Forging the Dragon Blade", "Mysteries of the Elven Forest",
	"Surviving the Darklands", "The Wizard's Quest", "Tales from the Dwarven Halls",
	"Secrets of the Arcane", "Journeys in the Shadow Realm", "The Ranger's Path",
	"Alchemy and Potions Mastery", "Battles of the Forgotten Kingdom", "The Bard's Ballad",
	"Runes of Power", "Ancient Dungeons Unveiled", "Taming the Mystic Beasts",
	"The Rogue's Gambit", "The Art of Spellcraft", "Adventures in the Astral Plane",
	"Crafting Legendary Weapons", "The Shaman's Vision", "Exploring the Enchanted Isles",
}

var contents = []string{
	"Discover the path of a young sorcerer unlocking their hidden powers and the challenges they must overcome.",
	"Learn the ancient art of crafting legendary weapons and the journey to create the mythical Dragon Blade.",
	"Uncover the secrets hidden deep within the enchanted Elven Forest, where magic and mystery abound.",
	"Explore strategies and tips for navigating the perilous Darklands and surviving its deadly threats.",
	"Join a wizard's epic quest to master powerful spells and battle dark forces threatening the realm.",
	"Delve into the rich history and lore of the dwarven halls, where ancient treasures await discovery.",
	"Unlock the arcane secrets of magic and how it shapes the world in this in-depth exploration of spellcraft.",
	"Journey into the Shadow Realm, a place of darkness and danger, where only the brave dare to tread.",
	"Follow the ranger’s path as they protect the wilds, tracking foes and mastering survival in untamed lands.",
	"Master the art of potion-making and alchemy, discovering rare ingredients and powerful brews.",
	"Relive epic battles and heroic tales from the forgotten kingdom, where legends were born.",
	"Embark on a musical adventure with the bard, learning the power of storytelling through song.",
	"Harness the power of ancient runes to unlock magical abilities and reshape destiny.",
	"Explore the depths of forgotten dungeons filled with traps, treasures, and long-lost relics.",
	"Learn how to tame and bond with mystic beasts, unlocking their strength and loyalty for your quests.",
	"Step into the shadows and embrace the rogue’s gambit, using cunning and stealth to outwit enemies.",
	"Discover the intricacies of spellcrafting, from beginner enchantments to powerful arcane rituals.",
	"Venture into the astral plane, where physical and magical worlds collide in a battle for supremacy.",
	"Learn the skills needed to craft legendary weapons, from gathering materials to perfecting your technique.",
	"Tap into the shaman’s vision and uncover hidden truths about the spiritual world and its influence.",
	"Explore the enchanted isles, a realm of beauty and danger, where powerful magic lies dormant.",
}

var tags = []string{
	"Magic Mastery", "Ancient Artifacts", "Spellcraft", "Enchanted Realms", "Arcane Power",
	"Hero's Journey", "Dungeon Exploration", "Quests", "Epic Battles", "Legendary Weapons",
	"Beast Taming", "Mystical Creatures", "Wilderness Survival", "Darklands", "Shadow Realm",
	"Guild Strategy", "Alchemy", "Potion Crafting", "Runes of Power", "Sorcery",
	"Elven Mysteries", "Dwarven Lore", "Mythic Lore", "Kingdoms and Castles", "Fantasy Warfare",
	"Dragon Slaying", "Rogue Tactics", "Shaman Wisdom", "Bard's Tales", "Astral Adventures",
	"Ancient Runes", "Magical Beasts", "Forgotten Kingdoms", "Ranger Skills", "Mystic Vision",
}

var comments = []string{
	"Awesome journey! Thanks for taking us through it.",
	"I love the idea of forging a legendary weapon!",
	"Interesting take on the Elven Forest, never thought about it that way.",
	"Great survival tips for the Darklands, thanks for sharing!",
	"I really enjoyed reading about the wizard's quest, very inspiring.",
	"This exploration of dwarven halls is fascinating, thanks!",
	"Thanks for diving deep into the mysteries of arcane magic!",
	"The Shadow Realm sounds intense—great advice on navigating it.",
	"Fantastic insight into the ranger's path. Very useful for adventurers!",
	"Thanks for the potion tips, I’ll definitely be trying those!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
