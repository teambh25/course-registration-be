# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a course registration backend API built with Go, Gin web framework, GORM, and PostgreSQL (Supabase). It provides endpoints for students to register for courses and administrators to manage students and courses. The system uses session-based authentication with role-based access control (Admin and Student roles).

## Development Commands

### Setup and Running
```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go
# Server runs on port 3000 by default
```

### Testing
- **Claude는 테스트를 직접 실행하지 말 것** - 사용자가 직접 실행함
```bash
go test ./...
go test -v ./...
go test -v ./_test
go test -v ./_test -run TestAPI
```

### Build
```bash
go build -o course-reg main.go
```

### Load Testing
- **Claude는 로드 테스트를 직접 실행하지 말 것** - 사용자가 직접 실행함
```bash
cd test/performance
python3 generate_test_data.py --num_students 1000 --num_courses 100
python3 register_test_data.py
locust -f locustfile.py
```

## Architecture

### Project Structure

The codebase follows a clean architecture pattern with clear separation of concerns:

```
internal/
├── app/
│   ├── handler/      # HTTP handlers (presentation layer)
│   ├── service/      # Business logic layer
│   ├── repository/   # Data access layer
│   ├── models/       # Domain models (GORM entities)
│   ├── middleware/    # HTTP middleware (auth, CORS)
│   ├── routers/      # Route definitions
│   └── domain/       # Domain-specific logic
│       ├── cache/         # In-memory cache (NOT thread-safe)
│       ├── constants/     # Constants (user roles, course status)
│       ├── dto/           # Data transfer objects
│       ├── e/             # Custom error definitions
│       ├── export/        # Static file export (courses.json)
│       ├── registration/  # Registration state management
│       └── worker/        # Concurrent request processing with goroutines
└── pkg/
    ├── database/     # Database setup and connection
    ├── setting/      # Configuration management (env vars)
    ├── utils/        # Utility functions (schedule parsing, time conflicts)
    ├── session/      # Session management
    └── file/         # File operations (JSON utilities)
```

### Dependency Injection Pattern

The application uses manual dependency injection initialized in main.go:

1. **Database** → **Repositories** → **EnrollmentWorker** → **Services** → **Handlers** → **Router**
2. Dependencies flow from main.go downward through constructor functions
3. Each layer depends only on interfaces defined in the layer below

Example flow:
```
app.NewApplication(cfg)
  → database.Setup() creates *gorm.DB
  → repository.New*Repository(db) creates repositories
  → export.ExportCoursesToJson(courseRepo) exports initial courses.json
  → worker.NewEnrollmentWorker(queueSize, enrollRepo) creates worker
  → loadRegistrationState(regConfigRepo) creates registration state
  → service.New*Service(repos..., enrollWorker, regState, ...) creates services
  → handler.New*Handler(services...) creates handlers
  → routers.InitRouter(runMode, sessionKey, handlers) wires up HTTP routes

Note: Worker is NOT automatically started at startup. Admin must call StartRegistration() endpoint,
which loads students/courses/enrollments from DB and calls enrollmentWorker.Start(...).
```

### Concurrency Architecture (CRITICAL)

**EnrollmentWorker** (`internal/app/domain/worker/`):
- Single goroutine processes all enrollment requests sequentially via channel
- Buffered channel with queue size 1000
- **Thread-safety**: 쓰기 연산은 worker goroutine만 cache에 접근 (lock 불필요)
- 읽기 연산(`GetAllCourseStatus`)은 atomic counter를 사용하므로 채널 없이 직접 호출 가능
- Request types: `ENROLL` (구현), `READ_ALL` (TODO: admin 수강 신청 현황 조회용), `CANCEL`/`ADMIN_ENROLL`/`ADMIN_CANCEL` (미구현)
- 에러 처리: `domain/e` 패키지의 커스텀 에러 반환 (`ErrCourseNotFound`, `ErrCourseFull` 등)
- **Initialization pattern**:
  - `NewEnrollmentWorker(queueSize, enrollRepo)` - Sets immutable dependencies only
  - `Start(students, courses, enrollments)` - Creates channel and cache, spawns goroutine
  - `Stop()` - Closes channel, waits for goroutine, sets requestChan to nil
  - Can be restarted: Stop() then Start() creates fresh channel and cache

**EnrollmentCache** (`internal/app/domain/cache/enrollment.go`):
- 쓰기는 worker goroutine에서만 접근, 읽기는 atomic counter(`atomic.Int32`)로 thread-safe
- `EnrolledCount`, `WaitingCount`: `atomic.Int32` 기반 - 채널 없이 직접 읽기 가능
- `StudentCourses`, `ConflictGraph` 등 map 기반 필드: worker goroutine에서만 접근
- Precomputes time conflicts between all course pairs at startup for O(1) lookups
- **Initialization**: `NewEnrollmentCacheWithData(students, courses, enrollments)` returns `(*EnrollmentCache, error)`
  - Private methods: `loadInitStudents()`, `loadInitCourses()`, `loadEnrollments()`, `buildConflictGraph()`
  - Cache is always in complete, valid state after construction

**Request/Response Pattern**:
```go
// Main thread sends request
err := enrollmentWorker.Enroll(studentID, courseID)
  ↓
// Create request with response channel
req := EnrollmentRequest{
  Type: ENROLL,
  StudentID: studentID,
  CourseID: courseID,
  Response: make(chan error, 1),
}
  ↓
// Send to worker via channel
w.requestChan <- req
  ↓
// Worker processes sequentially
for req := range w.requestChan {
  err := w.processEnroll(req)  // Access cache here
  req.Response <- err
}
  ↓
// Main thread receives response (blocks until worker responds)
return <-req.Response  // NOTE: No timeout currently - blocks indefinitely
```

**Why This Matters**:
- Services call `enrollmentWorker.Enroll()` which sends requests via channel
- Worker goroutine is the ONLY goroutine that writes to the cache
- Atomic counters allow lock-free reads from any goroutine (`GetAllCourseStatus`)

### Authentication and Authorization

- **Session-based authentication** using `gin-contrib/sessions` with in-memory store (`memstore`)
- Session에 `role` 저장 (1: admin, 2: student), session key: `course_reg_session`
- **Two user roles** defined in `internal/app/domain/constants/role.go`:
  - `RoleAdmin` (value: 1): Full access to admin endpoints
  - `RoleStudent` (value: 2): Limited access to course registration
- **Admin credentials** are stored in environment variables (`SECRET_ADMIN_ID`, `SECRET_ADMIN_PW`)
- **Student authentication** uses phone number as username and birth date as password (stored in Student model)
- **Middleware**:
  - `middleware.Auth()`: Requires any authenticated user
  - `middleware.AuthAdmin()`: Requires admin role specifically

### Key Models

- **Student**: Uses phone number (unique) for authentication, stores name and birth date
- **Course**: Has name, instructor, description, schedules (text field), capacity, and is_special flag
  - Schedules format: `"월 09:10~11:30, 수 17:10~19:20"` (Korean day names with time ranges)
  - Schedule parsing in `internal/pkg/utils/schedule_parser.go`
- **Enrollment**: Relationship table with student-course relationship, waitlist flag, and position

### Conflict Detection System

**Schedule Parsing** (`internal/pkg/utils/schedule_parser.go`):
- Parses Korean schedule strings: `"월 09:10~11:30, 수 17:10~19:20"`
- Extracts day of week (월=Monday, 화=Tuesday, etc.) and time ranges
- Converts to minute-based representation for efficient comparison

**Conflict Graph** (Precomputed at startup):
- O(n²) build time: Compare all course pairs for time conflicts
- O(1) lookup during enrollment: `ConflictGraph[courseID][otherCourseID]`
- Two courses conflict if they share any overlapping time on the same day
- Built in `EnrollmentCache.BuildConflictGraph()` called during worker initialization

**Enrollment Flow with Conflict Check**:
1. Validate student exists in cache (`StudentExists`) — session 기반이라 실패 시 내부 에러
2. Validate course exists in cache (`CourseExists`)
3. Check time conflicts: For each enrolled course, check `ConflictGraph[newCourse][enrolledCourse]`
4. Check if student already enrolled in course
5. Check if course has capacity available
6. If all checks pass → Insert to DB → update cache → return success
7. Else → Return appropriate error message

**DB/Cache 일관성**: DB write first, then cache update 순서.
- DB 실패 → cache 미변경 → 일관성 유지
- DB 성공 후 cache 업데이트 → validation을 통과했으므로 panic 없음
- DB 성공 후 서버 crash → 재시작 시 `Start()`에서 DB 기준으로 cache 재구축

### Configuration

The application reads configuration from environment variables (`.env` file via `godotenv`):
- `APP_*`: Logging configuration
- `SERVER_*`: Server settings (port, timeouts, run mode)
- `SECRET_*`: Admin credentials and session key
- `DATABASE_*`: Database connection (URL, pool size, connection lifetime)

Configuration is loaded via `setting.Load()` in main.go. See `.env` for all available variables.
- 영구 디스크가 없는 배포 환경이라 파일 로깅 제거, 표준 `log` 패키지로 콘솔 출력만 사용

### API Routes

Base path: `/api/v1`

**Authentication** (`/api/v1/auth`):
- POST `/login` - Login with credentials
- POST `/logout` - Logout current session
- GET `/check` - Check authentication status

**Admin** (`/api/v1/admin`) - Requires AuthAdmin middleware:
- GET `/registration/state` - Get registration enabled status
- POST `/registration/start` - Start registration (loads data into worker cache)
- POST `/registration/pause` - Pause registration
- PUT `/registration/period` - Set registration period
- GET `/registration/period` - Get registration period
- Setup (`/api/v1/admin/setup`):
  - POST `/students/register` - Bulk register students
  - DELETE `/students/reset` - Delete all students
  - POST `/courses` - Create a single course
  - DELETE `/courses/:course_id` - Delete a course
  - POST `/courses/register` - Bulk register courses
  - DELETE `/courses/reset` - Delete all courses
  - DELETE `/enrollments/reset` - Delete all enrollments

**Courses** (`/api/v1/courses`) - Requires AuthUser middleware:
- GET `/` - Get all courses (served from static/courses.json)
- GET `/status` - Get all course status (available/waitlist/full)

**Course Registration** (`/api/v1/course-reg`) - Requires AuthStudent middleware:
- POST `/enrollment` - Enroll in course (goes through EnrollmentWorker)

### Database

- **PostgreSQL (Supabase)** - 배포 환경에 영구 디스크가 없어서 외부 DB 사용
- **GORM** ORM with auto-migration on startup
- `PrepareStmt: false` - prepared statement 캐싱이 연결 레벨이라 단일 워커 구조에서 오히려 악영향
- Migrations run in `database.Setup()` for Student, Course, Enrollment, RegistrationConfig models
- Connection pool: `MaxIdleConns == MaxOpenConns` (환경 변수 `DATABASE_POOL_SIZE`로 설정)
- **Warmup**: 수강 신청 시작 시 `WarmupConnectionPool()`로 커넥션을 미리 생성하여 초기 트래픽 지연 방지

### Static File Caching

- `export.ExportCoursesToJson()`로 전체 과목을 `static/courses.json`에 export
- `GET /api/v1/courses/`에서 static file로 직접 서빙
- 과목 변경 시 자동 재생성, mutex로 thread-safe

### Registration Control System

**State** (`internal/app/domain/registration/registration.go`):
- Thread-safe component (uses `sync.RWMutex`) that controls whether enrollment is enabled
- Key methods:
  - `IsEnabled()`: Check if registration is currently enabled
  - `ChangeEnabledAndAct(enabled, act)`: 상태 변경 + 액션을 atomic하게 실행 (Start/Pause에 사용)
  - `RunIfEnabled(enabled, act)`: 특정 상태일 때만 액션 실행 (admin setup 작업에 사용, `TryRLock` 사용)

**Registration Lifecycle**:
1. Admin calls `POST /admin/registration/start`
2. `AdminService.StartRegistration()` calls `regState.ChangeEnabledAndAct(true, ...)`
3. Setup function: warmup → load students/courses/enrollments from DB
4. `enrollmentWorker.Start(students, courses, enrollments)` spawns worker goroutine
5. Worker builds conflict graph and begins processing enrollment requests
6. Enabled 상태가 DB(`RegistrationConfig`)에 저장되어 재시작 시 복원

## Development Notes

### Adding New Endpoints

1. Define interface method in `internal/app/service/service.go`
2. Implement business logic in appropriate service file (e.g., `service/admin.go`)
3. Add repository method if database access needed in `repository/repository.go` and implementation
4. Create handler method in appropriate handler file (e.g., `handler/admin.go`)
5. Register route in `internal/app/routers/router.go`
6. Apply appropriate middleware (`Auth()` or `AuthAdmin()`)

### Test Conventions

- **공유 테스트 데이터**: helper 함수로 다양한 상태(empty/partial/full)를 미리 구성, `createTestCache` 하나로 재사용
- **table-driven test**: 케이스가 여러 개면 `t.Run` + 테이블로 구조화, 2개 이하면 인라인
- **magic number 제거**: table에서는 필드명이 설명하므로 그대로, 그 외는 역할을 설명하는 지역 변수 사용

## Important Architectural Constraints

### DB/Cache 동기화 이슈
- Admin operations (`CreateCourse`, `DeleteCourse`, `RegisterStudents` 등)은 **DB만** 수정
- EnrollmentWorker cache는 admin 작업 시 갱신되지 않음
- **현재 방식**: Pause → Start로 cache를 DB에서 다시 로드
- **개선안**: worker에 `ADD_COURSE`, `REMOVE_COURSE` 등 request type 추가

### Worker 사용 규칙
- 수강 신청은 반드시 `enrollmentWorker.Enroll()`을 통해 채널로 전달
- `worker.cache`에 직접 접근 금지
- No timeout 구현 - worker가 멈추면 요청이 무한 블록

### Service 독립성
- Service 간 직접 호출 금지, repository나 domain 컴포넌트를 통해 공유
- `export.ExportCoursesToJson(courseRepo)`: 과목 변경 후 서비스에서 직접 호출

### TimeProvider
- `internal/pkg/utils/utils.go`에 `TimeProvider` 인터페이스 정의
- `KoreaTimeProvider`: Asia/Seoul 타임존 사용
- 시간 의존 로직에는 `time.Now()` 대신 주입된 `TimeProvider` 사용

### Known TODOs

1. **Cancel enrollment**: `CANCEL` request type 정의됨, 미구현 (대기열 승격 로직 필요)
2. **Admin force enroll**: `ADMIN_ENROLL` request type 정의됨, 미구현
3. **READ_ALL**: admin 수강 신청 현황 조회용으로 활용 검토 중
4. **Waitlist**: 데이터 구조 존재하나 enrollment flow에서 미사용 (정원 초과 시 거절)
5. **수강 신청 예약**: `startTime`/`endTime` 필드 존재, 미구현