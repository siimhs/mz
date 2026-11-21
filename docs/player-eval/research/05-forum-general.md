# Position-independent factors and existing rating heuristics (forum DB + manual)

Scope notes:
- The DB actually spans **2011-02-17 .. 2012-03** (posts ~23k). The new sim went live for
  cups/friendlies on **2011-10-18** and for all matches from Season 40 [post 35303637, 2011-10-18, patrick7777].
  "pre" = old sim (in use since about 2008), "post" = new sim. The new sim is closer to the
  current one but was rebalanced after that, so treat everything here as dated.
- **patrick7777** is MZ staff (CEO / sim producer; quoted as "patrick7777 (Crew) (GPA)" in
  [post 35172210]) even though the users table does not flag him as Crew. His posts are tagged
  **official**.
- Match ratings, MVP, scout reports and "player evaluation" features **are not discussed at
  all** in the forum. They are later features. Only the current manual covers them.

## Summary (the answer)

1. **Old sim (pre-S40): Speed and Stamina were "king"**, followed by Ball Control. Heading and Set
   Plays were widely called useless. **consensus** [34003124; 35001651; 35205487; 33883982; 33886074]
2. **New sim (S40+): staff say Aerial Passing and Heading are "far more important" than
   before**, and shooting weight was cut by 50%. Posters report that speed and stamina matter
   less while PI, heading and AP matter more. **official + consensus** [34779168; 35073428; 35405495; 35542098]
3. Experience was a strong multiplier-like "bonus" in the old sim (**consensus**, though
   contested on how strong). In the new sim **staff say its effect was reduced** and moved
   partly into avoiding bad tackles and cards [35367521, patrick7777]. **official**
4. Manual: Form and Experience "help to maximise a player's other attributes". Experience also
   means "less likely to make mistakes". **official** [manual A05 Training]. Nobody outside
   staff knows the numbers. One poster's guess is +0.25x per Form ball and 0.05–0.15x per
   Experience ball. **opinion, explicitly speculative** [35709160]
5. Form: posters see **no noticeable difference between about 8 and 10 Form**. Only gaps of 3 or
   more balls are felt. **consensus (weak)** [34221116; 34216764; 34225242]. After S40 several
   posters say Form "doesn't seem to matter" much. **opinion** [35424425; 35345043]
6. Stamina: **official** that low stamina means a player tires sooner (keepers included, since
   they tire even without running) [35597407; 35598946]. Pressing "Committed" uses more stamina
   [manual A04]. Whether tiredness lowers *all* other skills is **contested**. Many say yes, with
   "a 9-ball striker with 3 stamina plays like 4–5 late on" [35675259; 35709160], but one
   poster calls this misinformation [35711025].
7. **Diminishing returns at the top are widely perceived.** The step from 7 to 10 balls feels
   small ("almost no difference between 7 and 10"), in both the old and new sims. **consensus
   (as player perception)** [34829306; 35371425; 35252444; 35368493; 34858233]. Staff say each
   keeping step matters ("10 ball keeper still a lot better than 1 ball … and each step in
   between") [35073428]. **official**
8. **Hard floor (official): any goalkeeper gets at least 4 balls of keeping.** "If your player
   has '0' balls in keeping, the Sim gives them 4 free balls" [35248205, 2011-10-12, test sim]
9. **Training boost (official):** while training a skill, a player gets up to +1 ball in that
   skill in matches. The boost is smaller the more balls he has, and a 10-ball player may not
   notice it [34507883]. This is also evidence that the gain per ball shrinks near the top.
10. Partial skills count: "Players benefit from each part of their skill, not only the full skill
    balls" [manual A05]. **official**. So a continuous skill value beats an integer one.
11. **Not used by the sim (official): foot, height, weight** [35382524, 2011-10-28]
12. **Team context (official):** "Weaker/younger players will receive a small boost when playing
    in a strong team" [35567056, 2011-11-23]. A player's output is therefore not purely his own
    skills.
13. **Player value ≠ playing strength.** Value sets wages and tax. It weights skills unevenly:
    Passing, Tackling, BC and Keeping raise value a lot, while Speed, Stamina and PI raise it
    little. **consensus** [34522360; 34527629; 36301421; 34526404; 36300430; 36301019]. Staff:
    "There is more to this game than team value" [35302436]. **official**
14. "Total balls" (the sum of skills) is the common community heuristic. Many believed the old sim
    effectively worked on totals or averages and that positions barely mattered (keeper/striker
    swap experiments) [33956388; 33957206; 33954233]. For the new sim, others say individual,
    position-relevant skills matter [35560227]. **contested**
15. Old-sim balance stated by staff: **"5/95 coach/player"**, meaning players' skills dominated
    tactics [35207160]. **official**

## Findings

### 1. Overall skill importance rankings

Old sim (before 2011-10-18):
- "The most important skills are Speed, Stamina and ball control" [post 34003124, 2011-05-18, neilc], from a widely used youth guide, "Rolling 9". **consensus**
- "speed and stamina are the most important skills on this game" [post 35001651, 2011-09-09, steviedee]; "In previous sim the most important skills was speed and stamina" [post 35205487, 2011-10-06, shadowwolf666]. **consensus**
- Training order "Speed, Stamina and Ball Control, then PI and the positional skill" [post 34828513, 2011-08-18, richardmbutler]; "tackling, speed and stamina always come before [aerial pass]" [post 33883982, 2011-05-04, alx460]. **consensus**
- "Imo ball control is the most important skill … average of 6-7 in BC you will never make it in the top divisions" [post 34810289, 2011-08-16, superdad123]. **opinion**
- PI is seen as an "amplifier" or complement skill, "much like Form and Exp" [post 34811106, 2011-08-16, jaysargey]; "PI is somewhat of a compliment skill … however it is just as vital" [post 34827274, sikent]. High PI helps a player pick up the ball on the run and beat defenders [post 34810976, mosmosmos]. **consensus (qualitative)**
- Heading: "Aerial passing is not completely useless, heading however is" [post 33886074, 2011-05-05, gebruikersnaam_]. Set plays: "still a tremendous waste of time" [post 33883982, alx460]. "I only train players in set plays if they have maxed in all the other useful skills" [post 33895270, wanttostrokemypuppy]. **consensus (old sim)**, with a dissent: "Heading is actually important as is PI … no one really trains heading" [post 33889174, jamie1976]. **contested, minority**
- Speed vs PI: "The 3 difference in speed is a lot more significant than the PI" [post 33871237, 2011-05-03, tosspot]; "speed is more important than PI in that position" [post 33873056, jaysargey]. **consensus (old sim)**
- Staff on the old sim: "The current Sim is 5/95 coach/player" [post 35207160, 2011-10-07, patrick7777]. **official**. A user's estimate was "5 tactics/65 skills and 30 luck" [post 35238908, trisi]. **opinion**

New sim (H2H test versions mid-2011, live from 2011-10-18):
- "Air passing and heading are now far more important skills than before" [post 34779168, 2011-08-12, patrick7777]. **official**
- "Reduced the weighting of the shooting skill for all players by 50%" (H2H V1.49) [post 35073428, 2011-09-19, patrick7777]. **official** (a relative change only, no absolute weight)
- "Players will head the ball more accurately, based on their heading ability"; "improved the accuracy of passing, based on the passing ability" [post 34449451, 2011-07-08, patrick7777]. **official**
- Wingers should have "pace, stamina, ball control, aerial passing/passing" [post 34451557, patrick7777]. **official**
- "Speed and stamina seem less important than before … heading, play intelligence and aerial passing seem much more important" [post 35405495, 2011-10-31, innuendo]. **consensus** (also [post 35205487, shadowwolf666]; [post 35365975, gorgio] "speed is not as dominant a factor anymore")
- Speed for strikers: "To say Speed is useless is wrong … it just isn't as effective as it used to be" [post 35542098, 2011-11-19, mosmosmos]. A striker with speed 6 "performs equal to" one with speed 9, other skills equal [post 35542375, tomamientras] (**opinion**). "A striker with 7/8 speed but 6/7 bc/pi scores more than one with 9 speed and 3/4 bc/pi" [post 35542695, nomercyurz]. **consensus: speed is down but not zero; BC/PI matter for strikers**
- Strikers: "SH, PI, BC, ST, SP in that order" [post 35418681, 2011-11-01, egtor]. **opinion**. "Attackers need BC, Speed, Stamina and PI while shooting seems diminished" [post 35407726, prednibe]. **opinion**
- Better players (faster, better PI) attempt more tackles [post 35172210, 2011-10-02, patrick7777]. **official** (Speed and PI act on involvement, not only on quality)
- Proxy for interest: posts mentioning "heading" rose from 0.53% (pre) to 1.08% (post) of all posts. Mentions of "speed and stamina" fell from 0.33% to 0.27% (my count via SQL LIKE).

### 2. Experience, Form, Age

Manual (current): "Form and Experience can both help to maximise a player's other attributes." Experience "ensures a player is more comfortable on the pitch and less likely to make mistakes" and grows with games played and age [manual A05 Training]. **official**. The wording suggests a multiplier or boost on other skills, but gives no numbers.

Experience:
- Old sim: experience was a big factor. "The difference between 25 y old players and 30y old players with the same skills is just too huge, just because of experience" [post 34648349, 2011-07-28, mwosle]; "experience skill is too important in the current sim and older players are too strong" [post 35188147, clabbem]; "massive experience bonus = buy old players regardless" [post 35589537, mosmosmos]; "grossly overweighted 'bonus' skills (Experience)" [post 35550601, mosmosmos]. **consensus**
- Counterpoint: "I also don't think there's that much difference between 8 and 9/10 experience" [post 34650714, 2011-07-28, civilu]. **opinion**. This suggests saturation at the top.
- New sim: "Yes, the effect of experience has been reduced in this Simulator, but also used in different areas than before. For example, experience now gives players a chance to avoid receiving cards by reducing the number of bad tackles they make" [post 35367521, 2011-10-25, patrick7777]. **official**. Also "their intelligence and experience will allow them to avoid making tackles that would have resulted in a card" [post 35194786, patrick7777]. **official**
- "experience is less important on this SIM" [post 35534472, 2011-11-18, darkline]; "an inexperienced player is more prone to get a card … so experience is still important, it just gets a different kind of influence" [post 35586019, darkline]. **consensus**
- Speculative multiplier: Experience at "around 0.05 - 0.15X" per ball versus Form at 0.25X per ball. The author admits the "actual figures … appear to be much lower" [post 35709160, 2011-12-15, matttombstone]. **opinion, do not use as fact**

Form:
- "I think significant difference only if 3 or more balls difference in form, so 8 to 10 relatively equal" [post 34221116, 2011-06-14, coldfire01]; reaching 10 form "never made any significant difference" [post 34216764, john_b_wirral]; "no real difference … with 10 rather than 9" [post 34225242, iahawkeye53]. **consensus (weak, anecdotal, old sim)**
- "form is vital" [post 33768716, 2011-04-23, bluedean]. **opinion (old sim)**
- After S40: "the difference between the different skills … is that low that form really doesn't seem to matter lots of times" [post 35424425, 2011-11-02, baltovenaar]; "What happened for the 'form' skill? That was a very important skill earlier" [post 35345043, gorgio]. **opinion**
- Staff: "a top team that doesn't pay any attention to tactics or form or experience can LOSE to a weaker team" [post 35163795, 2011-10-01, patrick7777]. **official (qualitative: Form and Exp matter)**
- Form mechanics (games played, tired/sleep status) are in [manual A05]. Form Package sets Form to 7 or 9 depending on the source [manual A05 vs beginners_guide manager_basics] (the two sources disagree).

Age:
- Manual: Speed and Stamina deteriorate from age 30 and fastest at 34. Other skills can still be gained. A maxed player who starts deteriorating "will quickly lose a ball … [but] will performance almost as good" [manual A05]. **official**. Age itself is not described as a match factor, only through deterioration and Experience.
- "at the moment, the older the player, the more total balls" [post 35589537, mosmosmos]. **opinion**

### 3. Stamina / fatigue

- Manual: high stamina players "are less likely to tire towards the end of games and less prone to injury" [manual A05]. Committed pressing uses more stamina [manual A04]. **official**
- "A keeper with 1 ball in stamina will soon tire" [post 35597407, 2011-11-28, patrick7777]; "Keepers need stamina, they move and run, they become tired. The lower their stamina, the greater the risk they will become tired or exhausted"; his keepers have 5–7 stamina, "not a magic number" [post 35598946, patrick7777]. **official**
- "a crew member said that stamina is tied to the distance the player travels" (old) [post 35597963, darkline]. **opinion (second-hand)**. Keepers' shot-stopping visibly drops late in games when stamina is low [post 35616125, hanzinho]. **opinion**
- Effect on other skills: "Lower stamina means other skills deteriorate faster as the game goes on … a 9 ball striker with a 3 stamina would be closer to a 4 or 5 ball striker in the second half" [post 35675259, 2011-12-10, tree_amigos]; "commonly acknowledged and observable … players with lower stamina play worse the longer they are on the pitch" [post 35682044, wanttostrokemypuppy]. **consensus**. Disputed by "posts stating that stamina effects all attributes [are] misinformation" [post 35711025, cyberman]. Overall **contested, majority yes**
- Detailed model (**opinion**, [post 35709160, matttombstone]): 0 stamina tires in about 15 minutes and 10 stamina lasts the full match plus extra time. Tiredness reduces sprint speed, PI and all technical skills ("10 tackling will become 9…"). With Committed pressing a 5-stamina player tires at 50–60 minutes; with Normal, 6 stamina lasts to about 80 minutes. A central midfielder needs about 8 stamina and a full-back about 6 under Committed. Half-time gives partial recovery.
- New sim complaint: "They dont really run as they should nor do they get tired" [post 35571131, 2011-11-23, bema]; low-stamina players "run pretty damn good the whole game" (test sim) [post 35195012, capcash]. **opinion**. Stamina may have mattered less right after S40.

### 4. Existing evaluation methods

- **Player value**: recalculated each season, and salaries are based on it [manual A03]. The sale price cap is 3x value [manual A06]. Value tracks some skills much more than others: "Passing has a much larger effect on value than BC, Speed or Stamina" [post 34522360, mosmosmos]; "keeping has probably the highest value per ball" [post 34527629, jamie1976]; "Skills like speed and Play intelligence does not pump up the value very much. Skills like tackling, Ball control and passing do make a significant difference" [post 36301421, 2012-02-28, baltovenaar]. **consensus**. A better indicator for well-rounded top players than for lopsided low-level ones [post 34522360]. "Value pretty useless" [post 36301019, mrbling]; a 750k player "is rubbish. Slow, can't pass…" [post 36297332, mosmosmos]. **consensus: value is a poor proxy for playing strength, since it under-weights Speed, Stamina and PI**
- Staff: "If team value was the only factor, we could add up the $$ value … and declare the winner right away" [post 35302436, patrick7777]. **official**. A user model: the higher-value squad wins a single match about 75–80% of the time [post 35324604, adv41449]. **opinion**
- **Total skill balls**: an MZ transfer-search field ("Total Skill Balls") exists [post 35575537, 2011-11-24, blackwater, Crew]. Community thresholds include "18 & 19 y/old … with at least 30 total" [post 34009640, darkline]; U21 scouting asked for "Speed 5/6, Stamina 5/6, Combined amount of balls 30+" [post 34574740, tonnerwhu]; a "good" youth has 23 total balls "all on the right places" [post 36176796, darkline].
- **Position minimums posted (NT selection, opinion)**: "minimum 9 on primary skills for field players. 10 for keeper"; stamina 7–8; CM speed ≥7; strikers and defenders speed 9–10; wingers high speed + AP, with strikers heading 5+ [post 36078858, 2012-01-31, sephiroth]. U21 USA: "Speed 6/7 Stamina 5/6 BC 5/6 Shooting 8/9 Tackling 8/9 Keepers 8/9" [post 34577478, kbssa-costa]; mids "Speed 6 Stamina 6 PI 3/4" [post 34609500].
- No numeric weight formulas for positions were found in the DB. Posters say only Crew knows the internals and that Crew never confirms them [posts 33901527, 35682044, 35675664].
- **Match ratings / MVP**: these exist only in the current manual. The rating is based on visible statistics, team performance, the player's actions and position, is open-ended and can be negative. The MVP is the highest-rated player [manual B01 Fixtures and Results]. Not discussed in the 2011 forum.
- **Scout reports / player evaluation tool**: not found in the DB ("scout" hits are only NT scouting threads). One wish for "buying a scout report" [post 36026344, 2012-01-25, mawi] suggests the feature did not exist then.
- **Third-party tools named**: mzstuff.net/.com "Skiller" (skill card images, estimated partial balls) [posts 34114187, 35675136]; mzplayer.se (player cards); mzplus.com.ar (training camp calculator, player cards) [post 35241091]; MZ Assistant, MZManager / "MZZ Manager" [posts 34526177, 35887925]; CoZone addon [topic 9098097]; mzstats.net (stats graphs) [post 34311400]. None is described as computing a player rating.

### 5. Interaction effects / non-linearity

- Diminishing returns above about 7 balls: "the difference between ball 7 and 8 (or 8 to 9 or 9 to 10) is less than the difference between balls up to 7" [post 34829306, 2011-08-18, civilu]. He recalls this from an older official statement; patrick denies authoring it [post 34829566]. **consensus as perception (old sim); source unverified**
- New sim: "almost no difference between 7 and 10 IMO … not that much gap between a 4 shooting striker and a 10 shooting striker" [post 35371425, 2011-10-26, darkline]; "There just isn't a big enough difference between 3 balls and 10 balls, and that is in EVERY skill" [post 35368493, mosmosmos]; "almost no difference between having strikers of 10 or 7 balls" [post 35252444, neilillo]; "No noticeable difference between a 6 and 10 balls in tackling" [post 35157670, trisi]. **consensus that the skill spread was compressed right after release**. Counter: darkline has "data … that proves the difference between a 10, 9 & 8 ball striker" [post 35534858]. Patrick then promised "a significant difference between the best and worst strikers" after shooting bug fixes [post 35266669, 2011-10-13]. **official**
- A wish for "skill-evaluation exponential instead of linear as it feels like it is today" [post 35221028, capcash]. **opinion** (a sign that it *felt* linear to concave)
- "little to no difference between 7 balls with 99% training and 8 balls maxed" [post 34858233, mosmosmos]. This is consistent with the manual saying partial skills count. **consensus**
- Floors: the keeper minimum of 4 is **official** [35248205]. "Keepers are over-weighted, even poor ones" [post 35562356, mosmosmos]. **opinion**
- Skills interact: BC alone does not evade defenders and needs PI and Speed [posts 34810546, 34810976]. A 10 BC attacker against a 10 BC defender is "50-50" [post 34810289]. This implies relative (attacker vs defender) skill comparisons in duels. **opinion/consensus**
- Lopsided players: "10 Passing and little else can really throw the balance" for value [post 34522360]. Speed-only strikers no longer work in the new sim [post 35542098]. **consensus**: well-rounded > one maxed skill.
- Old sim totals vs individual skills: keeper/striker swap tests gave "more or less the same" results [posts 33957206, 33959210]; "averages are more important than individual balls" [post 33954233, civilu]. New sim: a striker in goal plays as a 4-ball keeper plus his secondaries; a keeper up front plays as a striker "as much as his stats allow" [post 35560227, mosmosmos]. **contested (old sim), official-backed for new sim** (patrick: "play best in the position(s) that match their skills" [post 35567056])

## Open questions / contradictions

- No numeric weights or multipliers from staff for Form or Experience. The only numbers (0.25x and 0.05–0.15x per ball) are one poster's admitted guess.
- Whether stamina/fatigue scales all skills down (majority belief) or only physical ones (cyberman). Also how much stamina mattered right after S40, when posters said players stopped tiring.
- Diminishing returns above 7: widely perceived in both sims, but it is unclear how much of this is real concavity versus sim randomness or the compressed spread in early S40 versions (later bug fixes widened the spread).
- Value under-weights Speed, Stamina and PI. If the sim weights them highly (old sim), value and strength diverge most for fast, high-stamina players.
- The "small boost" for weak players in strong teams (official) and the "superboost" for weaker teams (user complaint, old sim) mean individual evaluation cannot fully predict team output.
- The Form Package level is 7 in manual A05 and 9 in beginners_guide/manager_basics.
- Everything here predates later sim changes. Match ratings and scout reports must come from the manual or newer sources.

## Queries used

- FTS: `"most important skill" OR "most important skills"`, `"speed is"`, `"play intelligence" NEAR(important,8)`, `experience NEAR(important,10)`, `"form is" OR "form affects" OR "form matters"`, `stamina NEAR(tired,10)`, `stamina NEAR(minute*,10)`, `NEAR(useless heading|"set plays"|"aerial passing"|speed, 8-12)`, `"total balls" OR "total skill" OR "skill sum" OR "total skills"`, `NEAR(value skills,6)`, `"player value" OR "players value"`, `skiller OR mzstuff OR mzplayer OR mztools`, `calculator OR evaluator OR spreadsheet`, `formula`, `NEAR("7 balls" difference,15)`, `NEAR(ball 10 difference,6)`, `diminishing`, `NEAR(experience bonus,8)`, `superboost`, `NEAR(minimum balls,8)`, `scout*`, `MVP`, `linear OR exponential`, `"0 form" OR "low form"`.
- SQL: all posts by patrick7777 grepped for skill terms; Crew-role posts (`users.roles LIKE '%Crew%'`: blackwater, mikelony, marcus, thunderexpress, powdersnow). Nothing in them about skills except blackwater's "Total Skill Balls" search field.
- Topics read (partly): 9383947 "Bye bye experience", 9674840 (stamina), 9181750 "Aerial Pass & Set Plays", 9634719 "Speed is now useless in strikers", 9595256 "New Sim Think Tank", 9430344 (ball control), 9352904 "Players Value", 9842140 (value), 9273461 "10 form", 9200785 "Switched keeper and forwards", 9263270 "Team value in SIM?".
- Manual: A05 Training (skills, Form/Exp, partial bar, ageing), A03 Squad (values), A04 Tactics (pressing/stamina), B01 (Rating System, MVP), A06 (value cap).
