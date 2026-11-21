# Hockey manual and Beginners Guide: what MZ officially says

Sources: `docs/hockey-eval/sources/manual/` (40 chapters) and
`docs/hockey-eval/sources/beginners_guide/` (4 chapters), fetched 2026-11-21.
Every claim below is **official** (MZ's own text), **era: current site text**
(undated; it carries football leftovers, so read some lines with care).
Citations: [manual A05 Training], [guide manager_basics].

## Summary (≤15 bullets, the answer)

- **Ten trainable skills**, shown as 0–10 "skill pucks": Play Intelligence, Power, Skating, Passing,
  Quickness, Shooting, Keeping, Puck Control, Checking, Stamina. **Form** is untrainable and
  "has as much of an impact on their performances as the attributes". Experience is named but never defined.
- **The sim uses fractional skills.** The display rounds them down to quarter pucks (Partial Skill Bar), and
  every daily training change counts in matches. Our model should take fractional input.
- **Profile data**: age, shoots (handedness), height, weight, birth season, nationality, value, salary,
  Loyal flag. No source ties height, weight or handedness to the match engine.
- **Rink**: 6 zones per line: 1 goalkeeper, 2 defenders, 2 wingers, 1 centre. Players can be moved
  inside their zone. The guide's shirt types are Goalkeeper, Defender, Center, Forward.
- **Lines**: at least 10 skaters + 1 GK in "3 lines, including powerplay and boxplay lines". There is
  no text on line ice time, rotation, substitutions or backup goalies.
- **Tactic options described**: Formation (player placement), Mentality (Defensive/Balanced/Attacking),
  Pressing (Passive/Normal/Committed). The A04 intro names "Formation, Passing style and team Mentality",
  but **no section anywhere describes a Passing style** (probably left over from football; Pressing
  has a section the intro doesn't list).
- **The only official tactic-to-skill link**: Committed pressing "use[s] up more Stamina" and causes
  more fouls/penalties; Passive saves Stamina. Mentality changes positioning and risk. Its only
  skill-related hint is that Attacking makes defenders join attacks.
- **Position-to-skill hints** come from the camp packages (counts below). Forward camps: Shooting in 10/13.
  Defender camps: Checking in 14/14. Center camps: Puck Control and Power 6/11 each, Play Intelligence 5/11.
  Keeping is only in Wild Card camps (3/7), and **there is no goalie camp**.
- Skill definitions point the same way: Shooting gets you "into goal scoring positions"; Keeping is
  needed for GKs but "good goalkeepers need other attributes as well" (not named).
- **Aging**: Quickness and Stamina deteriorate from about 28–30, fastest at 34, at most −2 pucks each.
  Retirement can come from age 31 and is forced by 38.
- **Per-player match data** (calibration): match page shows position, penalty minutes, goals and a
  **Rating**. Full stats add shifts, shots, assists, points. League top-scorer table: GP, G, A, P, +/-, PIM,
  FOL, FOW, S, Sft, TOI. The Rating is "open ended, and could even be negative". It uses visible stats,
  team result, the player's actions **and his position**. MVP = highest Rating.
- **Youth scout report** names the likely position (centre/forward/defender/goal keeper) and
  "main strengths and weaknesses", **not maxings**. It exists for hockey youths from season 62 on.
- **No official star rating, team rating, or per-position score** appears in these sources. Player value
  is recalculated each season (salary follows value). MZ Rank = achievement points over 8 weeks.
- Form: playing 5–6 matches per gameweek keeps a player "match fit". Too many League/Friendly matches
  cause exhaustion and a Form drop. Injuries slow training. A clinic speeds recovery.
- Several numbers disagree between the manual and the guide (Form Package 7 vs 9, walkover threshold,
  tactic slots, youth arrival age, camp skill count). See "Open questions".

## Findings

### 1. Skills and attributes

Exact definitions. Manual [manual A05 Training] first, then guide [guide get_started] where it differs:

| Skill | Manual | Guide |
|---|---|---|
| Play Intelligence | "ability to find good positions on the ice, both when attacking and defending" | same |
| Power | "ability to make, and withstand, tackles" | same |
| Skating | "improves a Player's balance and helps them turn faster on the ice" | same |
| Passing | "more likely to retain possession … finding his team-mates with accurate passes" | "ability to find team-mates with accurate passes" |
| Quickness | "How fast a player can skate, both with and without the puck" | same |
| Shooting | "more likely to get into goal scoring positions, and to exploit them" | "ability to shoot the puck accurately with their stick" |
| Keeping | "Although good goalkeepers need other attributes as well, without high Keeping … less effective shot stoppers" | "ability to stop shots and anticipate an opposing Player's actions" |
| Puck Control | "more likely to trap incoming passes and … less likely to lose possession" | named "Puck Control / Technique": "control passes with their stick and skate with the puck" |
| Checking | "apply pressure to opponents and win the puck in 1 vs. 1 situations" | same |
| Stamina | "physical endurance … less likely to tire towards the end of games" | "How far a Player can skate before they become tired and their performance drops" |

- **Scale**: 0–10 pucks. The total "can increase or decrease at different rates depending on his
  training efficiency, injuries, age and how close he is to reaching his limit in an attribute
  (… maxing)" [manual A05 Training]. Each skill has a hidden per-player **limit** (max).
- **Fractions**: "In matches Players benefit from each part of their skill, not only the full skill
  pucks". The Partial Skill Bar is "rounded down to the nearest quarter" [manual A05 Training].
  "each daily change is reflected in a Player's match performances" [manual A05 Training; guide manager_basics].
- **Training boost**: moving a player to a new training zone gives "a short-term boost in that
  attribute" (small green puck) [manual A05 Training]. Its match effect is not stated.
- **Form**: "cannot be influenced by training". It rises with matches played and falls from exhaustion if
  he plays too many League and Friendly matches [manual A05 Training]. Form has "as much of an impact on
  their performances as the attributes" [guide get_started]. Profile icons: barbell = needs more
  games, sleeping = needs rest, happy face = right load [manual A03 Squad]. "Match fit" = 5–6 matches per
  gameweek [guide get_started]. The Form Package sets Form to **7** [manual A05 Training] or
  **9** [guide manager_basics], plus "match fit". This hints at a Form scale running to at least 9 (likely 0–10).
- **Experience**: A02/A03 link to a "Form and Experience" section of Training that **does not exist** in
  the A05 text. Friendlies "affect players' experience and form levels" [guide competitions]. World League
  teams "earn team money and experience" [manual B11 World League]. No scale, display or effect is given.
- **Other profile fields**: age, "how he shoots", height, weight, birth season, nationality, current
  market value, salary, goals and "cards" for the league season [manual A03 Squad]. The guide says
  "strongest Shoots" [guide get_started]. **None is tied to match performance.** Loyal Players
  (6–7 per new team) can't be sold [manual A03 Squad]. No special traits or personalities are mentioned.
- **Age**: +1 on day 1 of each season [manual A03 Squad]. Youths arrive at 15–16 [manual A08 Youth]
  ("16" [guide manager_basics]) and graduate after their 18th birthday. Retirement is announced from age 31
  and comes no later than 38 [manual A03 Squad].

### 2. Positions, lines, goalies

- "The Rink is divided into six zones, one for each player. There are zones for one goal keeper, two
  defenders, two wingers and one centre. After placing a player … you can fine-tune his position by
  moving him around within the zone" [manual A04 Tactics]. The guide names four areas (goalkeeper,
  defender, forward, center). A player dragged into an area "will automatically adopt its associated
  position to the best of their abilities" [guide manager_basics]. Wingers = "Forward" in guide/camps.
- "At least 10 players and 1 goal keeper has to be placed in 3 lines, including powerplay and boxplay
  lines. It is not possible to save an incorrect formation" [manual A04 Tactics]. It is unclear
  whether "3 lines" means 3 even-strength lines plus PP/BP lines, or 3 lines in total.
- Fewer than 11 eligible players (incl. 1 GK) = 5-0 walkover [manual A04 Tactics]. The guide says
  "nine or fewer eligible Players, or have gaps in your formations" [guide manager_basics].
- Auto Fill exists, but "players won't necessarily be placed in their best positions" [manual A04 Tactics].
  This implies a "best position" per player, but no definition is given.
- **Not described**: line ice-time shares, line rotation, PP/BP selection, backup goalie, substitutions,
  injuries during a match. Match stats do record per-player **shifts (Sft)** and **TOI**, so lines get
  different ice time [manual B02 League].

### 3. Tactic options (exhaustive from these sources)

- Intro: "There are three main elements to setting tactics: Formation, Passing style and team
  Mentality" [manual A04 Tactics]. **Passing style: mentioned only there. No section in the manual or
  the guide describes it, and it doesn't appear in the guide's tactic chapter.**
  Pressing has its own section but isn't in the intro list.
- **Formation**: placing players in zones and lines (above).
- **Mentality** (dropdown): Defensive / Balanced / Attacking. "affects their positioning and how they
  react to different in-game circumstances" [guide manager_basics].
  - Defensive: "more conservative and rarely attack. They focus on retaining the puck and not giving away
    too many scoring chances."
  - Balanced: "position themselves neutrally … If they are losing, they will be more likely to attack,
    but without taking excessive risks in defence."
  - Attacking: "Defenders will support attacks, even if it means leaving gaps in defence … more likely to
    create chances, but it also increases the likelihood of conceding" [guide manager_basics].
  - Manual: Defensive avoids committing "too many players forward and risk being caught on the counter
    attack" [manual A04 Tactics].
- **Pressing**: Passive / Normal / Committed. "A Committed approach means that your players will harass
  the puck carrier more but they will also use up more Stamina … and run the risk of committing more
  fouls" ("penalties" in the guide). "A Passive policy … will save on Stamina and reduce the chance of
  disciplinary problems but could give the opposition too much of an easy time" [manual A04 Tactics].
- Other: 3 tactic slots (manual: "Club Members can set up to ten"; guide: "a further six"). Tactics lock
  2 h before the match. XML import/export [manual A04 Tactics]. **No** aggression, line-time, PP/BP
  tactic, or shooting/passing-style options are described.
- MZ's own ranking of what matters: "The skills of your Players have the largest impact on results,
  but ineffective Tactics will lead to more losses than wins" [guide your_first_season].

### 4. Links tactic ↔ skill, position ↔ skill

- Pressing ↔ **Stamina** (Committed drains it), and Pressing ↔ **penalties** [manual A04 Tactics].
  Stamina's definition ("less likely to tire towards the end of games") makes it a late-game modifier.
- Mentality Attacking ↔ defenders join attacks. That suggests offensive skills (Shooting, Passing)
  matter more for D under Attacking. This is an inference, not stated.
- Position ↔ skill, stated directly: Keeping for goalkeepers, plus unnamed "other attributes"
  [manual A05 Training]. Nothing else is stated. The camps below are the strongest official hint.

### 5. Training and camps

- Training days: Mon, Tue, Thu, Fri, Sat (none on league days Wed/Sun). Camps train every day except Sunday
  [manual A05 Training; guide manager_basics]. Senior field holds 60, youth field 32. Youths improve faster.
- Coaches raise efficiency, more so when the coach's skill exceeds the player's. One coach per zone
  helps up to 7 seniors / 9 youths [manual A09 Staff]. "a Coach will always improve a Player, even if
  his skill level is lower" [guide manager_basics]. 7–8 meter segments with a coach = "a great trainer".
- Training speed differs per player and doesn't predict final level [guide manager_basics].
- Camp efficiency, weeks 1–5: 100/150/200/250/300% ("normal training efficiency x 2.0 … x 4.0").
  One regular camp per season, at most 15 players at once, 8 free youth places [manual A05 Training].
- **Camp packages: how often each skill appears** (from the tables in [manual A05 Training]):

| Camp (n) | Sho | Pow | Ska | Pas | PC | Qui | PI | Chk | Sta | Kee |
|---|---|---|---|---|---|---|---|---|---|---|
| Forward (13) | **10** | 5 | 5 | 5 | 4 | 4 | 4 | 1 | 1 | 0 |
| Center (11) | 2 | **6** | 1 | 3 | **6** | 3 | **5** | 4 | 2 | 0 |
| Defender (14) | 3 | 5 | 5 | 1 | 5 | 2 | 3 | **14** | 3 | 0 |
| Wild Card (7) | 0 | 2 | 3 | 1 | 0 | 4 | 1 | 0 | **6** | 3 |

  So MZ's view: forwards = Shooting first. Centers = Puck Control, Power, Play Intelligence, then
  Checking (a two-way role). Defenders = Checking always, then Power/Skating/Puck Control. No camp
  combines Keeping with Shooting or Checking. The GK-ish Wild Card packs pair Keeping with Stamina,
  Quickness, Play Intelligence and Passing.
- Youth Package: "Tries to pinpoint the most probable role … trains him in three useful skills for that
  role" [manual A05 Training]. The guide says youth packages contain "only two attributes per package,
  which makes it easier to establish roles" [guide manager_basics].
- **Deterioration**: from about 30 (a heading says "Speed and Stamina", the body says Quickness and Stamina).
  "Players who start to deteriorate at 28 will do so slower". The rate peaks at 34, at most −2 pucks per
  attribute. Each training day one of the two is picked at random. Training Quickness/Stamina then
  "will not help" [manual A05 Training].
- Maxing: each skill's limit is hidden. The youth scout report gives no hint of it [manual A10 Youth Exchange].

### 6. Ratings, values, reports

- **Match Rating**: "based on visible statistics and the formulas … takes into account both how well
  the team did in general, the actions each Player undertook, and the position of the Player. The
  Rating is open ended, and could even be negative". The player with the highest Rating is MVP
  [manual B01 Fixtures and Results]. This is a good calibration target: one number per player per match,
  adjusted for position.
- **Per-player stats**: Match Page lists position, penalty minutes, goals, Rating. "Full Match
  Statistics" list shifts, shots, assists, points per player and team [manual B01]. The league/youth top-scorer
  table has GP, G, A, P, +/- (on-ice goal difference), PIM, FOL, FOW (face-offs, so centres), S, Sft,
  TOI [manual B02 League; B03 Youth Leagues]. Club Members get a Squad Stats screen with averages [manual A03 Squad].
  Opponent stats: the "%-symbol" next to a match shows "how each of their players performed" [manual A04 Tactics].
- **Player value**: recalculated at season start (and on mid-season transfers, from the fee). Salary is
  based on value. Asking price ≤ 8× value [manual A03 Squad; A06 Transfers]. The formula is not given.
  It's a possible proxy for overall skill. Team value > €2.7M unlocks unlimited bids [manual A06 Transfers].
- **Youth scout report** (€872): "gives an idea of which position the player will be good at
  (centre/forward/defender/goal keeper), but does not give any indications of potential maxings"
  [manual A10 Youth Exchange]. The report "will pinpoint the main strengths and weaknesses". It exists for
  hockey youths from season 62 [guide manager_basics]. A "Scout report" icon appears on player profiles
  [guide get_started]. **No star levels or high/low scale are described.**
- **MZ Rank**: points from "a variety of different achievements … over the previous eight weeks"
  [manual C05 Ranking]. This is team-level and achievement-based, not a skill rating.
- **Not in sources**: star ratings, team strength rating, a per-position player rating, or a training
  "talent" indicator (beyond meter segments).

### 7. Injuries, fatigue, form, experience

- Injuries slow skill development ("depending on … injuries") [manual A05 Training]. A stadium clinic
  speeds "recovery time for your injured players" [manual C03 Stadium]. Friendly injuries "apply to all
  competitions" [guide competitions]. The Availability screen lists injuries and suspensions
  [manual A04 Tactics]. No match-time effect of injuries is described.
- Fatigue in matches: Stamina governs tiring "towards the end of games", and Committed pressing burns
  it faster [manual A05 Training; A04 Tactics].
- Form influence by match type: Senior/Youth/World League and Friendly can raise or lower Form. Cups and
  Friendly League can only raise it ("a player that have the status tired will not loose form"). Only the first
  two league matches count [manual A05 Training].
- Experience: its effect is undocumented (see §1).

## Open questions / contradictions

- Passing style: does it exist in hockey? It is only in the A04 intro sentence.
- "3 lines, including powerplay and boxplay lines": how many even-strength lines, and how ice time splits.
- Form Package target: Form 7 [manual A05] vs 9 [guide manager_basics]. Form and Experience scales are not given.
- Walkover: < 11 eligible [manual A04] vs ≤ 9 [guide manager_basics].
- Tactic slots for Club Members: 10 total [manual] vs 3 + 6 [guide].
- Youth Package skill count: three [manual A05] vs two [guide manager_basics].
- Deterioration start age: 30 vs "start to deteriorate at 28". The heading says "Speed", the body Quickness.
- Camp efficiency table: "100% (x2.0)" labels are inconsistent with each other.
- Leftovers that look like football: "cards" in the profile, "Puck Control / Technique", "pitch", "kick-off".

## Gaps the forum must answer

1. Which skills matter per position (GK, D, C, W), with any ranking or numbers. Above all, the GK's
   "other attributes" beyond Keeping.
2. Is there a Passing style in hockey? Full list of tactic-screen options (line ice time, PP/BP, aggression…).
3. How lines work: number of lines, ice-time split, PP/BP usage, what a weak 3rd line costs.
4. Pressing and Mentality effects on skill value (e.g. does Committed raise Checking's weight?).
5. Effect sizes of Form and Experience, and their scales. Does height, weight or handedness matter?
6. Hockey sim engine versions and dates, to tag eras.
7. The Rating formula per position, and how faceoffs (centres) and +/- feed into it.
8. Are there youth scout report star levels / high-low skills in hockey, and how reliable are they?
9. Effect of training boost and partial pucks in matches. Is the value formula skill-weighted by position?

## Queries used

- Read in full: manual A02, A03, A04, A05, A08, A09, A10, B01, C05, C06. Guide: all 4 chapters.
- Skimmed: B02/B03 (stats tables), A06 (value), C03 (clinic), B10/B11.
- Whole-tree grep for: passing style, powerplay, boxplay, star, experience, ice time, shift, lines, wing,
  centre, aggress, height, shoot, scout, injur, clinic.
- Camp skill counts: python over the A05 tables (Defender/Wild Card hand-checked).
