# Forum: goalkeeper and defender skills, and how tactics shift them

Scope: local forum DB (forum 10, 2011-02 … 2012-02) plus the manual and Beginners Guide.
Timeline: the "old sim" ran until about 2011-10-19. H2H/test sims V1.x–V2.0 ran from 2011-07 to 2011-10.
The new (Season 40) sim went live on 2011-10-19/21 (topic 9568328 "New Simulator live now").
Each claim is marked **old sim**, **test sim** or **new sim** where it matters.
`patrick7777` is MZ Crew and the sim developer (shown as "patrick7777 (Crew) (GPA)", quoted in
[post 35172210]), so his posts are tagged **official**.

## Summary
- **GK, Keeping is dominant.** Posters treat Keeping as the primary skill by a wide margin. Typical
  cut-offs are "discard keepers under 8" and a 10 for national teams. Crew confirms that "higher keeper
  score → more shots saved". The new sim floors Keeping at 4 (a 0-ball keeper gets 4 free balls), which
  compresses the low end. **official + consensus**
- **GK secondary skills.** Stamina, Ball Control (= handling), Speed and Play Intelligence are the core
  secondaries. Aerial Passing comes next (goal kicks and distribution) and Set Plays last. Passing,
  Shooting, Tackling and Heading are useless for a GK. **consensus** (6+ posters). Crew: "Keepers need
  stamina … the lower their stamina, the greater the risk they will become tired". The Crew developer uses
  keepers with 5–7 Stamina. **official**
- **GK numeric heuristics from posters:**
  - Keeping ≥7–9, Speed ≥5, Stamina ≈5–6, BC ≥6–7, PI "as high as you can".
  - matttombstone's grouping: Speed/Stamina/BC (tier 1), PI (tier 1.5), Aerial Passing (tier 2),
    Set Plays (tier 3).
  - The game's own youth camp for keepers is Keeping + Speed, PI, Aerial Passing, Set Plays.
- **Defenders, core skills.** Tackling is the primary skill, with Speed and Stamina just behind. Then come
  Ball Control and Play Intelligence, then Passing and Aerial Passing. Heading is weak to useless (in the
  new sim defenders rarely headed until a late-2011 patch).
  - Heuristics: Tackling 9, Speed 7, Stamina 7, BC 5. Or "≥7 in Speed, Stamina, Tackling and BC".
  - National-team level: Speed 9–10 for defenders "to eliminate breakaway goals".
  - The youth defender camp is Tackling + Speed, Stamina, PI, Aerial Passing. **consensus**
- **New sim (S40) shift for defenders.** Posters report:
  - Tackling matters relatively less: "breaks" (interceptions) outnumber tackles about 3:1.
  - PI, Experience and Passing matter more.
  - Speed and Stamina were reported both as less important and as still important. **contested**
  - Crew: better (faster, higher-PI) players make more tackles. Experience and PI reduce bad tackles and
    cards, and the general effect of Experience was reduced. **official**
- **Central vs wide defenders: the game does not label them.** The manual has only GK/Defender/Midfielder/Attacker
  areas and no positional roles. **official**
- **The new sim's AI does treat defenders differently by where you place them.** Crew (V1.63): "Wide
  DEFENDERS (fullbacks and wing backs) will try to support central defenders more, but … their primary role
  is to defend against wide attackers." **official**
- **Full-backs had little impact in the old sim, more in the new one.** In the old sim they averaged about
  1 pass or tackle a game and many teams played 3 CBs + 2 DMs. In the new sim they are active, carry the
  ball forward, and are exposed to fast, high-BC wingers when the opponent plays Wings, which brings fouls
  and red cards. Teams "put their worst defenders at full back". **consensus**
- **CB vs FB profile.** "Not quick enough for CB" → move him to full back. So Speed is weighted more for
  central defenders. **opinion**
- **Sweeper and DM.** There is no sweeper or libero setting. Placement creates the role: a "superfast
  sweeper" behind a back four handled a strong striker (test sim). DMs need the same skills as defenders
  (Tackling + PI). DMs were card-prone in the old sim (low PI or low Tackling), and in the new sim "anchor
  men" were said to help the defence little. **opinion**
- **Passing style.**
  - Long ball (and Wings) makes defenders' Aerial Passing matter. In the old sim: "aerial pass is hugely
    important to defenders if you play long ball" (but behind Tackling, Speed and Stamina).
  - In the new sim, Aerial Passing was reported as "absolutely massive". One poster instead saw almost no
    effect: his 2.5-ball AP players completed 90% of long balls. **contested**
  - Crew says Wings "tests the opposition full backs and goalkeeper with crosses". That points to full-back
    quality and GK mattering more against Wings teams. **official**
- **Mentality, pressing, lines.**
  - Crew: an Attacking mentality makes "defenders take more risks", so you concede more. A Defensive
    mentality means fewer risks.
  - A relayed patch note: defenders no longer press as early as midfielders and forwards.
  - Committed pressing burns stamina and causes fouls (manual). That raises the value of defenders' Stamina,
    PI and Experience.
  - There are no offside-trap, defensive-line or man-marking settings. Posters control the "line" only by
    where they place players.
- **No explicit formula or weighted skill sum exists in the forum.** Posters argue only through rankings and
  per-skill thresholds. Crew states that "stats + tactics = result" and that skill weights are rebalanced
  between sim versions. In V1.49, for example, Shooting's weight was cut by 50% and all GK skills were
  boosted while keeping the relative differences between levels. **official**

## Findings

### 1. Goalkeeper

**Official (Crew / manual)**

- **Keeping floor in the new sim.** "If your player has '0' balls in keeping, the Sim gives them 4 free
  balls … The higher the keeper score, the more shots you will save" [post 35248205, 2011-10-12, patrick7777].
  **official, test sim (carried to the live sim)**. This implies effective_keeping = max(keeping, 4).
- **Stamina.** "Keepers need stamina, they move and run, they become tired. The lower their stamina, the
  greater the risk they will become tired or exhausted." Patrick7777's own keepers have 5–7 Stamina, "not a
  magic number" [post 35598946, 2011-11-28, patrick7777].
- **Secondary stats count.** "A keeper with 1 ball in stamina will soon tire, the rest of the secondary stats
  are also quite poor" [post 35597407, 2011-11-28, patrick7777]. **official, new sim**
- **V1.49 rebalance.** "Improved all Goalkeepers skills, but maintained the relative differences … a 10 ball
  Goalkeeper is still a lot better than 1 ball", and "Reduced the weighting of the shooting skill for all
  players by 50%" [post 35073428, 2011-09-19, patrick7777]. **official, test sim**
- **GK too strong in tests.** Keepers were "much better, and probably too good" in the test sim
  [post 35124101, 2011-09-26, patrick7777]. In V1.63, "Goalkeeping abilities have been adjusted in line with
  the changes to DEFENDERS" [post 35157183, 2011-09-30, patrick7777]. **official**
- **Manual and guide definitions.**
  - Keeping: "without high Keeping skill levels they will be less effective shot stoppers", and "good
    goalkeepers need other attributes as well" [manual A05 Training].
  - Stamina "includes goalkeepers" [Beginners Guide get_started].
- **Camp packages.** Clubhouse and Wild Card camps that include Keeping pair it with Stamina, PI, Speed or
  Passing [manual A05 Training].
- **Youth camp.** The youth "KEEPER" camp is Keeping + Speed, PI, crossing (= Aerial Passing) and Set Plays
  [post 36248128, 2012-02-21, civilu]. **official-by-design** (relayed game data).

**Community rankings of GK secondary skills.** Topic 9532156 "keeper secondary skills" was written in
2011-10, during the S40 test and launch.

| Poster | Order (best first) | Cite |
|---|---|---|
| egtor | BC > Speed > PI > Stamina > Set Plays > Aerial Passing > Passing; Tackling, Shooting, Heading "totally un-necessary" | [post 35174244, 2011-10-02, egtor] |
| dowopado | Stamina, BC, PI, Speed, Aerial Passing; "the rest who cares" | [post 35174788, 2011-10-02, dowopado] |
| marciocha_20 | Stamina, BC, PI, Aerial Passing | [post 35175914, 2011-10-03, marciocha_20] |
| bestplaymaker2134 | Stamina, BC, Speed, PI; Set Plays "no point" | [post 35184745, 2011-10-04, bestplaymaker2134] |
| dagriggstar | BC, Stamina, Speed, Aerial Passing, PI | [post 35185990, 2011-10-04, dagriggstar] |
| civilu | PI & Stamina, Speed, BC, Aerial Passing, Passing | [post 35186004, 2011-10-04, civilu] |
| matttombstone | Tier 1: Speed, Stamina, BC; tier 1.5: PI; tier 2: Aerial Passing; tier 3: Set Plays. Passing, Shooting, Tackling and Heading are useless | [post 35208087, 2011-10-06, matttombstone] |

- **Overall picture.** Stamina, BC, Speed and PI make up the top four in nearly every list. Aerial Passing
  comes fifth, and Set Plays and Passing come last. **consensus**
- **matttombstone's reasoning** [post 35208087]. **opinion**
  - "BC can be translated to handling."
  - Speed helps positioning, rushing out and diving.
  - Stamina of "about 5 or 6" is enough unless your defence is weak and the keeper faces many shots.
  - PI improves judgement on rushing out and positioning.
  - Passing is useless because keepers only play aerial passes and set pieces.
- **Older thresholds (old sim).**
  - "Speed at least 5. Stamina about 5 … PI as high as you can get it. BC at least 6 … Aerial Passing isn't
    a must but it helps", because the keeper "will always lug it upfield" [post 34663925, 2011-07-30, jaysargey].
  - Keepers need "at least 5 Speed and Stamina and 7 Keeping and Ball Control"
    [post 33933740, 2011-05-10, jaysargey].
  - "Keepers max keeping (9) then speed (5) and stamina (5)" [post 34136545, 2011-06-03, acidburns85].
  - "All keepers under 8 get discarded" [post 34792387, 2011-08-14, civilu].
- **Contested: Aerial Passing for a GK.** "Not that important … It never helps"
  [post 34665712, 2011-07-30, xellu] vs "it helps" [post 34663925].
- **Keeping dominance.** "Keeping as a skill is so so important/dominant in that position than any other
  skill" [post 34907394, 2011-08-29, ade-pm]. **opinion, echoed widely**
- **New-sim perception.** Posters felt the GK decides matches: "Only the Goalkeeper and the tactic?"
  [post 35549883, 2011-11-21, konstantinidisk]. Crew agrees that GK influence was too large and reduced it
  [post 35248205]. **contested** (perception vs Crew's intent)
- **Experience for a GK.** One poster's keeper had Keeping 10, Exp 10, Speed 8, Stamina 8 and still
  "ships goals" in the old sim [post 33892389, 2011-05-05, jamie1976]. There is no evidence that experience
  matters specially for keepers. Crew says experience's general effect was reduced in the new sim
  [post 35367521, 2011-10-25, patrick7777]. **official** (general)

### 2. Defenders

**Official**

- **Manual.** Tackling: "Good tacklers are more likely to win the ball … and less likely to foul". PI: "more
  likely to be in the right place at the right time, whether defending or attacking" [manual A05 Training].
  Speed: "Faster Players will be able to intercept more passes if they see them in time"
  [Beginners Guide get_started].
- **Defender camp packages.** All 15 include Tackling. The partner skills are Speed, Stamina, PI, Passing,
  Heading, Aerial Passing and BC; Set Plays appears only in package 15, which has no Tackling
  [manual A05 Training]. The youth defender camp is Tackling + Speed, Stamina, PI, crossing
  [post 36248128, 2012-02-21, civilu].
- **Tackle volume.** "Better players (faster, better PI) are able to make more tackles against lower skill
  players" (more tackles also means more cards) [post 35172210, 2011-10-02, patrick7777].
- **Cards.** A relayed patch note says "the most Intelligent or experience players will now receive fewer
  cards" [post 35194706, 2011-10-05, capcash, quoting the patch notes]. Crew confirms that experience reduces
  bad tackles [post 35367521, 2011-10-25, patrick7777].
- **Defender AI, V1.63.** "Better DEFENDERS will carry out these duties more efficiently than worst
  DEFENDERS" [post 35157183, 2011-09-30, patrick7777].
- **Height, weight and feet** are not used by the sim [post 35382524, 2011-10-28, patrick7777].

**Community: skill priorities**

- **Tackling, Speed and Stamina come first.**
  - "Tackling (9) speed (7) stamina (7) control (5)" [post 34136545, 2011-06-03, acidburns85].
  - "At least 7 Speed, Stamina, Tackling and Ball Control"; Passing and Aerial Passing help but are "not
    compulsary" [post 33933740, 2011-05-10, jaysargey].
  - "Tackling, speed and stamina always come before" Aerial Passing
    [post 33883982, 2011-05-04, alx460].
  - **consensus, old sim**
- **"Best defender" archetype.** A national-team defender is quoted with Speed 10, Stamina 9, PI 4,
  Passing 4, BC 10, Tackling 10 [post 34858421, 2011-08-22, jere92]. Another poster would take a "9-9-9-9
  with 6 passing" [post 34845836, 2011-08-20, ivorylegend]. Elite = maximum Tackling, Speed, Stamina, BC/PI.
  **opinion**
- **Speed for national-team defenders.** "Strikers and Defenders should have high speeds of 9-10 … having
  every defender speed at 9 should eliminate breakaway goals (of course tackling and stamina have to measure
  up as well)". Midfielders can be slower [post 36078858, 2012-01-31, sephiroth]. **opinion, new sim**
- **Ball Control.** A low-BC defender is "many times being tackled by the opposition's attackers"
  [post 34811106, 2011-08-16, jaysargey]. **opinion**
- **PI** is "more of an amplifier", like form and experience [post 34811106]. **opinion**
- **Market value vs match value.** BC and Passing drive transfer value most, then Shooting and Tackling.
  Speed, Stamina and PI "don't really carry much value" in price [post 34846085, 2011-08-20, mosmosmos].
  Note that this is market price, not match effect.
- **Heading.**
  - Old sim: "heading however is [useless]" [post 33886074, 2011-05-04, gebruikersnaam_].
  - New sim at launch: defenders never headed the ball away [post 35362346, 2011-10-25, jamie1976]
    [post 35300899, 2011-10-18, fathertaylor] [post 35326679, 2011-10-21, capcash].
  - After an update: "Defenders does Headng now!" [post 35580629, 2011-11-25, simbo737].
  - Later: "defenders head with long ball and wings" [post 35540436, 2011-11-19, steviedee] vs "heading is
    zero importance" [post 35538815, 2011-11-19, heddo].
  - **contested, low weight**

**New sim (S40) shift**

- **Tackling is less central.** "Tackling is now a secondary skill, because there are like 300% more breaks
  than tackles" [post 35283585, 2011-10-16, huangho, test V2.0]. "The main way of getting the ball is
  actually by interception rather than tackling", and "Passing is a huge primary skill now"
  [post 35283648, 2011-10-16, mosmosmos]. **consensus (2+), test sim**
- **Experience and PI.** "D need experience and PI. I guess stamina/speed is more important than on
  midfield … Still D needs tackling" [post 35407726, 2011-10-31, prednibe]. **opinion, new sim**
- **Speed and Stamina.** "Speed and stamina seem less important than before … heading, play intelligence
  and aerial passing seem much more important" [post 35405495, 2011-10-31, innuendo]. In the opposite
  direction, a slow defender cannot catch a fast striker [post 35314258, 2011-10-19, darkline]. A later post
  complains of "levelling of speed": a 5-speed defender catches an 8-speed winger
  [post 36168470, 2012-02-11, adv41449]. **contested**

### 2b. Central vs wide defenders, sweeper, DM

- **The game has no positional roles.** Players dragged into the defensive third "automatically adopt that
  role (e.g. defender…)". You can only fine-tune their position within the third, so there is no CB, FB, SW
  or DM label [manual A04 Tactics]. A poster asks Crew for role-based positions (SWP, RB, CB, LB, LWB, DM…),
  which confirms that none exist [post 35359828, 2011-10-24, madmonk]. **official**
- **The sim's AI separates wide from central defenders by placement.** "Wide DEFENDERS (fullbacks and wing
  backs) will try to support central defenders more, but will also remember that their primary role is to
  defend against wide attackers" [post 35157183, 2011-09-30, patrick7777]. **official, test V1.63 onward**
- **Patch note on pressing.** "Players considered DEFENDERS in your tactics will not press as early as
  before. MIDFIELDERS and FORWARDS will continue to press" [post 35136061, 2011-09-27, nemgabo, quoting the
  patch notes]. So the sim distinguishes the defender line from a DM placed in midfield. **official (relayed)**
- **Full-backs in the old sim.** "The full-backs never get involved", so this poster plays 3 CBs + 2 DMs
  [post 34831911, 2011-08-19, dance_of_the_dead]. In the test sim, "Side defenders are not playing … Only
  central defenders participate" [post 35139156, 2011-09-28, mikelony, relaying Spanish community feedback].
  **consensus, old/test sim**
- **Full-backs in the new sim.**
  - "Full backs are actually active … [in] the last 40 seasons [they made] about 1 pass/tackle a game"
    [post 35283662, 2011-10-16, mosmosmos].
  - Other posters still see side defenders "do absolutely nothing" in mirror 4-3-3s
    [post 35331798, 2011-10-21, jauho], or drop them [post 35410648, 2011-10-31, prednibe].
  - Others went from 3 to 4 at the back "because many people now play wings"
    [post 35410979, 2011-10-31, konstas].
  - **contested**
- **Full-back weakness is exploitable.** "If you have poor full backs against a team who has fast wingers
  with High BC and Speed and play the use wings tactic, you are more likely to commit fouls … Teams generally
  put their worst defenders at full back, it never used to matter in the old sim"
  [post 35853073, 2012-01-02, mosmosmos]. **opinion, new sim**. This means full-backs need Speed, Tackling
  and PI/Experience (to avoid cards) against wing teams.
- **CB vs FB speed.** "Deakin should make it now, as a Full Back. Not quick enough for CB and not good
  enough Passing for midfield" [post 35960806, 2012-01-16, mosmosmos]. **opinion**. This implies that CBs
  need more Speed than FBs, and that the defender → midfielder step is driven by Passing.
- **Sweeper.** "I then put a superfast sweeper behind my 4 back line, and it dealt with that [strong
  striker]" [post 34854348, 2011-08-21, hanzinho, H2H V1.34]. Crew's own example tactic used "3 Central
  defenders, one holding midfielder/forward defender" [post 35163795, 2011-10-01, patrick7777]. **opinion /
  official example**
- **DM.**
  - "DMs have always had trouble with refs … It could be as much a problem of low PI as low TA [Tackling]"
    [post 33919356, 2011-05-09, frankzappa, old sim]. The poster with the DM wanted the defender with
    "better passing and far better tackling" in that role [post 33919153, 2011-05-08, alx460].
  - In the new sim, "Defensive midfielders need the same [as D: experience, PI, stamina/speed]"
    [post 35407726, 2011-10-31, prednibe]. "Anchor men doesnt seem to help much my defense"
    [post 35408812, 2011-10-31, konstas].
  - "Defensive midfielders are drawn out of position very easily" [post 35223923, 2011-10-09, dagriggstar].
  - **opinion**

### 3. Tactic dependence

- **Passing style (manual).**
  - Long ball means "less precise passes" up to the forwards.
  - Wings "tests the opposition full backs and goalkeeper with crosses to your strikers".
  - Long balls = "longer ground and air passes to Players closest to the opposition's goal".
  - [manual A04 Tactics] [Beginners Guide manager_basics]. **official**
- **Long ball → defenders' Aerial Passing (old sim).** "Aerial pass is hugely important to defenders if you
  play long ball" [post 33883982, 2011-05-04, alx460]. "Got to have fast strikers if you are playing long
  balls" [post 34772326, 2011-08-11, willo1904]. **opinion (+1 agree, post 33885838)**
- **Long ball and Aerial Passing in the new sim.**
  - Matters more: "Certainly if your tactics are crossing or long ball then it figures that its more
    important … heading seems to be more significant too" [post 35621217, 2011-12-02, webregas]. "Absolutely
    massive" [post 35630660, 2011-12-03, mosmosmos].
  - Hardly matters: "the long ball skill doesnt really make any difference. I have players with an average
    of 2,5 balls … and 90% of their long balls reach correctly their target" [post 35408812, 2011-10-31,
    konstas]. "Long balls is aerial passing" [post 35354025, 2011-10-24, safieeamat].
  - **contested**
- **Wings attack → you need strong CBs and FBs.** "Defense requires strong, quality centre backs to allow for
  such an attacking set up" (wings/attacking/committed) [post 35405495, 2011-10-31, innuendo]. Wings beat a
  narrow tactic without side-backs [post 34854348, 2011-08-21, hanzinho]. Crew beat a dense central formation
  with a wide 4-4-2 [post 35163795, 2011-10-01, patrick7777]. **consensus**. Against Wings opponents,
  full-back quality matters (see 2b).
- **Mentality (Crew).** "Attacking: … your defenders will take more risks, as a result, you are more likely to
  concede. Defensive: … Your defenders will take less risks" [post 35252335, 2011-10-12, patrick7777]. Under
  Attacking mentality "Defenders will support attacks, even if it means leaving gaps"
  [Beginners Guide manager_basics]. **official**. Attacking mentality should therefore weight defenders'
  Passing and BC up, and Defensive mentality should weight pure defending.
- **Pressing.** Committed pressing burns more Stamina and causes more fouls [manual A04 Tactics]. So defenders'
  Stamina and PI/Experience (fewer card tackles) are worth more under Committed. In the old sim, "your
  defenders usually go out of position if you play aggressive … just play passive"
  [post 33748924, 2011-04-21, darkline]. "Forgot to reset … committed to normal and had a defender sent off"
  [post 35551806, 2011-11-21, tree_amigos]. **official + opinion**
- **Offside, defensive line, man marking.** There are no settings for these in the manual's six tactic
  elements [manual A04 Tactics]. The old-sim offside thread advises "set them [defenders] lower or switch to
  Normal play" [post 34001100, 2011-05-18, overthetop3]. Posters asked whether a zonal vs man-marking setting
  exists, and none is mentioned [post 35094199, 2011-09-22, dagriggstar]. **official (absence)**. Line height
  is purely formation placement, so it has no skill-weight effect I can source.

### 4. Formulas and skill-sum heuristics

- **No poster gives a numeric weighted formula** for GK or defenders. All heuristics are thresholds or
  ordered lists (see the tables above). One poster explicitly says only Crew knows the weights
  [post 33901527, 2011-05-06, matttombstone].
- **Crew on how results are computed.** Crew describes results as "stats + tactics = result". The "influence
  of each stat" is tuned between versions [post 35248205]. Known adjustments:
  - Shooting weight −50% (V1.49).
  - GK skills uniformly boosted (V1.49).
  - GK Keeping floor of 4.
  - Small boost for weak or young players in a strong team [post 35567056, 2011-11-23, patrick7777].
- **A poster test of skill sums.** He switched GK and striker and found similar results, which "makes me
  think rather than summing up individual skills…" [post 35218015, 2011-10-08, oner7, test sim]. This is
  anecdotal evidence against a plain sum. Crew rebutted the GK-as-striker results [post 35207160, 2011-10-07,
  patrick7777].
- **Implied heuristics I can derive** (my synthesis, not from posters; not cited as fact):
  - GK ≈ Keeping dominant (≥50%), then Stamina ≈ BC ≈ Speed ≈ PI, then Aerial Passing, then Set Plays.
    Passing, Shooting, Tackling and Heading count as 0.
  - Defender ≈ Tackling highest, then Speed ≈ Stamina, then PI ≈ BC (and Experience), then Passing, then
    Aerial Passing. Heading is minor and Shooting/Keeping count as 0.
  - Modifiers: under Long ball, Aerial Passing goes up. For a CB, Speed goes up. For an FB against Wings
    opponents, Speed and Tackling go up. Under Committed pressing, Stamina and PI go up. Under Attacking
    mentality, Passing and BC go up.

## Open questions / contradictions
- **How much a GK matters after S40.** Posters say GKs decide matches, while Crew says the influence was
  reduced and floored at 4. Nothing quantitative.
- **Aerial Passing for long-ball defenders and keepers:** "massive" vs "no difference" (90% completion at
  2.5 balls).
- **Speed for defenders in the new sim:** "less important" and "levelled" vs a 9–10 national-team ideal.
- **Tackling vs PI/Experience in the new sim.** Interceptions dominate, which suggests PI and Speed matter
  more than Tackling, but nobody has measured this.
- **Whether full-backs matter** in the new sim (active vs "do absolutely nothing").
- **No data on Experience specifically for keepers.**
- **No Crew statement on relative skill weights** for any position.
- **Everything here is from 2011–12.** The current engine may differ further.

## Queries used
- `goalkeeper OR keeper OR goalie`
- `(keeper OR goalkeeper OR goalie) AND (stamina OR intelligence OR "aerial passing" OR speed OR experience)`
- `(keeper OR goalkeeper) AND (skills OR important OR matters)`
- `defender* AND (tackling OR tackle) AND (speed OR heading OR intelligence)`
- `"full back" OR fullback OR "wing back" OR wingback …`
- `sweeper OR libero`
- `"defensive midfielder" OR DMF OR CDM …`
- `offside OR "man marking" OR "defensive line"`
- `"marking machines"`
- `"long ball" OR longball OR "long balls"`
- `formula OR weighting OR "skill sum"`
- `(defender OR defence) AND (pressing OR committed OR passive OR mentality)`
- `"centre back" OR "center back" OR CB …`
- `NEAR(defender* speed,10)…`
- `NEAR(defender* heading,8)…`
- SQL: all posts by patrick7777 filtered by skill keywords.
- Topics read: 9532156 (keeper secondary skills), 9389666 (Training keepers), 9452602 (How to spot a
  keeper), 9181750 (Aerial Pass & Set Plays), 9431177 (Best defender ever), 9191331 (Defensive midfielder
  dilemma), 9659131 (Aerial Passing / Crossballs), 9595256 (New Sim Think Tank), 9574849 (Tactical Advice
  Thread, partial).
