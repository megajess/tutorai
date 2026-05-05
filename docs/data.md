# Data

## Overview

The chatbot draws from three distinct data sources, each serving a different type of query.

---

## 1. Scryfall Card Data

**Source:** [Scryfall Bulk Data API](https://api.scryfall.com/bulk-data)  
**File:** Oracle Cards (`oracle_cards` type)  
**Size:** ~270MB JSON  
**Update frequency:** Daily

### What gets stored

**SQLite (structured fields — for filtering):**
- `name` — card name
- `mana_cost` — e.g. `{2}{U}{B}`
- `cmc` — converted mana cost
- `type_line` — e.g. `Legendary Creature — Zombie Warrior`
- `oracle_text` — rules text
- `colors` — array of color symbols
- `color_identity` — for Commander legality
- `legalities` — format legality map (standard, modern, commander, etc.)
- `prices_usd` — current market price in USD

**Chroma (embedded — for semantic search):**
- `oracle_text` is embedded and stored with card name as metadata
- Used for natural language queries like "draw a card when a creature dies"

### Retrieval strategy

1. Filter SQLite by format legality and color identity first
2. Apply price ceiling if budget was specified
3. Semantic search over Chroma using oracle text embeddings
4. Return top-k results ranked by relevance

---

## 2. WotC Comprehensive Rules

**Source:** [Wizards of the Coast](https://magic.wizards.com/en/rules)  
**Format:** Plain text  
**Size:** ~250KB

### Chunking strategy

The rules are numbered hierarchically (e.g. 702, 702.2, 702.2a). Chunks are assembled by top-level rule section so related sub-rules stay together. For example, all of rule 702.2 (Deathtouch and its sub-rules) becomes one chunk, rather than splitting each lettered sub-rule apart.

This means a question like "how does deathtouch interact with trample" retrieves the full deathtouch rule block, which contains everything needed to answer it.

### What gets stored

**Chroma:**
- Each chunk is embedded with the rule number range as metadata
- e.g. `{"rule_range": "702.2–702.2b", "section": "Deathtouch"}`

---

## 3. Slang & Terminology Glossary

**Source:** Hand-curated (`data/slang_glossary.json`)

Community language that doesn't appear in the rules or oracle text but is essential for natural conversation.

### Color Identity

```json
{
  "term": "golgari",
  "colors": ["B", "G"],
  "aliases": ["bg", "black green", "green black"]
}
```

Guild names, shard names, and wedge names are stored as a lookup table — not embedded — since they map deterministically to color identities. Checked before any vector retrieval.

**Guilds (2-color):** Azorius (WU), Dimir (UB), Rakdos (BR), Gruul (RG), Selesnya (GW), Orzhov (WB), Izzet (UR), Golgari (BG), Boros (RW), Simic (GU)

**Shards (3-color):** Bant (GWU), Esper (WUB), Grixis (UBR), Jund (BRG), Naya (RGW)

**Wedges (3-color):** Mardu (RWB), Temur (GUR), Abzan (WBG), Jeskai (URW), Sultai (BGU)

### Archetypes

```json
{
  "term": "aristocrats",
  "description": "A strategy based on sacrifice synergies, death triggers, and life drain",
  "oracle_text_signals": ["sacrifice", "dies", "when a creature you control dies", "drain"]
}
```

| Term | Description |
|---|---|
| Aristocrats | Sacrifice synergies, death triggers, drain |
| Stax | Resource denial, symmetrical hate pieces |
| Voltron | Single creature with equipment/auras for commander damage |
| Pillow Fort | Redirect damage, protection enchantments |
| Goodstuff | No synergy theme, just high-power individual cards |
| Combo | Specific card interactions that create an infinite loop or win condition |
| Tokens | Wide board of small creatures, anthem effects |

### Community Shorthand

Terms not in the rulebook that players use when asking for recommendations:

| Term | Maps to |
|---|---|
| Mana rock | Artifact that taps for mana |
| Ramp piece | Any card that accelerates mana production |
| Draw engine | Repeatable card draw effect |
| Board wipe / Wrath | Mass creature removal |
| Staple | Card that goes in almost every deck of its color |
| Tutor | Card that searches your library for another card |
| Value piece | Card that generates card advantage or resources |
| Hate piece | Card targeting a specific strategy or card type |

---

## Adding New Data

To add new glossary terms, edit `data/slang_glossary.json` and re-run:

```bash
python scripts/ingest_slang.py --refresh
```
