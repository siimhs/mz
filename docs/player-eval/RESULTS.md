# Player evaluation – results (v2)

v2 adds the round-2 title sweep (`research/06-sweep-*.md`) over all 10,928
topics, 2011-02 … 2026-11, from forums 10 and 12. v1 used only `research/01…05`
(2011–12).

**No source gives numeric skill weights**: not the manual, not Crew, not posters.
The numbers below are ours. They convert rankings, polls, release notes and
posters' tests into weights, and must be calibrated on match data (see "Next steps").

When sources disagree, this order decides: Crew release notes > current manual >
post-2016 consensus > 2012–16 consensus > 2011 posts > single opinions.

## Key facts

| Fact | Confidence | Source |
|---|---|---|
| Current engine is **sim V2.40** (live 2016-04-13), an update of the 2011 S40 2.x sim. Posters see no major engine change since | official + consensus | 06-11081198, 06-12109520 |
| 11 skills, 0–10; **fractional values count** | official | 01 |
| Field has only GK/DEF/MID/ATT thirds. Position comes from placement: wide defenders mark wingers, and mid-third players press differently from defenders | official | 01, 02 |
| Passing styles: **Short, Long balls, Wings** only | official | 01, 04 |
| **Pressing = marking**: Passive, nobody man-marks; Normal, defenders man-mark; Committed, defenders + midfielders man-mark. "A slower player cannot catch up to a faster player." | official (V2.29) | 06-11081198 |
| **AP decides whether an aerial pass succeeds**, for every position incl. GK goal kicks. A failed aerial pass goes to the opponent; pass length doesn't matter | official (V2.32/2.40) + consensus | 06-11081198 |
| **Heading**: helps strikers under Wings; near-useless for defenders | consensus (3 eras) | 06-9942556, 06-11081198, 06-12109520 |
| Wings dominant 2013–2024, near balance with Short by 2026. **Long balls weak** in every era ("never really worked") | consensus | all 06 |
| **Experience**: Crew said reduced (2011); every later era says it is one of the strongest factors (seniors need Exp 9–10) | contested → modern consensus | 05, 06-* |
| **Form**: only the form of the players on the pitch counts (team-form graph does not); posters avoid fielding form < 9 | Crew second-hand + opinion | 06-9942556, 06-12109520 |
| **Stamina** lowers all skills as the match goes on; Committed pressing needs about 8–9 | consensus | 06-9942556, 06-11081198 |
| Overall quality beats single skills: 10 vs 8 Tackling won 49% vs 48% of tackles; 10 vs 8 Shooting showed no difference against good keepers | posters' tests, small samples | 06-9942556, 06-12109520 |
| Tackle success tracks **PI**, not BC | consensus | 06-9942556 |
| No effect: Set Plays (barely), height, weight, foot, captain | Crew (foot/height/weight) + consensus | 05, 06-9942556 |
| Since 2024-11-25 match pages show passes, failed passes, tackles, assists, saves, distance and avg speed per player | official | 06-12109520 |

## Positions: 7 (GK, CD, FB, DM, CM, W, ST)

Your 5 positions, with Def split into CD + FB and DM split out of Mid:

| Code | Meaning | Why |
|---|---|---|
| GK | goalkeeper | – |
| CD | central defender | man-marks (Normal/Committed), tackles, intercepts |
| FB | wide defender | marks wingers; needs more Speed; idle against narrow opponents |
| DM | holding midfielder | defender-like skills, but sits in the mid third (presses/marks as mid) |
| CM | central/attacking mid | generalist; PI and BC above Passing (2018 poll) |
| W  | wide mid or wide forward (must be *really* wide) | the crosser in Wings |
| ST | central forward | – |

FB and DM can be merged into CD and CM if 5 positions are wanted: their
weights are close (FB ≈ CD with more Speed, DM ≈ CM with Tackling).

## Base weights (each row sums to 1 → 0–10 scale)

SPD speed, STA stamina, PI play intelligence, PAS passing, SHO shooting,
HEA heading, KEE keeping, BC ball control, TAC tackling, AP aerial passing.
Set Plays goes only into `setPieceScore`.

| | SPD | STA | PI | PAS | SHO | HEA | KEE | BC | TAC | AP |
|---|---|---|---|---|---|---|---|---|---|---|
| GK | .05 | .11 | .12 | – | – | – | .58 | .07 | – | .07 |
| CD | .17 | .15 | .18 | .07 | – | .02 | – | .10 | .24 | .07 |
| FB | .20 | .16 | .14 | .08 | – | .02 | – | .12 | .20 | .08 |
| DM | .13 | .14 | .18 | .12 | .04 | – | – | .11 | .20 | .08 |
| CM | .14 | .14 | .20 | .16 | .03 | – | – | .17 | .06 | .10 |
| W  | .20 | .16 | .10 | .12 | .07 | – | – | .17 | – | .18 |
| ST | .15 | .15 | .12 | .07 | .19 | .10 | – | .19 | – | .03 |

Changes vs v1, and why:
- **GK**: PI and Stamina up, Speed and BC down, AP up (goal kicks are aerial
  passes). The MZ Assistant tool uses KEE .55 plus .09 each, which is close.
- **CD**: Heading .05 → .02 (consensus in three eras); PI up (tackles follow PI).
  Speed kept, since the official marking mechanic makes it matter.
- **CM**: PI is #1 and BC above Passing (2018 poll; PI top-3 for every
  non-attacker).
- **ST**: Shooting .22 → .19: gains above 8 are small and capped by good keepers.

## Tactic modifiers

Multiply, then re-normalise each row to sum 1.

| Tactic | GK | CD / FB | DM / CM | W | ST |
|---|---|---|---|---|---|
| **Short** | AP×0.8 | PAS×1.2, AP×0.8 (FB also SPD×0.9) | PAS×1.3, BC×1.1, AP×0.7 | PAS×1.2, AP×0.7 | PAS×1.2, BC×1.1, HEA×0.6 |
| **Long balls** ⚠ weak in practice | AP×1.25 | AP×1.3 (FB ×1.2) | AP×1.15, PAS×0.9 | AP×1.1 | SPD×1.15, BC×1.1 |
| **Wings** | AP×1.2 | AP×1.3 | AP×1.3, PAS×0.85 | AP×1.4, STA×1.1 | HEA×1.8, SHO×0.9 |
| Pressing **passive** | – | SPD×0.9, STA×0.85 | STA×0.85 | STA×0.85 | STA×0.85 |
| Pressing **normal** | – | (baseline: defenders man-mark) | – | – | – |
| Pressing **committed** | STA×1.1 | SPD×1.1, STA×1.2 | SPD×1.15, TAC×1.15, STA×1.2 | SPD×1.15, STA×1.2 | STA×1.2 |
| Opponent plays **Wings** (optional) | – | TAC×1.1, PI×1.1, FB SPD×1.1 | DM TAC×1.2 | – | – |

Mentality: no modifiers. No source shows a skill effect; Crew only describes
how much risk defenders take.

**Your hypothesis**: long passes do use Aerial Passing (AP decides whether an
aerial pass succeeds), and Long balls shifts Def/Mid toward AP and away from
midfield Passing. But Long balls has been the weak style in every era, so its
modifiers are kept small. Where your idea really applies is **Wings**: it
raises AP's weight for the whole team, not just the wingers.

## Pseudo code

```
// skills fractional 0..10
evaluate(player, pos, tactic, opponent=null) -> double:
    w = BASE[pos]
    for mod in modifiers(pos, tactic.passing, tactic.pressing, opponent?.passing):
        w[mod.skill] *= mod.factor
    w = normalise(w)

    s = player.skills
    if pos == GK:
        s.KEE = max(s.KEE, 4)                 // Crew: keeper floor
        s.STA = min(s.STA, 7)                 // stamina saturates for keepers

    posScore  = Σ w[k] * shape(s[k])                        // 0..10
    coreScore = mean(shape(s.SPD), shape(s.STA), shape(s.PI), shape(s.BC))
    skill = pos == GK ? posScore
                      : 0.8 * posScore + 0.2 * coreScore    // "overall quality" term

    return skill * fatigue(s.STA, pos, tactic.pressing)
                 * condition(player.form, player.exp)

shape(x) = 10 * (x/10)^0.8            // diminishing returns: 8→10 is small

fatigue(sta, pos, pressing):
    T = {GK:6, CD:8, FB:8, DM:8, CM:8, W:8, ST:7}[pos]
    if pressing == committed: T += 1
    return 1 - 0.04 * max(0, T - sta)  // sta 4 as CD → 0.84

condition(form, exp):                  // present strength; use 1.0 for long-term/transfer views
    mf = 1 - 0.05 * max(0, 9 - form)   // form 7 → 0.90
    me = 0.80 + 0.02 * exp             // exp 0 → 0.80, exp 10 → 1.00
    return mf * me

// one set-piece taker, not part of the position score (low confidence)
setPieceScore(p) = 0.35*SP + 0.35*AP + 0.30*SHO

bestPosition(player, tactic) = argmax over pos of evaluate(player, pos, tactic)
```

All constants (0.8/0.2, 0.04, 0.05, 0.80 + 0.02·exp, exponent 0.8) are
placeholders, ordered by evidence but not measured.

## Next steps

1. **Implement** in Go next to the players scraper, with weights and modifiers
   in a JSON table.
2. **Calibrate** with the per-player match stats (since 2024-11): fit tackles
   won vs TAC/PI/SPD, pass success vs PAS/AP by passing style, saves vs KEE/PI,
   and goals per shot vs SHO against the opposing KEE. Scrape our own and
   opponents' matches.
3. Test the open hypotheses on that data: the total-balls/core term, the size
   of Experience, and whether Speed flattens above 8.

## Still open

- **Experience**: how big it is, and whether it counts per player or per team
  average (modern consensus says big; Crew 2011 said reduced).
- **Speed**: the poll ranks it top-tier and the marking mechanic makes it
  matter, but the 2024 avg-speed stat shows little gap between 6 and 10.
- **Shooting**: a high rank, but low-SHO midfielders convert about as well.
- **Set Plays**: barely any effect seen; penalties mostly saved after 2024
  (unconfirmed engine tweak).
- Formation limits (how many defenders, midfielders or attackers trigger a
  penalty), and which skills reduce stunning after a failed tackle.
