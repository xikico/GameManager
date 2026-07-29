# AGENTS.md
## Project Overview
- **Tech stack**: Go backend (hexagonal architecture) + Vue3 frontend (HTTP API)
- **Purpose**: Local game management system
- **Root dir**: `E:\Program\GO_Program\GameManager`

## Directory Structure (Hexagonal Architecture)
| Directory | Role |
|-----------|------|
| `/domain/entity` | Core business entities (Game, Category etc.) |
| `/domain/service` | **All business logic resides here** |
| `/domain/ports/in` | Input port interfaces (API boundary, DTO definitions) |
| `/domain/ports/out` | Output port interfaces (DB, utils, unzip etc.) |
| `/domain/utils` | DTO ↔ Domain entity conversion utils |
| `/adapters/in/gin_support` | Gin HTTP input adapter (routes, VO, API handlers) |
| `/adapters/in/utils` | VO ↔ DTO conversion utils |
| `/adapters/out/db/sqlite` | SQLite DB output adapter |
| `/adapters/out/db/fake` | Fake DB for testing |

## Key Commands
- Build: `go build` (run from root)
- No configured test/lint commands as of now

## Non-Negotiable Code Conventions
1. **Game list sorting**: All game list responses are sorted **ONLY in domain/service layer** (not DB, not adapter layer) by `InsertTime` descending (newest first)
2. **Hexagonal architecture rules**: Domain layer **MUST NOT** depend on any adapter implementations, only on port interfaces
3. **Entity conversion flow**: VO (HTTP layer) → DTO (port layer) → Domain entity (core) → DB entity (adapter layer), all conversions use dedicated utils
4. **Adapter layer rules**: Adapters only handle protocol/IO conversion, **NO business logic allowed**
5. Spelling conventions: Use `GetAllCategory` (not `GatAllCategory`), `condition_vo.go` (not `conditon_vo.go`)

## Important Notes
- Main entrypoint: `/main.go`
- DB port interface: `/domain/ports/out/db/db.go`
- HTTP API routes: `/adapters/in/gin_support/app.go`
- Business logic entrypoint: `/domain/service/service.go`
- All game list methods use shared `sortAndConvertToDTO` utility for sorting
