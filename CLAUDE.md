# CLAUDE.md

@.claude/skills/golang-pro.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`tuto-api` is the backend for the Tuto (Lumio) Flutter app at `../tuto/`. Go + chi.

```
cmd/api/            # main entrypoint
internal/auth/      # auth handlers (parent + kid sessions)
internal/kid/       # kid-side handlers (map, lessons, library, badges, me)
internal/contest/   # …
migrations/         # SQL migrations (apply against PostgreSQL)
pkg/                # reusable helpers
schema.dbml         # canonical data model — see "Database schema" below
```

## Database schema

The canonical data model lives in `schema.dbml` (DBML — paste into [dbdiagram.io](https://dbdiagram.io) to visualize). It is derived from every screen in the Flutter app (`_Screen` enum in `../tuto/lib/main.dart:114-140`) plus the notification opt-in screen.

### Screen → table mapping

| Screen | Tables |
|---|---|
| `kickoff` | — (pure routing choice) |
| `pAccount` | `parents` |
| `pKid` | `kids` |
| `pGoals` | `goals` + `parent_goals` |
| `pPlan` | `learning_plans` |
| `pConnect` | `pairing_codes` |
| `kEnterCode` | `pairing_codes` (consumed_at) |
| `kHello` | — (greets using `kids.name`) |
| `kAvatar` | `kids.avatar_emoji` |
| `kInterests` | `interests` + `kid_interests` |
| `kReady` | — (transition screen) |
| `map` | `lessons`, `lesson_progress`, `subjects` |
| `lesson` (intro) | `lessons` |
| `miniGame` | `mini_games`, `lesson_attempts` |
| `quiz` | `quizzes`, `quiz_options`, `lesson_attempts` |
| `reward` | `lesson_progress.stars_earned`, `kid_badges` |
| `badges` (tab) | `badges`, `kid_badges` |
| `library` (tab) | `library_collections`, `library_items`, `library_progress` |
| `me` (tab) | `kids`, `kid_settings` |
| `parentGate` | `parent_gate_attempts` |
| `parentDash` | `daily_activity`, `subject_progress`, `kids` |
| `weekly` | `weekly_reports`, `weekly_report_highlights` |
| `notification_screen` | `kid_settings.notifications_on`, `notifications` |

`kickoff`, `kHello`, `kReady` are intentionally stateless — they render text / route, using data already stored elsewhere. Cross-cutting tables `auth_sessions` and `kid_devices` support sign-in and device pairing that the screens imply but don't show directly.

### Tables (full list, grouped)

**1. Identity & pairing**
| Table | Purpose | Key columns |
|---|---|---|
| `parents` | Parent account (pAccount) | `email` (unique), `password_hash`, `display_name`, `marketing_opt_in` |
| `kids` | Kid profile under a parent | `parent_id → parents`, `name`, `age`, `avatar_emoji`, `level`, `stars`, `hearts`, `streak_days` |
| `pairing_codes` | 6-digit code for pConnect → kEnterCode | `kid_id → kids`, `code` (unique), `expires_at`, `consumed_at` |
| `kid_devices` | Devices a kid is signed in on | `kid_id → kids`, `platform`, `push_token`, `last_seen_at` |
| `parent_gate_attempts` | Math-puzzle audit (parentGate) | `parent_id → parents`, `a`, `b`, `submitted`, `succeeded` |

**2. Onboarding selections**
| Table | Purpose | Key columns |
|---|---|---|
| `goals` | Catalogue of 6 parent goals (pGoals) | `title_key` (unique), `emoji`, `sub_key`, `sort_order` |
| `parent_goals` | M:N parent ↔ goals | PK `(parent_id, goal_id)` |
| `interests` | Catalogue of 9 kid interests (kInterests) | `title_key` (unique), `emoji`, `sort_order` |
| `kid_interests` | M:N kid ↔ interests | PK `(kid_id, interest_id)` |
| `learning_plans` | pPlan output | `kid_id` (unique), `minutes_per_day`, `days_of_week` (bitmask Mon=1…Sun=64, default 31 = Mon–Fri) |
| `kid_settings` | me_screen toggles + push opt-in | PK `kid_id`, `sounds_on`, `bedtime_on`, `notifications_on` |

**3. Content catalogue**
| Table | Purpose | Key columns |
|---|---|---|
| `subjects` | Subject tags (For You / Reading / Math / Science / World / Art / Bonus) | `key` (unique), `emoji`, `sort_order` |
| `lessons` | A MapNode lesson | `subject_id → subjects`, `emoji`, `label_key` (unique), `x_offset`, `sort_order`, `est_minutes`, `stars_reward` |
| `mini_games` | mini_game_screen step inside a lesson | `lesson_id → lessons`, `prompt_key`, `config_json`, `sort_order` |
| `quizzes` | quiz_screen step inside a lesson | `lesson_id → lessons`, `question_key`, `sort_order` |
| `quiz_options` | A/B/C/D options | `quiz_id → quizzes`, `letter`, `label_key`, `is_correct`, unique `(quiz_id, letter)` |
| `library_collections` | Library "row" grouping | `title_key` (unique), `subtitle_key`, `sort_order` |
| `library_items` | Library card | `collection_id → library_collections`, `subject_id`, `title_key`, `level_key`, `minutes` |

**4. Kid progress**
| Table | Purpose | Key columns |
|---|---|---|
| `lesson_progress` | Drives MapNode.state per kid | `(kid_id, lesson_id)` unique, `state` (enum `locked`/`current`/`done`), `stars_earned`, `best_score`, `attempts`, `first_started_at`, `completed_at` |
| `lesson_attempts` | Quiz/mini-game attempt audit | `kid_id`, `lesson_id`, `quiz_id?`, `selected_option_id?`, `is_correct`, `duration_ms`, `attempted_at` |
| `library_progress` | Drives "Continue" rail | PK `(kid_id, library_item_id)`, `in_progress`, `minutes_left`, `finished_at` |
| `badges` | Catalogue of badges (badges_screen) | `name_key` (unique), `emoji`, `sub_key`, `sort_order` |
| `kid_badges` | M:N kid ↔ badges | PK `(kid_id, badge_id)`, `earned_at`, `fresh` (unseen) |
| `daily_activity` | 7-day chart on parentDash | unique `(kid_id, activity_date)`, `minutes`, `lessons_done`, `stars_earned`, `hearts_lost` |
| `subject_progress` | Per-subject rings on parentDash | PK `(kid_id, subject_id)`, `value` (0.0–1.0) |

**5. Parent reporting**
| Table | Purpose | Key columns |
|---|---|---|
| `weekly_reports` | weekly_report_screen header | unique `(kid_id, week_start)`, `minutes_total`, `lessons_total`, `stars_total`, `badges_earned` |
| `weekly_report_highlights` | Bullet highlights | `weekly_report_id → weekly_reports`, `emoji`, `text_key`, `sort_order` |

**6. Notifications & auth**
| Table | Purpose | Key columns |
|---|---|---|
| `notifications` | Outbound push log | `parent_id?`, `kid_id?`, `kind` (enum: `daily_reminder` / `weekly_report` / `badge_earned` / `bedtime` / `parent_summary`), `title`, `body`, `payload_json`, `sent_at`, `read_at` |
| `auth_sessions` | JWT sessions for both actor types | `actor_type` (enum `parent`/`kid`), `actor_id`, `token_hash` (unique), `device_id?`, `expires_at`, `revoked_at` |

### Enums

* `lesson_state` — `locked`, `current`, `done` (used by `lesson_progress`)
* `actor_type` — `parent`, `kid` (used by `auth_sessions`)
* `notification_kind` — `daily_reminder`, `weekly_report`, `badge_earned`, `bedtime`, `parent_summary`

### Conventions

* snake_case columns, plural table names
* uuid PKs for entities; small int PKs only for catalogue tables (`goals`, `interests`, `subjects`)
* `created_at` / `updated_at` audit columns on every mutable entity
* l10n keys stored as `*_key` columns — see "Localization" below

### Localization

Content tables keep the existing l10n keys as columns (`label_key`, `title_key`, `sub_key`, etc.) — these match the keys already used by the Flutter app (`mapNodeWhales`, `subjectScience`, `badgeStarCounter`…), so translations continue to flow through the Flutter ARB files instead of being duplicated in the DB.

### Editing the schema

When the Flutter app gains a new screen or data field, update `schema.dbml` first (it is the source of truth), then generate / hand-write the corresponding SQL migration under `migrations/`. Keep the screen → table table above in sync.
