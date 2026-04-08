package generator

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"blockchain-bot/internal/model"
)

type Generator struct{}

func (g *Generator) GenerateHash(data model.Data, time time.Time, payload string) string {
	r := data.PrevHash + strconv.FormatInt(data.BlockNumber, 10) + time.String() + payload

	h := sha256.New()

	h.Write([]byte(r))
	hash := h.Sum(nil)

	return hex.EncodeToString(hash)
}

func (g *Generator) GeneratePost(data model.Data, time time.Time, payload string) (string, *model.Data, error) {
	hash := g.GenerateHash(data, time, payload)
	s := "🧬 BLOCK #" + strconv.FormatInt(data.BlockNumber, 10) + "\n" +
		"⏰ Timestamp:" + time.String() + "\n" +
		"🔐 Hash:" + hash + "\n" +
		"🔗 Prev Hash:" + data.PrevHash + "\n" +
		"━━━━━━━━━━━━━━━━━━━━━━\n" +
		payload + "\n" +
		"━━━━━━━━━━━━━━━━━━━━━━\n" +
		"✔️Status: Mining | Difficulty: High |Miner: [БL0KЧ3ЙНИ5Т]"

	return s, &model.Data{
		BlockNumber: data.BlockNumber + 1,
		PrevHash:    hash,
	}, nil
}
