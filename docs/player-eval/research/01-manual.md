# Official sources: manual + Beginners Guide

Scope: only `docs/player-eval/sources/manual/*.txt` and
`docs/player-eval/sources/beginners_guide/*.txt`. Every claim here is
**official**, because these are MZ's own docs. Where I go beyond the text I
mark it *(inference)*. Citation keys: `[manual A05 Training]`,
`[BG get_started]`, and so on.

## Summary (the answer)

- There are **11 trainable skills**, each 0–10 balls: Speed, Stamina, Play
  Intelligence, Passing, Shooting, Heading, Keeping, Ball Control, Tackling,
  Aerial Passing, Set Plays. **Experience and Form** can't be trained
  [manual A05 Training].
- **The match engine uses fractional skill values**, not just whole balls.
  The partial bar is only a display, rounded down to the nearest quarter
  [manual A05 Training]. So an evaluator should use decimals when they're
  available.
- **Form and Experience "have as much of an impact on their performances as
  the attributes"** [BG get_started]. The manual says more softly that they
  "help to maximise a player's other attributes" [manual A05 Training]. This
  suggests multiplying (skills × form/exp factor), which is *(inference)*.
- **The tactics field has only thirds**: GK, defender, midfielder,
  attacker/striker. A player's role comes from which third he's dropped into.
  Inside a third he can be moved around freely [manual A04 Tactics]. The docs
  define **no separate wing, full-back or DM positions**. "Wings" exists only
  as a passing style ("attacking colleagues who operate on the wings"), so a
  winger is just a wide midfielder or attacker *(inference)*.
- The **position families the game itself uses** are Defender, Midfielder and
  Striker training camps, plus Keeping inside "Wild Card" camps. The youth
  scout report names forward, midfielder, defender and goal keeper
  [manual A05, A10].
- **No official per-position skill weights exist.** The best official proxy
  is what each position's training camps contain (counts below). Defender
  camps centre on Tackling (13/15). Striker camps centre on Shooting (10/13)
  and Ball Control (8/13). Midfielder camps centre on Play Intelligence
  (9/16), Passing (8/16) and Aerial Passing (9/16). Keeping appears only with
  Stamina, Play Intelligence, Speed and Passing.
- **Passing styles**: Short, Long balls, Wings [manual A04; BG manager_basics].
  Long balls = "longer ground and air passes to Players closest to the
  opposition's goal". Wings = crosses to the strikers, which "tests the
  opposition full backs and goalkeeper". The skill links (Aerial Passing,
  also called Crossing, and Heading) are implied by the skill definitions,
  not stated in the tactic text.
- **Mentality**: Defensive, Balanced, Attacking. It changes positioning.
  Under Attacking, "Defenders will support attacks" [BG manager_basics].
- **Pressing**: Passive, Normal, Committed. Committed uses more Stamina and
  causes more fouls; Passive saves Stamina [manual A04]. This makes Stamina
  matter more under Committed.
- **Roles**: penalty-taker order, free-kick-taker order and captain
  [manual A04]. Set Plays makes "teams" better at corners, free kicks and
  penalties [manual A05]. The manual never says what the captain does.
- **Stamina**: tired or exhausted players perform worse. Players regain some
  stamina at half time, and the same applies to goalkeepers. High Stamina
  also means fewer injuries [BG get_started; manual A05].
- **Age**: Speed and Stamina start to decline at 30, fastest at 34, by at
  most 2 balls each. Other skills can still be trained. Players retire at
  31–38 [manual A05, A03].
- **Post-match Rating**: open-ended and can be negative. It is computed from
  visible statistics, team result, the player's actions and **his
  position**. The MVP is the player with the highest rating
  [manual B01 Fixtures and Results]. There are **no star ratings** or
  overall-skill numbers for players anywhere in the official docs.
- **Scout report (youth exchange)**: shows which position the youth will be
  good at, and his main strengths and weaknesses. It explicitly gives **no
  indication of potential maxings** [manual A10; BG manager_basics].
- "The skills of your Players have the largest impact on results, but
  ineffective Tactics will lead to more losses than wins" [BG your_first_season].

## Findings

### 1. Skills: what the official docs say each one does

Both sources give definitions. The Beginners Guide (BG) versions are often
more specific about match behaviour.

| Skill | Manual A05 Training | BG get_started |
|---|---|---|
| Speed | "How fast a player can run, both with and without the ball." | "Faster Players will be able to intercept more passes if they see them in time." (So it combines with PI for defending.) |
| Stamina | "less likely to tire towards the end of games and less prone to injury." | "How far a Player can run before they are tired or exhausted. This includes goalkeepers. At half time, Players will regain some stamina. When Players become tired or exhausted their performance will drop." |
| Play Intelligence | "more likely to be in the right place at the right time, whether defending or attacking." | "ability to make footballing decisions and to avoid making tackles that result in a card. Intelligent Players will see more passes. Intelligence helps Players adapt to different roles on the field." |
| Passing | "more likely to retain possession … accurate passes." | Same text. |
| Shooting | "more likely to get into goal scoring positions and find the back of the net." | "ability to strike the ball accurately with their feet." |
| Heading | "comfortable when dealing with aerial balls." | "ability to pass or shoot accurately with their head." |
| Keeping | "Although good goalkeepers need other attributes as well, without high Keeping … less effective shot stoppers." | "ability to stop shots and anticipate an opposing Player's actions, and therefore make better saves." |
| Ball Control (BG: "Ball Control / Technique") | "more likely to trap incoming passes and less likely to lose possession." | "control incoming passes with their feet. **This is an important attribute for all Players.**" |
| Tackling | "more likely to win the ball … and less likely to foul them." | "apply pressure to a ball-carrying opponent and to win the ball in a 1 vs. 1 situation." |
| Aerial Passing (BG: "Aerial Passing / Crossing") | "finding team-mates with long, aerial balls from anywhere on the pitch." | "pass the ball through the air to a teammate and to cross the ball into the opposition's penalty area." |
| Set Plays | "Teams containing a player with good skill … more effective from corners, free kicks and penalties." | Same. |

Notes, all **official**:
- Ball Control is the only skill called important "for all Players"
  [BG get_started].
- Keeping alone isn't enough: "good goalkeepers need other attributes as
  well" [manual A05]. The docs don't say which ones. The Wild Card camps pair
  Keeping with Stamina, Play Intelligence, Speed and Passing.
- Play Intelligence affects several things: positioning in attack and
  defence, seeing passes, avoiding cards, and adapting to different roles
  [BG get_started]. *(inference)* That makes it a reasonable universal skill
  for every outfield position.
- Speed works with PI on interceptions [BG get_started].
- Tackling, PI and Experience all reduce fouls or bad tackles
  [manual A05; BG get_started].
- Set Plays is described at team level ("Teams containing a player"). It
  matters through the designated set-piece takers (see Roles). *(inference)*
  So it should count for only one or two players in a lineup, not as a
  general skill.

### 2. Skill scale, partial skills, maxing, age

- The scale is 0–10 balls. Balls can rise or fall depending on training
  efficiency, injuries, age and "maxing" (reaching his limit in an
  attribute) [manual A05].
- "In matches Players benefit from each part of their skill, not only the
  full skill balls." The partial bar is "rounded down to the nearest
  quarter". "The SIM recognise the partial skills" [manual A05]. Each daily
  training change shows up in match performance [manual A05; BG manager_basics].
- Deterioration: it starts at 30 and affects only Speed and Stamina. It
  peaks at 34. At most 2 balls are lost per attribute, so 4 in total. A
  player who just lost a ball from max is "still very close to the previous
  level" [manual A05].
- Age goes up on day 1 of each season. Players may retire from 31 and must
  by 38. Retiring players take on average 50% longer to recover from
  injuries [manual A03].
- Youths arrive at 15 or 16 and graduate after their 18th birthday. They
  train faster than seniors [manual A08; BG manager_basics]. Fast early
  training doesn't mean a higher ceiling [BG manager_basics].

### 3. Form, Experience, Stamina in matches

- Form and Experience "have as much of an impact on their performances as
  the attributes mentioned above" [BG get_started]. The manual says they
  "can both help to maximise a player's other attributes" [manual A05].
- Form goes up with matches played, on the pitch or on the bench. Too many
  League or Friendly matches lower it through exhaustion. Cups and Friendly
  Leagues can only raise it [manual A05]. Staying "match fit" takes 5–6
  matches per gameweek [BG get_started]. The squad icons show a barbell
  (needs games), a sleeping icon (needs rest) or a happy face (right amount)
  [manual A03].
- The docs never state the Form scale. The Form Package sets Form to **7**
  according to the manual [manual A05], but to **9** according to the BG
  [BG manager_basics]. That contradicts, but either way it implies a scale
  of roughly 0–10.
- Experience grows with games played and age. It "ensures a player is more
  comfortable on the pitch and less likely to make mistakes" [manual A05],
  for example by avoiding bad tackles [BG get_started].
- Stamina in matches: see the table above. Committed pressing uses more of
  it [manual A04].

### 4. Positions and the tactics field

- "The Field is divided into defensive, midfield and attacking thirds and
  any player you drag into them will automatically adopt that role (e.g.
  defender, midfielder, attacker). Any player you drop close to your goal
  will become your goalkeeper. You can then fine-tune players' positions by
  moving them around their third" [manual A04]. The BG says the same with
  the names goalkeeper, defender, midfielder and striker areas, and adds
  that a player "will automatically adopt its associated position to the
  best of their abilities" [BG manager_basics].
- The minimum valid formation is 1 GK, 2 DEF, 2 MID and 1 ATT. Eleven
  players are expected; seven are needed to save a tactic. Light green areas
  are off-limits, and players can't stand too close to each other
  [manual A04].
- The preset formations are 4-4-2, 4-3-3 and 5-3-2 [manual A04].
- The bench has 5 slots: GK reserve, defender reserve, then midfielder and
  striker reserves, and an extra. Up to 3 substitutions are allowed
  [manual A04].
- The docs never name a **wide position or a winger**. Lateral placement
  inside a third is free-form. *(inference)* Wide versus central is set by
  the x-coordinate within the third. Whether the engine treats wide players
  differently is unstated, apart from the Wings passing style sending the
  ball "to their attacking colleagues who operate on the wings".
- The scout report and the training camps use the same four families:
  goalkeeper, defender, midfielder, forward/striker [manual A05, A10].

### 5. Official proxy for position skill sets: training camp contents

Camps are "aimed at different Player positions" [BG manager_basics]. Below
is how often each skill appears across a position's camps [manual A05
Training camp tables]. These are counts, **not weights**. They show which
skills MZ thinks belong together for a position.

| Skill | Defender (15 camps) | Midfielder (16) | Striker (13) | Wild Card (7) |
|---|---|---|---|---|
| Tackling | **13** | 2 | 0 | 0 |
| Speed | 5 | 6 | 5 | 2 |
| Stamina | 4 | 4 | 3 | 5 |
| Play Intelligence | 3 | **9** | 3 | 4 |
| Passing | 4 | **8** | 4 | 1 |
| Aerial Passing | 4 | **9** | 2 | 0 |
| Heading | 3 | 1 | 3 | 0 |
| Ball Control | 2 | 5 | **8** | 3 |
| Shooting | 0 | 1 | **10** | 0 |
| Set Plays | 1 | 1 | 0 | 0 |
| Keeping | 0 | 0 | 0 | 3 |

What the counts show:
- **DEF**: Tackling dominates, then Speed, Stamina, Passing and Aerial
  Passing. Heading appears in 3 camps, all with Tackling. Shooting never
  appears.
- **MID**: the core is Play Intelligence, Passing and Aerial Passing, with
  Speed and Ball Control next. Tackling appears only twice and Shooting
  once.
- **ST**: the core is Shooting and Ball Control, with Speed next.
  Heading, Passing and PI each appear 3–4 times. Tackling never appears.
- **GK** (Wild Card 5–7): Keeping comes with Stamina, Play Intelligence,
  Speed and Passing.
- Set Plays appears only in combination camps: DEF 15 (Set Plays, Passing,
  Aerial Passing) and MID 16 (Heading, Set Plays, Passing).
- There is no winger camp. The closest are MID 12 (Aerial Passing, Ball
  Control, Speed) and ST 3 (Ball Control, Passing, Aerial Passing)
  *(inference)*.
- The Youth Package "Tries to pinpoint the most probable role for your Youth
  Player and trains him in three useful skills for that role" [manual A05].
  The BG says youth packages have two attributes each "to establish roles"
  [BG manager_basics]. Either way, MZ evaluates roles internally, but the
  docs don't say how.

### 6. Tactic options (as written)

The manual lists six elements: Formation, tactical substitutions, Passing
style, Mentality, Pressing and Roles. Rules come on top [manual A04].

**Passing** [manual A04; BG manager_basics]
- **Short**: "patient … good possession stats but sometimes at the expense
  of a cutting edge in front of goal" [manual]. The BG adds "more runs on
  and off the ball", "higher possession and shooting stats, but the
  positions from which the shots occur may not be ideal".
- **Long balls**: "knock Long balls up to the forwards (allowing your team
  to quickly get into shooting positions but increasing the risk of losing
  possession with less precise passes)" [manual]. The BG: "longer ground and
  air passes to Players closest to the opposition's goal".
- **Wings**: "channel their attacks down the Wings (which tests the
  opposition full backs and goalkeeper with crosses to your strikers)"
  [manual]. The BG: "prefer to pass to their attacking colleagues who
  operate on the wings".
- "it's up to you to gauge which most suits your players" [manual A04]. No
  skill is named for any style. *(inference, from the skill definitions)*
  - Short leans on Passing, Ball Control and PI.
  - Long leans on Aerial Passing, which is long aerial balls "from anywhere
    on the pitch". It also needs forwards' Speed, Heading and Ball Control,
    and opponents' Heading.
  - Wings leans on Aerial Passing in its "Crossing" role, the BG's name for
    the skill, plus the strikers' Heading and the wide players' Speed.

**Mentality** [manual A04; BG manager_basics]
- **Defensive**: "more conservative and rarely attack. They focus on
  retaining the ball, and keep more Players between their goal and the
  ball."
- **Balanced**: neutral positioning. More likely to attack when losing,
  "without taking excessive risks".
- **Attacking**: "Defenders will support attacks, even if it means leaving
  gaps in defence … more likely to create chances … increases the likelihood
  of conceding." *(inference)* This could add some attacking value to
  defenders (Passing, Ball Control, Speed).

**Pressing** [manual A04; BG manager_basics]
- **Committed**: harasses the ball carrier more, uses more Stamina, more
  fouls.
- **Normal**.
- **Passive**: saves Stamina, fewer disciplinary problems, "could give the
  opposition too much of an easy time". *(inference)* Committed should give
  Stamina, and probably Tackling and PI, more weight. Tackling is also the
  skill that is "less likely to foul".

**Roles** [manual A04]
- A priority order for **penalties** and a priority order for **free
  kicks**, plus a **captain**. The docs mention no corner-taker role, though
  Set Plays covers corners [manual A05]. The docs never describe what the
  captain does.

**Rules** [manual A04; BG manager_basics]
- In-play changes at 5-minute intervals: substitutions, Passing style,
  Mentality, Formation, or a different tactic. They can be triggered by the
  score, injuries or cards.

**Tactics Test** is an official tool for "analysing your players'
performances in certain positions" [manual A04]. Scout matches and
challenges are recommended for "gaining a better understanding of your
Players" [BG competitions].

### 7. Match engine, ratings and scouting

- Rating System: "After every match each participating Player is given a
  Rating … based on visible statistics and the formulas … takes into account
  both how well the team did in general, the actions each Player undertook,
  and the position of the Player. The Rating is open ended, and could even
  be negative" [manual B01]. The MVP is the highest-rated player [manual B01].
  *(inference)* Because the rating depends on team results and on position,
  it's a weak target for checking a skill-based evaluator.
- Full match statistics per player: shots, shots on target, passes, tackles,
  saves [manual B01].
- Scouting opponents: the %-symbol on the opponent's last game shows their
  formation and "how each of their players performed" [manual A04].
- Youth scout report: it costs €872. It "gives an idea of which position the
  player will be good at (forward/midfielder/defender/goal keeper), but does
  not give any indications of potential maxings" [manual A10]. It "will
  pinpoint the main strengths and weaknesses" [BG manager_basics]. It's only
  available for players who were youths in season 70 or later
  [BG manager_basics].
- The docs contain no "stars", overall rating or position-suitability score
  for senior players.

### 8. Value, salary, other attributes

- Value is recalculated at the start of each season. Salary is based on
  value. A mid-season sale triggers a recalculation based on the transfer
  fee [manual A03]. The maximum asking price is 3× value [manual A06]. The
  docs don't give the value formula. *(inference)* Value could serve as a
  rough market-based overall rating, but how it weights skills, age and
  position is unknown.
- Height, weight and strongest foot appear on the profile [manual A03;
  BG get_started]. The docs never say whether they affect matches.
- Player items are cosmetic only [manual E03].
- Coaches have a "class" (A-class and so on), which shows coaching ability,
  not player rating [manual A09; BG your_first_season].

## Open questions and contradictions

- Form Package: the manual says Form 7, the BG says Form 9. The Form scale
  itself is never stated.
- How exactly do Form and Experience combine with skills? The BG says "as
  much impact as the attributes", which is vague. It could be additive or
  multiplicative.
- Is there any wide or central difference in the engine? There are no
  winger, full-back or DM positions officially, only thirds and free
  placement.
- What does the captain do? The docs never say.
- Who takes corners? Only penalty and free-kick orders exist. Is it the
  free-kick taker *(inference)*?
- Do height, weight or foot affect heading or matches? The docs don't say.
- No official skill weights per position or per passing style exist. The
  camp counts in section 5 are the only official quantitative proxy, and
  forum or empirical research needs to supply real weights.
- Does the Long balls style use Aerial Passing (air) or Passing (ground)?
  The BG says "longer ground and air passes", which suggests both.
- Were these texts updated after the Season 40 sim change? The README says
  the manual is current. Mentions of seasons 51 and 70 support that.

## Queries used

- Read in full: manual A03 Squad, A04 Tactics, A05 Training. Read in part:
  B01 (Rating/MVP), A08 Youth, A10 Youth Exchange, A06 (value/salary lines),
  C05, E03.
- BG sections: get_started lines 75–183 (attributes, form, experience) and
  manager_basics lines 20–210 (tactics, mentality, passing, training,
  youth).
- grep -i over all sources: scout, rating, star, position, defender,
  striker, wing, goalkeeper, keeper, midfield, long ball, short, counter,
  experience, form, stamina, potential, skill, captain, penalt, corner, free
  kick, value, max, height, weight, foot, class.
