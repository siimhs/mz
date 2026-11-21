# Forum: central midfielders, wingers, strikers — skills and tactic effects

Source: local forum DB (forum 10, posts 2011-02 … early 2012) + manual text.
Era labels used below:
- **old sim**: the live engine before Season 40 (posts before ~2011-10-18).
- **H2H/test**: test versions V1.1–V2.02 of the new engine (Jul–Oct 2011, Head2Head / test page).
- **S40 sim**: the new engine, live from ~2011-10-18/19 ("New Simulator live now", topic 9568328).
The engine has changed again since then, so treat all forum weights as directional only.
`patrick7777` posts as Crew: other posters quote him as "patrick7777 (Crew) (GPA)" [post 35172210, 2011-10-02]. That makes his posts **official**.

## Summary (the answer)

1. **The game has no explicit position roles.** Which third of the pitch a player is dropped into makes him a defender, midfielder or attacker. Everything else (winger, anchor, deep forward) comes from where he stands on the pitch, and from the passing style for wingers [manual A04 Tactics] **official**. Crew still talk about "wingers", "CM", "deep forward" and "centre forward" as positions you create by placement [post 35163795, 2011-10-01, patrick7777] **official**.
2. **A "winger" is a player placed near the touchline, and he only behaves like one under the Wings passing style.** Crew: "Wingers like to be positioned near the right and/or left touchlines" and "wingers enjoy playing really wide". Players placed too narrow "play more like inside forwards" [posts 34779168, 35159177, 35194786, patrick7777] **official**.
3. **Winger skills according to Crew:** "pace, stamina, ball control, aerial passing/passing" [post 34451557, 2011-07-08, patrick7777] **official**. Posters say speed + AP [post 35407726, prednibe] or speed + AP for crosses [post 36078858, sephiroth] **opinion**. Several posters say AP made no visible difference on S40 wingers: 0–4 AP crossers hit their target ~90% of the time **contested**.
4. **Striker core skills: Shooting, Ball Control, Speed, Stamina, plus Play Intelligence.** This is the near-universal list in both sims **consensus**. In camp packages, Shooting appears in 10 of the 13 striker camps, BC in 7, Speed in 5, Passing in 4, PI/Heading/Stamina in 3 each and AP in 2 [manual A05 Training].
5. **Striker weight shift in the S40 sim:** Speed dropped in importance, and BC and PI matter more. One poster ranks them "SH, PI, BC, ST, SP in that order" [post 35418681, egtor] **opinion**. Others saw no difference between 6 and 9 speed [post 35542375, tomamientras]. Crew cut shooting weight by 50% in H2H V1.49 [post 35073428] **official**.
6. **Old-sim striker:** speed was king. "Speed is often more important than shooting in the lower leagues" [post 33919438, frankzappa]. Speed/Stamina/BC were seen as the "most important skills" overall [posts 34003124, 35205487] **consensus (old sim)**.
7. **Heading for strikers matters mainly with Wings plus wingers who have good AP** [posts 35359536, 35364885, patrick7777] **official**. Posters' tests found no difference between 0, 3 and 7 heading [posts 35362547, 35514157, capcash; 35372616, civilu] **contested**. Player height is not used [post 36241601, quoting crew] **opinion relaying crew**.
8. **No "target man" concept appears in the forum DB.** The phrase gets zero hits. A fast striker suits Long ball ("Got to have fast strikers if you are playing long balls!" [post 34772326]) **opinion**. A heading striker suits Wings with AP crossers (point 7).
9. **Central midfield skills:** Passing, PI, Aerial Passing, Ball Control, Speed and Stamina. The midfielder camp counts are Passing 9/16, PI 8, AP 8, Speed 7, BC 6, Stamina 5, Tackling 2, Shooting 1 [manual A05]. For the youth package, the midfielder primary is PI, with stamina, passing, BC, tackling and crossing (AP) as secondaries [post 36248128, civilu] **consensus**. Posters' rule of thumb is that midfielders need "a bit of it all" [posts 34136545, 33933740] **consensus**.
10. **Defensive vs attacking mid:** posters do not separate them clearly. A DM is like a CM with more tackling and PI, because low TA/PI DMs collect cards [posts 33919356, 35367120, 36255118]. An AM is like a CM with more shooting: midfield shots from distance score well in the S40 sim [posts 35367120, 35892148] **opinion/consensus-lite**.
11. **S40 sim:** central midfielders mostly do one-touch passing, so "won't need anything more than passing, stamina and playing intelligence" [post 35132019, nikos6, test sim] **opinion**. Another view: the CM "is very demanding now (in skills)" [post 35285181, mihairo] **contested**.
12. **Tactic shifts:**
    - **Short passing** rewards Passing/BC/PI across the midfield [post 34428793; manual A04].
    - **Wings** rewards winger Speed/BC/AP and striker Heading [Crew posts above].
    - **Long ball** rewards AP, receiver BC and striker Speed [posts 35533817, 34772326, 33883982].
    - The S40 sim strongly favoured **Wings**. Crew said short and long "can beat wings ... it depends on the type of players you have" [post 35299246] **official**.
13. **Mentality/pressing:** the forum gives nothing skill-specific. The manual says Committed pressing uses more Stamina and risks more fouls [manual A04] **official**. So Stamina weight should rise under Committed pressing, and Tackling/PI weight for midfielders too (card avoidance).
14. **Explicit heuristics posters use are minimum thresholds, not formulas** (see Findings §5). No poster gives a weighted formula. Several suggest the sum of all balls matters more than single skills [posts 34273623, 35367120] **opinion**.
15. **PI is controversial.** Some treat it as a multiplier on other skills, like form/experience [posts 34811106, 36316670]. Others say high-PI strikers pass instead of shooting [posts 36298244, 36256283, 35283585] **contested**.

## Findings

### 1. Is there a winger role? (positions vs placement)

- The tactics field is split into thirds. A player dropped into a third "will automatically adopt that role (e.g. defender, midfielder, attacker)", and is then fine-tuned within it [manual A04 Tactics] **official**. No winger or DM role exists in the UI. The Roles screen is only for penalty/free-kick takers and the captain [manual A04 Tactics] **official**.
- Bench slots: GK, defensive players, and two slots for "midfielders and strikers" [manual A04 Tactics] **official**. So the engine knows only 3 outfield categories.
- Crew on wing placement: "The width is important as the player feels like a Winger. The distance from the left and right backs ... encourage my wingers to run more" [post 34451557, 2011-07-08, patrick7777, H2H] **official**.
- "Wingers like to be positioned near the right and/or left touchlines to get the best use out of their wing skills" [post 34779168, 2011-08-12, patrick7777, H2H V1.34] **official**.
- Players that are too narrow "are not really acting as wingers" [post 35159177, 2011-09-30, patrick7777]. They "play more like inside forwards" [post 35194786, 2011-10-05, patrick7777] **official**.
- Crew notation for a 4-4-2 Wings: "2 x CM 2 x wingers 1 x deep forward 1 x centre forward" [post 35163795, 2011-10-01, patrick7777] **official**. This is a useful position taxonomy: CM, wide mid/winger, deep forward, centre forward.
- "Players can operate in different positions on the pitch, but play best in the position(s) that match their skills" [post 35567056, 2011-11-23, patrick7777, S40] **official**. This implies a hidden mapping from skills to positions. It is not specified.
- Wingers placed in the attacking third receive passes "not as midfielders, but as attackers" [post 35138843, baltovenaar]. "Wingers has to be in the forwards' zone to be effective" [post 35354025, safieeamat]. The opposite view is that wide mids are better than wide attackers [post 35367120, 2011-10-26, prednibe]. Another poster advises placing them "deep and wide (against side line)" [post 35387202, baltovenaar] **contested** (where to put the winger).
- A winger behaves like a winger mainly under the **Wings** passing style. V1.2: "When playing wings tactic, the wing players will never make the collective decision to start a passing sequence. They will act instinctively instead" [post 34449451, patrick7777] **official**. In the S40 sim, "A wing player is the only player on the pitch that will run with the ball" [post 35367120, prednibe] **opinion**.
- **Implication for the evaluator:** "Wing" is a valid evaluation position, but it is a midfielder or attacker placed wide. Its weight profile should apply most strongly when passing = Wings.

### 2. Winger / wide midfielder skills

- Crew: "think about the kind of skills a 'Winger' would have, pace, stamina, ball control, aerial passing/passing, these skills will also affect how well a player and team can play with Wings" [post 34451557, 2011-07-08, patrick7777] **official** (H2H era, but stated as design intent).
- Crew, V1.34: "Air passing and heading are now far more important skills than before" [post 34779168, 2011-08-12, patrick7777] **official**. V1.2: "improved the accuracy of passing, based on the passing ability" and "Players are better able to control air passes" [post 34449451] **official**.
- Crew, V2.02: "Winger will cross if they can ... Wingers will use ground passes to create scoring chances ... cut inside and shoot" [post 35299246, 2011-10-17, patrick7777] **official**. So wingers use Passing and Shooting as well as AP.
- Posters, S40 sim:
  - "wingers should need speed and AP" [post 35407726, 2011-10-31, prednibe] **opinion**.
  - "if your winger has a decent speed he can run many yards before he passes/shoots" [post 35408812, konstas] **opinion**.
  - "Outside midfielder should have high speeds and aerial passing to deliver good crosses, and have to be accompanied by strikers with heading skill 5+" [post 36078858, 2012-01-31, sephiroth] **opinion**.
  - Poor full backs facing "fast wingers with High BC and Speed" with Wings commit more fouls [post 35853073, 2012-01-02, mosmosmos] **opinion**.
- Against AP mattering (S40 sim):
  - "even with low aerial passing the wingers are still consistently getting the ball to the strikers" [post 35405495, innuendo].
  - "my winger with zero in crossballs have 100% passing" [post 35357205, cannava].
  - Wingers with 7–9 passing/AP replaced by 0–3 ones showed "ANY ... difference" none [posts 35514157, 35533511, capcash].
  - "I have 2 balls on 1 winger and 0 on the other and they are making about 90%+ of Aerial Passing" [post 35368490, mosmosmos].
  - Players with ~2.5 balls of AP land 90% of long balls [post 35408812, konstas].
  - The same mosmosmos post says that in test V1.99, low-AP crossers "kept on doing crosses over the head of the striker".
  - **contested**. Crew insist the skills are used [post 35364885]. Later crew wanted defenders to head away more crosses [post 35567056], which would make crossing quality matter more.
- For AP mattering:
  - One poster moved a midfielder with 8 AP to the right wing and won 0-4, with 2 headed goals [post 35366025, baltovenaar] **opinion** (n=1).
  - "Absolutely massive IMO" [post 35630660, 2011-12-03, mosmosmos]. That is the same poster as above, later.
  - AP "Seems to be more important than previously ... if your tactics are crossing or long ball" [post 35621217, webregas] **opinion**.
- Skill naming: "Aerial Passing" was earlier called "Cross balls/Crossballs". The skill is unchanged; "Aerial passing just better describes what the skill actually affects ... i.e. all passes in the air" [post 35628640, 2011-12-03] **consensus**. Older posts saying "crossballs" or "LB" mean AP.
- H2H V1.2: one poster put "someone out there with 9 Passing and not much Speed" instead of a quick winger, and it worked better as an outlet [post 34474111, 2011-07-11, mosmosmos] **opinion**. On non-Wings styles, Passing may matter more than Speed for wide players.
- Old sim: "wings aren't effective if you're playing 'Use wings', they can work with short passing, however the winger will utilise more passing than crossing". Train passing to 7–8 before AP [post 33887996, 2011-05-05, matttombstone] **opinion (old sim)**.

### 3. Strikers / attackers

- Manual: "Players with high Shooting skills are more likely to get into goal scoring positions and find the back of the net" [manual A05 Training] **official**. So Shooting also affects positioning.
- Camp package frequency for strikers (13 packages): Shooting 10, BC 7, Speed 5, Passing 4, Stamina 3, PI 3, Heading 3, AP 2 [manual A05 Training] **official** (as a proxy for designer intent, not engine weights).
- Youth Training Camp combos: "ATTACKER: shoot + speed, intelligence, ball control, tackling" [post 36248128, 2012-02-21, civilu] **consensus** (describes the game's YTC). Tackling shows up oddly for strikers. Another poster: "Tackling seems to be irrelevant to strikers" [post 35516769, gorgio].
- Old-sim minimums:
  - "Strikers should have at least 7 Speed, Stamina, Shooting and Ball Control" and PI ≥5 [post 33933740, 2011-05-10, jaysargey].
  - "Shooting (9) speed (7) stamina (7) control (5)" [post 34136545, 2011-06-03, acidburns85].
  - A good old-sim striker: "9 pi, 7 speed, 8 stamina, 10 BC and 9 shoot" [post 34273623, baltovenaar].
  - **consensus (old sim)**: Shooting + Speed + Stamina + BC are the core, with PI as a support skill.
- Old sim: "Speed is often more important than shooting in the lower leagues" [post 33919438, 2011-05-09, frankzappa] **opinion**. Also "a striker with 7/8 speed but with 6/7 bc/pi scores more than a striker with 9 speed and 3/4 bc/pi" [post 35542695, nomercyurz] **opinion**.
- S40 sim, speed falls:
  - "today forwards with high shooting, control and Intelligence are enough". A 6-speed and a 9-speed striker with otherwise equal skills "performs equal" [posts 35541939/35542375, 2011-11-19, tomamientras].
  - Reply: speed "isn't as effective as it used to be" [post 35542098, mosmosmos].
  - "speed is nowhere near as important as it used to be" [post 35348078, matttombstone].
  - **consensus (S40)**.
- S40 sim striker rankings:
  - "most important skills for strikers now are SH, PI, BC, ST, SP in that order" [post 35418681, 2011-11-01, egtor] **opinion**.
  - "Attackers need BC, Speed, Stamina and PI while shooting seems diminished" [post 35407726, prednibe] **opinion**.
  - The two agree on BC + PI + ST and disagree on Shooting vs Speed. **contested**
- Crew, V1.49 (H2H): "Reduced the weighting of the shooting skill for all players by 50% - In the new Sim, our attackers are just too good" [post 35073428, 2011-09-19, patrick7777] **official**. V1.98: "players with low shooting (and connected secondary) skills are able to score so many goals, this was not our intent" [post 35266669] **official**. Shooting stays the primary skill, but its spread was compressed in the S40 sim.
- Striker PI (S40) **contested**:
  - For PI: "I won't touch a striker unless he is above 6 PI" [post 36291423, nav]. "Not enough BC..." in reply to strikers with 3–5 PI who can't score [post 35457336].
  - Against: "having LESS PI balls as a Striker actually helps as they just decide to shoot not to pass" [post 36298244, mosmosmos]. A 9-PI striker shot less than a 6-PI one [post 35283585, huangho]. "my strikers with 8 Pi likes to pass ... my striker with 5 Pi likes to take shots" [post 36256283, darkline].
- Heading:
  - Crew: "Heading is a useful skill for strikers, particularly in teams playing wings, where the wingers have good aerial passing" [posts 35359536, 35364885, 2011-10-24/25, patrick7777] **official**.
  - Posters' tests found no difference [posts 35362547, 35514157, capcash; 35358753, 35372616, civilu] **contested**.
  - "Heading only come from high crossing by winger" [post 35369401, coldfire01] **opinion**. Header goals come almost only from Wings crosses.
- Target man vs fast striker:
  - No "target man" hits.
  - Fast striker: suits Long ball [post 34772326, willo1904] **opinion**. The S40 sim reduced the value of pure pace: "you can't just stick a player with 10Sp 9 Shooting with nothing else and expect him to run through teams" [post 35542098, mosmosmos] **opinion**.
  - Aerial striker: Heading + Wings + AP wingers (Crew view above).
  - Placement: in the S40 sim, attackers placed slightly deep shoot from the edge of the box, which converts better than 1-on-1s or headers [post 35408812, konstas; post 35407911, ozteam] **opinion**.
- Weighing a real choice [post 35885893, 2012-01-06, mosmosmos]: Speed + Stamina + Shooting with poor BC, vs 10 Shooting with low fitness, vs decent pace, great Stamina, average Shooting, great BC and passing. A "perfect striker" needs all of Speed, Stamina, Shooting and BC **opinion**.

### 4. Central midfielders (CM / DM / AM)

- Camp package frequency for midfielders (16 packages): Passing 9, PI 8, AP 8, Speed 7, BC 6, Stamina 5, Tackling 2, Shooting 1, Heading 1, Set Plays 1 [manual A05 Training] **official** (proxy).
- YTC: "MIDFIELDER: intelligence + stamina, passing, ball control, tackling, crossing". PI is the midfielder primary [post 36248128, civilu; post 34481636, jaysargey "PI as a main skill is seen as a midfielder primary in YTC"] **consensus**.
- Old-sim minimums:
  - "Midfielders should have at least 7 Speed, Stamina, Passing, BC, Tackling and Shooting", with PI ≥5 [post 33933740, jaysargey].
  - "Midfield needs a bit of it all: speed (7) stamina (7) control (6) tackling (5) passing (5) shooting (5) Int (5)" [post 34136545, acidburns85].
  - **consensus (old sim)**: midfielders are generalists.
- "first passing then ball control" for wingers and midfielders [post 34787314, jamiescfc2010] **opinion**.
- National-team view: CMs can have lower speed (≥6–7) than strikers and defenders (9–10). Primary skill (passing) ≥9 [post 36078858, 2012-01-31, sephiroth] **opinion**.
- S40/test sim:
  - Midfielders "won't need anything more than passing, stamina and playing intelligence ... they don't keep the ball" [post 35132019, 2011-09-27, nikos6, test sim] **opinion**.
  - Against that: "The Central Midfielder is very demanding now (in skills)". An attacking CM with 7 Sp, 7 St, 5 PI, 5 Pa, 5 Sh, 6 BC, 6 Ta, 4 AP and 9 exp did well [post 35285181, 2011-10-16, mihairo] **opinion**.
  - "Midfielders don't have much use" [post 35311963, dagriggstar] **opinion**.
  - Crew rebutted that last view with an example where a missing central midfield conceded goals [post 35303055, patrick7777] **official** (midfield presence matters).
- DM:
  - Card problems come from low tackling, and "could be as much a problem of low PI as low TA" [posts 33919153, 33919356, frankzappa].
  - Crew: "Their intelligence and experience will allow them to avoid making tackles that would have resulted in a card" [post 35194786, patrick7777] **official**.
  - A DM/defender with 5 PI was sold for missed tackles [post 36255118, civilu] **opinion**.
  - So for the DM, Tackling + PI rise.
  - The S40 DM/anchor takes long balls: a DM with 8 AP mostly played long balls [post 35286610, huangho] **opinion**. AP matters for a deep playmaker.
- AM: "Long shots are important for CAMs" [post 35350704, matttombstone, as a proposal]. A CM who can "put the pass through to the Striker, or hit the 30 yard shot" [post 35892148, mosmosmos]. In the S40 sim, "Central defensive midfielders will take part in the offensive ... and score from the distance" [post 35367120, prednibe]. **opinion**: AM = CM + Shooting.
- One poster's path to DM vs AM: train speed/stamina/BC to 6+, then at about 22 "decide to go for a defensive midfielder or attacking midfielder and train PI along with your chosen stat" (tackling vs shooting/passing) [post 35661163, ashleyhughes] **opinion**.

### 5. Explicit heuristics / formulas / thresholds

None of these are weighted formulas. All are threshold rules.
- jaysargey [post 33933740, 2011-05-10, old sim]:
  - GK ≥5 Sp/St, ≥7 Keeping/BC
  - Def ≥7 Sp/St/Ta/BC
  - **St ≥7 Sp/St/Sh/BC**
  - **Mid ≥7 Sp/St/Pa/BC/Ta/Sh**
  - all players PI ≥5, "the more, the better"
- acidburns85 [post 34136545, 2011-06-03, old sim]:
  - St Sh 9, Sp 7, St 7, BC 5
  - Mid Sp 7, St 7, BC 6, Ta 5, Pa 5, Sh 5, PI 5
- Youth cut rules (all positions): sell or convert youths who max ≤6 in Speed or Stamina, or below ~6–7 BC [posts 34789814, 34790059, 34792387] **consensus**. These posters treat Speed/Stamina/BC as universal base skills, so they belong in every outfield position's score.
- "Total balls" view: "strikers with low shooting and lots of overall balls are doing a better job. This made me think if the total number of balls is important" [post 34273623, baltovenaar]. "most likely aggregated attributes are decisive" [post 35367120, prednibe] **opinion**. This supports a base term of the skill sum plus position-specific weights.
- PI as a multiplier: "more of an amplifier on how your player reads the game ... much like Form and Exp" [post 34811106, jaysargey]. "PI acts as a multiplier on the players other skills in a similar way to how form and experience work" [post 36316670, wanttostrokemypuppy] **opinion**.

### 6. How importance shifts with tactic

- **Short passing** [manual A04]: patient possession, possibly at the cost of cutting edge. With better passing accuracy, "teams with high (average) passing skills, BC skills and possibly PI" benefit [post 34428793, casper8, H2H] **opinion**. The old sim was dominated by short passing through the middle [post 35367120, prednibe] **consensus**. In the S40 sim, short passing is seen as weak: "short passes dont work well" [post 35408812, konstas; post 35437439, fva2101] **consensus (early S40)**. One poster found it raised strikers' shot accuracy [post 35450955, chelseax1] **opinion**.
  - Weights: Passing, BC and PI up for all midfielders. Winger-specific Speed/AP down.
- **Long ball** [manual A04]: quick to shooting positions, less precise.
  - Strikers: "Got to have fast strikers if you are playing long balls!" [post 34772326, willo1904] **opinion**.
  - Passers: AP "hugely important to defenders if you play long ball" [post 33883982, alx460] **opinion**. A low-AP passer gives the receiver trouble controlling the ball [post 35533817, darkline] **opinion**. So receiver BC matters too.
  - "Long balls is aerial passing" [post 35354025, safieeamat] **opinion**.
  - Several posters said long ball played like short passing in the S40 sim [post 35455768, dagriggstar] **opinion**.
  - Weights: AP up for DM/CM/defenders. Speed and BC up for strikers.
- **Wings** [manual A04]: "tests the opposition full backs and goalkeeper with crosses to your strikers".
  - Crew: winger skills are pace, stamina, BC, AP/passing; heading matters for strikers; wingers must be wide [posts 34451557, 34779168, 35359536] **official**.
  - The S40 sim was widely seen as Wings-dominated [posts 35408812, 35418681, 35367120] **consensus**. Crew: short or long "can beat wings ... it depends on the type of players you have" [post 35299246] **official**.
  - Wings stats run high possession with "a lot of useless shots" [post 35407911, ozteam] **opinion**.
  - Weights: winger Speed, AP and BC up. Striker Heading up. CM Passing slightly down.
- **Mentality** [manual A04]: Attacking/Balanced/Defensive changes how many players commit forward. The forum has no data on skill weights by mentality.
- **Pressing** [manual A04]: Committed "will also use up more Stamina ... and run the risk of committing more fouls" **official**. Weights: Stamina up for all. Tackling/PI up for DM and defenders (cards).
- Counter attack is not a separate passing style. It appears only as a proposal [post 35348078, matttombstone] and in the manual's description of Defensive mentality.

## Open questions / contradictions

- **AP and Heading in the S40 sim:** Crew say they matter (Wings), while several careful posters measured no difference [capcash, civilu]. The likely truth is a small effect per ball ("There just isn't a big enough difference between 3 balls and 10 balls, and that is in EVERY skill" [post 35368490, mosmosmos]). Crew later planned for defenders to head crosses away more often [post 35567056], which would increase the value of AP.
- **Striker Shooting vs Speed** after S40: egtor ranks Shooting first and Speed last. prednibe says shooting is diminished and Speed matters. Crew halved shooting weight in test.
- **Striker PI:** amplifier or "passes instead of shoots"? Unresolved. A modest weight is the safest choice.
- **Where wingers should sit:** midfield-wide vs attack-wide. Both are reported to work. This matters for whether "Wing" pulls from the midfield or the striker weight profile.
- **All of this is from 2011–12.** The engine has been revised since. Crew's statement that players "play best in the position(s) that match their skills" [post 35567056] means a hidden mapping exists, but it is never specified.
- No forum data separates DM from AM numerically. Weights for them would be extrapolated.

## Queries used

- FTS: `"target man"` (0 hits), `winger NEAR(speed, 10)`, `striker NEAR(shooting, 8)`, `"key skills"`, `midfielder NEAR(passing, 6)`, `"defensive midfielder" OR "attacking midfielder" OR "holding midfielder"`, `"long ball" OR "long balls"`, `"aerial passing" NEAR(winger, 10)`, `crossball* OR "cross balls"`, `"counter attack" OR counterattack`, `"strikers need" OR "midfielders need" OR "wingers need" ...`, `"most important skill(s)"`, `"short passing" skills`, `wings NEAR(wingers speed, 15)`, `NEAR(midfielder intelligence, 12)`, `NEAR(midfield PI, 10)`, `NEAR(striker heading, 6)`, `anchorman OR "anchor man"`, `"wing forward" OR wingforward OR "inside forward" OR "deep forward"`, `NEAR(speed striker, 5)`, `NEAR(height heading, 15)`.
- Topics read: 9595256 (New Sim Think Tank), 9634719 (Speed is now useless in strikers), 9423402 (Speed and Stamina first), 9430344 (Does BC matter), 9582428 (Heading a useful skill?), 9831050 (Player Intelligence), 9659131 (Aerial Passing / Crossballs), 9668707, 9191331, 9272613, 9760927, 9731231 (Who for Centre Forward), 9823434 (Improved header), 9417827, 9579300 (partial), 9574849 (partial).
- SQL: every post by Crew users (patrick7777, marcus, mikelony, blackwater, powdersnow, thunderexpress) mentioning skill, striker, winger or midfield. Key Crew posts: 34449451, 34451557, 34779168, 35073428, 35159177, 35163795, 35194786, 35266669, 35299246, 35359536, 35364885, 35567056.
- Manual: A04 Tactics (placement, passing, mentality, pressing, roles), A05 Training (skill definitions, camp packages, counted per skill).
