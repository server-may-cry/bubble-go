package models

type Test struct {
	Ttt int
}

type User struct {
	id                      int
	sysId                   uint8
	extId                   string
	reachedStage01          uint8
	reachedSubStage01       uint8
	ignoreSavePointBlock    bool
	inifinityExtra00        bool
	inifinityExtra01        bool
	inifinityExtra02        bool
	inifinityExtra03        bool
	inifinityExtra04        bool
	inifinityExtra05        bool
	inifinityExtra06        bool
	inifinityExtra07        bool
	inifinityExtra08        bool
	inifinityExtra09        bool
	remainingTries          uint8
	restoreTriesAt          int
	credits                 uint16
	friendsBonusCreditsTime int
}
