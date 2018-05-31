package application

type gameConfig struct {
	DefaultRemainingTries         int8
	IntervalTriesRestoration      int
	FriendsBonusCreditsMultiplier int
	DefaultCredits                platformsBonus
	InitProgress                  [7][]int8
}

type platformsBonus struct {
	Vk int
	Ok int
}

var defaultConfig = gameConfig{
	DefaultCredits: platformsBonus{
		Vk: 1000,
		Ok: 3000,
	},
	DefaultRemainingTries:         5,
	IntervalTriesRestoration:      1800,
	FriendsBonusCreditsMultiplier: 40,
	InitProgress: [7][]int8{
		{-1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	},
}
