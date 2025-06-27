package domain

type Player struct {
	Name          string `bson:"Name"`
	Age           int    `bson:"Age"`
	Photo         string `bson:"Photo"`
	Nationality   string `bson:"Nationality"`
	Flag          string `bson:"Flag"`
	Overall       int    `bson:"Overall"`
	Club          string `bson:"Club"`
	ClubLogo      string `bson:"Club Logo"`
	Value         int    `bson:"Value"`
	Special       int    `bson:"Special"`
	PreferredFoot string `bson:"Preferred Foot"`
	WeakFoot      int    `bson:"Weak Foot"`
	SkillMoves    int    `bson:"Skill Moves"`
	WorkRate      string `bson:"Work Rate"`
	RealFace      string `bson:"Real Face"`
	Position      string `bson:"Position"`
	Height        int    `bson:"Height"`
}
