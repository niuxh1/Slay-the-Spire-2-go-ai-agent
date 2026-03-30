package models

import "encoding/json"

// Response 是 API 的通用响应格式
type Response struct {
	Ok        bool            `json:"ok"`
	RequestID string          `json:"request_id"`
	Data      json.RawMessage `json:"data"`
	Error     *APIError       `json:"error,omitempty"`
}

type APIError struct {
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Details   json.RawMessage `json:"details"`
	Retryable bool            `json:"retryable"`
}

// GameState 代表 /state 返回的完整快照
type GameState struct {
	StateVersion     int             `json:"state_version"`
	RunID            string          `json:"run_id"`
	Screen           string          `json:"screen"`
	InCombat         bool            `json:"in_combat"`
	Turn             *int            `json:"turn"`
	AvailableActions []string        `json:"available_actions"`
	Combat           *CombatState    `json:"combat"`
	Run              *RunState       `json:"run"`
	Map              *MapState       `json:"map"`
	Reward           *RewardState    `json:"reward"`
	Selection        *SelectionState `json:"selection"`
	Event            *EventState     `json:"event"`
	Shop             *ShopState      `json:"shop"`
	Rest             *RestState      `json:"rest"`
	Chest            *ChestState     `json:"chest"`
}

// CombatState 战斗状态
type CombatState struct {
	Player  PlayerState   `json:"player"`
	Hand    []Card        `json:"hand"`
	Enemies []EnemyState  `json:"enemies"`
}

type PlayerState struct {
	CurrentHP int     `json:"current_hp"`
	MaxHP     int     `json:"max_hp"`
	Block     int     `json:"block"`
	Energy    int     `json:"energy"`
	Stars     int     `json:"stars"`
	Powers    []Power `json:"powers"`
}

type Card struct {
	Index            int     `json:"index"`
	CardID           string  `json:"card_id"`
	Name             string  `json:"name"`
	Upgraded         bool    `json:"upgraded"`
	TargetType       string  `json:"target_type"`
	RequiresTarget   bool    `json:"requires_target"`
	EnergyCost       int     `json:"energy_cost"`
	StarCost         int     `json:"star_cost"`
	Playable         bool    `json:"playable"`
	UnplayableReason *string `json:"unplayable_reason"`
}

type EnemyState struct {
	Index     int      `json:"index"`
	EnemyID   string   `json:"enemy_id"`
	Name      string   `json:"name"`
	CurrentHP int      `json:"current_hp"`
	MaxHP     int      `json:"max_hp"`
	Block     int      `json:"block"`
	IsAlive   bool     `json:"is_alive"`
	IsHittable bool    `json:"is_hittable"`
	Powers    []Power  `json:"powers"`
	Intents   []Intent `json:"intents"`
}

type Power struct {
	Index   int    `json:"index"`
	PowerID string `json:"power_id"`
	Name    string `json:"name"`
	Amount  *int   `json:"amount"`
	IsDebuff bool  `json:"is_debuff"`
}

type Intent struct {
	Index       int    `json:"index"`
	IntentType  string `json:"intent_type"`
	Label       string `json:"label"`
	TotalDamage *int   `json:"total_damage"`
}

// RunState 局内全局状态
type RunState struct {
	Floor     int      `json:"floor"`
	CurrentHP int      `json:"current_hp"`
	MaxHP     int      `json:"max_hp"`
	Gold      int      `json:"gold"`
	MaxEnergy int      `json:"max_energy"`
	Deck      []Card   `json:"deck"`
	Relics    []Relic  `json:"relics"`
	Potions   []Potion `json:"potions"`
}

type Relic struct {
	Index   int    `json:"index"`
	RelicID string `json:"relic_id"`
	Name    string `json:"name"`
}

type Potion struct {
	Index    int    `json:"index"`
	PotionID string `json:"potion_id"`
	Name     string `json:"name"`
	Occupied bool   `json:"occupied"`
	CanUse   bool   `json:"can_use"`
}

// MapState 地图状态
type MapState struct {
	AvailableNodes []MapNode  `json:"available_nodes"`
	Nodes          []FullNode `json:"nodes"`
}

type MapNode struct {
	Index    int    `json:"index"`
	Row      int    `json:"row"`
	Col      int    `json:"col"`
	NodeType string `json:"node_type"`
}

type FullNode struct {
	Row      int           `json:"row"`
	Col      int           `json:"col"`
	NodeType string        `json:"node_type"`
	State    string        `json:"state"`
	Children []Coordinates `json:"children"`
}

type Coordinates struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// RewardState 奖励状态
type RewardState struct {
	PendingCardChoice bool         `json:"pending_card_choice"`
	Rewards           []RewardItem `json:"rewards"`
	CardOptions       []Card       `json:"card_options"`
}

type RewardItem struct {
	Index       int    `json:"index"`
	RewardType  string `json:"reward_type"`
	Description string `json:"description"`
	Claimable   bool   `json:"claimable"`
}

// 其他状态简化定义
type SelectionState struct {
	Prompt string `json:"prompt"`
	Cards  []Card `json:"cards"`
}

type EventState struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Options     []EventOption `json:"options"`
}

type EventOption struct {
	Index       int    `json:"index"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ShopState struct {
	IsOpen bool       `json:"is_open"`
	Cards  []ShopItem `json:"cards"`
}

type ShopItem struct {
	Index     int    `json:"index"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Available bool   `json:"available"`
}

type RestState struct {
	Options []RestOption `json:"options"`
}

type RestOption struct {
	Index     int    `json:"index"`
	OptionID  string `json:"option_id"`
	Title     string `json:"title"`
	IsEnabled bool   `json:"is_enabled"`
}

type ChestState struct {
	IsOpened     bool    `json:"is_opened"`
	RelicOptions []Relic `json:"relic_options"`
}
