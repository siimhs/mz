# Hockey player evaluation – results (v1)

v1 is based on `research/01-manual.md` (manual and Beginners Guide) and the
title sweep `research/02-sweep-*.md`. The sweep covers all 1,191 hockey topics
in forums 138 and 140 (11.7k posts, 2009-12 … 2026-11) in four topic-id
quartiles. Cite as `01`, `7848005` (2009–12), `10149335` (2012–14),
`11255602` (2014–18), `12340151` (2018–26).

**No source gives numeric skill weights**: not the manual, not Crew, not
posters. Third-party tools (Demerzel's DZMZHV, MZ Manager) never published
their formulas. The one poster formula found (nav, 2013, `10149335`) was
rejected by the replies. **Every number below is ours.** The numbers turn
rankings, thresholds and one Crew statement into weights, and must be
calibrated on match data (see "Next steps").

When sources disagree, this order decides: Crew posts and release notes >
current manual > post-2015 (sim 1.2) consensus > 2009–15 consensus > single
opinions.

## Key facts

| Fact | Confidence | Source |
|---|---|---|
| Ten trainable skills, 0–10 pucks: Play Intelligence, Power, Skating, Passing, Quickness, Shooting, Keeping, Puck Control, Checking, Stamina. Plus Form. Experience is not shown (maybe hidden, effect unknown) | official / contested (Exp) | 01, 12340151 |
| **Partial skills count** in the sim ("Players benefit from each part of the skill"); the display rounds down to quarter pucks | official (news 2015-04) | 01, 11255602 |
| Engine: 1.0 (to ~2009) → 1.1 → **1.2** (Crew beta 2015-06-17, live ~S47, mid-2015) → stamina/line-change update (~2015–16) → away-boost fix (2020-10-06). Release notes for 1.2 were not found. No reported change to skills since then: the "no changes" guide thread was re-confirmed in 2020, 2022 and 2023 | official (dates) + consensus | 11255602, 12340151 |
| **Power is by far the most important skater skill**: "a flaw in the game that we … will have to live with" (powdersnow, Crew, 2019-02-13) | **official** + consensus in every era | 12340151 [44719547], all |
| Power thresholds: < 8 is "disposable", 8 is the floor, top teams want 9–10 | consensus (all eras) | all 02 |
| Stamina is #2 overall: P8/C8/S8 beats P7/C10/S6, and 9 Pow/7 Sta beats 10/6 | consensus | 10149335, 11255602 |
| **Skating is the only counter to Power** (absorbs or avoids body checks); P8/Sk8 beats P9/Sk4 | consensus | 10149335, 11255602, 12340151 |
| Face-offs: **Puck Control** is the main skill, Power is not a factor; PI helps | **official** (mikelony, Crew, 2018) + consensus | 12340151 [44384899], 7848005 |
| Passing ≈ useless for skaters ("pass fine with 0"); PI helps a little (fewer wasted chances) | consensus, a few dissent | all 02 |
| GK: Keeping ≫ Stamina > PI > Quickness > Passing/PC/Skating; Power, Shooting and Checking useless; keep only keepers who reach Keeping 10 | consensus (moderator guides) | all 02 |
| Left/right (LW/RW, LD/RD) and stick hand are cosmetic | consensus | 10149335, 12340151 |
| Height, weight: no effect. No injuries in hockey | consensus | 7848005, 12340151 |
| Training a skill on match day plays it as s + 0.1·(10 − s) (9 → 9.1) | official-ish (moderator relays changelog, unverified) | 7848005, 12340151 |
| Since 2015-04, ageing lowers only Quickness and Stamina (from 27–30, at most 2 pucks each) | official | 01, 11255602 |

## Positions: 4 (GK, D, C, W)

| Code | Meaning | Why |
|---|---|---|
| GK | goalkeeper | one per match, plays all 60 min; no backup, no 6th skater |
| D | defender (LD/RD) | Power + Checking; under Attacking, joins the attack |
| C | centre | takes face-offs (Puck Control); needs Shooting *and* Checking; "best player on the line" |
| W | winger (LW/RW) | Power + Shooting; under Defensive, tracks back and checks |

A line is **C, LW, RW, LD, RD**. The tactic has 21 slots: regular lines
L1–L4, powerplay PP1–PP2, boxplay BP1–BP2, and the GK. The minimum is 10
skaters + GK; 20 + GK is recommended. A player may appear in more than one
line, which costs him stamina.

## Base weights (each row sums to 1 → 0–10 scale)

POW power, STA stamina, SKA skating, QUI quickness, SHO shooting, CHK
checking, PC puck control, PI play intelligence, PAS passing, KEE keeping.

| | POW | STA | SKA | QUI | SHO | CHK | PC | PI | PAS | KEE |
|---|---|---|---|---|---|---|---|---|---|---|
| GK | – | .13 | .03 | .07 | – | – | .02 | .11 | .02 | .62 |
| D  | .28 | .15 | .12 | .11 | .06 | .18 | .03 | .05 | .02 | – |
| C  | .25 | .14 | .10 | .09 | .13 | .12 | .09 | .06 | .02 | – |
| W  | .28 | .15 | .12 | .12 | .18 | .05 | .03 | .05 | .02 | – |

Why:
- **POW** is the largest weight for every skater (Crew 2019). It is about 1.5× the
  role skill, not 2×, because the role skill is what separates D from W
  (Checking vs Shooting). C's Power is a little lower because C has to carry
  both role skills and PC.
- **STA, SKA, QUI** ("athletic core") come next and are nearly equal. Skating is
  slightly above Quickness, as in the current-era advice (scruttino
  2025–26); malkin78 (2019) put Quickness first, so this is contested.
- **PC** gets a medium weight only for C (face-offs, Crew). **PAS** gets a
  token .02 everywhere. **KEE** gets 0 for skaters (it inflates value, not
  play).
- **GK**: the order is consensus in all eras; the size of the Keeping gap
  (10 ≫ 8) justifies .62.

## Tactic modifiers

Hockey tactics are **set per line**: Mentality (Defensive / Balanced /
Attacking) and Pressing (Passive / Normal / Committed), plus where each player
sits inside his zone. **No passing style exists**: the manual's A04 intro names
one but no section describes it, and no forum post in 17 years mentions it.
There is also no aggression, ice-time or rotation setting, no substitution, no
pulled goalie and no Test Tactic. Placement inside the zone shows almost no
effect, so it is not modelled.

Multiply, then re-normalise each row to sum 1. All factors are ours. The
direction comes from the GPA tactics guide (agapitod/floodish 2016, topic
11868027, never contradicted since); the sizes are not measured.

| Line setting | D | C | W | GK |
|---|---|---|---|---|
| Mentality **Attacking** | SHO×1.6, QUI×1.1, STA×1.1, CHK×0.9 | SHO×1.1 | – | – |
| Mentality **Balanced** | (baseline) | | | |
| Mentality **Defensive** | CHK×1.1 | CHK×1.1 | CHK×2.0, QUI×1.1, STA×1.1, SHO×0.9 | – |
| Pressing **Passive** | STA×0.85 | STA×0.85 | STA×0.85 | – |
| Pressing **Normal** | (baseline) | | | |
| Pressing **Committed** | STA×1.25 | STA×1.25 | STA×1.25 | – |
| Slot **PP** (powerplay) | SHO×1.5 | SHO×1.3, PC×1.1 | SHO×1.3 | – |
| Slot **BP** (boxplay) | CHK×1.2, POW×1.1 | CHK×1.3, SHO×0.7 | CHK×1.5, SHO×0.7 | – |

Committed also draws more penalties, which is a team-level cost, not a skill
weight (manual + consensus).

## Pseudo code

```
// skills fractional 0..10 (partial skills count, official)
evaluate(player, pos, lineTactic, slot = L1, nLines = 4) -> double:
    w = BASE[pos]
    for mod in modifiers(pos, lineTactic.mentality, lineTactic.pressing, slot.kind):
        w[mod.skill] *= mod.factor
    w.STA *= {4: 1.0, 3: 1.3, 2: 1.6}[nLines]   // fewer lines → more ice time each
    w = normalise(w)

    s = player.skills                           // optionally s[trained] += 0.1*(10 - s[trained])
    score = Σ w[k] * shape(s[k])                // 0..10
    if pos != GK:
        score *= powerFloor(s.POW)
    return score * condition(player.form)

shape(x)      = 10 * (x/10)^0.8                // diminishing returns, as in football
powerFloor(p) = 1 - 0.04 * max(0, 8 - p)       // "below 8 is disposable": p 6 → 0.92
condition(f)  = 1 - 0.05 * max(0, 9 - f)       // form 7 → 0.90; use 1.0 for transfer views
// Experience: no factor (hidden; no evidence of effect)

bestPosition(player, tactic) = argmax over {D, C, W} (or GK) of evaluate(...)

// Lineup level: ice time is not a property of the player
lineStrength(line) = Σ over its 5 skaters of evaluate(...)
teamStrength(t)    = Σ_line share(line) * lineStrength(line)  + GK term
share(line): set by stamina. Since ~2015-16 the line with the most stamina left goes
    on at each stoppage, so lines of equal stamina get about equal time; one
    low-stamina player cuts his whole line's shifts. Pre-update rule of thumb: 21/18/13/8 min.
```

All constants (exponent 0.8, 0.04, 0.05, nLines factors, tactic factors) are
placeholders, ordered by evidence but not measured.

## Where hockey differs from football (matters for code)

- **Tactics are per line, not per team.** `tactic` must be a per-line struct
  (mentality, pressing), and the caller must pass the line the player is in.
- **The slot changes the weights**: the same player is worth different things
  on L1, PP1 and BP1. The signature gets `slot` (L1–L4, PP1–2, BP1–2). A player
  listed in two slots gets two scores, plus extra stamina drain (not modelled).
- **Ice time is a lineup property.** It depends on every line-mate's stamina,
  so it belongs in a team or lineup function, not in `evaluate`. Only
  `nLines` (how many regular lines are filled) goes into the player score, via
  the Stamina weight.
- **No passing style**, so the modifier table's main axis in football has no hockey
  equivalent. Mentality *does* shift weights here; in football it did not.
- **One skill dominates.** Power needs a floor term, not just a large weight.
  Football had no comparable rule.
- **No injuries; skill decline hits only QUI and STA.** Useful for age
  projection, not the score.
- **Positions**: 4 (GK, D, C, W), not 7; there is no wide/central split beyond
  C vs W, and left and right are the same.

## Calibration data available

- Per player per match: G, A, P, +/- (incl. PP/PK goals, goalies), PIM,
  face-offs won/lost, shots, shifts, TOI, and an open-ended **Rating** (can be
  negative; the highest is MVP; formula undisclosed) [manual A03/B01].
- Per team per match: goals, shots, face-offs won and %, saves, penalties, and
  scoring chances (since about 2014).
- Old-sim baseline (2014): about 16 shots per team per match, 14–26 % scored.
  Matches are very random (lost 0–1 with 25–10 shots).

## Next steps

1. **Implement** in Go next to the football evaluator, with weights and
   modifiers in a JSON table keyed by sport.
2. **Calibrate**: fit face-off % vs centre PC/PI/CHK; shots and goals vs
   SHO/POW against GK KEE; TOI/shifts vs line stamina (this tests the rotation
   rule); and +/- and Rating vs the position scores. Scrape our own and
   opponents' matches.
3. Test the Power floor (is there a knee at 8?), and Skating vs Quickness.

## Still open

- **What 1.2 changed**: no release note found. The away boost (2020)
  changes results, not player value.
- **Mentality and placement**: real but small, or no effect? There is no Test
  Tactic in hockey to check.
- **Skating vs Quickness** order, and whether there is a **Power × Skating**
  interaction (a player's Power is worth less against a high-Skating opponent).
- **Face-offs**: Crew says Puck Control; one small 2013 sample (8 centres)
  said Checking correlates more (0.81 vs 0.48).
- **PI**: small effect or none for skaters; its size for keepers (third)
  rests on posters' preference, not tests.
- **Experience**: hidden or non-existent; effect unknown.
- **Fatigue**: does a tired player play worse, or only get fewer shifts?
  Contested.
